package main

import (
	"fmt"
	glut "github.com/Ysgard/goglutils"
)

func main() {
	mesh := glut.NewMesh("TestMesh")
	err := mesh.LoadGLUTMesh("world_tut/UnitCylinderTint.xml")
	if err != nil {
		fmt.Printf("meshtest: Cannot load glutmesh:\n\t%s\n", err)
	}
	mesh.Debug()
}
