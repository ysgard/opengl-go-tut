// Matrix driver
package main

import (
//"fmt"
//gl "github.com/chsc/gogl/gl33"
)

//Test matrix multiplication

func main() {
	// Create two matrices
	var m1 = Mat4{
		{1, 2, 3, 4}, // col1
		{1, 2, 3, 4}, // col2
		{1, 2, 3, 4}, // col3
		{1, 2, 3, 4}, // col4
	}

	var m2 = Mat4{
		{1, 1, 1, 1}, // col1
		{1, 1, 1, 1}, // col2
		{1, 1, 1, 1}, // col3
		{1, 1, 1, 1}, // col4
	}

	m3 := m1.MulM(&m2)
	m3.Print("Matrix Mul")
	m4 := m1.MulM(&m2).MulM(m3)
	m4.Print("Chained Mul")
	m5 := m4.Transpose()
	m5.Print("Transpose of Last")
	var ms MatrixStack
	ms.Init()
	ms.currMat.Print("Current Matrix Stack")
	ms.Push()
	ms.RotateX(45.0)
	ms.currMat.Print("Rotated 45 degrees")
	ms.Push()
	ms.RotateX(45.0)
	ms.currMat.Print("Rotated another 45 degrees")
	ms.Push()
	vs := &Vec4{444.0, 555.0, 444.0, 1.0}
	ms.Scale(vs)
	cp := ms.currMat.Copy()
	ms.currMat.Print("Scaled by 4,5,4")
	ms.Push()
	ms.Invert()
	ms.currMat.Print("Inverse of last matrix")
	ms.Push()
	ms.MulM(cp)
	ms.currMat.Print("Multiply Inverse by previous, should be Ident")
	ms.Pop()
	ms.currMat.Print("Back to Inverse")
	ms.Pop()
	ms.currMat.Print("Back to scaled by 4,5,5")
	ms.Pop()
	ms.currMat.Print("Back to 90 degrees")
	ms.Pop()
	ms.currMat.Print("Back to 45 degrees")
	ms.Pop()
	ms.currMat.Print("Should be Ident")

}
