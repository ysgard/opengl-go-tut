package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"github.com/go-gl/glfw"
	"math"
	"os"
	"runtime"
	"time"
	"unsafe"
)

const (
	Width  = 500
	Height = 500
	Title  = "Hierarchy Demo"
)

// Shader vars
var theProgram gl.Uint
var positionAttrib gl.Uint
var colorAttrib gl.Uint

var modelToCameraMatrixUnif gl.Uint
var cameraToClipMatrixUnif gl.Uint

var cameraToClipMatrix Mat4

var shaders = []string{
	"shaders/PosColorLocalTransform.vert",
	"shaders/ColorPassthrough.frag",
}

// OpenGL Buffers

var vertexBufferObject gl.Uint
var indexBufferObject gl.Uint
var vao gl.Uint

// Data definitions
// Some simple definitions
type Color struct {
	R gl.Float
	G gl.Float
	B gl.Float
	W gl.Float
}

var GREEN_COLOR = Color{0.0, 1.0, 0.0, 1.0}
var BLUE_COLOR = Color{0.0, 0.0, 1.0, 1.0}
var RED_COLOR = Color{1.0, 0.0, 0.0, 1.0}
var YELLOW_COLOR = Color{1.0, 1.0, 0.0, 1.0}
var CYAN_COLOR = Color{0.0, 1.0, 1.0, 1.0}
var MAGENTA_COLOR = Color{1.0, 0.0, 1.0, 1.0}

// Vertex Data
var numberOfVertices = 24

var vertexData = []gl.Float{
	// Front
	1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,

	// Top
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,

	// Left
	1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,

	// Back
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,

	// Bottom
	1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,

	// Right
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,

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
	RED_COLOR.R, RED_COLOR.G, RED_COLOR.B, 1.0,

	YELLOW_COLOR.R, YELLOW_COLOR.G, YELLOW_COLOR.B, 1.0,
	YELLOW_COLOR.R, YELLOW_COLOR.G, YELLOW_COLOR.B, 1.0,
	YELLOW_COLOR.R, YELLOW_COLOR.G, YELLOW_COLOR.B, 1.0,
	YELLOW_COLOR.R, YELLOW_COLOR.G, YELLOW_COLOR.B, 1.0,

	CYAN_COLOR.R, CYAN_COLOR.G, CYAN_COLOR.B, 1.0,
	CYAN_COLOR.R, CYAN_COLOR.G, CYAN_COLOR.B, 1.0,
	CYAN_COLOR.R, CYAN_COLOR.G, CYAN_COLOR.B, 1.0,
	CYAN_COLOR.R, CYAN_COLOR.G, CYAN_COLOR.B, 1.0,

	MAGENTA_COLOR.R, MAGENTA_COLOR.G, MAGENTA_COLOR.B, 1.0,
	MAGENTA_COLOR.R, MAGENTA_COLOR.G, MAGENTA_COLOR.B, 1.0,
	MAGENTA_COLOR.R, MAGENTA_COLOR.G, MAGENTA_COLOR.B, 1.0,
	MAGENTA_COLOR.R, MAGENTA_COLOR.G, MAGENTA_COLOR.B, 1.0,
}

var indexData = []gl.Short{
	0, 1, 2,
	2, 3, 0,

	4, 5, 6,
	6, 7, 4,

	8, 9, 10,
	10, 11, 8,

	12, 13, 14,
	14, 15, 12,

	16, 17, 18,
	18, 19, 16,

	20, 21, 22,
	22, 23, 20,
}

// Matrix Stack
type MatrixStack struct {
	mCurrMat  *Mat4
	mMatrices []*Mat4
}

func (m *MatrixStack) Init() {
	m.mCurrMat = IdentMat4()
}

func (m *MatrixStack) Top() *Mat4 {
	return m.mCurrMat
}

func (m *MatrixStack) RotateX(fAngDeg gl.Float) {
	m.mCurrMat = m.mCurrMat.Mul
}

// Frustum scale
var fFrustumScale gl.Float

func CalcFrustumScale(fFovDeg gl.Float) gl.Float {
	fFovRad := fFovDeg * degToRad
	return 1.0 / TanGL(fFovRad/2.0)
}

// Initialize vertex array objects
func InitializeVAO() {
	gl.GenBuffers(1, &vertexBufferObject)

	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	bufferLen := unsafe.Sizeof(gl.Float(0)) * (uintptr)(len(vertexData))
	gl.BufferData(
		gl.ARRAY_BUFFER,
		gl.Sizei(bufferLen),
		&vertexData[0],
		gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.GenBuffers(1, &indexBufferObject)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBufferObject)
	bufferLen = unsafe.Sizeof(gl.Short(0)) * (uintptr)(len(indexData))
	gl.BufferData(
		gl.ELEMENT_ARRAY_DATA,
		gl.Sizei(bufferLen),
		&indexData[0],
		gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	colorDataOffset := gl.Offset(nil, unsafe.Sizeof(gl.Float(0)*3*numberOfVertices))
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.EnableVertexAttribArray(positionAttrib)
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(positionAttrib, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.VertexAttribPointer(colorAttrib, 4, gl.FLOAT, gl.FALSE, 0, colorDataOffset)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBufferObject)

	gl.BindVertexArray(0)
}

// Initialize
func InitializeShaders() {

	// Shader program creation, bind attributes
	theProgram = CreateShaderProgram(shaderFiles)
	positionAttrib = gl.GetAttribLocation(theProgram, "position")
	colorAttrib = gl.GetAttribLocation(theProgram, "color")

	modelToCameraMatrixUnif = gl.GetUniformLocation(theProgram, "modelToCameraMatrix")
	cameraToClipMatrixUnif = gl.GetUniformLocation(theProgram, "cameraToClipMatrix")

	fzNear := gl.Float(1.0)
	fzFar := gl.Float(100.0)

	cameraToClipMatrix[0].x = fFrustumScale
	cameraToClipMatrix[1].y = fFrustumScale
	cameraToClipMatrix[2].z = (fzFar + fzNear) / (fzNear - fzFar)
	cameraToClipMatrix[2].w = -1.0
	cameraToClipMatrix[3].z = (2 * fzFar * fzNear) / (fzNear - fzFar)

	gl.UseProgram(theProgram)
	gl.UniformMatrix4fv(cameraToClipMatrixUnif, 1, gl.FALSE, &cameraToClipMatrix[0])
	gl.UseProgram(0)
}

func Initialize() {
	fFrustumScale = CalcFrustumScale(45.0)

	glfwInitWindow()
	gl.Init()

	InitializeShaders()
	InitializeVAO()
}

func reshape(w, h int) {
	cameraToClipMatrix[0].x = fFrustumScale * (gl.Float)(h) / (gl.Float)(w)
	cameraToClipMatrix[1].y = fFrustumScale

	gl.UseProgram(currentShader)
	gl.UniformMatrix4fv(cameraToClipMatrixUnif, 1, gl.FALSE, &cameraToClipMatrix[0].x)
	gl.UseProgram(0)

	gl.Viewport(0, 0, (gl.Sizei)(w), (gl.Sizei)(h))
}

func keyboard(key, state int) {
	if state == glfw.KeyPress {
		switch key {
		case glfw.KeyEsc:
			shutdown()
		}
		return
	}
}

func shutdown() {
	// Delete all buffers, shut down glfw
	// gl.DeleteBuffers(1, &vertexBufferObject)
	// gl.DeleteBuffers(1, &indexBufferObject)
	// gl.DeleteProgram(currentShader)
	// gl.DeleteVertexArrays(1, &vao)
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
		fmt.Fprintf(os.Stdout, "*** Frame: %f ***\n", glfw.Time())
		time.Sleep(time.Millisecond * 100)
		display()
	}

}
