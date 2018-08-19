import React, { Component } from 'react'

import * as webgl from './webgl'
import * as scene from './scene'

let gl = null // It doesn't hurt to make it a bit global.

class Preview extends Component {
  constructor(props) {
    super(props)

    this.width = window.innerWidth
    this.height = window.innerHeight
    this.scalingFactor = window.devicePixelRatio || 1

    if (!gl) {
      const canvas = webgl.setUpCanvas(
        this.width,
        this.height,
        this.scalingFactor,
      )
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

    const vertexShader = webgl.createShader(
      gl, gl.VERTEX_SHADER,
      require('./shaders/vertex.glsl')
    )

    const fragmentShader = webgl.createShader(
      gl, gl.FRAGMENT_SHADER,
      require('./shaders/fragment.glsl')
    )

    // Compile shaders.
    this.program = webgl.createProgram(gl, vertexShader, fragmentShader)

    // Normalize data.
    const { data } = this.props
    Object.keys(data).forEach(type => {
      for (let i = 0; i < data[type].value.length; i += 2) {
        data[type].value[i] *= this.width
        data[type].value[i + 1] *= this.width
      }
    })

    // Setup scene.
    scene.setup(gl, this.program)
  }

  shouldComponentUpdate(nextProps) {
    return true
  }

  render() {
    const { data } = this.props
    const values = Object.values(data)

    const objects = scene.setBuffers(gl, values)
    scene.draw(gl, this.program, objects, this.constants)

    return null
  }
}

export default Preview
