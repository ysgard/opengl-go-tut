// Matrix - A collection of simple routines and structures to help with matrices 

package main

import (
	"fmt"
	"github.com/Jragonmiris/mathgl"
)





func Mat4fDebug(m mathgl.Mat4f) {
	fmt.Printf("\t-------------------------Mat4f---------------------------\n")
	for i := 0; i < 4; i++ {
		fmt.Printf("\t%f\t%f\t%f\t%f\n", m[i*4], m[i*4+1], m[i*4+2], m[i*4+3])
	}
	fmt.Printf("\t--------------------------------------------------------\n")
}