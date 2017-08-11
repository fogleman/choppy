package choppy

func segmentsIntersect(v1x1, v1y1, v1x2, v1y2, v2x1, v2y1, v2x2, v2y2 float64) bool {
	a1 := v1y2 - v1y1
	b1 := v1x1 - v1x2
	c1 := (v1x2 * v1y1) - (v1x1 * v1y2)
	d1 := (a1 * v2x1) + (b1 * v2y1) + c1
	d2 := (a1 * v2x2) + (b1 * v2y2) + c1
	if d1 > 0 && d2 > 0 {
		return false
	}
	if d1 < 0 && d2 < 0 {
		return false
	}
	a2 := v2y2 - v2y1
	b2 := v2x1 - v2x2
	c2 := (v2x2 * v2y1) - (v2x1 * v2y2)
	d1 = (a2 * v1x1) + (b2 * v1y1) + c2
	d2 = (a2 * v1x2) + (b2 * v1y2) + c2
	if d1 > 0 && d2 > 0 {
		return false
	}
	if d1 < 0 && d2 < 0 {
		return false
	}
	if (a1*b2)-(a2*b1) == 0 {
		return false
	}
	return true
}
