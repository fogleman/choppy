package main

import (
	"log"

	"github.com/fogleman/choppy"
	"github.com/fogleman/fauxgl"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	z      = kingpin.Flag("z", "Z offset for slicing.").Short('z').Required().Float64()
	input  = kingpin.Flag("input", "Input STL file.").Short('i').Required().ExistingFile()
	output = kingpin.Flag("output", "Output STL file.").Short('o').Required().String()
)

func main() {
	kingpin.Parse()

	mesh, err := fauxgl.LoadMesh(*input)
	if err != nil {
		log.Fatal(err)
	}

	z0 := mesh.BoundingBox().Min.Z
	p := fauxgl.Vector{0, 0, z0 + *z}
	n := fauxgl.Vector{0, 0, -1}

	choppedMesh := choppy.Chop(mesh, p, n)
	choppedMesh.SaveSTL(*output)
}
