package choppy

import (
	"github.com/fogleman/fauxgl"
)

type Plane struct {
	Point  fauxgl.Vector
	Normal fauxgl.Vector
	U, V   fauxgl.Vector
}

func MakePlane(point, normal fauxgl.Vector) Plane {
	u := normal.Perpendicular().Normalize()
	v := u.Cross(normal).Normalize()
	return Plane{point, normal, u, v}
}

func (p Plane) Project(point fauxgl.Vector) fauxgl.Vector {
	d := point.Sub(p.Point)
	x := d.Dot(p.U)
	y := d.Dot(p.V)
	return fauxgl.Vector{x, y, 0}
}

func (p Plane) Unproject(point fauxgl.Vector) fauxgl.Vector {
	return p.Point.Add(p.U.MulScalar(point.X)).Add(p.V.MulScalar(point.Y))
}

func (p Plane) ClipMesh(m *fauxgl.Mesh) *fauxgl.Mesh {
	var triangles []*fauxgl.Triangle
	for _, t := range m.Triangles {
		if t.IsDegenerate() {
			continue
		}
		f1 := p.pointInFront(t.V1.Position)
		f2 := p.pointInFront(t.V2.Position)
		f3 := p.pointInFront(t.V3.Position)
		if f1 && f2 && f3 {
			triangles = append(triangles, t)
		} else if f1 || f2 || f3 {
			triangles = append(triangles, p.clipTriangle(t)...)
		}
	}
	return fauxgl.NewTriangleMesh(triangles)
}

func (p Plane) SliceMesh(m *fauxgl.Mesh) *fauxgl.Mesh {
	p.Point = p.Point.RoundPlaces(9)
	var paths []Path
	for _, t := range m.Triangles {
		if v1, v2, ok := p.intersectTriangle(t); ok {
			paths = append(paths, Path{v1, v2})
		}
	}
	paths = joinPaths(paths)
	paths = projectPaths(paths, p)
	polygons := pathsToPolygons(paths)

	// im := renderPolygons(polygons)
	// gg.SavePNG("out.png", im)

	mesh := fauxgl.NewEmptyMesh()
	for _, polygon := range polygons {
		mesh.Add(polygon.Triangulate(p))
	}
	return mesh
}

func (p Plane) clipTriangle(t *fauxgl.Triangle) []*fauxgl.Triangle {
	p1 := t.V1.Position
	p2 := t.V2.Position
	p3 := t.V3.Position
	points := []fauxgl.Vector{p1, p2, p3}
	newPoints := sutherlandHodgman(points, []Plane{p})
	var result []*fauxgl.Triangle
	for i := 2; i < len(newPoints); i++ {
		b1 := fauxgl.Barycentric(p1, p2, p3, newPoints[0])
		b2 := fauxgl.Barycentric(p1, p2, p3, newPoints[i-1])
		b3 := fauxgl.Barycentric(p1, p2, p3, newPoints[i])
		v1 := fauxgl.InterpolateVertexes(t.V1, t.V2, t.V3, b1)
		v2 := fauxgl.InterpolateVertexes(t.V1, t.V2, t.V3, b2)
		v3 := fauxgl.InterpolateVertexes(t.V1, t.V2, t.V3, b3)
		result = append(result, fauxgl.NewTriangle(v1, v2, v3))
	}
	return result
}

func (p Plane) pointInFront(v fauxgl.Vector) bool {
	return v.Sub(p.Point).Dot(p.Normal) > 0
}

func (p Plane) intersectSegment(v0, v1 fauxgl.Vector) (fauxgl.Vector, bool) {
	// TODO: do slicing in Z, rotate mesh to plane
	v0 = v0.RoundPlaces(9).Add(p.Normal.MulScalar(5e-10))
	v1 = v1.RoundPlaces(9).Add(p.Normal.MulScalar(5e-10))
	u := v1.Sub(v0)
	w := v0.Sub(p.Point)
	d := p.Normal.Dot(u)
	if d > -1e-9 && d < 1e-9 {
		return fauxgl.Vector{}, false
	}
	n := -p.Normal.Dot(w)
	t := n / d
	if t < 0 || t > 1 {
		return fauxgl.Vector{}, false
	}
	return v0.Add(u.MulScalar(t)), true
}

func (p Plane) intersectTriangle(t *fauxgl.Triangle) (fauxgl.Vector, fauxgl.Vector, bool) {
	v1, ok1 := p.intersectSegment(t.V1.Position, t.V2.Position)
	v2, ok2 := p.intersectSegment(t.V2.Position, t.V3.Position)
	v3, ok3 := p.intersectSegment(t.V3.Position, t.V1.Position)
	var p1, p2 fauxgl.Vector
	if ok1 && ok2 {
		p1, p2 = v1, v2
	} else if ok1 && ok3 {
		p1, p2 = v1, v3
	} else if ok2 && ok3 {
		p1, p2 = v2, v3
	} else {
		return fauxgl.Vector{}, fauxgl.Vector{}, false
	}
	p1 = p1.RoundPlaces(8)
	p2 = p2.RoundPlaces(8)
	if p1 == p2 {
		return fauxgl.Vector{}, fauxgl.Vector{}, false
	}
	n := p2.Sub(p1).Cross(p.Normal)
	if n.Dot(t.Normal()) < 0 {
		return p1, p2, true
	} else {
		return p2, p1, true
	}
}

func sutherlandHodgman(points []fauxgl.Vector, planes []Plane) []fauxgl.Vector {
	output := points
	for _, plane := range planes {
		input := output
		output = nil
		if len(input) == 0 {
			return nil
		}
		s := input[len(input)-1]
		for _, e := range input {
			if plane.pointInFront(e) {
				if !plane.pointInFront(s) {
					x, _ := plane.intersectSegment(s, e)
					output = append(output, x)
				}
				output = append(output, e)
			} else if plane.pointInFront(s) {
				x, _ := plane.intersectSegment(s, e)
				output = append(output, x)
			}
			s = e
		}
	}
	return output
}
