/**
 * Set up <canvas> element for later use
 * @param {Number} width expected width of the viewport
 * @param {Number} height expected height of the viewport
 * @param {Number} scalingFactor targets mostly window.devicePixelRatio
 * @returns {Object} prepared DOM canvas element
 */
const setUpCanvas = (width, height, scalingFactor) => {
  const canvas = document.createElement('canvas')
  canvas.setAttribute('width', scalingFactor * width)
  canvas.setAttribute('height', scalingFactor * height)
  canvas.setAttribute('style', `width: ${width}px; height: ${height}px`)
  document.body.appendChild(canvas)

  return canvas
}

/**
 * Compile shader
 * @param {Object} gl WebGL context
 * @param {Object} type shader type, one of: gl.VERTEX_SHADER, gl.FRAGMENT_SHADER
 * @param {string} source source of the program
 * @returns {Object|undefined} compiled shader on success, does not return on
 * failure
 */
const createShader = (gl, type, source) => {
  const shader = gl.createShader(type)
  gl.shaderSource(shader, source)
  gl.compileShader(shader)
  const success = gl.getShaderParameter(shader, gl.COMPILE_STATUS)
  if (success) return shader

  console.error(gl.getShaderInfoLog(shader))
  gl.deleteShader(shader)
}

/**
 * Creates final program by combining two shaders
 * @param {Object} gl WebGL context
 * @param {Object} vertexShader compiled vertex shader
 * @param {Object} fragmentShader compiled fragment shader
 * @returns {Object|undefined} linked program on success, does not return on
 * failure
 */
const createProgram = (gl, vertexShader, fragmentShader) => {
  const program = gl.createProgram()
  gl.attachShader(program, vertexShader)
  gl.attachShader(program, fragmentShader)
  gl.linkProgram(program)
  const success = gl.getProgramParameter(program, gl.LINK_STATUS)
  if (success) return program

  console.error(gl.getProgramInfoLog(program))
  gl.deleteProgram(program)
}

export {
  setUpCanvas,
  createShader,
  createProgram,
}
