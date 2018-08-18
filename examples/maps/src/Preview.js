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
    const { vertices } = this.props
    for (let i = 0; i < vertices.length; i += 2) {
      vertices[i] *= this.width
      vertices[i + 1] *= this.width
    }

    // Setup scene.
    scene.setup(gl, this.program)
  }

  shouldComponentUpdate(nextProps) {
    return false
  }

  render() {
    const { vertices } = this.props
    // const vertices = [0.15068493783473969, -0.534246563911438, 0.5890411138534546, -0.8082191944122314, 1, -0.232876718044281, 0.15068493783473969, -0.534246563911438, 1, -0.232876718044281, 0.5068492889404297, -0.2876712381839752, 0.15068493783473969, -0.534246563911438, 0.5068492889404297, -0.2876712381839752, 0.2876712381839752, 0.15068493783473969, 0.15068493783473969, -0.534246563911438, 0.2876712381839752, 0.15068493783473969, -0.12328767031431198, -0.6164383292198181, -0.6438356041908264, -0.4794520437717438, -0.6164383292198181, 0.013698630034923553, -0.8904109597206116, -0.5616438388824463, -0.6438356041908264, -0.4794520437717438, -0.8904109597206116, -0.5616438388824463, -0.34246575832366943, -1, -0.34246575832366943, -1, 0.15068493783473969, -0.534246563911438, -0.12328767031431198, -0.6164383292198181, -0.34246575832366943, -1, -0.12328767031431198, -0.6164383292198181, -0.6438356041908264, -0.4794520437717438]
    const triangles = new Float32Array(vertices)
    const objects = scene.setBuffers(gl, [triangles])

    scene.draw(gl, this.program, objects, this.constants)

    return null
  }
}

export default Preview
