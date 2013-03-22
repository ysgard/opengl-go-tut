/*
Orthographic Cube
*/

package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"github.com/go-gl/glfw"
	"os"
	"runtime"
	"unsafe"
)

const (
	Width  = 500
	Height = 500
	Title  = "Orthographic Cube"
)

// Shader program and uniorms
var currentShader gl.Uint
var offsetUniform gl.Int

// Shader ilenames
var shaders = []string{
	"shaders/orthographic.vertexshader",
	"shaders/orthographic.fragmentshader",
}

// Vertex arrays and buffers
var vertexBufferObject gl.Uint
var vao gl.Uint

var vertexData = []gl.Float{
	0.25, 0.25, 0.75, 1.0,
	0.25, -0.25, 0.75, 1.0,
	-0.25, 0.25, 0.75, 1.0,

	0.25, -0.25, 0.75, 1.0,
	-0.25, -0.25, 0.75, 1.0,
	-0.25, 0.25, 0.75, 1.0,

	0.25, 0.25, -0.75, 1.0,
	-0.25, 0.25, -0.75, 1.0,
	0.25, -0.25, -0.75, 1.0,

	0.25, -0.25, -0.75, 1.0,
	-0.25, 0.25, -0.75, 1.0,
	-0.25, -0.25, -0.75, 1.0,

	-0.25, 0.25, 0.75, 1.0,
	-0.25, -0.25, 0.75, 1.0,
	-0.25, -0.25, -0.75, 1.0,

	-0.25, 0.25, 0.75, 1.0,
	-0.25, -0.25, -0.75, 1.0,
	-0.25, 0.25, -0.75, 1.0,

	0.25, 0.25, 0.75, 1.0,
	0.25, -0.25, -0.75, 1.0,
	0.25, -0.25, 0.75, 1.0,

	0.25, 0.25, 0.75, 1.0,
	0.25, 0.25, -0.75, 1.0,
	0.25, -0.25, -0.75, 1.0,

	0.25, 0.25, -0.75, 1.0,
	0.25, 0.25, 0.75, 1.0,
	-0.25, 0.25, 0.75, 1.0,

	0.25, 0.25, -0.75, 1.0,
	-0.25, 0.25, 0.75, 1.0,
	-0.25, 0.25, -0.75, 1.0,

	0.25, -0.25, -0.75, 1.0,
	-0.25, -0.25, 0.75, 1.0,
	0.25, -0.25, 0.75, 1.0,

	0.25, -0.25, -0.75, 1.0,
	-0.25, -0.25, -0.75, 1.0,
	-0.25, -0.25, 0.75, 1.0,
	// The Color data
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,

	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,

	0.8, 0.8, 0.8, 1.0,
	0.8, 0.8, 0.8, 1.0,
	0.8, 0.8, 0.8, 1.0,

	0.8, 0.8, 0.8, 1.0,
	0.8, 0.8, 0.8, 1.0,
	0.8, 0.8, 0.8, 1.0,

	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,

	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,

	0.5, 0.5, 0.0, 1.0,
	0.5, 0.5, 0.0, 1.0,
	0.5, 0.5, 0.0, 1.0,

	0.5, 0.5, 0.0, 1.0,
	0.5, 0.5, 0.0, 1.0,
	0.5, 0.5, 0.0, 1.0,

	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,

	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,

	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,

	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
}

func glfwInitWindow() {
	// Initialize glfw
	glfw.Init()
	// Set some basic params or the window
	glfw.OpenWindowHint(glfw.FsaaSamples, 4) // 4x antialiasing
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 2)
	// Core, not compat
	glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// Open a window and initialize its OpenGL content
	if err := glfw.OpenWindow(Width, Height, 0, 0, 0, 0, 32, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "glfw.OpenWindow failed: %s\n", err)
		os.Exit(1)
	}

	// Set the Window title
	glfw.SetWindowTitle(Title)

	// Make sure we can capture the escape key
	glfw.Enable(glfw.StickyKeys)
}

// glInit - initialize OpenGL context and vertex arrays
func glInit() {
	gl.Init()

	// Load the shaders
	currentShader = CreateShaderProgram(shaders)
	offsetUniform = gl.GetUniformLocation(currentShader, gl.GLString("offset"))
	gl.UseProgram(currentShader)

}

// Create necessary buffers
func initializeVertexBuffer() {
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	buflen := unsafe.Sizeof(vertexData[0]) * (uintptr)(len(vertexData))
	gl.BufferData(gl.ARRAY_BUFFER,
		gl.Sizeiptr(buflen),
		gl.Pointer(&vertexData[0]),
		gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

// Called after window creation and OpenGL initialization
// Once before main loop
func programInit() {
	initializeVertexBuffer()
	// Init vertex arrays
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Cull faces
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	// Winding order - determine which face is the 'front'
	// of the triangle.  In this case, the vertices were
	// defined in vertexData such on order that a clockwise 
	// winding defines their front.
	gl.FrontFace(gl.CW)

	// Set the key handler for the main loop
	glfw.SetKeyCallback(keyboard)

	// Set the resize handler
	glfw.SetWindowSizeCallback(reshape)
}

// display - render an OpenGL frame, this function should be
// called once per frame
func display() {
	// Clear the background, then draw black
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(currentShader)
	// Pass a vector of two floats to the shader using
	// the offsetUniform handle
	gl.Uniform2f(offsetUniform, 0.5, 0.25)

	// Append the 
	offset := (len(vertexData) / 2) * 4
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, gl.FALSE, 0, nil)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, gl.FALSE, 0,
		gl.Offset(nil, (uintptr)(offset)))

	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.UseProgram(0)

	glfw.SwapBuffers()
}

// Called whenever the window is resized.  The new window size is
// is given in pixels.  This is an opportunity to call glViewport
// or glScissor to keep up with the change in size.
func reshape(w, h int) {
	gl.Viewport(0, 0, (gl.Sizei)(w), (gl.Sizei)(h))
}

func shutdown() {
	// Delete all buffers, shut down glfw
	gl.DeleteBuffers(1, &vertexBufferObject)
	gl.DeleteProgram(currentShader)
	gl.DeleteVertexArrays(1, &vao)
	glfw.Terminate()
}

// Called whenever a key is pressed or released.  Esc terminates the 
// program.
func keyboard(key, state int) {
	if state == glfw.KeyPress {
		switch key {
		case glfw.KeyEsc:
			shutdown()
		}

	}
	return
}

// Main loop
func main() {
	// Sit. Down. Good boy.
	runtime.LockOSThread()

	// Initialize subsystems
	glfwInitWindow()
	glInit()
	programInit()

	// Main loop.  Run until it dies, or we find someone better.
	for glfw.WindowParam(glfw.Opened) == 1 {
		display()
	}

}
