/* 
Driver for testing ReadSourceFile(), a function that returns a string representing
all the text in a file, with the lines delineated by '\n'.
*/
package main 

import (
	"fmt"
	"os"
	"bufio"
)

func ReadSourceFile(filename string) (string, error) {
	if fp, err := os.Open(filename); err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %s!\n", filename)
		fmt.Fprintf(os.Stderr, "os.Open: %e\n", err)
		return "", err
	}
	var lines []string
	for line, err := 

}

func main() {
	if codeString, err := ReadSourceFile("simple_vertex_shader.glsl"); err != nil {
		fmt.Fprintf(os.Stderr, "Could not read source code! : %s\n", err)
		return
	}
	fmt.Fprintf(os.Stdout, codeString)
}