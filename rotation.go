// go build rotation.go matrix.go shader.go
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
	Title  = "Rotation Demo"
)

// 
const degToRad = math.Pi * 2.0 / 360

var fFrustumScale gl.Float

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
var cameraToClipMatrix = IdentMat4()
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

var GREEN_COLOR = Color{0.0, 1.0, 0.0, 1.0}
var BLUE_COLOR = Color{0.0, 0.0, 1.0, 1.0}
var RED_COLOR = Color{1.0, 0.0, 0.0, 1.0}
var GREY_COLOR = Color{0.8, 0.8, 0.8, 1.0}
var BROWN_COLOR = Color{0.5, 0.5, 0.0, 1.0}

// Vertex data
var numberOfVertices = 8
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
var indexData = []gl.Ushort{
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
	name       string
	RotateFunc func(gl.Float) *Mat4
	offset     Vec4
}

func (i Instance) constructMatrix(fElapsedTime gl.Float) *Mat4 {
	theMat := i.RotateFunc(fElapsedTime)
	theMat[3] = i.offset
	return theMat
}

var instanceList = []Instance{
	{"NullRotation", NullRotation, Vec4{0.0, 0.0, -25.0, 1.0}},
	{"RotateX", RotateX, Vec4{-5.0, -5.0, -25.0, 1.0}},
	{"RotateY", RotateY, Vec4{-5.0, 5.0, -25.0, 1.0}},
	{"RotateZ", RotateZ, Vec4{5.0, 5.0, -25.0, 1.0}},
	{"RotateAxis", RotateAxis, Vec4{5.0, -5.0, -25.0, 1.0}},
}

func CalcLerpFactor(fElapsedTime, fLoopDuration gl.Float) gl.Float {
	fValue := ModGL(fElapsedTime, fLoopDuration) / fLoopDuration
	if fValue > 0.5 {
		fValue = 1.0 - fValue
	}
	return fValue * 2.0
}

func ComputeAngleRad(fElapsedTime, fLoopDuration gl.Float) gl.Float {
	fScale := Pi * 2.0 / fLoopDuration
	fCurrTimeThroughLoop := ModGL(fElapsedTime, fLoopDuration)
	return fCurrTimeThroughLoop * fScale
}

func NullRotation(_ gl.Float) *Mat4 {
	return IdentMat4()
}

func RotateX(fElapsedTime gl.Float) *Mat4 {
	fAngRad := ComputeAngleRad(fElapsedTime, 3.0)
	fCos := CosGL(fAngRad)
	fSin := SinGL(fAngRad)
	theMat := IdentMat4()
	theMat[1].y = fCos
	theMat[2].y = -fSin
	theMat[1].z = fSin
	theMat[2].z = fCos
	return theMat
}

func RotateY(fElapsedTime gl.Float) *Mat4 {
	fAngRad := ComputeAngleRad(fElapsedTime, 2.0)
	fCos := CosGL(fAngRad)
	fSin := SinGL(fAngRad)
	theMat := IdentMat4()
	theMat[0].x = fCos
	theMat[2].x = fSin
	theMat[0].z = -fSin
	theMat[2].z = fCos
	return theMat
}

func RotateZ(fElapsedTime gl.Float) *Mat4 {
	fAngRad := ComputeAngleRad(fElapsedTime, 2.0)
	fCos := CosGL(fAngRad)
	fSin := SinGL(fAngRad)
	theMat := IdentMat4()
	theMat[0].x = fCos
	theMat[1].x = -fSin
	theMat[0].y = fSin
	theMat[1].y = fCos
	return theMat
}

func RotateAxis(fElapsedTime gl.Float) *Mat4 {
	fAngRad := ComputeAngleRad(fElapsedTime, 2.0)
	fCos := CosGL(fAngRad)
	fSin := CosGL(fAngRad)
	fInvCos := 1.0 - fCos
	//fInvSin := 1.0 - fSin
	v := Vec4{1.0, 1.0, 1.0, 1.0}
	v.Normalize()
	m := IdentMat4()
	m[0].x = (v.x * v.x) + ((1 - v.x*v.x) * fCos)
	m[1].x = v.x*v.y*fInvCos - v.z*fSin
	m[2].x = v.x*v.z*fInvCos + v.y*fSin
	m[0].y = v.x*v.y*fInvCos + v.z*fSin
	m[1].y = v.y*v.y + (1-v.y*v.y)*fCos
	m[2].y = v.y*v.z*fInvCos - v.x*fSin
	m[0].z = v.x*v.z*fInvCos - v.y*fSin
	m[1].z = v.y*v.z*fInvCos + v.x*fSin
	m[2].z = v.z*v.z + (1-v.z*v.z)*fCos
	return m
}

func DynamicNonUniformScale(fElapsedTime gl.Float) Vec4 {
	fXLoopDuration := gl.Float(3.0)
	fZLoopDuration := gl.Float(5.0)
	mixx := 1.0 + 4.0*CalcLerpFactor(fElapsedTime, fXLoopDuration)
	mixz := 1.0 + 9.0*CalcLerpFactor(fElapsedTime, fZLoopDuration)
	return Vec4{mixx, 1.0, mixz, 1.0}
}

func CalcFrustumScale(fFovDeg gl.Float) gl.Float {
	fFovRad := fFovDeg * degToRad
	return (gl.Float)(1.0 / math.Tan((float64)(fFovRad/2.0)))
}

func InitializeVertexBuffers() {
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	if gl.IsBuffer(vertexBufferObject) == gl.FALSE {
		fmt.Fprintf(os.Stderr, "Cannot initialize vertexBufferObject!\n")
	}
	bufferLen := unsafe.Sizeof(gl.Float(0)) * (uintptr)(len(vertexData))
	gl.BufferData(
		gl.ARRAY_BUFFER,
		gl.Sizeiptr(bufferLen),
		gl.Pointer(&vertexData[0]),
		gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.GenBuffers(1, &indexBufferObject)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBufferObject)
	if gl.IsBuffer(indexBufferObject) == gl.FALSE {
		fmt.Fprintf(os.Stderr, "Cannot initialize indexBufferObject!\n")
	}
	bufferLen = unsafe.Sizeof(gl.Short(0)) * (uintptr)(len(indexData))
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER,
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
	if modelToCameraMatrixUnif == -1 {
		fmt.Fprintf(os.Stderr, "Invalid value error from glGetUniformLocation: modelToCameraMatrix\n")
	}
	cameraToClipMatrixUnif = gl.GetUniformLocation(currentShader, gl.GLString("cameraToClipMatrix"))
	if cameraToClipMatrixUnif == -1 {
		fmt.Fprintf(os.Stderr, "Invalid value error from glGetUniformLocation: cameraToClipMatrix\n")
	}

	//DebugMat(cameraToClipMatrix, "Camera Matrix")

	cameraToClipMatrix[0].x = fFrustumScale
	cameraToClipMatrix[1].y = fFrustumScale
	cameraToClipMatrix[2].z = (fzFar + fzNear) / (fzNear - fzFar)
	cameraToClipMatrix[2].w = -1.0
	cameraToClipMatrix[3].z = (2 * fzFar * fzNear) / (fzNear - fzFar)

	gl.UseProgram(currentShader)
	gl.UniformMatrix4fv(cameraToClipMatrixUnif, 1, gl.FALSE, &cameraToClipMatrix[0].x)
	gl.UseProgram(0)
}

func Initialize() {

	fFrustumScale = CalcFrustumScale(45.0)

	glfwInitWindow()
	gl.Init()

	InitializeProgram()
	InitializeVertexBuffers()

	gl.GenVertexArrays(1, &vao)
	if gl.IsVertexArray(vao) == gl.FALSE {
		fmt.Fprintf(os.Stderr, "Cannot initialize vao!")
	}
	gl.BindVertexArray(vao)

	colorDataOffset := gl.Offset(nil, unsafe.Sizeof(gl.Float(0))*(uintptr)(3*numberOfVertices))
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, gl.FALSE, 0, colorDataOffset)
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

	fElapsedTime := glfw.Time()
	for i := 0; i < len(instanceList); i++ {
		xform := instanceList[i].constructMatrix((gl.Float)(fElapsedTime))
		//xformT := ToColumnMajor(xform)
		//DebugMat(xform, instanceList[i].name)
		gl.UniformMatrix4fv(modelToCameraMatrixUnif, 1, gl.FALSE, &xform[0].x)
		fmt.Fprintf(os.Stderr, "Drawing %d elements\n", gl.Sizei(len(indexData)))
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
		fmt.Fprintf(os.Stdout, "*** Frame: %f ***\n", glfw.Time())
		time.Sleep(time.Millisecond * 100)
		display()
	}

}
