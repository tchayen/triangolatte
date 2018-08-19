import * as matrix from './matrix'

/**
 * A bunch of initialization commands. It's not a real setup, you know...
 * @param {Object} gl WebGL context
 * @param {Object} program linked shaders
 */

let positionLocation, resolutionLocation, colorLocation, matrixLocation

const setup = (gl, program) => {
  positionLocation = gl.getAttribLocation(program, 'a_position')

  // Set uniforms.
  resolutionLocation = gl.getUniformLocation(program, "u_resolution")
  colorLocation      = gl.getUniformLocation(program, "u_color")
  matrixLocation     = gl.getUniformLocation(program, "u_matrix")

  // Set up data in buffers.
  gl.clearColor(0, 0, 0, 0)
}

/**
 * Sets data in GL buffers.
 * @param {Object} gl
 * @param {Object[]} objects objects (flat arrays of vertices) to draw
 * @returns {Object[]} objects containing WebGL buffers and matching vertex arrays
 */
const setBuffers = (gl, objects) => {
  const resultObjects = []
  objects.forEach(object => {
    const buffer = gl.createBuffer()
    gl.bindBuffer(gl.ARRAY_BUFFER, buffer)
    gl.bufferData(gl.ARRAY_BUFFER, object.value, gl.STATIC_DRAW)
    resultObjects.push({
      buffer,
      color: object.color.map(c => c / 255.0),
      triangles: object.value,
    })
  })

  return resultObjects
}

/**
 * Draw scene – enable vertex attribute, calculate scale-rotate-translate-projection
 * matrix, call `gl.drawArrays`.
 * @param {Object} gl WebGL context
 * @param {Object} program linked shaders
 * @param {Object} objects objects containg triangle data and initialized buffers
 * @param {Object} constants set of configuration constants to use for rendering
 */
const draw = (gl, program, objects, constants) => {
  gl.viewport(0, 0, gl.drawingBufferWidth, gl.drawingBufferHeight)
  gl.clear(gl.COLOR_BUFFER_BIT)
  gl.useProgram(program)

  objects.forEach(object => {
    const { buffer, triangles, color } = object

    gl.uniform2f(resolutionLocation, gl.canvas.width, gl.canvas.height)
    gl.uniform3f(colorLocation, ...color)
    const m = matrix.calculateSRTP([gl.canvas.clientWidth, gl.canvas.clientHeight], [0, 0], [1, 1], 0)

    gl.uniformMatrix3fv(matrixLocation, false, m)
    gl.bindBuffer(gl.ARRAY_BUFFER, buffer)
    gl.enableVertexAttribArray(positionLocation)

    gl.vertexAttribPointer(
      positionLocation,
      constants.size,
      constants.type,
      constants.normalize,
      constants.stride,
      constants.offset,
    )

    gl.drawArrays(
      constants.primitiveType,
      constants.arrayOffset,
      triangles.length / 2,
    )
  })
}

export {
  setup,
  setBuffers,
  draw,
}
