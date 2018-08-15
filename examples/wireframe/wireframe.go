package main

import (
	"fmt"
	"math"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	. "github.com/tchayen/triangolatte"
)

const (
	windowWidth  = 500
	windowHeight = 500
	previewSize  = 1.5 // Anything in (0, 2]

	vertexShaderSource = `
		#version 410
		in vec3 position;
		// in vec3 barycentric;
		out vec3 vbc;
		void main() {
			vbc = vec3(1, 1, 1);//barycentric;
			gl_Position = vec4(position, 1.0);
		}
` + "\x00"

	fragmentShaderSource = `
		#version 410
		in vec3 vbc;
		out vec4 color;
		void main() {
			// if (any(lessThan(vbc, vec3(0.02)))) {
			// 	color = vec4(0.0, 0.0, 0.0, 1.0);
			// } else {
			// 	color = vec4(0.5, 0.5, 0.5, 1.0);
			// }
			color = vec4(0, 0, 0, 1);
		}
` + "\x00"
)

var (
	window      *glfw.Window
	program     uint32
	vao         uint32
	triangles   []float32
	barycentric []float32
	points      = []Point{{X: 0, Y: 0}, {X: 5, Y: 0}, {X: 2, Y: 2}, {X: 3, Y: 4}, {X: 0, Y: 4}}
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// findMinMax takes array of points and finds min and max coordinates.
func findMinMax(points []Point) (xMin, yMin, xMax, yMax float64) {
	xMin, yMin, xMax, yMax = math.MaxFloat64, math.MaxFloat64, 0.0, 0.0
	for _, p := range points {
		if p.X < xMin {
			xMin = p.X
		}

		if p.X > xMax {
			xMax = p.X
		}

		if p.Y < yMin {
			yMin = p.Y
		}

		if p.Y > yMax {
			yMax = p.Y
		}
	}
	return
}

// toFloat32 takes array of float64 elements and changes them to float32.
func toFloat32(array []float64) []float32 {
	converted := make([]float32, len(array))
	for i, v := range array {
		converted[i] = float32(v)
	}
	return converted
}

// normalize puts array of vertices in range [-1, 1] (stretching to rectangle).
func normalize(points []Point) []Point {
	xMin, yMin, xMax, yMax := findMinMax(points)
	for i := range points {
		points[i] = Point{
			X: ((points[i].X-xMin)/(xMax-xMin))*previewSize - previewSize/2,
			Y: ((points[i].Y-yMin)/(yMax-yMin))*previewSize - previewSize/2,
		}
	}
	return points
}

// triangulate triangules array of points returning array of triangles.
func triangulate() {
	triangulated, err := Polygon(normalize(points))
	check(err)

	triangles = toFloat32(triangulated)
}

func calculateBarycentric() {
	// 2 coordinates times 3 vertices in a triangle.
	n := len(triangles) / 6

	// 3 barycentric vertices times 3 coordinates.
	barycentric = make([]float32, 9*n)

	for i := 0; i < 9*n; i += 9 {
		barycentric[i+0], barycentric[i+1], barycentric[i+2] = 1, 0, 0
		barycentric[i+3], barycentric[i+4], barycentric[i+5] = 0, 1, 0
		barycentric[i+6], barycentric[i+7], barycentric[i+8] = 0, 0, 1
	}
}

// initGlfw calls GLFW startup functions and does some initialization.
func initGlfw() {
	err := glfw.Init()
	check(err)

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

// initWindow creates new window.
func initWindow() {
	var err error
	window, err = glfw.CreateWindow(
		windowWidth,
		windowHeight,
		"Triangolatte – wireframe example",
		nil,
		nil,
	)
	check(err)

	window.MakeContextCurrent()
}

// initOpenGL initializes OpenGL.
func initOpenGL() {
	err := gl.Init()
	check(err)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL version: %s\n", version)

	gl.ClearColor(1, 1, 1, 1.0)
}

// createProgram links two shaders into a program.
func createProgram() {
	// compileShader takes shader source and compiles it.
	compileShader := func(source string, shaderType uint32) (uint32, error) {
		shader := gl.CreateShader(shaderType)

		csources, free := gl.Strs(source)
		gl.ShaderSource(shader, 1, csources, nil)
		free()
		gl.CompileShader(shader)

		var status int32
		gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

			return 0, fmt.Errorf("Failed to compile %v: %v", source, log)
		}

		return shader, nil
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	check(err)
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	check(err)

	program = gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
}

// makeVao creates a Vertex Array Object from given array of points.
func makeVao() {
	getBuffer := func(vbo *uint32, data []float32) {
		gl.GenBuffers(1, vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, *vbo)
		gl.BufferData(
			gl.ARRAY_BUFFER,
			4*len(data), // Size of buffer in bytes (float32 has 4 bytes).
			gl.Ptr(data),
			gl.STATIC_DRAW,
		)
	}

	getAttribute := func(vbo uint32, index uint32) {
		gl.EnableVertexAttribArray(index)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.VertexAttribPointer(index, 2, gl.FLOAT, false, 0, nil)
	}

	var trianglesVbo uint32
	getBuffer(&trianglesVbo, triangles)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	getAttribute(trianglesVbo, 0)
}

// draw takes vao, count of vertices to draw and program to use.
func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangles)/2))

	glfw.PollEvents()
	window.SwapBuffers()
}

// cleanUp cleans just before program exit.
func cleanUp() {
	glfw.Terminate()
}

func main() {
	runtime.LockOSThread()

	initGlfw()
	initWindow()
	initOpenGL()
	triangulate()
	calculateBarycentric()
	createProgram()
	makeVao()

	for !window.ShouldClose() {
		draw()
	}

	cleanUp()
}
