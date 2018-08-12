import React, { Component } from 'react'

import * as webgl from './webgl'
import * as scene from './scene'

let gl = null // It doesn't hurt to make it a bit global.

class Preview extends Component {
  constructor(props) {
    console.log('Preview_constructor')
    super(props)

    this.width = window.innerWidth
    this.height = window.innerHeight
    this.scalingFactor = window.devicePixelRatio || 1

    if (!gl) {
      const canvas = webgl.setUpCanvas(this.width, this.height, this.scalingFactor)
      gl = canvas.getContext('webgl')
    }

    this.constants = {
      type: gl.FLOAT,
      // Normalization means translating value in any type to [-1.0, 1.0] range
      // based on the range this given type has.
      normalize: false,
      // Start at the beginning of the buffer.
      offset: 0,
      // 2 components per iteration, i.e. for
      // a {x, y, z, w} vector we provide only {x, y}, z
      // and w will default to 0 and 1 respectively.
      size: 2,
      // move forward size * sizeof(type) each iteration to get the next position
      stride: 0,
      arrayOffset: 0,
      primitiveType: gl.TRIANGLES,
    }
  }

  componentDidMount() {
    if (this.state.initialized) return

    console.log('gl_componentDidMount')

    const vertexShader = webgl.createShader(
      gl, gl.VERTEX_SHADER,
      require('./shaders/vertex.glsl')
    )

    const fragmentShader = webgl.createShader(
      gl, gl.FRAGMENT_SHADER,
      require('./shaders/fragment.glsl')
    )

    this.program = webgl.createProgram(gl, vertexShader, fragmentShader)

    // Normalize data.
    for (let b = 0; b < this.props.triangleData.buildings.length; b++) {
      if (this.props.triangleData.buildings[b] === null) continue

      for (let i = 0; i < this.props.triangleData.buildings[b].length; i += 2) {
        this.props.triangleData.buildings[b][i] *= this.width * 0.8
        this.props.triangleData.buildings[b][i + 1] *= this.height * 0.8
      }
    }

    scene.setup(gl, this.program)
    this.setState({ initialized: true })
  }

  shouldComponentUpdate(nextProps) {
    console.log('Preview_shouldComponentUpdate')
    const next = nextProps.triangleData.selected
    const current = this.props.triangleData.selected

    if (next !== current) return true
    return false
  }

  render() {
    console.log('Preview_render')
    const { selected, buildings } = this.props.triangleData
    const objects = scene.setBuffers(gl, [new Float32Array(buildings[selected])])
    scene.draw(gl, this.program, objects, this.constants)

    return <div></div>
  }
}

export default Preview
