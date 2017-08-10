package choppy

import (
	"github.com/fogleman/fauxgl"
	"github.com/fogleman/triangle"
)

type Path []fauxgl.Vector

func joinPaths(paths []Path) []Path {
	lookup := make(map[fauxgl.Vector]Path, len(paths))
	for _, path := range paths {
		lookup[path[0]] = path
	}
	var result []Path
	for len(lookup) > 0 {
		var v fauxgl.Vector
		for v = range lookup {
			break
		}
		var path Path
		for {
			path = append(path, v)
			if p, ok := lookup[v]; ok {
				delete(lookup, v)
				v = p[len(p)-1]
			} else {
				break
			}
		}
		result = append(result, path)
	}
	return result
}

func triangulatePolygon(exterior Path, interiors []Path, plane Plane) *fauxgl.Mesh {
	exterior = exterior[1:]
	if len(exterior) < 3 {
		return fauxgl.NewEmptyMesh()
	}
	var points [][2]float64
	var segments [][2]int32
	n := len(exterior)
	for i, p := range exterior {
		p = plane.Project(p)
		points = append(points, [2]float64{p.X, p.Y})
		j := (i + 1) % n
		segments = append(segments, [2]int32{int32(i), int32(j)})
	}
	in := triangle.NewTriangulateIO()
	in.SetPoints(points)
	in.SetSegments(segments)
	opts := triangle.NewOptions()
	opts.ConformingDelaunay = false
	opts.SegmentSplitting = triangle.NoSplitting
	out := triangle.Triangulate(in, opts, false)
	points = out.Points()
	var triangles []*fauxgl.Triangle
	for _, t := range out.Triangles() {
		point1 := points[t[0]]
		point2 := points[t[1]]
		point3 := points[t[2]]
		p1 := plane.Unproject(fauxgl.Vector{point1[0], point1[1], 0}).RoundPlaces(8)
		p2 := plane.Unproject(fauxgl.Vector{point2[0], point2[1], 0}).RoundPlaces(8)
		p3 := plane.Unproject(fauxgl.Vector{point3[0], point3[1], 0}).RoundPlaces(8)
		ft := fauxgl.NewTriangleForPoints(p1, p2, p3)
		triangles = append(triangles, ft)
	}
	triangle.FreeTriangulateIO(in)
	triangle.FreeTriangulateIO(out)
	return fauxgl.NewTriangleMesh(triangles)
}
