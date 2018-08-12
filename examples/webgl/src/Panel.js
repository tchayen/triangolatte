import React, { Component } from 'react'

class Panel extends Component {
  constructor(props) {
    super(props)
  }

  button = (label = '', classes = [], action = () => {}) =>
    <div className={`${['button', ...classes].join(' ')}`}>
      {label}
    </div>

  render() {
    return (
      <div className="panel">
        <div className="buttons">
          {this.button('Correct', ['correct'])}
          {this.button('Not sure')}
          {this.button('Incorrect', ['incorrect'])}
        </div>
      </div>
    )
  }
}

export default Panel
