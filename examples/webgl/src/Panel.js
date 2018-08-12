import React, { Component } from 'react'

class Panel extends Component {
  constructor(props) {
    super(props)
  }

  renderButton = (label, classes, action, index) =>
    <div
      className={`${['button', ...classes].join(' ')}`}
      onClick={action}
      key={`button-${index}`}
    >
      {label}
    </div>

  render() {
    return (
      <div className="panel">
        <div className="buttons">
          {this.props.buttons.map((b, i) => this.renderButton(...b, i))}
        </div>
      </div>
    )
  }
}

export default Panel
