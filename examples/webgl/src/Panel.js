import React, { Component } from 'react'
import Button from './Button'

class Panel extends Component {
  constructor(props) {
    super(props)
  }

  render() {
    return (
      <div className="panel">
        <div className="buttons">
          {this.props.buttons.map((props, i) =>
            <Button
              {...props}
              key={`button-${i}`}
            />)}
        </div>
      </div>
    )
  }
}

export default Panel
