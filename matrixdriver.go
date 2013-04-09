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
	m3.PrintMat4("Matrix Mul")
	m4 := m1.MulM(&m2).MulM(m3)
	m4.PrintMat4("Chained Mul")
	m5 := m4.Transpose()
	m5.PrintMat4("Transpose of Last")

}
