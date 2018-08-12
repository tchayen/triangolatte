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
      const buttons = 100 // Top offset.
      const canvas = webgl.setUpCanvas(
        this.width / 2.0,
        this.height - buttons,
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
    const { buildings } = this.props.triangleData
    for (let b = 0; b < buildings.length; b++) {
      if (buildings[b] === null) continue

      for (let i = 0; i < buildings[b].length; i += 2) {
        buildings[b][i] *= this.width / 2.0
        buildings[b][i + 1] *= this.height / 2.0
      }
    }

    // Setup scene.
    scene.setup(gl, this.program)
  }

  shouldComponentUpdate(nextProps) {
    const next = nextProps.triangleData.selected
    const current = this.props.triangleData.selected

    if (next !== current) return true
    return false
  }

  render() {
    const { selected, buildings } = this.props.triangleData
    const triangles = new Float32Array(buildings[selected])
    const objects = scene.setBuffers(gl, [triangles])
    scene.draw(gl, this.program, objects, this.constants)

    return null
  }
}

export default Preview
