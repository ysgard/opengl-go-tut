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
	warnfunc := func(s string, e error) (uint) {
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

// Dump the output of LoadOBJ
func dumpOBJ( filePath string ) {
	vtx, uv, nrm := loadOBJ(filePath)
	fmt.Printf("*** VERTEXES ***\n\n")
	for i, val := range vtx {
		fmt.Printf(" %f\t%f\t%f\n", val[0], val[1], val[2])
	}
	fmt.Printf("*** UVS ***\n\n")
	for i, val := range vtx {
		fmt.Printf(" %f\t%f\n", val[0], val[1])
	}
	fmt.Printf("*** NORMALS ***\n\n")
	for i, val := range vtx {
		fmt.Printf(" %f\t%f\t%f\n", val[0], val[1], val[2])
	}
}


func loadOBJ( filePath string ) ([]mathgl.Vec3f, []mathgl.Vec2f, []mathgl.Vec3f) {

	var vertexIndices, uvIndices, normalIndices []uint

	var vertices = make([]mathgl.Vec3f, 1, 50)
	var uvs = make([]mathgl.Vec2f, 1, 50)
	var normals = make([]mathgl.Vec3f, 1, 50)

	// Open the file, and prepare a buffer to read in lines
	fp, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not open file %s\n", filePath)
		return nil, nil, nil
	}
	defer fp.Close()
	fileBuf := bufio.NewReader(fp)

	// Loop through the lines, 'till the file is done
	for {
		line, err := fileBuf.ReadString('\n')
		if line == "" && err == io.EOF { 
			break 
		} // EOF on line by itself
		words := strings.Fields(line)
		switch words[0] {
		case "v":
			var vtx = mathgl.Vec3f{
				strconv.ParseFloat(words[1], 32), 
				strconv.ParseFloat(words[2], 32), 
				strconv.ParseFloat(words[3], 32)}
			vertices = append(vertices, vtx)
		case "vt":
			var uv = mathgl.Vec2f{words[1], words[2]}
			uvs = append(uvs, uv)
		case "vn":
			var normal = mathgl.Vec3f{
				strconv.ParseFloat(words[1], 32), 
				strconv.ParseFloat(words[2], 32), 
				strconv.ParseFloat(words[3], 32)}
			normals = append(normals, normal)
		case "f":
			vi1, uvi1, ni1 := decipherFace(words[1])
			vi2, uvi2, ni2 := decipherFace(words[2])
			vi3, uvi3, ni3 := decipherFace(words[3])
			vertexIndices = append(vertexIndices, vi1, vi2, vi3)
			uvIndices = append(uvIndices, uvi1, uvi2, uvi3)
			normalIndices = append(normalIndices, ni1, ni2, ni3)
		default:
			continue
		}
		if err == io.EOF { 
			break 
		} // EOF at end of line
	}

	// Translate object file vertices into an order more preferable to 
	// OpenGL.  OpenGL expects triplets of vertexes, one after another,
	// defining each triangle, but the OBJ format uses the 'f' (face)
	// line to define a triangle and maps it to a listed vertex.
	realVertices := make([]mathgl.Vec3f, 1, 50)
	realUVs := make([]mathgl.Vec2f, 1, 50)
	realNormals := make([]mathgl.Vec3f, 1, 50)

	// Now, for each vertex of each triangle
	for i := 0; i < len(vertexIndices); i++ {
		// Get the indices of its attributes
		vertexIndex := vertexIndices[i]
		uvIndex := uvIndices[i]
		normalIndex := normalIndices[i]

		// Get the attributes thanks to the index
		realVertices = append(realVertices, vertices[vertexIndex - 1])
		realUVs = append(realUVs, uvs[uvIndex - 1])
		realNormals = append(realNormals, normals[normalIndex - 1])

	}
	return realVertices, realUVs, realNormals

}

// Driver, remove once loader has been tested.
func main() {
	dumpOBJ("art/cylinder.obj")
}