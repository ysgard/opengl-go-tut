/*
A simple example of how to use glfw and gogl to open a window and deliver
an OpenGL context to it.
*/
package main

import (
	"fmt"
	//gl "github.com/chsc/gogl"
	"github.com/jteeuwen/glfw"
	"os"
)

const (
	Title  = "Tutorial 01"
	Width  = 800
	Height = 600
)

func main() {
	// Always call init first
	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot init: glfw: %s\n", err)
		return
	}

	// Set Window hints - necessary information before we can
	// call the underlying OpenGL context.
	glfw.OpenWindowHint(glfw.FsaaSamples, 4)        // 4x antialiasing
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3) // OpenGL 3.2
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 2)
	// We want the new OpenGL
	glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// Open a window and initialize its OpenGL context
	if err := glfw.OpenWindow(Width, Height, 0, 0, 0, 0, 32, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "OpenWindow failed: glfw: %s\n")
		return
	}
	// Go idiom - this function will be called right after its enclosing function (main)
	// finishes
	defer glfw.CloseWindow()

	// Set the Window title
	glfw.SetWindowTitle(Title)

	// Make sure we can capture the escape key
	glfw.Enable(glfw.StickyKeys)

	// Main loop - run until it dies, or we find something better
	for (glfw.Key(glfw.KeyEsc) != glfw.KeyPress) && 
		(glfw.WindowParam(glfw.Opened) == 1) {
		glfw.SwapBuffers()
	}

}
