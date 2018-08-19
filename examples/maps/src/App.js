import React, { Component } from 'react'
import Preview from './Preview'

import './styles.scss'

const workerTask = () => {
  Promise
    .all([{
      type: 'buildings',
      color: [150, 150, 150],
    }, {
      type: 'parks',
      color: [210, 210, 210],
    }, {
      type: 'roads',
      color: [240, 240, 240],
    }].map(({ type, color }) => new Promise((resolve, reject) => {
    const request = new XMLHttpRequest()
      request.open('GET', `${SERVER}/${type}_tmp`, true)
      request.responseType = 'arraybuffer'
      request.onload = event => {
        const arrayBuffer = request.response

        if (!arrayBuffer) {
          reject('Array buffer conversion failed')
        }

        resolve({
          type,
          color,
          value: new Float32Array(arrayBuffer),
        })
      }
      request.send()
    })))
    .then(values => postMessage(values))
    .catch(error => postMessage({ error }))
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
      const renderingData = {}
      data.forEach(object => {
        const { color, value } = object
        renderingData[object.type] = { color, value }
      })
      this.setState({ data: renderingData, error: null })
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

    console.log(data)

    return (
      <div className="container">
        <div className="navigation">
          <img src="logo.png" className="logo" />
        </div>
        {data
          ? <Preview data={data} />
          : <div className="loading">loading...</div>}
        {error && <div>{JSON.stringify(error)}</div>}
      </div>
    )
  }
}

export default App
