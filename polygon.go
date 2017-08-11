package choppy

import (
	"fmt"
	"image"

	"github.com/fogleman/fauxgl"
	"github.com/fogleman/gg"
	"github.com/fogleman/triangle"
)

type Polygon struct {
	Exterior  Path
	Interiors []Path
}

func renderPolygons(polygons []Polygon) image.Image {
	dc := gg.NewContext(1024, 1024)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.Translate(512, 512)
	dc.Scale(384, 384)
	for _, polygon := range polygons {
		for _, p := range polygon.Exterior {
			dc.LineTo(p.X, p.Y)
		}
		dc.ClosePath()
		dc.SetRGB(0, 0, 0)
		dc.SetLineWidth(3)
		dc.Stroke()
		for _, path := range polygon.Interiors {
			h := path.HolePoint()
			dc.DrawPoint(h.X, h.Y, 3)
			dc.NewSubPath()
			for _, p := range path {
				dc.LineTo(p.X, p.Y)
			}
			dc.ClosePath()
		}
		dc.SetRGB(0.4, 0.4, 0.4)
		dc.SetLineWidth(3)
		dc.Stroke()
	}
	return dc.Image()
}

func (polygon Polygon) Triangulate(plane Plane) *fauxgl.Mesh {
	paths := make([]Path, len(polygon.Interiors)+1)
	paths[0] = polygon.Exterior
	copy(paths[1:], polygon.Interiors)

	var points [][2]float64
	var segments [][2]int32
	for _, path := range paths {
		path = path[1:]
		start := len(points)
		for i, p := range path {
			points = append(points, [2]float64{p.X, p.Y})
			if i == len(path)-1 {
				segments = append(segments, [2]int32{int32(len(points) - 1), int32(start)})
			} else {
				segments = append(segments, [2]int32{int32(len(points) - 1), int32(len(points))})
			}

		}
	}

	var holes [][2]float64
	for _, hole := range polygon.Interiors {
		p := hole.HolePoint()
		holes = append(holes, [2]float64{p.X, p.Y})
	}

	// fmt.Println(points)
	// fmt.Println(segments)
	// fmt.Println(holes)

	if len(segments) < 3 {
		return fauxgl.NewEmptyMesh()
	}

	in := triangle.NewTriangulateIO()
	in.SetPoints(points)
	in.SetSegments(segments)
	if len(holes) > 0 {
		in.SetHoles(holes)
	}
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

func pathsToPolygons(paths []Path) []Polygon {
	var result []Polygon
	seen := make([]bool, len(paths))
	done := false
	for !done {
		done = true
		for i, p := range paths {
			if seen[i] {
				continue
			}
			// see if p is a top-level contour (no others contain it)
			ok := true
			for j, q := range paths {
				if i != j && !seen[j] && q.Contains(p[0]) {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
			seen[i] = true
			// see which contours q are only contained by p (not contained by any other r)
			var holes []Path
			for j, q := range paths {
				if seen[j] || !p.Contains(q[0]) {
					continue
				}
				ok := true
				for k, r := range paths {
					if i != k && j != k && r.Contains(q[0]) {
						ok = false
						break
					}
				}
				if !ok {
					continue
				}
				seen[j] = true
				holes = append(holes, q)
			}
			// create polygon
			result = append(result, Polygon{p, holes})
			done = false
		}
	}
	fmt.Println()
	return result
}
