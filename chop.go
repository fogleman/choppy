package choppy

import (
	"github.com/fogleman/fauxgl"
)

func Chop(mesh *fauxgl.Mesh, point, normal fauxgl.Vector) *fauxgl.Mesh {
	plane := MakePlane(point, normal)
	clipped := plane.ClipMesh(mesh)
	sliced := plane.SliceMesh(mesh)
	clipped.Add(sliced)
	return clipped
}
