import React, { Component } from 'react'

import * as webgl from './webgl'
import * as scene from './scene'

let width, height, scalingFactor, canvas, gl, constants, vertexShader, fragmentShader, program

class Preview extends Component {
  constructor(props) {
    super(props)
  }

  componentDidMount() {
    width = window.innerWidth
    height = window.innerHeight
    scalingFactor = window.devicePixelRatio || 1
    canvas = webgl.setUpCanvas(width, height, scalingFactor)
    gl = canvas.getContext('webgl')

    if (!gl) throw 'WebGL is not supported'

    constants = {
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

    vertexShader = webgl.createShader(
      gl, gl.VERTEX_SHADER,
      require('./shaders/vertex.glsl')
    )

    fragmentShader = webgl.createShader(
      gl, gl.FRAGMENT_SHADER,
      require('./shaders/fragment.glsl')
    )

    program = webgl.createProgram(gl, vertexShader, fragmentShader)

    // Normalize data.
    for (let b = 0; b < this.props.triangleData.buildings; b++) {
      for (let i = 0; i < this.props.triangleData.buildings[b]; i += 2) {
        this.props.triangleData.buildings[b][i] *= width * 0.8
        this.props.triangleData.buildings[b][i + 1] *= height * 0.8
      }
    }

    objects = scene.setup(gl, program, [new Float32Array(triangles)])
    scene.draw(gl, program, objects, constants)
  }

  shouldComponentUpdate(nextProps) {
    const next = nextProps.triangleData.selected
    const current = this.props.triangleData.selected

    if (next !== current) return true
    return false
  }

  render() {
    const { triangleData } = this.props
    return (

    )
  }
}

export default Preview
