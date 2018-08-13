const _multiply = (a, b) => {
  const a00 = a[0 * 3 + 0]; const a01 = a[0 * 3 + 1]; const a02 = a[0 * 3 + 2]
  const a10 = a[1 * 3 + 0]; const a11 = a[1 * 3 + 1]; const a12 = a[1 * 3 + 2]
  const a20 = a[2 * 3 + 0]; const a21 = a[2 * 3 + 1]; const a22 = a[2 * 3 + 2]
  const b00 = b[0 * 3 + 0]; const b01 = b[0 * 3 + 1]; const b02 = b[0 * 3 + 2]
  const b10 = b[1 * 3 + 0]; const b11 = b[1 * 3 + 1]; const b12 = b[1 * 3 + 2]
  const b20 = b[2 * 3 + 0]; const b21 = b[2 * 3 + 1]; const b22 = b[2 * 3 + 2]

  return [
    b00 * a00 + b01 * a10 + b02 * a20,
    b00 * a01 + b01 * a11 + b02 * a21,
    b00 * a02 + b01 * a12 + b02 * a22,
    b10 * a00 + b11 * a10 + b12 * a20,
    b10 * a01 + b11 * a11 + b12 * a21,
    b10 * a02 + b11 * a12 + b12 * a22,
    b20 * a00 + b21 * a10 + b22 * a20,
    b20 * a01 + b21 * a11 + b22 * a21,
    b20 * a02 + b21 * a12 + b22 * a22,
  ]
}

/**
 * Multiply `n` matrices of size `3x3`.
 * @param {Number[]} args variable number of matrices
 */
const multiply = (...args) => {
  let matrix = _multiply(args[0], args[1])
  let i = 2
  while (i < args.length) {
    matrix = _multiply(matrix, args[i])
    i += 1
  }
  return matrix
}

/**
 * @param {Number} x pixels
 * @param {Number} y pixels
 * @returns {Number[]} `3x3` translation matrix
 */
const translation = (x, y) => [1, 0, 0, 0, 1, 0, x, y, 1]

/**
 * @param {Number} x scaling factor
 * @param {Number} y scaling factor
 * @returns {Number[]} `3x3` scale matrix by given factors
 */
const scaling = (x, y) => [x, 0, 0, 0, y, 0, 0, 0, 1]

/**
 * @param {Number} angleInRadians
 * @returns {Number[]} `3x3` rotation matrix
 */
const rotation = (angleInRadians) => {
    const c = Math.cos(angleInRadians)
    const s = Math.sin(angleInRadians)
    return [c,-s, 0, s, c, 0, 0, 0, 1]
}

/**
 * **Note:** This matrix flips the Y axis so that 0 is at the top.
 * @param {Number} width pixels
 * @param {Number} height pixels
 */
const projection = (width, height) => [2 / width, 0, 0, 0, 2 / height, 0, -1, -1, 1]

/**
 * Calculate matrix of scale-rotation-translation-projection.
 * @param {Number[]} size canvas size
 * @param {Number[]} translate 2D translation vector in pixels
 * @param {Number[]} scale 2D scale vector in floats
 * @param {Number} angle in radians
 * @returns {Number[]} scale-rotation-translation-projection matrix
 */
const calculateSRTP = (size, translate, scale, angle) => multiply(
  projection(...size),
  translation(...translate),
  rotation(angle),
  scaling(...scale),
)

export {
  multiply,
  translation,
  scaling,
  rotation,
  projection,
  calculateSRTP,
}
