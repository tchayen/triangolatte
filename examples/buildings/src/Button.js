import React, { Component } from 'react'

const noop = () => {}

class Button extends Component {
  constructor(props) {
    super(props)
  }

  state = { waiting: false }

  handleAction = async action => {
    this.setState({ waiting: true })
    // Artificial timeout to show loading animation.
    setTimeout(() => {
      this.setState({ waiting: false })
      action()
    }, 300)
  }

  renderSpinner = () => <div className="spinner"><div></div><div></div><div></div><div></div></div>

  render() {
    const { label, action, waiting } = this.props

    const classes = waiting ? ['waiting', 'button'] : ['button']

    return (
      <div
        className={`${classes.join(' ')}`}
        onClick={!waiting
          ? () => this.handleAction(action)
          : noop
        }
      >
        {this.state.waiting ? this.renderSpinner() : label}
      </div>
    )
  }
}

export default Button
