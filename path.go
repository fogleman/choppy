package choppy

import (
	"math"

	"github.com/fogleman/fauxgl"
)

type Path []fauxgl.Vector

func (a Path) BoundingBox() fauxgl.Box {
	x0 := a[0].X
	y0 := a[0].Y
	z0 := a[0].Z
	x1 := a[0].X
	y1 := a[0].Y
	z1 := a[0].Z
	for _, p := range a {
		x0 = math.Min(x0, p.X)
		y0 = math.Min(y0, p.Y)
		z0 = math.Min(z0, p.Z)
		x1 = math.Max(x1, p.X)
		y1 = math.Max(y1, p.Y)
		z1 = math.Max(z1, p.Z)
	}
	return fauxgl.Box{fauxgl.Vector{x0, y0, z0}, fauxgl.Vector{x1, y1, z1}}
}

func (a Path) HolePoint() fauxgl.Vector {
	p1 := a[0]
	p2 := a[1]
	mid := p1.Add(p2).DivScalar(2)
	dir := p1.Sub(p2).Perpendicular().Normalize()
	off := 1.0
	for {
		p := mid.Add(dir.MulScalar(off))
		if a.Contains(p) {
			return p
		}
		off /= 2
	}
}

func (a Path) Project(plane Plane) Path {
	b := make(Path, len(a))
	for i, p := range a {
		b[i] = plane.Project(p)
	}
	return b
}

func (a Path) Contains(p fauxgl.Vector) bool {
	outside := a.BoundingBox().Min.Sub(fauxgl.Vector{1, 1, 1})
	v1x1 := p.X
	v1y1 := p.Y
	v1x2 := outside.X
	v1y2 := outside.Y
	n := len(a)
	count := 0
	for i := 0; i < n-1; i++ {
		j := (i + 1) % n
		v2x1 := a[i].X
		v2y1 := a[i].Y
		v2x2 := a[j].X
		v2y2 := a[j].Y
		if segmentsIntersect(v1x1, v1y1, v1x2, v1y2, v2x1, v2y1, v2x2, v2y2) {
			count++
		}
	}
	return count%2 == 1
}

func projectPaths(paths []Path, plane Plane) []Path {
	result := make([]Path, len(paths))
	for i, path := range paths {
		result[i] = path.Project(plane)
	}
	return result
}

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
