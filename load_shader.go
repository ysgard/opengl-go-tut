/* Loads fragment and vertex shader code from the supplied file. */

package main 

import (
	gl "github.com/chsc/gogl/gl33"
	//"fmt"
	//"unsafe"
	"bufio"
	"os"
)

func LoadShaders(vertexShaderFilePath, fragmentShaderFilePath string) gl.Uint {

	// Create the shaders
	vertexShaderID := gl.CreateShader(gl.VERTEX_SHADER)
	fragmentShaderID := gl.CreateShader(gl.FRAGMENT_SHADER) 

	// Read the Vertex shader code from the file
	vertexShaderCode :=

	return 0
}

// Returns a '\n'-delimited string