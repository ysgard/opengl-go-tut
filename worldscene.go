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

type TreeData struct {
	fXPos, fZPos, fTrunkHeight, fConeHeight gl.Float
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

var g_bDrawLookatPoint = bool(false)
var g_camTarget = &glut.Vec3{0.0, 0.4, 0.0}
// Spherical coordinates
var g_sphereCamRelPos = &glut.Vec3{67.5, -46.0, 150.0}

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

var g_forest []TreeData {

	{-45.0, -40.0, 2.0, 3.0},
	{-42.0, -35.0, 2.0, 3.0},
	{-39.0, -29.0, 2.0, 4.0},
	{-44.0, -26.0, 3.0, 3.0},
	{-40.0, -22.0, 2.0, 4.0},
	{-36.0, -15.0, 3.0, 3.0},
	{-41.0, -11.0, 2.0, 3.0},
	{-37.0, -6.0, 3.0, 3.0},
	{-45.0, 0.0, 2.0, 3.0},
	{-39.0, 4.0, 3.0, 4.0},
	{-36.0, 8.0, 2.0, 3.0},
	{-44.0, 13.0, 3.0, 3.0},
	{-42.0, 17.0, 2.0, 3.0},
	{-38.0, 23.0, 3.0, 4.0},
	{-41.0, 27.0, 2.0, 3.0},
	{-39.0, 32.0, 3.0, 3.0},
	{-44.0, 37.0, 3.0, 4.0},
	{-36.0, 42.0, 2.0, 3.0},

	{-32.0, -45.0, 2.0, 3.0},
	{-30.0, -42.0, 2.0, 4.0},
	{-34.0, -38.0, 3.0, 5.0},
	{-33.0, -35.0, 3.0, 4.0},
	{-29.0, -28.0, 2.0, 3.0},
	{-26.0, -25.0, 3.0, 5.0},
	{-35.0, -21.0, 3.0, 4.0},
	{-31.0, -17.0, 3.0, 3.0},
	{-28.0, -12.0, 2.0, 4.0},
	{-29.0, -7.0, 3.0, 3.0},
	{-26.0, -1.0, 2.0, 4.0},
	{-32.0, 6.0, 2.0, 3.0},
	{-30.0, 10.0, 3.0, 5.0},
	{-33.0, 14.0, 2.0, 4.0},
	{-35.0, 19.0, 3.0, 4.0},
	{-28.0, 22.0, 2.0, 3.0},
	{-33.0, 26.0, 3.0, 3.0},
	{-29.0, 31.0, 3.0, 4.0},
	{-32.0, 38.0, 2.0, 3.0},
	{-27.0, 41.0, 3.0, 4.0},
	{-31.0, 45.0, 2.0, 4.0},
	{-28.0, 48.0, 3.0, 5.0},

	{-25.0, -48.0, 2.0, 3.0},
	{-20.0, -42.0, 3.0, 4.0},
	{-22.0, -39.0, 2.0, 3.0},
	{-19.0, -34.0, 2.0, 3.0},
	{-23.0, -30.0, 3.0, 4.0},
	{-24.0, -24.0, 2.0, 3.0},
	{-16.0, -21.0, 2.0, 3.0},
	{-17.0, -17.0, 3.0, 3.0},
	{-25.0, -13.0, 2.0, 4.0},
	{-23.0, -8.0, 2.0, 3.0},
	{-17.0, -2.0, 3.0, 3.0},
	{-16.0, 1.0, 2.0, 3.0},
	{-19.0, 4.0, 3.0, 3.0},
	{-22.0, 8.0, 2.0, 4.0},
	{-21.0, 14.0, 2.0, 3.0},
	{-16.0, 19.0, 2.0, 3.0},
	{-23.0, 24.0, 3.0, 3.0},
	{-18.0, 28.0, 2.0, 4.0},
	{-24.0, 31.0, 2.0, 3.0},
	{-20.0, 36.0, 2.0, 3.0},
	{-22.0, 41.0, 3.0, 3.0},
	{-21.0, 45.0, 2.0, 3.0},

	{-12.0, -40.0, 2.0, 4.0},
	{-11.0, -35.0, 3.0, 3.0},
	{-10.0, -29.0, 1.0, 3.0},
	{-9.0, -26.0, 2.0, 2.0},
	{-6.0, -22.0, 2.0, 3.0},
	{-15.0, -15.0, 1.0, 3.0},
	{-8.0, -11.0, 2.0, 3.0},
	{-14.0, -6.0, 2.0, 4.0},
	{-12.0, 0.0, 2.0, 3.0},
	{-7.0, 4.0, 2.0, 2.0},
	{-13.0, 8.0, 2.0, 2.0},
	{-9.0, 13.0, 1.0, 3.0},
	{-13.0, 17.0, 3.0, 4.0},
	{-6.0, 23.0, 2.0, 3.0},
	{-12.0, 27.0, 1.0, 2.0},
	{-8.0, 32.0, 2.0, 3.0},
	{-10.0, 37.0, 3.0, 3.0},
	{-11.0, 42.0, 2.0, 2.0},


	{15.0, 5.0, 2.0, 3.0},
	{15.0, 10.0, 2.0, 3.0},
	{15.0, 15.0, 2.0, 3.0},
	{15.0, 20.0, 2.0, 3.0},
	{15.0, 25.0, 2.0, 3.0},
	{15.0, 30.0, 2.0, 3.0},
	{15.0, 35.0, 2.0, 3.0},
	{15.0, 40.0, 2.0, 3.0},
	{15.0, 45.0, 2.0, 3.0},

	{25.0, 5.0, 2.0, 3.0},
	{25.0, 10.0, 2.0, 3.0},
	{25.0, 15.0, 2.0, 3.0},
	{25.0, 20.0, 2.0, 3.0},
	{25.0, 25.0, 2.0, 3.0},
	{25.0, 30.0, 2.0, 3.0},
	{25.0, 35.0, 2.0, 3.0},
	{25.0, 40.0, 2.0, 3.0},
	{25.0, 45.0, 2.0, 3.0},
}

func DrawForest(modelMatrix *glut.MatrixStack) {
	for iTree := 0; iTree < len(g_forest); iTree++ {
		currTree := g_forest[iTree]
		modelMatrix.Push()
		modelMatrix.Translate(&glut.Vec4{currTree.fXPos, 0.0, currTree.fZPos, 0.0})
		DrawTree(modelMatrix, currTree.fTrunkHeight, currTree.fConeHeight)
	}
}

func ResolveCamPosition() *glut.Vec3 {
	tempMat := GetMatrixStack()
	phi := glut.DegToRad(g_sphereCamRelPos.x)
	theta := glut.DegToRad(g_sphereCamRelPos.y + 90.0)
	fSinTheta := glut.SinGL(theta)
	fCosTheta := glut.CosGL(theta)
	fCosPhi := glut.CosGL(phi)
	fSinPhi := glut.SinGL(phi)
	dirToCamera := &glut.Vec3{fSinTheta * fCosPhi, fCosTheta, fSinTheta * fSinPhi}
	dirToCamera.MulS(g_sphereCamRelPos.z).Add(g_camTarget)
	return dirToCamera
}

// Called to update the display
// You should call glfw.SwapBuffers() after all your rendering to display what you rendered.
// If you need continuous updates of the screen, call glutPostRedisplay() at the end of the 
// function.
func display() {
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT)

	campPos := ResolveCamPosition()
	camMatrix := GetMatrixStack()
	camMatrix.Set(CalcLookAtMatrix(camPos, g_camTarget, &glut.Vec3{0.0, 1.0, 0.0}))

	gl.UseProgram(UniformColor.theProgram)
	gl.UniformMatrix4fv(UniformColor.worldToCameraMatrixUnif, 1, gl.FALSE, camMatrix.Top())
	gl.UseProgram(ObjectColor.theProgram)
	gl.UniformMatrix4fv(ObjectColor.worldToCameraMatrixUnif, 1, gl.FALSE, camMatrix.Top())
	gl.UseProgram(UniformColorTint.theProgram)
	gl.UniformMatrix4fv(UniformColorTint.worldToCameraMatrixUnif, 1, gl.FALSE, camMatrix.Top())
	gl.UseProgram(0)

	modelMatrix := GetMatrixStack()
	modelMatrix.Push()
	modelMatrix.Scale(&glut.Vec4(100.0, 1.0, 100.0, 1.0))
	gl.UseProgram(UniformColor.theProgram)
	gl.UniformMatrix4fv(UniformColor.modelToWorldMatrixUnif, 1, gl.FALSE, modelMatrix.Top())
	gl.Uniform4f(UniformColor.baseColorUnif, 0.302, 0.416, 0.0589, 1.0)
	g_pPlaneMesh.Render()
	gl.UseProgram(0)

	DrawForest(modelMatrix)

	// Draw the building
	modelMatrix.Push()
	modelMatrix.Translate(&glut.Vec4{20.0, 0.0, -10.0, 0.0})
	DrawParthenon(modelMatrix)

	if g_bDrawLookatPoint == true {
		gl.Disable(gl.DEPTH_TEST)
		identity := IdentMat4()
		modelMatrix.Push()
		cameraAimVec := g_camTarget.Sub(camPos)
		modelMatrix.Translate(&glut.Vec4{0.0, 0.0, -cameraAimVec.Length(), 0.0})
		modelMatrix.Scale(&glut.Vec4{1.0, 1.0, 1.0, 1.0})

		gl.UseProgram(ObjectColor.theProgram)
		gl.UniformMatrix4fv(ObjectColor.modelToWorldMatrixUnif, 1, gl.FALSE, modelMatrix.Top())
		gl.UniformMatrix4fv(ObjectColor.worldToCameraMatrixUnif, 1, gl.FALSE, identity)
		g_pCubeColorMesh.Render()
		gl.UseProgram(0)
		gl.Enable(gl.DEPTH_TEST)
	}

	glfw.SwapBuffers()
}

// Called whenever teh window is resized.  The new window size is given, in pixels.
// This is an opportunity to call glViewPort or glScissor to keep up with the change
// in size
func reshape(w, h int) {
	persMatrix := GetMatrixStack()
	persMatrix.Perspective()
}