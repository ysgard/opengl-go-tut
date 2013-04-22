package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
)

// shader program structure
type ProgramData struct {
	theProgram              gl.Uint
	modelToWorldMatrixUnif  gl.Uint
	worldToCameraMatrixUnif gl.Uint
	cameraToClipMatrixUnif  gl.Uint
	baseColorUnif           gl.Uint
}

func (p *ProgramData) LoadProgram(shaders []string) {
	p.theProgram = CreateShaderProgram(shaders)
	p.modelToWorldMatrixUnif = gl.GetUniformLocation(p.theProgram, gl.String("modelToWorldMatrix"))
	p.worldToCameraMatrixUnif = gl.GetUniformLocation(p.theProgram, gl.String("worldToCameraMatrix"))
	p.cameraToClipMatrixUnif = gl.GetUniformLocation(p.theProgram, gl.String("cameraToClipMatrix"))
	p.baseColorUnif = gl.GetUniformLocation(p.theProgram, gl.String("baseColor"))
}

// camera zoom
var fzNear = gl.Float(1.0)
var fzFar = gl.Float(1000.0)

var UniformColor ProgramData
var ObjectColor ProgramData
var UniformColorTint ProgramData

func InitializeProgram() {
	UniformColor.LoadProgram([]string{
		"world_tut/PosOnlyWorldTransform.vert",
		"world_tut/ColorUniform.frag",
	})
	ObjectColor.LoadProgram([]string{
		"world_tut/PosColorWorldTransform.vert",
		"world_tut/ColorPassthrough.frag",
	})
	UniformColorTint.LoadProgram([]string{
		"world_tut/PosColorWorldTransform.vert",
		"world_tut/ColorMultUniform.frag",
	})
}

func CalcLookAtMatrix(cameraPt, loopPt, upPt *Vec4) {
	
}