/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, { enumerable: true, get: getter });
/******/ 		}
/******/ 	};
/******/
/******/ 	// define __esModule on exports
/******/ 	__webpack_require__.r = function(exports) {
/******/ 		if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 			Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 		}
/******/ 		Object.defineProperty(exports, '__esModule', { value: true });
/******/ 	};
/******/
/******/ 	// create a fake namespace object
/******/ 	// mode & 1: value is a module id, require it
/******/ 	// mode & 2: merge all properties of value into the ns
/******/ 	// mode & 4: return value when already ns object
/******/ 	// mode & 8|1: behave like require
/******/ 	__webpack_require__.t = function(value, mode) {
/******/ 		if(mode & 1) value = __webpack_require__(value);
/******/ 		if(mode & 8) return value;
/******/ 		if((mode & 4) && typeof value === 'object' && value && value.__esModule) return value;
/******/ 		var ns = Object.create(null);
/******/ 		__webpack_require__.r(ns);
/******/ 		Object.defineProperty(ns, 'default', { enumerable: true, value: value });
/******/ 		if(mode & 2 && typeof value != 'string') for(var key in value) __webpack_require__.d(ns, key, function(key) { return value[key]; }.bind(null, key));
/******/ 		return ns;
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = "./src/index.js");
/******/ })
/************************************************************************/
/******/ ({

/***/ "./src/WebGLUtils.js":
/*!***************************!*\
  !*** ./src/WebGLUtils.js ***!
  \***************************/
/*! exports provided: setUpCanvas, createShader, createProgram */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, \"setUpCanvas\", function() { return setUpCanvas; });\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, \"createShader\", function() { return createShader; });\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, \"createProgram\", function() { return createProgram; });\n/**\n * Set up <canvas> element for later use\n * @param {Number} width expected width of the viewport\n * @param {Number} height expected height of the viewport\n * @param {Number} scalingFactor targets mostly window.devicePixelRatio\n * @returns {Object} prepared DOM canvas element\n */\nconst setUpCanvas = (width, height, scalingFactor) => {\n    const canvas = document.createElement('canvas');\n    canvas.setAttribute('width', scalingFactor * width);\n    canvas.setAttribute('height', scalingFactor * height);\n    canvas.setAttribute('style', `width: ${width}px; height: ${height}px`);\n    document.body.appendChild(canvas);\n\n    return canvas;\n};\n\n/**\n * Compile shader\n * @param {Object} gl WebGL context\n * @param {Object} type shader type, one of: gl.VERTEX_SHADER, gl.FRAGMENT_SHADER\n * @param {string} source source of the program\n * @returns {Object|undefined} compiled shader on success, does not return on\n * failure\n */\nconst createShader = (gl, type, source) => {\n    const shader = gl.createShader(type);\n    gl.shaderSource(shader, source);\n    gl.compileShader(shader);\n    const success = gl.getShaderParameter(shader, gl.COMPILE_STATUS);\n    if (success) return shader;\n\n    console.error(gl.getShaderInfoLog(shader));\n    gl.deleteShader(shader);\n};\n\n/**\n * Creates final program by combining two shaders\n * @param {Object} gl WebGL context\n * @param {Object} vertexShader compiled vertex shader\n * @param {Object} fragmentShader compiled fragment shader\n * @returns {Object|undefined} linked program on success, does not return on\n * failure\n */\nconst createProgram = (gl, vertexShader, fragmentShader) => {\n    const program = gl.createProgram();\n    gl.attachShader(program, vertexShader);\n    gl.attachShader(program, fragmentShader);\n    gl.linkProgram(program);\n    const success = gl.getProgramParameter(program, gl.LINK_STATUS);\n    if (success) return program;\n\n    console.error(gl.getProgramInfoLog(program));\n    gl.deleteProgram(program);\n};\n\n\n\n//# sourceURL=webpack:///./src/WebGLUtils.js?");

/***/ }),

/***/ "./src/index.js":
/*!**********************!*\
  !*** ./src/index.js ***!
  \**********************/
/*! no exports provided */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _WebGLUtils__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./WebGLUtils */ \"./src/WebGLUtils.js\");\n/* harmony import */ var _matrix__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./matrix */ \"./src/matrix.js\");\n\n\n\nlet positionLocation, resolutionLocation, matrixLocation;\n\n/**\n * A bunch of initialization commands. It's not a real setup, you know...\n * @param {Object} gl WebGL context\n * @param {Object} program linked shaders\n * @param {Object[]} objects objects (flat arrays of vertices) to draw\n * @returns {Object[]} objects containing WebGL buffers and matching vertex arrays\n */\nconst setupScene = (gl, program, objects) => {\n    positionLocation = gl.getAttribLocation(program, 'a_position');\n\n    // Set uniforms.\n    resolutionLocation = gl.getUniformLocation(program, \"u_resolution\");\n    matrixLocation = gl.getUniformLocation(program, \"u_matrix\");\n\n    // Set up data in buffers.\n    const resultObjects = [];\n    objects.forEach(object => {\n        const buffer = gl.createBuffer();\n        gl.bindBuffer(gl.ARRAY_BUFFER, buffer);\n        gl.bufferData(gl.ARRAY_BUFFER, object, gl.STATIC_DRAW);\n        resultObjects.push({\n            buffer,\n            triangles: object\n        });\n    });\n\n    gl.clearColor(0, 0, 0, 0);\n\n    return resultObjects;\n};\n\n/**\n * Draw scene – enable vertex attribute, calculate scale-rotate-translate-projection\n * matrix, call `gl.drawArrays`.\n * @param {Object} gl WebGL context\n * @param {Object} program linked shaders\n * @param {Object} objects objects containg triangle data and initialized buffers\n * @param {Object} constants set of configuration constants to use for rendering\n */\nconst drawScene = (gl, program, objects, constants) => {\n    gl.viewport(0, 0, gl.drawingBufferWidth, gl.drawingBufferHeight);\n    // gl.clear(gl.COLOR_BUFFER_BIT)\n    gl.useProgram(program);\n\n    objects.forEach(object => {\n        const { buffer, triangles } = object;\n\n        gl.uniform2f(resolutionLocation, gl.canvas.width, gl.canvas.height);\n\n        const matrix = _matrix__WEBPACK_IMPORTED_MODULE_1__[\"calculateSRTP\"]([gl.canvas.clientWidth, gl.canvas.clientHeight], [0, 0], [1, 1], 0);\n\n        gl.uniformMatrix3fv(matrixLocation, false, matrix);\n\n        gl.bindBuffer(gl.ARRAY_BUFFER, buffer);\n\n        gl.enableVertexAttribArray(positionLocation);\n        gl.vertexAttribPointer(positionLocation, constants.size, constants.type, constants.normalize, constants.stride, constants.offset);\n\n        gl.drawArrays(constants.primitiveType, constants.arrayOffset, object.triangles.length / 3);\n    });\n};\n\nconst width = window.innerWidth;\nconst height = window.innerHeight;\nconst scalingFactor = window.devicePixelRatio || 1;\n\nconst canvas = _WebGLUtils__WEBPACK_IMPORTED_MODULE_0__[\"setUpCanvas\"](width, height, scalingFactor);\nconst gl = canvas.getContext('webgl');\nif (!gl) throw 'WebGL is not supported';\n\nconst constants = {\n    type: gl.FLOAT,\n    // Normalization means translating value in any type to [-1.0, 1.0] range\n    // based on the range this given type has.\n    normalize: false,\n    // Start at the beginning of the buffer.\n    offset: 0,\n    // 2 components per iteration, i.e. for\n    // a {x, y, z, w} vector we provide only {x, y}, z\n    // and w will default to 0 and 1 respectively.\n    size: 2,\n    // 0 = move forward size * sizeof(type) each iteration to get the next position\n    stride: 0,\n    arrayOffset: 0,\n    primitiveType: gl.LINE_STRIP\n};\n\nconst vertexShader = _WebGLUtils__WEBPACK_IMPORTED_MODULE_0__[\"createShader\"](gl, gl.VERTEX_SHADER, __webpack_require__(/*! ./shaders/vertex.glsl */ \"./src/shaders/vertex.glsl\"));\n\nconst fragmentShader = _WebGLUtils__WEBPACK_IMPORTED_MODULE_0__[\"createShader\"](gl, gl.FRAGMENT_SHADER, __webpack_require__(/*! ./shaders/fragment.glsl */ \"./src/shaders/fragment.glsl\"));\n\nconst program = _WebGLUtils__WEBPACK_IMPORTED_MODULE_0__[\"createProgram\"](gl, vertexShader, fragmentShader);(async () => {\n    const response = await fetch('http://localhost:3000/polygon_tmp');\n    const data = await response.json();\n    const objects = setupScene(gl, program, [data]);\n\n    drawScene(gl, program, objects, constants);\n})();\n\n//# sourceURL=webpack:///./src/index.js?");

/***/ }),

/***/ "./src/matrix.js":
/*!***********************!*\
  !*** ./src/matrix.js ***!
  \***********************/
/*! exports provided: multiply, translation, scaling, rotation, projection, calculateSRTP */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, \"multiply\", function() { return multiply; });\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, \"translation\", function() { return translation; });\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, \"scaling\", function() { return scaling; });\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, \"rotation\", function() { return rotation; });\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, \"projection\", function() { return projection; });\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, \"calculateSRTP\", function() { return calculateSRTP; });\nconst _multiply = (a, b) => {\n    const a00 = a[0 * 3 + 0];const a01 = a[0 * 3 + 1];const a02 = a[0 * 3 + 2];\n    const a10 = a[1 * 3 + 0];const a11 = a[1 * 3 + 1];const a12 = a[1 * 3 + 2];\n    const a20 = a[2 * 3 + 0];const a21 = a[2 * 3 + 1];const a22 = a[2 * 3 + 2];\n    const b00 = b[0 * 3 + 0];const b01 = b[0 * 3 + 1];const b02 = b[0 * 3 + 2];\n    const b10 = b[1 * 3 + 0];const b11 = b[1 * 3 + 1];const b12 = b[1 * 3 + 2];\n    const b20 = b[2 * 3 + 0];const b21 = b[2 * 3 + 1];const b22 = b[2 * 3 + 2];\n\n    return [b00 * a00 + b01 * a10 + b02 * a20, b00 * a01 + b01 * a11 + b02 * a21, b00 * a02 + b01 * a12 + b02 * a22, b10 * a00 + b11 * a10 + b12 * a20, b10 * a01 + b11 * a11 + b12 * a21, b10 * a02 + b11 * a12 + b12 * a22, b20 * a00 + b21 * a10 + b22 * a20, b20 * a01 + b21 * a11 + b22 * a21, b20 * a02 + b21 * a12 + b22 * a22];\n};\n\n/**\n * Multiply `n` matrices of size `3x3`.\n * @param {Number[]} args variable number of matrices\n */\nconst multiply = (...args) => {\n    let matrix = _multiply(args[0], args[1]);\n    let i = 2;\n    while (i < args.length) {\n        matrix = _multiply(matrix, args[i]);\n        i += 1;\n    }\n    return matrix;\n};\n\n/**\n * @param {Number} x pixels\n * @param {Number} y pixels\n * @returns {Number[]} `3x3` translation matrix\n */\nconst translation = (x, y) => [1, 0, 0, 0, 1, 0, x, y, 1];\n\n/**\n * @param {Number} x scaling factor\n * @param {Number} y scaling factor\n * @returns {Number[]} `3x3` scale matrix by given factors\n */\nconst scaling = (x, y) => [x, 0, 0, 0, y, 0, 0, 0, 1];\n\n/**\n * @param {Number} angleInRadians\n * @returns {Number[]} `3x3` rotation matrix\n */\nconst rotation = angleInRadians => {\n    const c = Math.cos(angleInRadians);\n    const s = Math.sin(angleInRadians);\n    return [c, -s, 0, s, c, 0, 0, 0, 1];\n};\n\n/**\n * **Note:** This matrix flips the Y axis so that 0 is at the top.\n * @param {Number} width pixels\n * @param {Number} height pixels\n */\nconst projection = (width, height) => [2 / width, 0, 0, 0, -2 / height, 0, -1, 1, 1];\n\n/**\n * Calculate matrix of scale-rotation-translation-projection.\n * @param {Number[]} size canvas size\n * @param {Number[]} translate 2D translation vector in pixels\n * @param {Number[]} scale 2D scale vector in floats\n * @param {Number} angle in radians\n * @returns {Number[]} scale-rotation-translation-projection matrix\n */\nconst calculateSRTP = (size, translate, scale, angle) => multiply(projection(...size), translation(...translate), rotation(angle), scaling(...scale));\n\n\n\n//# sourceURL=webpack:///./src/matrix.js?");

/***/ }),

/***/ "./src/shaders/fragment.glsl":
/*!***********************************!*\
  !*** ./src/shaders/fragment.glsl ***!
  \***********************************/
/*! no static exports found */
/***/ (function(module, exports) {

eval("module.exports = \"precision mediump float;\\n\\nvarying float v_color;\\n\\nvoid main() {\\n   gl_FragColor = vec4(0, v_color, v_color, 1);\\n}\\n\"\n\n//# sourceURL=webpack:///./src/shaders/fragment.glsl?");

/***/ }),

/***/ "./src/shaders/vertex.glsl":
/*!*********************************!*\
  !*** ./src/shaders/vertex.glsl ***!
  \*********************************/
/*! no static exports found */
/***/ (function(module, exports) {

eval("module.exports = \"attribute vec3 a_position;\\n\\nuniform mat3 u_matrix;\\nvarying float v_color;\\n\\nvoid main() {\\n  gl_Position = vec4((u_matrix * vec3(a_position.xy, 1)).xy, 0, 1);\\n  v_color = a_position.z;\\n}\\n\"\n\n//# sourceURL=webpack:///./src/shaders/vertex.glsl?");

/***/ })

/******/ });