/*
Third tutorial in opengl-tutorials.org.

Matrix transformations.
*/
package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"github.com/jteeuwen/glfw"
	"os"
	"unsafe"
	"runtime"
	"github.com/jragonmiris/mathgl"

)

const (
	Title  = "Tutorial 03"
	Width  = 800
	Height = 600
)

const (
	VertexFile = "shaders/simple_transform.vertexshader",
	FragementFile = "shaders/simple_color.fragmentshader"
)


func main() {
	runtime.LockOSThread()
	// Always call init first
	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}

	// Set Window hints - necessary information before we can
	// call the underlying OpenGL context.
	glfw.OpenWindowHint(glfw.FsaaSamples, 4)        // 4x antialiasing
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3) // OpenGL 3.3
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 2)
	// We want the new OpenGL
	glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// Open a window and initialize its OpenGL context
	if err := glfw.OpenWindow(Width, Height, 0, 0, 0, 0, 32, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n")
		return
	}
	defer glfw.Terminate() // Make sure this gets called if we crash.

	// Set the Window title
	glfw.SetWindowTitle(Title)

	// Make sure we can capture the escape key
	glfw.Enable(glfw.StickyKeys)

	// Initialize OpenGL, make sure we terminate before leaving.
	gl.Init()

	// Dark blue background
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	// Load Shaders
	var programID gl.Uint = LoadShaders(
		VertexFile,
		FragementFile)
	gl.ValidateProgram(programID)
	var validationErr gl.Int 
	gl.GetProgramiv(programID, gl.VALIDATE_STATUS, &validationErr)
	if validationErr == gl.FALSE {
		fmt.Fprintf(os.Stderr, "Shader program failed validation!\n")
	}
	
	// Time to create some graphics!  
	var vertexArrayID gl.Uint = 0
	gl.GenVertexArrays(1, &vertexArrayID)
	gl.BindVertexArray(vertexArrayID)
	defer gl.DeleteVertexArrays(1, &vertexArrayID) // Make sure this gets called before main is done

	// Get a handle for our "MVP" uniform
	matrixID := gl.GetUniformLocation(programID, gl.GLString("MVP"));

	// Projection matrix: 45° Field of View, 4:3 ratio, display range : 0.1 unit <-> 100 units
	projection := mathgl.Perspective(45.0, 4.0/3.0, 0.1, 100.0)

	// Camera matrix
	view := mathgl.LookAt(
		4.0, 3.0, 3.0,
		0.0, 0.0, 0.0,
		0.0, 1.0, 0.0)

	// Model matrix: and identity matrix (model will be at the origin)
	model := mathgl.Ident4f() // Changes for each model!

	// Our ModelViewProjection : multiplication of our 3 matrices - remember, matrix mult is other way around
	MVP := projection.Mul4(view).Mul4(model) // projection * view * model



	// An array of 3 vectors which represents 3 vertices of a triangle
	vertexBufferData := [9]gl.Float{	// N.B. We can't use []gl.Float, as that is a slice
		-1.0, -1.0, 0.0,				// We always want to use raw arrays when passing pointers
		1.0, -1.0, 0.0,					// to OpenGL
		0.0, 1.0, 0.0,
	}

	// Time to draw this sucker.
	var vertexBuffer gl.Uint // id the vertex buffer
	gl.GenBuffers(1, &vertexBuffer) // Generate 1 buffer, grab the id
	defer gl.DeleteBuffers(1, &vertexBuffer) // Make sure we delete this, no matter what happens
	// The following commands will talk about our 'vertexBuffer'
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer) 
	// Give our vertices to OpenGL
	// WARNING!  This looks EXTREMELY fragile
	gl.BufferData(
		gl.ARRAY_BUFFER, 
		gl.Sizeiptr(unsafe.Sizeof(vertexBufferData)), // Already pretty bad
		gl.Pointer(&vertexBufferData),  // SWEET ZOMBIE JESUS PLEASE DON'T CRASH MY MACHINE
		gl.STATIC_DRAW)

	// 

	// DEBUG - check MVP array
	for i, val := range MVP {
		fmt.Fprintf(os.Stdout, "%f ", val)
		if (i + 1) % 4 == 0 { fmt.Fprintf(os.Stdout, "\n") }
	}
	fmt.Fprintf(os.Stdout, "\n")
	

	// Main loop - run until it dies, or we find something better
	for (glfw.Key(glfw.KeyEsc) != glfw.KeyPress) && 
		(glfw.WindowParam(glfw.Opened) == 1) {

		// Clear the screen
		gl.Clear(gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT)

		// Want to use our loaded shaders
		gl.UseProgram(programID)

		// Perform the translation of the camera viewpoint
		// by sending the requested operation to the vertex shader
		//mvpm := [16]gl.Float{0.93, -0.85, -0.68, -0.68, 0.0, 1.77, -0.51, -0.51, -1.24, -0.63, -0.51, -0.51, 0.0, 0.0, 5.65, 5.83}
		gl.UniformMatrix4fv(matrixID, 1, gl.FALSE, (*gl.Float)(&MVP[0]))

		// 1st attribute buffer: vertices
		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
		gl.VertexAttribPointer(
			0,			// Attribute 0. No particular reason for 0, but must match layout in shader
			3,			// size
			gl.FLOAT,	// Type
			gl.FALSE,	// normalized?
			0,			// stride
			nil)	// array buffer offset

		// Draw the triangle!
		gl.DrawArrays(gl.TRIANGLES, 0, 3)	// Starting from vertex 0, 3 vertices total -> triangle

		gl.DisableVertexAttribArray(0)



		glfw.SwapBuffers()
	}



}
