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

	point := fauxgl.Vector{0, 0, 0}
	normal := fauxgl.Vector{1, 1, 1}.Normalize()

	m1 := choppy.Chop(mesh, point, normal)
	m2 := choppy.Chop(mesh, point, normal.Negate())

	m1.SaveSTL("out1.stl")
	m2.SaveSTL("out2.stl")

	m1.Transform(fauxgl.Translate(normal.MulScalar(30)))
	m2.Transform(fauxgl.Translate(normal.MulScalar(-30)))
	mesh = fauxgl.NewEmptyMesh()
	mesh.Add(m1)
	mesh.Add(m2)
	mesh.SaveSTL("out.stl")
}
