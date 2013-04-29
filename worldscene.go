package main

import (
	"fmt"
	glut "github.com/Ysgard/goglutils"
	gl "github.com/chsc/gogl/gl33"
	"github.com/go-gl/glfw"
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
	p.theProgram = glut.CreateShaderProgram(shaders)
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

func CalcLookAtMatrix(cameraPt, lookPt, upPt *glut.Vec3) *glut.Mat4 {
	lookDir := (lookPt.Sub(cameraPt)).Normalize()
	upDir := upPt.Normalize()

	rightDir := (lookDir.Cross(upDir)).Normalize()
	perpUpDir := rightDir.Cross(lookDir)

	rotMat := glut.IdentMat4()
	rotMat[0] = rightDir.V3to4(0.0)
	rotMat[1] = perpUpDir.V3to4(0.0)
	rotMat[2] = (lookDir.MulS(-1.0)).V3to4(0.0)

	rotMat = rotMat.Transpose()

	transMat = glut.IdentMat4()
	transMat[3] = (cameraPt.MulS(-1.0)).V3to4(1.0)

	return rotMat.MulM(transMat)
}

var g_fYAngle = gl.Float(0.0)
var g_fXAngle = gl.Float(0.0)

const g_fColumnBaseHeight = gl.Float(0.25)

const g_fParthenonWidth = gl.Float(14.0)
const g_fParthenonLength = gl.Float(20.0)
const g_fParthenonColumnHeight = 5.0
const g_fParthenonBaseHeight = 1.0
const g_fParthenonTopHeight = 2.0

var g_pConeMesh *glut.Mesh
var g_pCylinderMesh *glut.Mesh
var g_pCubeTintMesh *glut.Mesh
var g_pCubeColorMesh *glut.Mesh
var g_pPlaneMesh *glut.Mesh

func LoadMesh(file string) *glut.Mesh {
	mptr, err := glut.LoadMeshFromXML(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load UnitConeTint.xml, exiting...")
		exit()
	}
	return mptr
}

void Initialize() {
	InitializeProgram()
	err := nil
	g_pConeMesh = LoadMesh("UnitConeTint.xml")
	g_pCylinderMesh = LoadMesh("UnitCylinderTint.xml")
	g_pCubeTintMesh = LoadMesh("UnitCubeTint.xml")
	g_pCubeColorMesh = LoadMesh("UnitCubeColor.xml")
	g_pPlaneMesh = LoadMesh("UnitPlane.xml")

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(gl.TRUE)
	gl.DepthFunc(gl.LEQUAL)
	gl.DepthRange(0.0, 1.0)
	gl.Enable(gl.DEPTH_CLAMP)
}


// Trees are 3x3 in X/Y, and fTrunkHeight+fConeHeight in the Y
func DrawTree(modelMatrix *glut.MatrixStack, fTrunkHeight gl.Float, fConeHeight gl.Float) {
	if fTrunkHeight == -1.0 {
		fTrunkHeight = 2.0
	}
	if fConeHeight == -1.0 {
		fConeHeight = 3.0
	}

	// Draw the trunk
	modelMatrix.Push()
	modelMatrix.Scale(&glut.Vec4(1.0, fTrunkHeight, 1.0, 1.0))
	modelMatrix.Translate(&glut.Vec4(0.0, 0.5, 0.0, 0.0))

	gl.UseProgram(UniformColorTint.theProgram)
	gl.UniformMatrix4fv(UniformColorTint.modelToWorldMatrixUnif, 1, gl.FALSE, modelMatrix.Top())
	gl.Uniform4f(UniformColorTint.baseColorUnif, 0.694, 0.4, 0.106, 1.0)
	g_pCyclinderMesh.Render()
	gl.UseProgram(0)

	// Draw the treetop
	modelMatrix.Push()
	modelMatrix.Translate(&glut.Vec4(0.0, fTrunkHeight, 0.0, 1.0))
	modelMatrix.Scale(&glut.Vec4(3.0, fConeHeight, 3.0))
	gl.UseProgram(UniformColorTint.theProgram)
	gl.UniformMatrix4fv(UniformColorTint.modelToWorldMatrixUnif, 1, gl.FALSE, modelMatrix.Top())
	gl.Uniform4f(UniformColorTint.baseColorUnif, 0.0, 1.0, 0.0, 1.0)
	g_pConeMesh.Render()
	gl.UseProgram(0)
}

// Columns are 1x1 in the X/Z, and fHeight units in the Y
func DrawColumn(modelMatrix *glut.ModelStack, fHeight gl.Float) {
	if fHeight == -1.0 {
		fHeight = 5.0
	}

	// Draw the bottom of the column
	modelMatrix.Push()
	modelMatrix.Scale(&glut.Vec4(1.0, g_fColumnBaseHeight, 1.0, 1.0))
	modelMatrix.Translate(&glut.Vec4(0.0, 0.5, 0.0, 0.0))
	gl.UseProgram(UniformColorTint.theProgram)
	gl.UniformMatrix4fv(UniformColorTint.modelToWorldMatrixUnif, 1, gl.FALSE, modelMatrix.Top())
	gl.Uniform4f(UniformColorTint.baseColorUnif, 1.0, 1.0, 1.0, 1.0)
	g_pCubeColorMesh.Render()
	gl.UseProgram(0)

	// Draw the top of the column
	modelMatrix.Push()
	modelMatrix.Translate(&glut.Vec4(0.0, fHeight - g_fColumnBaseHeight, 0.0, 0.0))
	modelMatrix.Scale(&glut.Vec4(1.0, g_fColumnBaseHeight, 1.0, 1.0))
	modelMatrix.Translate(&glut.Vec4(0.0, 0.5, 0.0))
	gl.UseProgram(UniformColorTint.theProgram)
	gl.Uniform4fv(UniformColorTint.modelToWorldMatrixUnif, 1, gl.FALSE, modelMatrix.Top())
	gl.Uniform4f(UniformColorTint.baseColorUnif, 0.9, 0.9, 0.9, 0.9)
	g_pCubeTintMesh.Render()
	gl.UseProgram(0)

	// Draw the main column
	modelMatrix.Push()
	modelMatrix.Translate(&glut.Vec4(0.0, g_fColumnBaseHeight, 0.0, 0.0))
	modelMatrix.Scale(&glut.Vec4(0.8, fHeight - (g_fColumnBaseHeight * 2.0), 0.8, 1.0))
	modelMatrix.Translate(&glut.Vec4(0.0, 0.5, 0.0, 0.0))
	gl.UseProgram(UniformColorTint.theProgram)
	gl.UniformMatrix4fv(UniformColorTint.modelToWorldMatrixUnif, 1.0, gl.FALSE, modelMatrix.Top())
	gl.Uniform4f(UniformColorTint.baseColorUnif, 0.9, 0.9, 0.9, 0.9)
	g_pCylinderMesh.Render()
	gl.UseProgram(0)
}

func DrawParthenon(modelMatrix *glut.MatrixStack) {
	// Draw base.
	modelMatrix.Push()
	modelMatrix.Scale(&glut.Vec4{g_fParthenonWidth, g_fParthenonBaseHeight, g_fParthenonLength, 1.0})
	modelMatrix.Translate(&glut.Vec4{0.0, 0.5, 0.0, 0.0})
	gl.UseProgram(UniformColorTint.theProgram)
	gl.UniformMatrix4fv(UniformColorTint.baseColorUnif, 0.9, 0.9, 0.9, 0.9)
	g_pCubeTintMesh.Render()
	gl.UseProgram(0)

	// Draw Top
	modelMatrix.Push()
	modelMatrix.Translate(&glut.Vec4{0.0, g_fParthenonColumnHeight, g_fParthenonLength, 0.0})
	modelMatrix.Scale(&glut.Vec4{g_fParthenonWidth, g_fParthenonTopHeight, g_fParthenonLength, 1.0})
	modelMatrix.Translate(&glut.Vec4{0.0, 0.5, 0.0, 0.0})
	gl.UseProgram(UniformColorTint.theProgram)
	gl.UniformMatrix4fv(UniformColorTint.modelToWorldMatrixUnif, 1.0, gl.FALSE, modelMatrix.Top())
	gl.Uniform4f(UniformColorTint.baseColorUnif, 0.9, 0.9, 0.9, 0.9)
	g_pCubeTintMesh.Render()
	gl.UseProgram(0)

	// Draw Columns
	fFrontZVal := (g_fParthenonLength/2.0) - 1.0
	fRightXVal := (g_fParthenonWidth/2.0) -1.0

	for iColumnNum := 0; iColumnNum < (int)(g_fParthenonWidth / 2.0); iColumnNum++ {
		modelMatrix.Push()
		modelMatrix.Translate(&glut.Vec4{(2.0*iColumnNum)-(g_fParthenonWidth/2.0)+1.0, g_fParthenonBaseHeight, fFrontZVal, 0.0})
		DrawColumn(modelMatrix, g_fParthenonColumnHeight)

		modelMatrix.Push()
		modelMatrix.Translate(&glut.Vec4{(2.0*iColumnNum)-(g_fParthenonWidth/2.0)+1.0, g_fParthenonBaseHeight, -fFrontZVal, 0.0})
		DrawColumn(modelMatrix, g_fParthenonColumnHeight)
	}

	// Don't draw the first or last columns, since they've been drawn already.
	for iColumnNum := 1; iColumnNum < (int)((g_fParthenonLength - 2.0)/2.0); iColumnNum++ {
		modelMatrix.Push()
		modelMatrix.Translate(*glut.Vec4{fRightXVal, g_fParthenonBaseHeight, (2.0*iColumnNum) - (g_fParthenonLength / 2.0) + 1.0, 0.0})
		DrawColumn(modelMatrix, g_fParthenonColumnHeight)		
	}

	// Draw interior
	modelMatrix.Push()
	modelMatrix.Translate(&glut.Vec4{0.0, 1.0, 0.0, 0.0})
	modelMatrix.Scale(&glut.Vec4{g_fParthenonWidth - 6.0, g_fParthenonColumnHeight, g_fParthenonLength - 6.0, 1.0})
	modelMatrix.Translate(&glut.Vec4{0.0, 0.5, 0.0, 0.0})
	gl.UseProgram(ObjectColor.theProgram)
	gl.UniformMatrix4fv(ObjectColor.modelToWorldMatrixUnif, 1, gl.FALSE, modelMatrix.Top())
	g_pCubeColorMesh.Render()
	gl.UseProgram(0)

	// Draw Headpiece
	modelMatrix.Push()
	modelMatrix.Translate(&glut.Vec4{0.0, g_fParthenonColumnHeight + g_fParthenonBaseHeight + (g_fParthenonTopHeight / 2.0),
		g_fParthenonLength / 2.0, 0.0})
	modelMatrix.RotateX(-135.0)
	modelMatrix.RotateY(45.0)
	gl.UseProgram(ObjectColor.theProgram)
	gl.UniformMatrix4fv(ObjectColor.modelToWorldMatrixUnif, 1, gl.FALSE, modelMatrix.Top())
	g_pCubeColorMesh.Render()
	gl.UseProgram(0)
}



