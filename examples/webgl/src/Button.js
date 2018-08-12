import React, { Component } from 'react'

const noop = () => {
  console.log('noop')
}

class Button extends Component {
  constructor(props) {
    super(props)
  }

  state = { waiting: false }

  handleAction = async (action, postAction) => {
    this.setState({ waiting: true })
    await action()

    setTimeout(() => {
      this.setState({ waiting: false })
      postAction()
    }, 3000)
  }

  renderSpinner = () => <div className="spinner"><div></div><div></div><div></div><div></div></div>

  render() {
    const { label, classes, action, postAction, waiting } = this.props

    const applyClasses = ['button', ...classes]
    if (waiting) applyClasses.push('waiting')

    return (
      <div
        className={`${applyClasses.join(' ')}`}
        onClick={!waiting
          ? () => this.handleAction(action, postAction)
          : noop
        }
      >
        {this.state.waiting ? this.renderSpinner() : label}
      </div>
    )
  }
}

export default Button
