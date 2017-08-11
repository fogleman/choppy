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

// int areIntersecting(
//     float v1x1, float v1y1, float v1x2, float v1y2,
//     float v2x1, float v2y1, float v2x2, float v2y2
// ) {
//     float d1, d2;
//     float a1, a2, b1, b2, c1, c2;

//     // Convert vector 1 to a line (line 1) of infinite length.
//     // We want the line in linear equation standard form: A*x + B*y + C = 0
//     // See: http://en.wikipedia.org/wiki/Linear_equation
//     a1 = v1y2 - v1y1;
//     b1 = v1x1 - v1x2;
//     c1 = (v1x2 * v1y1) - (v1x1 * v1y2);

//     // Every point (x,y), that solves the equation above, is on the line,
//     // every point that does not solve it, is not. The equation will have a
//     // positive result if it is on one side of the line and a negative one
//     // if is on the other side of it. We insert (x1,y1) and (x2,y2) of vector
//     // 2 into the equation above.
//     d1 = (a1 * v2x1) + (b1 * v2y1) + c1;
//     d2 = (a1 * v2x2) + (b1 * v2y2) + c1;

//     // If d1 and d2 both have the same sign, they are both on the same side
//     // of our line 1 and in that case no intersection is possible. Careful,
//     // 0 is a special case, that's why we don't test ">=" and "<=",
//     // but "<" and ">".
//     if (d1 > 0 && d2 > 0) return NO;
//     if (d1 < 0 && d2 < 0) return NO;

//     // The fact that vector 2 intersected the infinite line 1 above doesn't
//     // mean it also intersects the vector 1. Vector 1 is only a subset of that
//     // infinite line 1, so it may have intersected that line before the vector
//     // started or after it ended. To know for sure, we have to repeat the
//     // the same test the other way round. We start by calculating the
//     // infinite line 2 in linear equation standard form.
//     a2 = v2y2 - v2y1;
//     b2 = v2x1 - v2x2;
//     c2 = (v2x2 * v2y1) - (v2x1 * v2y2);

//     // Calculate d1 and d2 again, this time using points of vector 1.
//     d1 = (a2 * v1x1) + (b2 * v1y1) + c2;
//     d2 = (a2 * v1x2) + (b2 * v1y2) + c2;

//     // Again, if both have the same sign (and neither one is 0),
//     // no intersection is possible.
//     if (d1 > 0 && d2 > 0) return NO;
//     if (d1 < 0 && d2 < 0) return NO;

//     // If we get here, only two possibilities are left. Either the two
//     // vectors intersect in exactly one point or they are collinear, which
//     // means they intersect in any number of points from zero to infinite.
//     if ((a1 * b2) - (a2 * b1) == 0.0f) return COLLINEAR;

//     // If they are not collinear, they must intersect in exactly one point.
//     return YES;
// }
