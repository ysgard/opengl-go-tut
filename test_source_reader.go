/* 
Driver for testing ReadSourceFile(), a function that returns a string representing
all the text in a file, with the lines delineated by '\n'.
*/
package main 

import (
	"fmt"
	"os"
	"bytes"
	"bufio"
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

func main() {
	
	if codeString, err := ReadSourceFile("simple_fragment_shader.glsl"); err == nil {
		fmt.Fprintf(os.Stdout, codeString)
	} else {
		fmt.Fprintf(os.Stderr, "Could not read source code! : %s\n", err)
	}
}