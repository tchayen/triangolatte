import React, { Component } from 'react'

import './styles.scss'

const workerTask = () => {
  fetch(`${SERVER}/data.geojson`)
    .then(value => value.json())
    .then(value => postMessage(value))
    .catch(error => postMessage({ error: 'data download failed' }))
}

class App extends Component {
  constructor(props) {
    super(props)
  }

  state = {
    error: null,
    data: null,
  }

  workerMessageHandler = event => {
    const { data } = event

    if (data.error) {
      this.setState({ data: null, error: data.error })
    } else {
      this.setState({ data, error: null })
    }
  }

  async componentDidMount() {
    const workerBlob = URL.createObjectURL(new Blob(
      ['(', workerTask.toString(), ')()'],
      { type: 'application/javascript' },
    ))
    const worker = new Worker(workerBlob)
    worker.onmessage = this.workerMessageHandler
    URL.revokeObjectURL(workerBlob)
  }

  render() {
    const { data, error } = this.state

    return (
      <div className="container">
        <div className="navigation">
          <img src="logo.png" className="logo" />
        </div>
        {data && <div>{JSON.stringify(data)}</div>}
        {error && <div>{JSON.stringify(error)}</div>}
      </div>
    )
  }
}

export default App
