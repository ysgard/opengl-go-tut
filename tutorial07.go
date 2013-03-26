/*
triangle.go - simple program to follow the tutorial at
http://arcsynthesis.org/gltut
*/
package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"github.com/go-gl/glfw"
	"os"
	"runtime"
	"unsafe"
	//"github.com/Jragonmiris/mathgl"
)

const (
	Width  = 1024
	Height = 768
	Title  = "Tutorial07"
)

// Default background color is black
var bgRed, bgGreen, bgBlue, bgAlpha gl.Float = 0.0, 0.0, 0.0, 0.0

// Current Shader program, set to nothing initially
var currentShader gl.Uint

// Shader filenames
var shaders = []string{
	"shaders/triangle.vertexshader",
	"shaders/triangle.fragmentshader",
}

// Object to load
var objFile = string("art/cylinder.obj")

func glfwInitWindow() {
	// Initialize glfw
	glfw.Init()
	// Set some basic parameters of the Window
	glfw.OpenWindowHint(glfw.FsaaSamples, 4)
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 2)
	// Core, not compat
	glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// Open a window and initialize its OpenGL context
	if err := glfw.OpenWindow(Width, Height, 0, 0, 0, 0, 32, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "OpenWindow failed: glfw: %s\n")
		os.Exit(1)
	}

	// Set the Window title
	glfw.SetWindowTitle(Title)

}

// glInit - initialize OpenGL context and Vertex Arrays
func glInit() {
	gl.Init()

	// Strange Vertex array init I don't understand yet
	var vertexArrayID gl.Uint = 0
	gl.GenVertexArrays(1, &vertexArrayID)
	gl.BindVertexArray(vertexArrayID)
}

// displayWindow - render an OpenGL frame, this function should be called
// from the main loop
func display(positionBuffer, colorBuffer gl.Uint, vertexCount gl.Sizei) {
	// Set the background
	gl.ClearColor(bgRed, bgGreen, bgBlue, bgAlpha)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Use the currently set shader program.  This usually never changes, 
	// which is why we set it in a global.
	gl.UseProgram(currentShader)

	// Bind the vertex array
	gl.BindBuffer(gl.ARRAY_BUFFER, positionBuffer)
	// Vertex position buffer
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(
		0,        // Vertex array to use (Position)
		4,        // number of floats per vertex
		gl.FLOAT, // Type of the value (32-bit float)
		gl.FALSE, // Normalized?
		0,        // Stride
		nil,      // Array buffer offset
	)

	// Bind the color array
	gl.BindBuffer(gl.ARRAY_BUFFER, colorBuffer)
	// Vertex color buffer
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(
		1,        // Vertex array to use (Color)
		4,        // Number of floats per color (R, G, B, A)
		gl.FLOAT, // type of value
		gl.FALSE, // Normalized?
		0,        // stride
		// Array buffer offset.  In this case, the color data is located
		// right after the vertex data. 3(#) * 4(XYZW) * 4(bytes/float) = 48 bytes
		nil,
	)

	// Draw the vertices
	gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)

	// Disable the vertex attribute arrays, reset shader
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.UseProgram(0)

	// Flip the screen
	glfw.SwapBuffers()

}

func initializeVertexBuffer(vertices []gl.Float) (gl.Uint, gl.Sizei) {
	// Create the vertex buffer object
	var buf gl.Uint
	gl.GenBuffers(1, &buf)

	// Now load the buffer with the data
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	bufferLen := unsafe.Sizeof(vertices[0]) * (uintptr)(len(vertices))
	gl.BufferData(
		gl.ARRAY_BUFFER,
		gl.Sizeiptr(bufferLen),
		gl.Pointer(&vertices[0]),
		gl.STATIC_DRAW)

	return buf, (gl.Sizei)(bufferLen)
}

func main() {
	// Sit. Good boy.
	runtime.LockOSThread()

	// Initialize the glfw subsystem and open a Window, then initialize OpenGL context
	glfwInitWindow()
	gl.Init()

	// Strange Vertex array init - required for OpenGL 3.1+
	// Apparently provides performance benefit as well
	var vertexArrayID gl.Uint = 0
	gl.GenVertexArrays(1, &vertexArrayID)
	gl.BindVertexArray(vertexArrayID)
	defer gl.DeleteVertexArrays(1, &vertexArrayID)

	// Make sure we can capture the escape key
	glfw.Enable(glfw.StickyKeys)

	// Change the OpenGL viewport when the window size changes -
	// this scales the contents of the window appropriately.
	glfw.SetWindowSizeCallback(func(w, h int) {
		gl.Viewport(0, 0, (gl.Sizei)(w), (gl.Sizei)(h))
	})

	// Data prep.  Make it a slice
	vertexPositions, _, _ := loadOBJ(objFile)

	vertexColors := []gl.Float{
		// Color data in RGBA format
		1.0, 0.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 1.0,
		0.0, 0.0, 1.0, 1.0,
	}

	// Load the shaders
	currentShader = CreateShaderProgram(shaders)

	// Main loop - run until it dies, or we find something better
	for (glfw.Key(glfw.KeyEsc) != glfw.KeyPress) &&
		(glfw.WindowParam(glfw.Opened) == 1) {

		vertexBuffer, vertexCount := initializeVertexBuffer(vertexPositions)
		vertexColors, _ := initializeVertexBuffer(vertexColors)
		display(vertexBuffer, vertexColors, vertexCount)
		gl.DeleteBuffers(1, &vertexBuffer)
		gl.DeleteBuffers(1, &vertexColors)
	}

	glfw.Terminate()
}
