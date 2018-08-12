import React, { Component } from 'react'
import Panel from './Panel'
import Preview from './Preview'

import './styles.scss'

class App extends Component {
  constructor(props) {
    super(props)
  }

  state = {
    dot: 0,
    loading: true,
    triangleData: {
      selected: 0,
      buildings: [],
    },
  }

  async componentDidMount() {
    // Begin UI animation.
    this.animation()

    // Create function for background task.
    const task = () => {
      fetch('http://localhost:3000/data')
        .then(response => response.json())
        .then(data => postMessage(data))
    }

    // Encode the function in a URL since Workers accept only them.
    const workerBlob = URL.createObjectURL(new Blob(
      ['(', task.toString(), ')()'],
      { type: 'application/javascript' },
    ))

    // Create worker with callback to send parsed data.
    const worker = new Worker(workerBlob)
    worker.onmessage = event => {
      this.setState({
        triangleData: {
          selected: 0,
          buildings: event.data,
        }
      })
    }

    // Clean up.
    URL.revokeObjectURL(workerBlob)
  }

  animation = () => {
    const { dot, loading, triangleData } = this.state

    // Proceed with animation if loading was faster than 5 dots or the data is
    // not yet here.
    if (dot <= 5 || triangleData.buildings.length === 0) {
      this.setState({ dot: dot + 1 })
      setTimeout(this.animation, 500)
    } else {
      this.setState({ loading: false })
    }
  }

  renderLoading = () =>
    <div className="loading">
      {[0, 1, 2].map(i =>
        <div
          key={`dot${i}`}
          className={`dot ${(this.state.dot % 3) === i ? 'selected' : ''}`}
        >.</div>
      )}
    </div>

  renderApp = () =>
    <div>
      <Panel />
      <Preview triangleData={this.state.triangleData} />
    </div>

  render() {
    return this.state.loading
      ? this.renderLoading()
      : this.renderApp()
  }
}

export default App
