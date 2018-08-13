import React, { Component } from 'react'
import Button from './Button'

class Panel extends Component {
  constructor(props) {
    super(props)
  }

  render() {
    const { previous, next, waiting } = this.props

    return (
      <div className="panel">
        <div className="buttons">
          <Button
            label="Previous"
            action={previous}
            waiting={waiting}
          />
          <h1>OSM buildings preview</h1>
          <Button
            label="Next"
            action={next}
            waiting={waiting}
          />
        </div>
      </div>
    )
  }
}

export default Panel
