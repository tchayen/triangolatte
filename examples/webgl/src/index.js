import * as webgl from './webgl'
import * as matrix from './matrix'

let positionLocation, resolutionLocation, matrixLocation

/**
 * A bunch of initialization commands. It's not a real setup, you know...
 * @param {Object} gl WebGL context
 * @param {Object} program linked shaders
 * @param {Object[]} objects objects (flat arrays of vertices) to draw
 * @returns {Object[]} objects containing WebGL buffers and matching vertex arrays
 */
const setupScene = (gl, program, objects) => {
    positionLocation = gl.getAttribLocation(program, 'a_position')

    // Set uniforms.
    resolutionLocation = gl.getUniformLocation(program, "u_resolution")
    matrixLocation     = gl.getUniformLocation(program, "u_matrix")

    // Set up data in buffers.
    const resultObjects = []
    objects.forEach(object => {
        const buffer = gl.createBuffer()
        gl.bindBuffer(gl.ARRAY_BUFFER, buffer)
        gl.bufferData(gl.ARRAY_BUFFER, object, gl.STATIC_DRAW)
        resultObjects.push({
            buffer,
            triangles: object,
        })
    })

    gl.clearColor(0, 0, 0, 0)

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
const drawScene = (gl, program, objects, constants) => {
    gl.viewport(0, 0, gl.drawingBufferWidth, gl.drawingBufferHeight)
    // gl.clear(gl.COLOR_BUFFER_BIT)
    gl.useProgram(program)

    objects.forEach(object => {
        const { buffer, triangles } = object

        gl.uniform2f(resolutionLocation, gl.canvas.width, gl.canvas.height)

        const m = matrix.calculateSRTP(
            [gl.canvas.clientWidth, gl.canvas.clientHeight], [0, 0], [1, 1], 0)

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
            object.triangles.length / 2,
        )
    })
}

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
    const response = await fetch('http://localhost:3000/polygon_tmp')
    const data = await response.json()

    const objects = setupScene(gl, program, [new Float32Array(data)])

    drawScene(gl, program, objects, constants)
})()
