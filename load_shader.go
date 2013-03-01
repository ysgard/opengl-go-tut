/* Loads fragment and vertex shader code from the supplied file. */

package main 

import (
	gl "github.com/chsc/gogl/gl33"
	"fmt"
	//"unsafe"
	"bufio"
	"os"
	"bytes"
	"io"
)

// Reads a file and returns its contents as a string.
func ReadSourceFile(filename string) (string, error) {
	
	fp, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ReadSourceFile: Could not open %s!\n", filename)
		fmt.Fprintf(os.Stderr, "os.Open: %e\n", err)
		return "", err
	}
	defer fp.Close()

	r := bufio.NewReaderSize(fp, 4*1024)
	var buffer bytes.Buffer
	for {
		line, err := r.ReadString('\n')
		buffer.WriteString(line)
		if err == io.EOF {
			// We've read the last string. Make sure there's an end of line.
			buffer.WriteByte('\n')
			break
		}
	}
	return buffer.String(), nil

}

func LoadShaders(vertexShaderFilePath, fragmentShaderFilePath string) gl.Uint {

	// Create the shaders, defer their deletion until the function quits.
	vertexShaderID := gl.CreateShader(gl.VERTEX_SHADER)
	defer gl.DeleteShader(vertexShaderID)
	fragmentShaderID := gl.CreateShader(gl.FRAGMENT_SHADER) 
	defer gl.DeleteShader(fragmentShaderID)

	// Read the Vertex shader code from the file
	vertexShaderCode, err := ReadSourceFile(vertexShaderFilePath)
	if err != nil { return 0 }

	// Read the Fragment shader code from the file
	fragmentShaderCode, err := ReadSourceFile(fragmentShaderFilePath)
	if err != nil { return 0 }

	var result gl.Int = gl.FALSE
	var infoLogLength gl.Int

	// Compile the Vertex Shader
	fmt.Fprintf(os.Stdout, "Compiling shader : %s\n", vertexShaderFilePath)
	glVertexCode := gl.GLString(vertexShaderCode)
	gl.ShaderSource(vertexShaderID, 1, &glVertexCode, nil)

	// Check Vertex Shader
	gl.GetShaderiv(vertexShaderID, gl.COMPILE_STATUS, &result)
	gl.GetShaderiv(vertexShaderID, gl.INFO_LOG_LENGTH, &infoLogLength)
	vertexErrorMsg := gl.GLStringAlloc(gl.Sizei(infoLogLength))
	gl.GetShaderInfoLog(vertexShaderID, gl.Sizei(infoLogLength), nil, vertexErrorMsg)


	// Compile the Fragment Shader
	fmt.Fprintf(os.Stdout, "Compiling shader : %s\n", fragmentShaderFilePath)
	glFragmentCode := gl.GLString(fragmentShaderCode)
	gl.ShaderSource(fragmentShaderID, 1, &glFragmentCode, nil)

	return 0
}

// Returns a '\n'-delimited string