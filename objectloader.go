/*
ObjectLoader - load .obj models for use by OpenGL
*/
package main 

import(
	"fmt"
	"io"
	"bufio"
	"os"
	"strings"
	"strconv"
	"github.com/Jragonmiris/mathgl"
)

// decipherFace - takes a v/uv/n, if a value is missing it 
// returns 0 for that value.
func decipherFace( face string ) (uint, uint, uint) {
	var vertexIndice, uvIndice, normalIndice uint
	ar := strings.Split(strings.TrimSpace(face), "/")
	warnfunc := func(s, e string) (uint) {
		fmt.Fprintf(os.Stderr, "Warning: Unable to parse %s into uint: %s\n", s, e)

	}
	
	if vertexIndice, err := strconv.ParseUint(ar[0], 0, 32); err != nil {
		warnfunc(ar[0], err)
		vertexIndice = 0
	}
	if uvIndice, err := strconv.ParseUint(ar[1], 0, 32); err != nil {
		warnfunc(ar[1], err)
		uvIndice = 0
	}
	if normalIndice, err := strconv.ParseUint(ar[2], 0, 32); err != nil {
		warnfunc(ar[2], err)
		normalIndice = 0
	}
	return vertexIndice, uvIndice, normalIndice
}


func loadOBJ( filePath string ) ([]mathgl.Vec3f, []mathgl.Vec2f, []mathgl.Vec3f) {

	var vertexIndices, uvIndices, normalIndices []uint

	var vertices = make([]mathgl.Vec3f)
	var uvs = make([]mathgl.Vec2f)
	var normals = make([]mathgl.Vec3f)

	// Open the file, and prepare a buffer to read in lines
	fp, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not open file %s\n", filePath)
		return
	}
	defer fp.Close()
	fileBuf := bufio.NewReader(fp)

	// Loop through the lines, 'till the file is done
	for {
		line, err := fileBuf.ReadString('\n')
		if line == "" && err == io.EOF { break } // EOF on line by itself
		words := strings.Fields(line)
		switch words[0] {
		case "v":
			var vtx = mathgl.Vec3f{words[1], words[2], words[3]}
			vertices = append(vertices, vtx)
		case "vt":
			var uv = mathgl.Vec2f{words[1], words[2]}
			uvs = append(uvs, uv)
		case "vn":
			var normal = mathgl.Vec3f{words[1], words[2], words[3]}
			normals = append(normals, normal)
		case "f":





	}

}