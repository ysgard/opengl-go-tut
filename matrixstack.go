// Matrix stack struct and supporting functions
package main

import ()

type MatrixStack struct {
	currMat  *Mat4
	matrices []*Mat4
}

func (ms *MatrixStack) Init() {
	currMat = IdentMat4()
}

// Return pointer to top matrix
func (ms *MatrixStack) Top() *Mat4 {
	return ms.currMat
}

func (ms *MatrixStack) RotateX(deg gl.Float)
