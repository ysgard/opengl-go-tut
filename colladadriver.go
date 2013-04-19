package main

import (
	"encoding/xml"
	//"errors"
	"fmt"
	"os"
)

func BuildModel(filename string) *Collada {
	file, err := os.Open(filename)
	fi, err := file.Stat()
	filelen := fi.Size()
	buf := make([]byte, filelen)
	read, err := file.Read(buf)
	if read != int(filelen) {
		fmt.Fprintf(os.Stderr, "Could not read complete contents of file: %d read vs %d size", read, filelen)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	c := new(Collada)
	err = xml.Unmarshal(buf, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	return c
}

func main() {
	c := BuildModel("world_tut/unitcone.dae")
	c.Debug()
}
