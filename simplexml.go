package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Fleet struct {
	Name string `xml:"name,attr"`
	Ship []Ship `xml:"ship"`
}

func (f *Fleet) Print() {
	fmt.Printf("*** Fleet * Name * %s\n", f.Name)
	for _, v := range f.Ship {
		v.Print()
	}
}

type Ship struct {
	XMLName xml.Name `xml:"ship"`
	Id      string   `xml:"id,attr"`
	Name    string   `xml:"name"`
	Type    string   `xml:"type"`
	Class   string   `xml:"class"`
	Turret  []Turret `xml:"turret"`
	Bay     []Bay    `xml:"bay"`
	Weapon  []Weapon `xml:"weapon"`
	Bombard Bombard  `xml:"bombard"`
	HP      string   `xml:"HP"`
	EW      bool     `xml:"EW"`
}

func (s *Ship) Print() {
	fmt.Printf("*** Ship * Id * %s\n", s.Id)
	fmt.Printf("*** Ship * Name * %s\n", s.Name)
	fmt.Printf("*** Ship * Type * %s\n", s.Type)
	fmt.Printf("*** Ship * Class * %s\n", s.Class)
	for _, v := range s.Turret {
		v.Print()
	}
	for _, v := range s.Bay {
		v.Print()
	}
	for _, v := range s.Weapon {
		v.Print()
	}
	if s.Bombard.Range != "" {
		s.Bombard.Print()
	}
	fmt.Printf("*** Ship * HP * %s\n", s.HP)
	fmt.Printf("*** Ship * EW * %v\n", s.EW)
}

type Turret struct {
	Arc    string `xml:"arc,attr"`
	Weapon Weapon `xml:"weapon"`
}

func (t *Turret) Print() {
	fmt.Printf("*** Turret * Arc * %s\n", t.Arc)
	t.Weapon.Print()
}

type Bay struct {
	Arc     string    `xml:"arc,attr"`
	Type    string    `xml:"type,attr"`
	Missile string    `xml:"missile"`
	Fighter []Fighter `xml:"fighter"`
	Count   string    `xml:"count"`
}

func (b *Bay) Print() {
	fmt.Printf("*** Bay * Arc * %s\n", b.Arc)
	fmt.Printf("*** Bay * Type * %s\n", b.Type)
	fmt.Printf("*** Bay * Missile * %s\n", b.Missile)
	for _, v := range b.Fighter {
		v.Print()
	}
}

type Weapon struct {
	Type   string `xml:"type"`
	Number string `xml:"number,attr"`
}

func (w *Weapon) Print() {
	fmt.Printf("*** Weapon * Type * %s\n", w.Type)
	fmt.Printf("*** Weapon * Number * %s\n", w.Number)
}

type Bombard struct {
	Range string `xml:"range"`
	Type  string `xml:"type"`
}

func (b *Bombard) Print() {
	fmt.Printf("*** Bombard * Range * %s\n", b.Range)
	fmt.Printf("*** Bombard * Type * %s\n", b.Type)
}

type Fighter struct {
	Name string `xml:"name"`
	Type string `xml:"type"`
}

func (f *Fighter) Print() {
	fmt.Printf("*** Fighter * Name * %s\n", f.Name)
	fmt.Printf("*** Fighter * Type * %s\n", f.Type)
}

func main() {

	// Open the XML file, and read its bytestream
	fp, err := os.Open("simple.xml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open simple.xml\n")
		return
	}

	fs, err := fp.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get stats on simple.xml\n")
		return
	}

	buflen := fs.Size()
	fmt.Fprintf(os.Stdout, "File has %d bytes\n", buflen)
	var bytes = make([]byte, buflen)
	numread, err := fp.Read(bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error, only %d bytes read from simple.xml\n", numread)
	}
	fmt.Fprintf(os.Stdout, "%d bytes read from simple.xml\n", numread)
	fp.Close()

	// Unmarshal the xml data
	f := new(Fleet)
	err = xml.Unmarshal(bytes, f)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	f.Print()
	fmt.Printf("Fleet: %v\n", f)

}
