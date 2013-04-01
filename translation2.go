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
	Title  = "Perspective Prism"
)

// Hardware vars
var mouseWheelPos int
var zOffset = gl.Float(0.2)
var bDepthClamping bool = false

// Shader program and uniorms
var currentShader gl.Uint
var offsetUniform gl.Int
var perspectiveMatrixUnif gl.Int

// Viewport matrix & constants
var theMatrix []gl.Float
var fFrustumScale = gl.Float(1.0)
var fzNear = gl.Float(0.5)
var fzFar = gl.Float(3.0)


// Shader ilenames
var shaders = []string{
	"shaders/standard.vert",
	"shaders/standard.frag",
}

// Vertex arrays and buffers
var vertexBufferObject gl.Uint
var indexBufferObject gl.Uint
var vao gl.Uint

var numberOfVertices = 36

// Some simple definitions
type Color struct {
	R gl.Float
	G gl.Float
	B gl.Float
	W gl.Float
}
var RIGHT_EXTENT = gl.Float(0.8)
var LEFT_EXTENT = -RIGHT_EXTENT
var TOP_EXTENT = gl.Float(0.2)
var MIDDLE_EXTENT = gl.Float(0.0)
var BOTTOM_EXTENT = -TOP_EXTENT
var FRONT_EXTENT = gl.Float(-1.25)
var REAR_EXTENT = gl.Float(-1.75)

var GREEN_COLOR = Color{0.75, 0.75, 1.0, 1.0}
var BLUE_COLOR = Color{0.0, 0.5, 0.0, 1.0}
var RED_COLOR = Color{1.0, 0.0, 0.0, 1.0}
var GREY_COLOR = Color{0.8, 0.8, 0.0, 1.0}
var BROWN_COLOR = Color{0.5, 0.5, 0.0, 1.0}

var indexData = [...]gl.Short{
	0, 2, 1,
	3, 2, 0,

	4, 5, 6,
	6, 7, 4, 

	8, 9, 10,
	11, 13, 12,

	14, 16, 15,
	17, 16, 14,
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
	perspectiveMatrixUnif = gl.GetUniformLocation(currentShader, gl.GLString("perspectiveMatrix"))
	gl.UseProgram(currentShader)

	// Initialize the uniforms
	theMatrix = make([]gl.Float, 16)
	theMatrix[0] = fFrustumScale
	theMatrix[5] = fFrustumScale
	theMatrix[10] = (fzFar + fzNear) / (fzNear - fzFar)
	theMatrix[14] = (2 * fzFar * fzNear) / (fzNear - fzFar)
	theMatrix[11] = -1.0
	// var theMatrix = [16]gl.Float{
	// 	fFrustumScale, 0, 0, 0,
	// 	0, fFrustumScale, 0, 0,
	// 	0, 0, (fzFar + fzNear) / (fzNear - fzFar), -1.0,
	// 	0, 0, (2 * fzFar * fzNear) / (fzNear - fzFar), 0,
	// }

	gl.UniformMatrix4fv(perspectiveMatrixUnif, 1, gl.FALSE, &theMatrix[0])

	gl.UseProgram(0)

	// Depth buffer setup
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(gl.TRUE)
	gl.DepthFunc(gl.LEQUAL)
	gl.DepthRange(0.0, 1.0)
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

	gl.GenBuffers(1, &indexBufferObject)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBufferObject)
	buflen = unsafe.Sizeof(indexData[0]) * (uintptr)(len(indexData))
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 
		gl.Sizeiptr(buflen),
		gl.Pointer(&indexData[0]),
		gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func initializeVertexArrayObjects() {
		// Create the first VAO
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Figure out the offset from the position data to the color data
	colorDataOffset := gl.Offset(nil, unsafe.Sizeof(gl.Float(0)) * (uintptr)(3 * numberOfVertices))

	// Attach attribute pointers to the data
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.EnableVertexAttribArray(0) // vertices
	gl.EnableVertexAttribArray(1) // colors
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, gl.FALSE, 0, colorDataOffset)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBufferObject)
	// Unbind 
	gl.BindVertexArray(0)
}

// Called after window creation and OpenGL initialization
// Once before main loop
func programInit() {

	initializeVertexBuffer()
	initializeVertexArrayObjects()

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

	// Set the mousewheel handler
	glfw.SetMouseWheelCallback(mouseWheelChange)
}

// display - render an OpenGL frame, this function should be
// called once per frame
func display() {
	// Clear the background, then draw black
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(currentShader)
	
	gl.BindVertexArray(vao)
	gl.Uniform3f(offsetUniform, 0.0, 0.0, zOffset)
	gl.DrawElements(gl.TRIANGLES, (gl.Sizei)(len(indexData)), gl.UNSIGNED_SHORT, nil)

	gl.Uniform3f(offsetUniform, 0.0, 0.0, -1.0)
	gl.DrawElementsBaseVertex(
		gl.TRIANGLES, 
		(gl.Sizei)(len(indexData)), 
		gl.UNSIGNED_SHORT, 
		nil,
		(gl.Int)(numberOfVertices / 2))

	gl.BindVertexArray(0)
	gl.UseProgram(0)

	glfw.SwapBuffers()
}

// Called whenever the window is resized.  The new window size is
// is given in pixels.  This is an opportunity to call glViewport
// or glScissor to keep up with the change in size.
func reshape(w, h int) {	
	theMatrix[0] = fFrustumScale / ((gl.Float)(w) / (gl.Float)(h))
	theMatrix[5] = fFrustumScale

	gl.UseProgram(currentShader)
	gl.UniformMatrix4fv(perspectiveMatrixUnif, 1, gl.FALSE, &theMatrix[0])
	gl.UseProgram(0)
	gl.Viewport(0, 0, (gl.Sizei)(w), (gl.Sizei)(h))
}

// Called when the mousewheel is changed.
func mouseWheelChange(pos int) {
	switch delta := mouseWheelPos - pos; {
	case delta < 0:
		if zOffset > -2.0 {
			zOffset -= 0.1 
		}
	case delta > 0:
		if zOffset < 2.0 {
			zOffset += 0.1
		} 
	}
	mouseWheelPos = pos
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
	
		case glfw.KeySpace:
			if bDepthClamping == true {
				gl.Disable(gl.DEPTH_CLAMP)
			} else {
				gl.Enable(gl.DEPTH_CLAMP)
			}
			bDepthClamping = !bDepthClamping
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



var vertexData = []gl.Float{
	//Object 1 positions
	LEFT_EXTENT,	TOP_EXTENT,		REAR_EXTENT,
	LEFT_EXTENT,	MIDDLE_EXTENT,	FRONT_EXTENT,
	RIGHT_EXTENT,	MIDDLE_EXTENT,	FRONT_EXTENT,
	RIGHT_EXTENT,	TOP_EXTENT,		REAR_EXTENT,

	LEFT_EXTENT,	BOTTOM_EXTENT,	REAR_EXTENT,
	LEFT_EXTENT,	MIDDLE_EXTENT,	FRONT_EXTENT,
	RIGHT_EXTENT,	MIDDLE_EXTENT,	FRONT_EXTENT,
	RIGHT_EXTENT,	BOTTOM_EXTENT,	REAR_EXTENT,

	LEFT_EXTENT,	TOP_EXTENT,		REAR_EXTENT,
	LEFT_EXTENT,	MIDDLE_EXTENT,	FRONT_EXTENT,
	LEFT_EXTENT,	BOTTOM_EXTENT,	REAR_EXTENT,

	RIGHT_EXTENT,	TOP_EXTENT,		REAR_EXTENT,
	RIGHT_EXTENT,	MIDDLE_EXTENT,	FRONT_EXTENT,
	RIGHT_EXTENT,	BOTTOM_EXTENT,	REAR_EXTENT,

	LEFT_EXTENT,	BOTTOM_EXTENT,	REAR_EXTENT,
	LEFT_EXTENT,	TOP_EXTENT,		REAR_EXTENT,
	RIGHT_EXTENT,	TOP_EXTENT,		REAR_EXTENT,
	RIGHT_EXTENT,	BOTTOM_EXTENT,	REAR_EXTENT,

	//Object 2 positions
	TOP_EXTENT,		RIGHT_EXTENT,	REAR_EXTENT,
	MIDDLE_EXTENT,	RIGHT_EXTENT,	FRONT_EXTENT,
	MIDDLE_EXTENT,	LEFT_EXTENT,	FRONT_EXTENT,
	TOP_EXTENT,		LEFT_EXTENT,	REAR_EXTENT,

	BOTTOM_EXTENT,	RIGHT_EXTENT,	REAR_EXTENT,
	MIDDLE_EXTENT,	RIGHT_EXTENT,	FRONT_EXTENT,
	MIDDLE_EXTENT,	LEFT_EXTENT,	FRONT_EXTENT,
	BOTTOM_EXTENT,	LEFT_EXTENT,	REAR_EXTENT,

	TOP_EXTENT,		RIGHT_EXTENT,	REAR_EXTENT,
	MIDDLE_EXTENT,	RIGHT_EXTENT,	FRONT_EXTENT,
	BOTTOM_EXTENT,	RIGHT_EXTENT,	REAR_EXTENT,
					
	TOP_EXTENT,		LEFT_EXTENT,	REAR_EXTENT,
	MIDDLE_EXTENT,	LEFT_EXTENT,	FRONT_EXTENT,
	BOTTOM_EXTENT,	LEFT_EXTENT,	REAR_EXTENT,
					
	BOTTOM_EXTENT,	RIGHT_EXTENT,	REAR_EXTENT,
	TOP_EXTENT,		RIGHT_EXTENT,	REAR_EXTENT,
	TOP_EXTENT,		LEFT_EXTENT,	REAR_EXTENT,
	BOTTOM_EXTENT,	LEFT_EXTENT,	REAR_EXTENT,

	//Object 1 colors
	GREEN_COLOR.R, GREEN_COLOR.G, GREEN_COLOR.B, 1.0,
	GREEN_COLOR.R, GREEN_COLOR.G, GREEN_COLOR.B, 1.0,
	GREEN_COLOR.R, GREEN_COLOR.G, GREEN_COLOR.B, 1.0,
	GREEN_COLOR.R, GREEN_COLOR.G, GREEN_COLOR.B, 1.0,

	BLUE_COLOR.R, BLUE_COLOR.G, BLUE_COLOR.B, 1.0,
	BLUE_COLOR.R, BLUE_COLOR.G, BLUE_COLOR.B, 1.0,
	BLUE_COLOR.R, BLUE_COLOR.G, BLUE_COLOR.B, 1.0,
	BLUE_COLOR.R, BLUE_COLOR.G, BLUE_COLOR.B, 1.0,

	RED_COLOR.R, RED_COLOR.G, RED_COLOR.B, 1.0,
	RED_COLOR.R, RED_COLOR.G, RED_COLOR.B, 1.0,
	RED_COLOR.R, RED_COLOR.G, RED_COLOR.B, 1.0,

	GREY_COLOR.R, GREY_COLOR.G, GREY_COLOR.B, 1.0,
	GREY_COLOR.R, GREY_COLOR.G, GREY_COLOR.B, 1.0,
	GREY_COLOR.R, GREY_COLOR.G, GREY_COLOR.B, 1.0,

	BROWN_COLOR.R, BROWN_COLOR.G, BROWN_COLOR.B, 1.0,
	BROWN_COLOR.R, BROWN_COLOR.G, BROWN_COLOR.B, 1.0,
	BROWN_COLOR.R, BROWN_COLOR.G, BROWN_COLOR.B, 1.0,
	BROWN_COLOR.R, BROWN_COLOR.G, BROWN_COLOR.B, 1.0,

	//Object 2 colors
	RED_COLOR.R, RED_COLOR.G, RED_COLOR.B, 1.0,
	RED_COLOR.R, RED_COLOR.G, RED_COLOR.B, 1.0,
	RED_COLOR.R, RED_COLOR.G, RED_COLOR.B, 1.0,
	RED_COLOR.R, RED_COLOR.G, RED_COLOR.B, 1.0,

	BROWN_COLOR.R, BROWN_COLOR.G, BROWN_COLOR.B, 1.0,
	BROWN_COLOR.R, BROWN_COLOR.G, BROWN_COLOR.B, 1.0,
	BROWN_COLOR.R, BROWN_COLOR.G, BROWN_COLOR.B, 1.0,
	BROWN_COLOR.R, BROWN_COLOR.G, BROWN_COLOR.B, 1.0,

	BLUE_COLOR.R, BLUE_COLOR.G, BLUE_COLOR.B, 1.0,
	BLUE_COLOR.R, BLUE_COLOR.G, BLUE_COLOR.B, 1.0,
	BLUE_COLOR.R, BLUE_COLOR.G, BLUE_COLOR.B, 1.0,

	GREEN_COLOR.R, GREEN_COLOR.G, GREEN_COLOR.B, 1.0,
	GREEN_COLOR.R, GREEN_COLOR.G, GREEN_COLOR.B, 1.0,
	GREEN_COLOR.R, GREEN_COLOR.G, GREEN_COLOR.B, 1.0,


	GREY_COLOR.R, GREY_COLOR.G, GREY_COLOR.B, 1.0,
	GREY_COLOR.R, GREY_COLOR.G, GREY_COLOR.B, 1.0,
	GREY_COLOR.R, GREY_COLOR.G, GREY_COLOR.B, 1.0,
	GREY_COLOR.R, GREY_COLOR.G, GREY_COLOR.B, 1.0,
}
