/* Loads fragment and vertex shader code from the supplied files. */

package main 

import (
	gl "github.com/chsc/gogl/gl33"
	"fmt"
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
			// We've read the last string. Make sure there's a null byte.
			buffer.WriteByte('\000')
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
	//fmt.Fprintf(os.Stdout, vertexShaderCode)
	//fmt.Fprintf(os.Stdout, vSrc)

	// Read the Fragment shader code from the file
	fragmentShaderCode, err := ReadSourceFile(fragmentShaderFilePath)
	if err != nil { return 0 }

	var result gl.Int = gl.TRUE
	var infoLogLength gl.Int

	// Compile the Vertex Shader
	fmt.Fprintf(os.Stdout, "Compiling shader : %s\n", vertexShaderFilePath)
	glslVertexCode := gl.GLStringArray(vertexShaderCode)
	defer gl.GLStringArrayFree(glslVertexCode)
	gl.ShaderSource(vertexShaderID, gl.Sizei(len(glslVertexCode)), &glslVertexCode[0], nil)
	gl.CompileShader(vertexShaderID)

	// Check Vertex Shader
	gl.GetShaderiv(vertexShaderID, gl.COMPILE_STATUS, &result)
	gl.GetShaderiv(vertexShaderID, gl.INFO_LOG_LENGTH, &infoLogLength)
	if infoLogLength > 0 {
		vertexErrorMsg := gl.GLStringAlloc(gl.Sizei(infoLogLength))
		defer gl.GLStringFree(vertexErrorMsg)
		gl.GetShaderInfoLog(vertexShaderID, gl.Sizei(infoLogLength), nil, vertexErrorMsg)
		fmt.Fprintf(os.Stdout, "Vertex Info: %s\n", gl.GoString(vertexErrorMsg))
	}
	if result == gl.FALSE {
		fmt.Fprintf(os.Stderr, "Vertex shader compile failed!\n")
	}

	// Compile the Fragment Shader
	fmt.Fprintf(os.Stdout, "Compiling shader : %s\n", fragmentShaderFilePath)
	glslFragmentCode := gl.GLStringArray(fragmentShaderCode)
	gl.ShaderSource(fragmentShaderID, gl.Sizei(len(glslFragmentCode)), &glslFragmentCode[0], nil)
	gl.CompileShader(fragmentShaderID)

	// Check the Fragment Shader
	gl.GetShaderiv(fragmentShaderID, gl.COMPILE_STATUS, &result)
	gl.GetShaderiv(fragmentShaderID, gl.INFO_LOG_LENGTH, &infoLogLength)
	if infoLogLength > 0 {
		fragmentErrorMsg := gl.GLStringAlloc(gl.Sizei(infoLogLength))
		defer gl.GLStringFree(fragmentErrorMsg)
		gl.GetShaderInfoLog(fragmentShaderID, gl.Sizei(infoLogLength), nil, fragmentErrorMsg)
		fmt.Fprintf(os.Stdout, "Fragment Info: %s\n", gl.GoString(fragmentErrorMsg))
	}

	// Link the shader program
	fmt.Fprintf(os.Stdout, "Linking program...\n")
	var ProgramID gl.Uint = gl.CreateProgram()
	gl.AttachShader(ProgramID, vertexShaderID)
	gl.AttachShader(ProgramID, fragmentShaderID)
	gl.LinkProgram(ProgramID)

	// Check the program
	gl.GetProgramiv(ProgramID, gl.LINK_STATUS, &result)
	gl.GetProgramiv(ProgramID, gl.INFO_LOG_LENGTH, &infoLogLength)
	if infoLogLength > 0 {
		programErrorMsg := gl.GLStringAlloc(gl.Sizei(infoLogLength))
		gl.GetProgramInfoLog(ProgramID, gl.Sizei(infoLogLength), nil, programErrorMsg)
		fmt.Fprintf(os.Stdout, "Program Info: %s\n", gl.GoString(programErrorMsg))
	}
	
	fmt.Fprintf(os.Stdout, "LoadShader completed, ProgramID: %d", ProgramID)
	return ProgramID
}
