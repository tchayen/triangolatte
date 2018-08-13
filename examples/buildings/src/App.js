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
      selected: null,
      currentId: '',
      buildings: [],
    },
  }

  async componentDidMount() {
    // Begin UI animation.
    this.animation()

    // Create function for background task.
    const task = () => {
      fetch(`${__SERVER__}/api/data`)
        .then(response => response.json())
        .then(data => {
          const shuffle = a => {
            var j, x, i;
            for (i = a.length - 1; i > 0; i--) {
                j = Math.floor(Math.random() * (i + 1))
                x = a[i]
                a[i] = a[j]
                a[j] = x
            }
            return a
          }

          const buildings = shuffle(data)
          postMessage({ buildings })
        })
        .catch(error => {
          console.log(error)
          postMessage({
            error: 'Data download failed. Make sure the server is up.'
          })
        })
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
        this.setState({ error: event.data.error })
      } else {
        const buildings = event.data
          .buildings
          .filter(b =>
            // Has some triangles at all.
            b.triangles !== null && b.triangles.length > 0 &&
            // Has more description than generic @id and "building": "yes".
            Object.values(b.properties).length > 2 &&
            // More complicated than rectangle.
            b.triangles.length > 12
          )

        this.setState({
          triangleData: {
            selected: 0,
            currentId: buildings[0].properties['@id'],
            buildings,
          }
        })
      }
    }

    // Clean up.
    URL.revokeObjectURL(workerBlob)
  }

  respond = async status => {
    const { currentId } = this.state.triangleData
    await fetch(`${__SERVER__}/api/report`, {
      method: 'POST',
      body: JSON.stringify({
        id: currentId,
        status,
      })
    })
  }

  acceptAction = async () => {
    this.setState({ waiting: true })
    this.respond('accepted')
  }

  passAction = async () => {
    this.setState({ waiting: true })
    this.respond('uncertain')
  }

  rejectAction = async () => {
    this.setState({ waiting: true })
    this.respond('rejected')
  }

  next = () => {
    const { selected, buildings } = this.state.triangleData

    const newSelected = (selected + 1) % buildings.length

    this.setState({
      waiting: false,
      triangleData: {
        selected: newSelected,
        currentId: buildings[newSelected].properties['@id'],
        buildings,
      },
    })
  }

  previous = () => {
    const { selected, buildings } = this.state.triangleData
    const n = buildings.length
    const i = selected - 1
    const newSelected = (i % n + n) % n

    this.setState({
      waiting: false,
      triangleData: {
        selected: newSelected,
        currentId: buildings[newSelected].properties['@id'],
        buildings,
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

  renderLabel = () => {
    const { buildings, selected } = this.state.triangleData
    const building = buildings[selected]
    return (
      <pre className="label">
        {JSON.stringify(building.properties, null, 4)}
      </pre>
    )
  }

  renderApp = () =>
    <div className="layout">
      <Panel
        previous={this.previous}
        next={this.next}
        waiting={this.state.waiting}
      />
      <Preview triangleData={this.state.triangleData} />
      {this.renderLabel()}
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
