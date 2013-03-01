/*
A simple example of how to use glfw and gogl to open a window and deliver
an OpenGL context to it.
*/
package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"github.com/jteeuwen/glfw"
	"os"
	"unsafe"
)

const (
	Title  = "Tutorial 02"
	Width  = 800
	Height = 600
)





func main() {
	// Always call init first
	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}

	LoadShaders("This", "that")

	// Set Window hints - necessary information before we can
	// call the underlying OpenGL context.
	glfw.OpenWindowHint(glfw.FsaaSamples, 4)        // 4x antialiasing
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3) // OpenGL 3.3
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 3)
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
	// Time to create some graphics!  
	var vertexArrayID gl.Uint = 0
	gl.GenVertexArrays(1, &vertexArrayID)
	gl.BindVertexArray(vertexArrayID)
	defer gl.DeleteVertexArrays(1, &vertexArrayID) // Make sure this gets called before main is done

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

	//fmt.Fprintf(os.Stdout, "Got this far.  vertexBufferData length: %d\n", len(vertexBufferData))
	//fmt.Fprintf(os.Stdout, "unsafe.Sizeof(vertexBufferData) - %v \n", unsafe.Sizeof(vertexBufferData))



	// Main loop - run until it dies, or we find something better
	for (glfw.Key(glfw.KeyEsc) != glfw.KeyPress) && 
		(glfw.WindowParam(glfw.Opened) == 1) {

		// Clear the screen
		gl.Clear(gl.COLOR_BUFFER_BIT)

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
