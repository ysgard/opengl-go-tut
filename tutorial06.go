/*
Third tutorial in opengl-tutorials.org.

Matrix transformations.
*/
package main

import (
	"fmt"
	"github.com/Jragonmiris/mathgl"
	gl "github.com/chsc/gogl/gl33"
	"github.com/go-gl/glfw"
	"os"
	"runtime"
	"unsafe"
)

const (
	Title  = "Tutorial 06"
	Width  = 800
	Height = 600
)

const (
	VertexFile   = "shaders/cube_texture.vertexshader"
	FragmentFile = "shaders/cube_texture.fragmentshader"
	TextureFile  = "art/liske.tga"
)

func loadTGA(imagePath string) gl.Uint {
	// Create one OpenGL texture
	var txid gl.Uint
	gl.GenTextures(1, &txid)

	// "Bind" the newly created texture: all future functions will modify this texture
	gl.BindTexture(gl.TEXTURE_2D, txid)

	// Read the file, call glTexImage2d with the right parameters
	glfw.LoadTexture2D(imagePath, 0)

	// Nice trilinear filtering.
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// Return the ID of the texture we just created
	return txid
}

func dump4f(m mathgl.Mat4f) {
	for i := 0; i < 4; i++ {
		fmt.Printf("{ %f, %f, %f, %f }\n",
			m[i*4], m[i*4+1], m[i*4+2], m[i*4+3])
	}
}

func main() {
	runtime.LockOSThread()
	// Always call init first
	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}

	// Set Window hints - necessary information before we can
	// call the underlying OpenGL context.
	glfw.OpenWindowHint(glfw.FsaaSamples, 4)        // 4x antialiasing
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3) // OpenGL 3.3
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 2)
	// We want the new OpenGL
	glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// Open a window and initialize its OpenGL context
	if err := glfw.OpenWindow(Width, Height, 0, 0, 0, 0, 32, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n")
		return
	}
	defer glfw.Terminate() // Make sure this gets called if we crash.

	// Set the Window title
	glfw.SetWindowTitle(Title)

	// Make sure we can capture the escape key
	glfw.Enable(glfw.StickyKeys)

	// Initialize OpenGL, make sure we terminate before leaving.
	gl.Init()

	// Dark blue background
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	// Load Shaders
	var programID gl.Uint = LoadShaders(
		VertexFile,
		FragmentFile)
	gl.ValidateProgram(programID)
	var validationErr gl.Int
	gl.GetProgramiv(programID, gl.VALIDATE_STATUS, &validationErr)
	if validationErr == gl.FALSE {
		fmt.Fprintf(os.Stderr, "Shader program failed validation!\n")
	}

	// Time to create some graphics!  
	var vertexArrayID gl.Uint = 0
	gl.GenVertexArrays(1, &vertexArrayID)
	gl.BindVertexArray(vertexArrayID)
	defer gl.DeleteVertexArrays(1, &vertexArrayID) // Make sure this gets called before main is done

	// Get a handle for our "MVP" uniform
	matrixID := gl.GetUniformLocation(programID, gl.GLString("MVP"))

	// Projection matrix: 45° Field of View, 4:3 ratio, display range : 0.1 unit <-> 100 units
	projection := mathgl.Perspective(45.0, 4.0/3.0, 0.1, 100.0)

	// Camera matrix
	view := mathgl.LookAt(
		4.0, 3.0, -3.0,
		0.0, 0.0, 0.0,
		0.0, 1.0, 0.0)

	// Model matrix: and identity matrix (model will be at the origin)
	model := mathgl.Ident4f() // Changes for each model!

	// Our ModelViewProjection : multiplication of our 3 matrices - remember, matrix mult is other way around
	MVP := projection.Mul4(view).Mul4(model) // projection * view * model

	// An array of 3 vectors which represents 3 vertices of a triangle
	/*vertexBufferData2 := [9]gl.Float{	// N.B. We can't use []gl.Float, as that is a slice
		-1.0, -1.0, 0.0,				// We always want to use raw arrays when passing pointers
		1.0, -1.0, 0.0,					// to OpenGL
		0.0, 1.0, 0.0,
	}*/

	// Load the texture
	texture := loadTGA(TextureFile)
	var textureID gl.Int = gl.GetUniformLocation(programID, gl.GLString("myTextureSampler"))

	// Three consecutive floats give a single 3D vertex
	// A cube has 6 faces with 2 triangles each, so this makes 6*2 = 12 triangles,
	// and 12 * 3 vertices
	vertexBufferData := [...]gl.Float{ // N.B. We can't use []gl.Float, as that is a slice
		-1, -1, -1, // face 1
		1, -1, 1,
		-1, -1, 1,
		-1, -1, -1,
		1, -1, 1,
		1, -1, -1,
		1, -1, -1, // face 2
		1, 1, 1,
		1, -1, 1,
		1, -1, -1,
		1, 1, 1,
		1, 1, -1,
		1, 1, -1, // face 3
		-1, 1, 1,
		1, 1, 1,
		1, 1, -1,
		-1, 1, 1,
		-1, 1, -1,
		-1, 1, -1, // face 4
		-1, -1, 1,
		-1, 1, 1,
		-1, 1, -1,
		-1, -1, 1,
		-1, -1, -1,
		1, -1, -1, // face 5
		-1, 1, -1,
		1, 1, -1,
		1, -1, -1,
		-1, 1, -1,
		-1, -1, -1,
		1, 1, 1, // face 6
		-1, -1, 1,
		1, -1, 1,
		1, 1, 1,
		-1, -1, 1,
		-1, 1, 1,
	}

	// Create a random number generator to produce colors
	//now := time.Now()
	//rnd := rand.New(rand.NewSource(now.Unix()))

	//var colorBufferData [3*12*3]gl.Float
	//for i := 0; i < 3*12*3; i += 3 {
	//	colorBufferData[i] = (gl.Float)(rnd.Float32())	// red
	//	colorBufferData[i+1] = (gl.Float)(rnd.Float32()) // blue
	//	colorBufferData[i+2] = (gl.Float)(rnd.Float32()) // green
	//}

	// Two UV coordinated for each vertex.  They were created with
	// Blender.  You'll learn shortly how to do this yourself.
	uvBufferData := [...]gl.Float{
		0, 0,
		1, 1,
		0, 1,
		0, 0,
		1, 1,
		1, 0,
		0, 0,
		1, 1,
		0, 1,
		0, 0,
		1, 1,
		1, 0,
		0, 0,
		1, 1,
		0, 1,
		0, 0,
		1, 1,
		1, 0,
		0, 0,
		1, 1,
		0, 1,
		0, 0,
		1, 1,
		1, 0,
		0, 0,
		1, 1,
		0, 1,
		0, 0,
		1, 1,
		1, 0,
		0, 0,
		1, 1,
		0, 1,
		0, 0,
		1, 1,
		1, 0,
	}

	// One color for each vertex. They were generated randomly.
	/*colorBufferData := [...]gl.Float{
		0.583, 0.771, 0.014,
		0.609, 0.115, 0.436,
		0.327, 0.483, 0.844,
		0.822, 0.569, 0.201,
		0.435, 0.602, 0.223,
		0.310, 0.747, 0.185,
		0.597, 0.770, 0.761,
		0.559, 0.436, 0.730,
		0.359, 0.583, 0.152,
		0.483, 0.596, 0.789,
		0.559, 0.861, 0.639,
		0.195, 0.548, 0.859,
		0.014, 0.184, 0.576,
		0.771, 0.328, 0.970,
		0.406, 0.615, 0.116,
		0.676, 0.977, 0.133,
		0.971, 0.572, 0.833,
		0.140, 0.616, 0.489,
		0.997, 0.513, 0.064,
		0.945, 0.719, 0.592,
		0.543, 0.021, 0.978,
		0.279, 0.317, 0.505,
		0.167, 0.620, 0.077,
		0.347, 0.857, 0.137,
		0.055, 0.953, 0.042,
		0.714, 0.505, 0.345,
		0.783, 0.290, 0.734,
		0.722, 0.645, 0.174,
		0.302, 0.455, 0.848,
		0.225, 0.587, 0.040,
		0.517, 0.713, 0.338,
		0.053, 0.959, 0.120,
		0.393, 0.621, 0.362,
		0.673, 0.211, 0.457,
		0.820, 0.883, 0.371,
		0.982, 0.099, 0.879,
	}*/

	// Time to draw this sucker.
	var vertexBuffer gl.Uint                 // id the vertex buffer
	gl.GenBuffers(1, &vertexBuffer)          // Generate 1 buffer, grab the id
	defer gl.DeleteBuffers(1, &vertexBuffer) // Make sure we delete this, no matter what happens
	// The following commands will talk about our 'vertexBuffer'
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	// Give our vertices to OpenGL
	// WARNING!  This looks EXTREMELY fragile
	gl.BufferData(
		gl.ARRAY_BUFFER,
		gl.Sizeiptr(unsafe.Sizeof(vertexBufferData)), // Already pretty bad
		gl.Pointer(&vertexBufferData),                // SWEET ZOMBIE JESUS PLEASE DON'T CRASH MY MACHINE
		gl.STATIC_DRAW)

	// Set up the UV buffer
	var uvBuffer gl.Uint
	gl.GenBuffers(1, &uvBuffer)
	defer gl.DeleteBuffers(1, &uvBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, uvBuffer)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		gl.Sizeiptr(unsafe.Sizeof(uvBufferData)),
		gl.Pointer(&uvBufferData),
		gl.STATIC_DRAW)

	// Enable Z-buffer
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	var iter uint = 0
	// Main loop - run until it dies, or we find something better
	for (glfw.Key(glfw.KeyEsc) != glfw.KeyPress) &&
		(glfw.WindowParam(glfw.Opened) == 1) {

		// Clear the screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Want to use our loaded shaders
		gl.UseProgram(programID)

		// Generate new MVP from mouse activity
		computeMatricesFromInputs()
		//model := mathgl.Ident4f()
		MVP = ProjectionMatrix.Mul4(ViewMatrix.Mul4(model)) // projection * view * mod

		// Perform the translation of the camera viewpoint
		// by sending the requested operation to the vertex shader
		//mvpm := [16]gl.Float{0.93, -0.85, -0.68, -0.68, 0.0, 1.77, -0.51, -0.51, -1.24, -0.63, -0.51, -0.51, 0.0, 0.0, 5.65, 5.83}
		gl.UniformMatrix4fv(matrixID, 1, gl.FALSE, (*gl.Float)(&MVP[0]))

		fmt.Printf("Iteration %d\n", iter)
		dump4f(MVP)
		//dump4f(MVP2)
		iter++

		// texture in Texture Unit 0
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.Uniform1i(textureID, 0)

		// 1st attribute buffer: vertices
		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
		gl.VertexAttribPointer(
			0,        // Attribute 0. No particular reason for 0, but must match layout in shader
			3,        // size
			gl.FLOAT, // Type
			gl.FALSE, // normalized?
			0,        // stride
			nil)      // array buffer offset

		// 2nd attribute buffer: UVs
		gl.EnableVertexAttribArray(1)
		gl.BindBuffer(gl.ARRAY_BUFFER, uvBuffer)
		gl.VertexAttribPointer(
			1,        // Attribute 1.  Again, no particular reason, but must match layout
			2,        // size
			gl.FLOAT, // Type
			gl.FALSE, // normalized?
			0,
			nil) // array buffer offset

		// Draw the cube!
		gl.DrawArrays(gl.TRIANGLES, 0, 12*3) // Starting from vertex 0, 3 vertices total -> triangle

		gl.DisableVertexAttribArray(0)
		gl.DisableVertexAttribArray(1)

		glfw.SwapBuffers()
	}

}
