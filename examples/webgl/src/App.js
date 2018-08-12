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
    waiting: false,
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
        .then(data => postMessage({ buildings: data }))
        .catch(error => postMessage({ error }))
    }

    // Encode the function in a URL since Workers accept only them.
    const workerBlob = URL.createObjectURL(new Blob(
      ['(', task.toString(), ')()'],
      { type: 'application/javascript' },
    ))

    // Create worker with callback to send parsed data.
    const worker = new Worker(workerBlob)
    worker.onmessage = event => {
      if (event.data.error) {
        this.setState({
          error: 'Data download failed. Make sure the server is up.',
        })
      } else {
        this.setState({
          triangleData: {
            selected: 0,
            buildings: event.data.buildings.filter(b => b !== null && b.length > 0),
          }
        })
      }
    }

    // Clean up.
    URL.revokeObjectURL(workerBlob)
  }

  acceptAction = async () => {
    this.setState({ waiting: true })
    fetch('http://localhost:3000/report/', {
      method: 'POST',
    })
  }

  passAction = async () => {
    this.setState({ waiting: true })
  }

  rejectAction = async () => {
    this.setState({ waiting: true })
  }

  next = () => {
    const { selected, ...rest } = this.state.triangleData
    this.setState({
      waiting: false,
      triangleData: {
        selected: selected + 1,
        ...rest,
      },
    })
  }

  animation = () => {
    const { dot, loading, triangleData } = this.state

    // Proceed with animation if loading was faster than 5 dots or the data is
    // not yet here.
    if (dot <= 5 || triangleData.buildings.length === 0) {
      this.setState({ dot: dot + 1 })
      setTimeout(this.animation, 300)
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

  renderError = () => <div className="loading">{this.state.error}</div>

  buttons = [{
    label: 'Incorrect',
    classes: ['incorrect'],
    action: this.rejectAction,
  }, {
    label: 'Not sure',
    classes: ['not-sure'],
    action: this.passAction,
  }, {
    label: 'Correct',
    classes: ['correct'],
    action: this.acceptAction,
  }]

  renderApp = () =>
    <div>
      <Panel
        buttons={this.buttons.map(button => ({
          ...button,
          postAction: this.next,
          waiting: this.state.waiting,
        }))}
      />
      <Preview triangleData={this.state.triangleData} />
    </div>

  render() {
    return this.state.loading
      ? this.state.error
        ? this.renderError()
        : this.renderLoading()
      : this.renderApp()
  }
}

export default App
