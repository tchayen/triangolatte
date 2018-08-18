import React, { Component } from 'react'
import Preview from './Preview'

import './styles.scss'

const workerTask = () => {
  const request = new XMLHttpRequest()
  request.open('GET', `${SERVER}/data_tmp`, true)
  request.responseType = 'arraybuffer'

  request.onload = event => {
    const arrayBuffer = request.response

    if (!arrayBuffer) {
      postMessage({ error: 'Array buffer conversion failed' })
      return
    }

    postMessage({ value: new Float32Array(arrayBuffer) })
  }
  request.send()
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
      this.setState({ data: data.value, error: null })
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
          ? <Preview vertices={data} />
          : <div>loading...</div>}
        {error && <div>{JSON.stringify(error)}</div>}
      </div>
    )
  }
}

export default App
