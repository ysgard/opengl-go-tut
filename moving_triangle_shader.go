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
)

const (
	Width  = 1024
	Height = 768
	Title  = "Moving Triangle"
)

// Default background color is black
var bgRed, bgGreen, bgBlue, bgAlpha gl.Float = 0.0, 0.0, 0.0, 0.0

// Current Shader program, set to nothing initially, and uniform shader vars
var currentShader gl.Uint
var elapsedTimeUniform gl.Int
var loopDurationUnf gl.Int

// Shader filenames
var shaders = []string{
	"shaders/moving_triangle_2.vertexshader",
	"shaders/moving_triangle.fragmentshader",
}

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

	// Load the shaders
	currentShader = CreateShaderProgram(shaders)
	elapsedTimeUniform = gl.GetUniformLocation(currentShader, gl.GLString("time"))
	loopDurationUnf = gl.GetUniformLocation(currentShader, gl.GLString("loopDuration"))
	fragLoopDurUnf := gl.GetUniformLocation(currentShader, gl.GLString("fragLoopDuration"))
	gl.UseProgram(currentShader)
	gl.Uniform1f(loopDurationUnf, 5.0)
	gl.Uniform1f(fragLoopDurUnf, 5.0)
	gl.UseProgram(0)
}

// displayWindow - render an OpenGL frame, this function should be called
// from the main loop
func display(positionBuffer gl.Uint, vertexCount gl.Sizei) {
	

	// Set the background
	gl.ClearColor(bgRed, bgGreen, bgBlue, bgAlpha)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Use the currently set shader program.  This usually never changes, 
	// which is why we set it in a global.
	gl.UseProgram(currentShader)
	gl.Uniform1f(elapsedTimeUniform, (gl.Float)(glfw.Time()))

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

	// Draw the vertices
	gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)

	// Second triangle
	gl.Uniform1f(elapsedTimeUniform, (gl.Float)(glfw.Time() + 1.0))
	gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)

	// Disable the vertex attribute arrays, reset shader
	gl.DisableVertexAttribArray(0)
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
	glInit()

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
	vertexPositions := []gl.Float{
		// vertex data in XYZW format
		0.0, 0.5, 0.0, 1.0,
		0.5, -0.366, 0.0, 1.0,
		-0.5, -0.366, 0.0, 1.0,
	}

	// Main loop - run until it dies, or we find something better
	for (glfw.Key(glfw.KeyEsc) != glfw.KeyPress) &&
		(glfw.WindowParam(glfw.Opened) == 1) {

		vertexBuffer, vertexCount := initializeVertexBuffer(vertexPositions)
		display(vertexBuffer, vertexCount)
		gl.DeleteBuffers(1, &vertexBuffer)
	}

	glfw.Terminate()
}
