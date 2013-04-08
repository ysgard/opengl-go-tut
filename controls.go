/*
Read controls and use them to generate coordinate mappings
*/
package main

import (
	"github.com/Jragonmiris/mathgl"
	"github.com/go-gl/glfw"
	"math"
)

// glfw.GetTime is called only once, the first time this function is called
var lastTime float64 = glfw.Time()

// Exported Views
var ViewMatrix mathgl.Mat4f
var ProjectionMatrix mathgl.Mat4f

// Initial position : on +Z
var position mathgl.Vec3f = mathgl.Vec3f{0.0, 0.0, 5.0}

// Initial horizontal angle : toward -Z
var horizontalAngle float32 = 3.14

// Initial vertical angle: none
var verticalAngle float32 = 0.0

// Initial field of view, or zoom
// 80° = very wide angle, huge deformations. 
// 60° – 45° : standard. 20° : big zoom.
var initialFoV float32 = 45.0

var speed float32 = 3.0 // 3 units / second
var mouseSpeed float32 = 0.005

// Low precision wrappers for the high-precision trig funcs
func lCos(op float32) float32 { return (float32)(math.Cos((float64)(op))) }
func lSin(op float32) float32 { return (float32)(math.Sin((float64)(op))) }

func computeMatricesFromInputs() {

	// Compute time difference between current and last frame
	currentTime := glfw.Time()
	deltaTime := (float32)(currentTime - lastTime)

	// Get mouse position
	xpos, ypos := glfw.MousePos()

	// Reset mouse position for next frame
	glfw.SetMousePos(Width/2, Height/2)

	// Compute new orientation
	horizontalAngle += mouseSpeed * (Width/2.0 - (float32)(xpos))
	verticalAngle += mouseSpeed * (Height/2.0 - (float32)(ypos))

	// Direction : Spherical coordinates to Cartesian coordinates conversion
	direction := mathgl.Vec3f{
		lCos(verticalAngle) * lSin(horizontalAngle),
		lSin(verticalAngle),
		lCos(verticalAngle) * lCos(horizontalAngle),
	}

	// Right vector
	right := mathgl.Vec3f{
		lSin(horizontalAngle - 3.14/2.0),
		0,
		lCos(horizontalAngle - 3.14/2.0),
	}

	// Up vector
	up := right.Cross(direction)

	// Move forward
	if glfw.Key(glfw.KeyUp) == glfw.KeyPress {
		position = position.Add(direction.Mul(deltaTime * speed))
	}

	// Move backwards
	if glfw.Key(glfw.KeyDown) == glfw.KeyPress {
		position = position.Sub(direction.Mul(deltaTime * speed))
	}

	// Strafe right
	if glfw.Key(glfw.KeyRight) == glfw.KeyPress {
		position = position.Add(right.Mul(deltaTime * speed))
	}

	// Strafe left
	if glfw.Key(glfw.KeyLeft) == glfw.KeyPress {
		position = position.Sub(right.Mul(deltaTime * speed))
	}

	FoV := initialFoV - 5.0*(float32)(glfw.MouseWheel())

	// Projection matrix : 45° Field of View, 4:3 ratio, display range : 0.1 unit <-> 100 units
	ProjectionMatrix = mathgl.Perspective((float64)(FoV), 4.0/3.0, 0.1, 100.0)
	// Camera matrix
	ViewMatrix = mathgl.LookAtV(
		position,                // Camera is here
		position.Add(direction), // And looks here - at the same position, plus
		up)                      // Head is up

	// For the next frame, the "last time" will be "now"
	lastTime = currentTime
}
