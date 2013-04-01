package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"github.com/go-gl/glfw"
	"github.com/Jragonmiris/mathgl"
	"os"
	"runtime"
	"unsafe"
	"math"
)

const (
	Width  = 500
	Height = 500
	Title  = "Perspective Prism"
)

// 
const degToRad = math.Pi * 2.0 / 360
var fFrustumScale = CalcFrustumScale(45.0)

const fLoopDuration = 3.0
const fScale = math.Pi * 2.0 / fLoopDuration

// Various GL
var currentShader gl.Uint
var vertexBufferObject gl.Uint
var indexBufferObject gl.Uint
var vao gl.Uint

// Shader vars
var shaders = []string{
	"shaders/PosColorLocalTransform.vert",
	"shaders/ColorPassthrough.frag",
}

var modelToCameraMatrixUnif gl.Int
var cameraToClipMatrixUnif gl.Int

// camera?
var cameraToClipMatrix mathgl.Mat4f
var fzNear = gl.Float(1.0)
var fzFar = gl.Float(45.0)

// Data definitions
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

// Vertex data
var numberOfVertices = int(8)
var vertexData = []gl.Float{
	+1.0, +1.0, +1.0,
	-1.0, -1.0, +1.0,
	-1.0, +1.0, -1.0,
	+1.0, -1.0, -1.0,

	-1.0, -1.0, -1.0,
	+1.0, +1.0, -1.0,
	+1.0, -1.0, +1.0,
	-1.0, +1.0, +1.0,

	GREEN_COLOR.R, GREEN_COLOR.G, GREEN_COLOR.B, 1.0,
	BLUE_COLOR.R, BLUE_COLOR.G, BLUE_COLOR.B, 1.0,
	RED_COLOR.R, RED_COLOR.G, RED_COLOR.B, 1.0,
	BROWN_COLOR.R, BROWN_COLOR.G, BROWN_COLOR.B, 1.0,

	GREEN_COLOR.R, GREEN_COLOR.G, GREEN_COLOR.B, 1.0,
	BLUE_COLOR.R, BLUE_COLOR.G, BLUE_COLOR.B, 1.0,
	RED_COLOR.R, RED_COLOR.G, RED_COLOR.B, 1.0,
	BROWN_COLOR.R, BROWN_COLOR.G, BROWN_COLOR.B, 1.0,
}

// Index data
var indexData = []gl.Short{
	0, 1, 2,
	1, 0, 3,
	2, 3, 0,
	3, 2, 1, 

	5, 4, 6,
	4, 5, 7,
	7, 6, 4,
	6, 7, 5,
}

type Instance struct {
	calcOffset func(gl.Float) (mathgl.Vec3f)
}

func (i Instance) constructMatrix(fElapsedTime gl.Float) mathgl.Mat4f {
	theMat := mathgl.Ident4f()
	co := i.calcOffset(fElapsedTime)
	theMat[3] = co[0]
	theMat[7] = co[1]
	theMat[11] = co[2]
	theMat[15] = 1.0
	return theMat
}


var instanceList = []Instance{
	{StationaryOffset},
	{OvalOffset},
	{BottomCircleOffset},
}


func CalcFrustumScale(fFovDeg gl.Float) gl.Float {
	fFovRad := fFovDeg * degToRad
	return 1.0 / (gl.Float)(math.Tan((float64)(fFovRad / 2.0)))
}

// arg: fElapsedTime gl.Float
func StationaryOffset(_ gl.Float) mathgl.Vec3f {
	return mathgl.Vec3f{0.0, 0.0, -20.0}
}

func OvalOffset(fElapsedTime gl.Float) mathgl.Vec3f {

	fCurrTimeThroughLoop := math.Mod((float64)(fElapsedTime), fLoopDuration)
	return mathgl.Vec3f{
		(float32)(math.Cos(fCurrTimeThroughLoop * fScale) * 4.0),
		(float32)(math.Sin(fCurrTimeThroughLoop * fScale) * 6.0),
		-20.0}
}

func BottomCircleOffset(fElapsedTime gl.Float) mathgl.Vec3f {

	fCurrTimeThroughLoop := math.Mod((float64)(fElapsedTime), fLoopDuration)	
	return mathgl.Vec3f{
		(float32)(math.Cos(fCurrTimeThroughLoop * fScale) * 5.0),
		-3.5,
		(float32)(math.Sin(fCurrTimeThroughLoop * fScale) * 5.0 - 20.0)}
}

func InitializeVertexBuffer() {
	gl.GenBuffers(1, &vertexBufferObject)
	
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	bufferLen := unsafe.Sizeof(vertexData[0]) * (uintptr)(len(vertexData))
	gl.BufferData(
		gl.ARRAY_BUFFER,
		gl.Sizeiptr(bufferLen),
		gl.Pointer(&vertexData[0]),
		gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.GenBuffers(1, &indexBufferObject)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBufferObject)
	bufferLen = unsafe.Sizeof(indexData[0]) * (uintptr)(len(indexData))
	gl.BufferData(
		gl.ARRAY_BUFFER,
		gl.Sizeiptr(bufferLen),
		gl.Pointer(&indexData[0]),
		gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
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


func InitializeProgram() {
	// Create shaders and bind their variables
	currentShader = CreateShaderProgram(shaders)
	modelToCameraMatrixUnif = gl.GetUniformLocation(currentShader, gl.GLString("modelToCameraMatrix"))
	cameraToClipMatrixUnif = gl.GetUniformLocation(currentShader, gl.GLString("cameraToClipMatrix"))

	cameraToClipMatrix[0] = (float32)(fFrustumScale)
	cameraToClipMatrix[5] = (float32)(fFrustumScale)
	cameraToClipMatrix[10] = (float32)((fzFar + fzNear) / (fzNear - fzFar))
	cameraToClipMatrix[14] = -1.0
	cameraToClipMatrix[11] = (float32)((2 * fzFar * fzNear) / (fzNear - fzFar))

	gl.UseProgram(currentShader)
	gl.UniformMatrix4fv(cameraToClipMatrixUnif, 1, gl.FALSE, (*gl.Float)(&cameraToClipMatrix[0]))
	gl.UseProgram(0)
}

func Initialize() {

	glfwInitWindow()
	gl.Init()
	InitializeProgram()
	InitializeVertexBuffer()

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	colorDataOffset := unsafe.Sizeof(gl.Float(0)) * (uintptr)(3 * numberOfVertices)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, gl.FALSE, 0, gl.Offset(nil, colorDataOffset))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBufferObject)

	gl.BindVertexArray(0)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(gl.TRUE)
	gl.DepthFunc(gl.LEQUAL)
	gl.DepthRange(0.0, 1.0)
}

func display() {

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(currentShader)

	gl.BindVertexArray(vao)

	fElapsedTime := glfw.Time() / (1000 * 1000)
	for i := 0; i < len(instanceList); i++ {
		xform := instanceList[i].constructMatrix((gl.Float)(fElapsedTime))
		Mat4fDebug(xform)
		gl.UniformMatrix4fv(modelToCameraMatrixUnif, 1, gl.FALSE, &(xform[0]))
		gl.DrawElements(
			gl.TRIANGLES, 
			gl.Sizei(len(indexData)), 
			gl.UNSIGNED_SHORT, 
			nil)
	}

	gl.BindVertexArray(0)
	gl.UseProgram(0)

	glfw.SwapBuffers()
}

func reshape(w, h int) {
	cameraToClipMatrix[0] = (float32)(fFrustumScale) * (float32)(h) / (float32)(w)
	cameraToClipMatrix[5] = (float32)(fFrustumScale)

	gl.UseProgram(currentShader)
	gl.UniformMatrix4fv(cameraToClipMatrixUnif, 1, gl.FALSE, (*gl.Float)(&cameraToClipMatrix[0]))
	gl.UseProgram(0)

	gl.Viewport(0, 0, (gl.Sizei)(w), (gl.Sizei)(h))
}

func keyboard(key, state int) {
	if state == glfw.KeyPress {
		switch key {
		case glfw.KeyEsc:
			shutdown()
	
		// case glfw.KeySpace:
		// 	if bDepthClamping == true {
		// 		gl.Disable(gl.DEPTH_CLAMP)
		// 	} else {
		// 		gl.Enable(gl.DEPTH_CLAMP)
		// 	}
		// 	bDepthClamping = !bDepthClamping
		// }
		}
	return
	}
}

func shutdown() {
	// Delete all buffers, shut down glfw
	gl.DeleteBuffers(1, &vertexBufferObject)
	gl.DeleteBuffers(1, &indexBufferObject)
	gl.DeleteProgram(currentShader)
	gl.DeleteVertexArrays(1, &vao)
	glfw.Terminate()
}

// Main loop
func main() {
	// Sit. Down. Good boy.
	runtime.LockOSThread()

	// Initialize subsystems
	Initialize()
	// Set the key handler for the main loop
	glfw.SetKeyCallback(keyboard)

	// Set the resize handler
	glfw.SetWindowSizeCallback(reshape)

	// Main loop.  Run until it dies, or we find someone better.
	for glfw.WindowParam(glfw.Opened) == 1 {

		display()
	}

}