import * as webgl from './webgl'
import * as scene from './scene'


const width = window.innerWidth
const height = window.innerHeight
const scalingFactor = window.devicePixelRatio || 1

const canvas = webgl.setUpCanvas(width, height, scalingFactor)
const gl = canvas.getContext('webgl')
if (!gl) throw 'WebGL is not supported'

const constants = {
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
  // 0 = move forward size * sizeof(type) each iteration to get the next position
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

const program = webgl.createProgram(gl, vertexShader, fragmentShader)

;(async () => {
  const response = await fetch('http://localhost:3000/data')
  const data = await response.json()

  const triangles = data[1]

  console.log(triangles)

  for (let i = 0; i < triangles.length; i += 2) {
    triangles[i] *= width * 0.8
    triangles[i + 1] *= height * 0.8
  }

  const objects = scene.setup(gl, program, [new Float32Array(triangles)])
  scene.draw(gl, program, objects, constants)
})()
