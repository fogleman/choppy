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

// func (a Path) SignedArea() float64 {
// 	var result float64
// 	for i, p1 := range a {
// 		p2 := a[(i+1)%len(a)]
// 		result += p1.X*p2.Y - p2.X*p1.Y
// 	}
// 	return result / 2
// }

func (a Path) IsHole() bool {
	_, ok := a.HolePoint()
	return ok
}

func (a Path) HolePoint() (fauxgl.Vector, bool) {
	p1 := a[0]
	p2 := a[1]
	mid := p1.Add(p2).DivScalar(2)
	dir := p1.Sub(p2).Perpendicular().Normalize()
	off := 1.0
	for i := 0; i < 32; i++ {
		p := mid.Add(dir.MulScalar(off))
		if a.ContainsPoint(p) {
			return p, true
		}
		off /= 2
	}
	return fauxgl.Vector{}, false
}

func (a Path) Project(plane Plane) Path {
	b := make(Path, len(a))
	for i, p := range a {
		b[i] = plane.Project(p)
	}
	return b
}

func (a Path) ContainsPath(b Path) bool {
	for _, p := range b {
		if !a.ContainsPoint(p) {
			return false
		}
	}
	return true
}

func (a Path) ContainsPoint(p fauxgl.Vector) bool {
	box := a.BoundingBox()
	if !box.Contains(p) {
		return false
	}
	outside := box.Min.Sub(fauxgl.Vector{1, 1, 1})
	v1x1 := p.X
	v1y1 := p.Y
	v1x2 := outside.X
	v1y2 := outside.Y
	n := len(a)
	count := 0
	for i := 0; i < n; i++ {
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
		if path[0] != path[len(path)-1] {
			continue
		}
		path = path[1:]
		if len(path) < 3 {
			continue
		}
		result = append(result, path)
	}
	return result
}
