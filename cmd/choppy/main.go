package main

import (
	"os"

	"github.com/fogleman/choppy"
	"github.com/fogleman/fauxgl"
)

func main() {
	mesh, err := fauxgl.LoadMesh(os.Args[1])
	if err != nil {
		panic(err)
	}

	mesh.Center()

	m1 := choppy.Chop(mesh, fauxgl.Vector{}, fauxgl.Vector{1, 0, 0})
	m2 := choppy.Chop(mesh, fauxgl.Vector{}, fauxgl.Vector{-1, 0, 0})

	// m1 := choppy.Chop(mesh, fauxgl.Vector{0, 0, 4.0011}, fauxgl.Vector{0, 0, 1})
	// m2 := choppy.Chop(mesh, fauxgl.Vector{0, 0, 4.0011}, fauxgl.Vector{0, 0, -1})

	m1.SaveSTL("out1.stl")
	m2.SaveSTL("out2.stl")

	m1.Transform(fauxgl.Rotate(fauxgl.Vector{0, 0, 1}, fauxgl.Radians(-30)).Translate(fauxgl.Vector{10, 0, 0}))
	m2.Transform(fauxgl.Rotate(fauxgl.Vector{0, 0, 1}, fauxgl.Radians(30)).Translate(fauxgl.Vector{-10, 0, 0}))
	mesh = fauxgl.NewEmptyMesh()
	mesh.Add(m1)
	mesh.Add(m2)
	mesh.SaveSTL("out.stl")
}
