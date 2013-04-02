// Matrix - A collection of simple routines and structures to help with matrices 

package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
)





func debugMat(m []gl.Float, s string) {
	fmt.Printf("\t-----------------------%s-------------------------\n", s)
	for i := 0; i < 4; i++ {
		fmt.Printf("\t%f\t%f\t%f\t%f\n", m[i*4], m[i*4+1], m[i*4+2], m[i*4+3])
	}
	fmt.Printf("\t--------------------------------------------------------\n")
}