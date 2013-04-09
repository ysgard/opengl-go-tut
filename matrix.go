// Matrix - A collection of simple routines and structures to help with matrices
// as well functions for common math operations 

package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"math"
	"os"
)

const degToRad = math.Pi * 2.0 / 360

// Change this to change where debug messages get sent
var debugOut = os.Stderr
var Pi = (gl.Float)(math.Pi)

// Struct that kinda, sorta represents a glm/glsl vector
type Vec4 struct {
	x, y, z, w gl.Float
}

// Struct that kinda, sorta represents a glm/glsl matrix
type Mat4 [4]Vec4

// MulV - multiply receiving matrix by given Vec4 and return
// the new Vec4
func (m *Mat4) MulV(v *Vec4) *Vec4 {
	rv := Vec4{0.0, 0.0, 0.0, 0.0}
	rv.x = m[0].x*v.x + m[1].x*v.y + m[2].x*v.z + m[3].x*v.w
	rv.y = m[0].y*v.x + m[1].y*v.y + m[2].y*v.z + m[3].y*v.w
	rv.z = m[0].z*v.x + m[1].z*v.y + m[2].z*v.z + m[3].z*v.w
	rv.w = m[0].w*v.x + m[1].w*v.y + m[2].w*v.z + m[3].w*v.w
	return &rv
}

// MulM - multiply receiving matrix by given Mat4 and return
// the new Mat.
func (m1 *Mat4) MulM(m2 *Mat4) *Mat4 {
	var rm = Mat4{
		{
			m1[0].x*m2[0].x + m1[1].x*m2[0].y + m1[2].x*m2[0].z + m1[3].x*m2[0].w,
			m1[0].y*m2[0].x + m1[1].y*m2[0].y + m1[2].y*m2[0].z + m1[3].y*m2[0].w,
			m1[0].z*m2[0].x + m1[1].z*m2[0].y + m1[2].z*m2[0].z + m1[3].z*m2[0].w,
			m1[0].w*m2[0].x + m1[1].w*m2[0].y + m1[2].w*m2[0].z + m1[3].w*m2[0].w,
		},
		{
			m1[0].x*m2[1].x + m1[1].x*m2[1].y + m1[2].x*m2[1].z + m1[3].x*m2[1].w,
			m1[0].y*m2[1].x + m1[1].y*m2[1].y + m1[2].y*m2[1].z + m1[3].y*m2[1].w,
			m1[0].z*m2[1].x + m1[1].z*m2[1].y + m1[2].z*m2[1].z + m1[3].z*m2[1].w,
			m1[0].w*m2[1].x + m1[1].w*m2[1].y + m1[2].w*m2[1].z + m1[3].w*m2[1].w,
		},
		{
			m1[0].x*m2[2].x + m1[1].x*m2[2].y + m1[2].x*m2[2].z + m1[3].x*m2[2].w,
			m1[0].y*m2[2].x + m1[1].y*m2[2].y + m1[2].y*m2[2].z + m1[3].y*m2[2].w,
			m1[0].z*m2[2].x + m1[1].z*m2[2].y + m1[2].z*m2[2].z + m1[3].z*m2[2].w,
			m1[0].w*m2[2].x + m1[1].w*m2[2].y + m1[2].w*m2[2].z + m1[3].w*m2[2].w,
		},
		{
			m1[0].x*m2[3].x + m1[1].x*m2[3].y + m1[2].x*m2[3].z + m1[3].x*m2[3].w,
			m1[0].y*m2[3].x + m1[1].y*m2[3].y + m1[2].y*m2[3].z + m1[3].y*m2[3].w,
			m1[0].z*m2[3].x + m1[1].z*m2[3].y + m1[2].z*m2[3].z + m1[3].z*m2[3].w,
			m1[0].w*m2[3].x + m1[1].w*m2[3].y + m1[2].w*m2[3].z + m1[3].w*m2[3].w,
		},
	}
	return &rm
}

// ToArray - produce a []gl.Float array from a given struct.
// Perhaps not necessary, doing &Mat4 should be sufficient!
func (m *Mat4) ToArray() []gl.Float {
	arr := make([]gl.Float, 16)
	for i, vec := range m {
		arr[i*4] = vec.x
		arr[i*4+1] = vec.y
		arr[i*4+2] = vec.z
		arr[i*4+3] = vec.w
	}
	return arr
}

// IdentMat4 - return a Mat4 with identity values
func IdentMat4() *Mat4 {
	var m Mat4
	m[0].x = 1.0
	m[1].y = 1.0
	m[2].z = 1.0
	m[3].w = 1.0
	return &m
}

// Normalize - normalizes a vector, doesn't include w 
func (v *Vec4) Normalize() {
	lenv := (gl.Float)(math.Sqrt((float64)(v.x*v.x + v.y*v.y + v.z*v.z)))
	v.x = v.x / lenv
	v.y = v.y / lenv
	v.z = v.z / lenv
}

// ModGL - Take two gl.Floats and return remainder as a gl.Float
func ModGL(a, b gl.Float) gl.Float {
	return (gl.Float)(math.Mod((float64)(a), (float64)(b)))
}

// LerpGL - Basic linear interpolation
func LerpGL(start, end, ratio gl.Float) gl.Float {
	return start + (end-start)*ratio
}

// CosGL - Cosine, in gl.Float
func CosGL(Rad gl.Float) gl.Float {
	return (gl.Float)(math.Cos((float64)(Rad)))
}

// SinGL - Sine, in gl.Float
func SinGL(Rad gl.Float) gl.Float {
	return (gl.Float)(math.Sin((float64)(Rad)))
}

// TanGL - Tan, in gl.Float
func TanGL(Rad gl.Float) gl.Float {
	return (gl.Float)(math.Tan((float64)(Rad)))
}

// Identity matrix, bare
func Ident4() []gl.Float {
	return []gl.Float{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
}

func DegToRad(fAngDeg gl.Float) gl.Float {
	return fAngDeg * degToRad
}

func Clamp(fValue, fMinValue, fMaxValue gl.Float) gl.Float {
	if fValue < fMinValue {
		return fMinValue
	} else if fValue > fMaxValue {
		return fMaxValue
	} else {
		return fValue
	}
}

// RotateX - returns a Mat4 representing a rotation matrix
// for the angle given in degrees
func RotateX(fAngDeg gl.Float) *Mat4 {
	fAngRad := DegToRad(fAngDeg)
	fCos := CosGL(fAngRad)
	fSin := SinGL(fAngRad)
	theMat := IdentMat4()
	theMat[1].y = fCos
	theMat[2].y = -fSin
	theMat[1].z = fSin
	theMat[2].z = fCos
	return theMat
}

// RotateY - returns a Mat4 representing a rotation matrix
// for the angle given in degree
func RotateY(fAngDeg gl.Float) *Mat4 {
	fAngRad := DegToRad(fAngDeg)
	fCos := CosGL(fAngRad)
	fSin := SinGL(fAngRad)
	theMat := IdentMat4()
	theMat[0].x = fCos
	theMat[2].x = fSin
	theMat[0].z = -fSin
	theMat[2].z = fCos
	return theMat
}

// RotateZ - returns a Mat4 representing a rotation matrix
// for the angle given in degrees
func RotateZ(fAngDeg gl.Float) *Mat4 {
	fAngRad := DegToRad(fAngDeg)
	fCos := CosGL(fAngRad)
	fSin := SinGL(fAngRad)
	theMat := IdentMat4()
	theMat[0].x = fCos
	theMat[1].x = -fSin
	theMat[0].y = fSin
	theMat[1].y = fCos
	return theMat
}

// DebugMat - Pretty-print a []gl.Float slice representing 
// a 16-item transformation matrix. 
func DebugMat(m []gl.Float, s string) {
	fmt.Fprintf(debugOut, "\t-----------------------%s-------------------------\n", s)
	for i := 0; i < 4; i++ {
		fmt.Fprintf(debugOut, "\t%f\t%f\t%f\t%f\n", m[i*4], m[i*4+1], m[i*4+2], m[i*4+3])
	}
	fmt.Fprintf(debugOut, "\t--------------------------------------------------------\n")
}

func (m *Mat4) PrintMat4(s string) {
	if s == "" {
		s = "Debugging Matrix"
	}
	fmt.Fprintf(debugOut, "\t-----------------------%s-------------------------\n", s)
	fmt.Fprintf(debugOut, "\t%f\t%f\t%f\t%f\n", m[0].x, m[1].x, m[2].x, m[3].x)
	fmt.Fprintf(debugOut, "\t%f\t%f\t%f\t%f\n", m[0].y, m[1].y, m[2].y, m[3].y)
	fmt.Fprintf(debugOut, "\t%f\t%f\t%f\t%f\n", m[0].z, m[1].z, m[2].z, m[3].z)
	fmt.Fprintf(debugOut, "\t%f\t%f\t%f\t%f\n", m[0].w, m[1].w, m[2].w, m[3].w)
	fmt.Fprintf(debugOut, "\t--------------------------------------------------------\n")
}
