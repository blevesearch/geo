// Copyright (c) 2026 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s2

import (
	"testing"
)

func pointFromDegrees(lat, lng float64) Point {
	return PointFromLatLng(LatLngFromDegrees(lat, lng))
}

func TestRectIntersectsEdge(t *testing.T) {
	// A rectangle 20 degrees wide whose top edge runs along the parallel at 60
	// degrees north. The geodesic joining its two top corners bulges north as
	// far as 60.363 degrees, so anything in that strip is a case that a polygon
	// built from the corners of this rectangle would get wrong.
	r := rectFromDegrees(50, -10, 60, 10)

	tests := []struct {
		a, b       Point
		intersects bool
		interior   bool
	}{
		{ // an edge along the parallel just north of the top edge stays outside
			a: pointFromDegrees(60.1, -1), b: pointFromDegrees(60.1, 1),
			intersects: false, interior: false,
		},
		{ // the same edge just inside the top edge
			a: pointFromDegrees(59.9, -1), b: pointFromDegrees(59.9, 1),
			intersects: true, interior: true,
		},
		{ // an edge crossing the top edge
			a: pointFromDegrees(59, 0), b: pointFromDegrees(61, 0),
			intersects: true, interior: true,
		},
		{ // an edge crossing the left edge
			a: pointFromDegrees(55, -11), b: pointFromDegrees(55, -9),
			intersects: true, interior: true,
		},
		{ // an edge whose endpoints are north of the top edge; the geodesic
			// bulges further north still, so it never comes back into the
			// rectangle
			a: pointFromDegrees(60.5, -30), b: pointFromDegrees(60.5, 30),
			intersects: false, interior: false,
		},
		{ // an edge whose endpoints are south of the bottom edge and whose
			// geodesic bulges north as far as 53 degrees, which brings it into
			// the rectangle even though neither endpoint is anywhere near it
			a: pointFromDegrees(49, -30), b: pointFromDegrees(49, 30),
			intersects: true, interior: true,
		},
		{ // the same, but only 10 degrees wide, so the bulge reaches 49.1
			// degrees and stays south of the rectangle
			a: pointFromDegrees(49, -5), b: pointFromDegrees(49, 5),
			intersects: false, interior: false,
		},
		{ // an edge wholly inside the rectangle
			a: pointFromDegrees(55, -5), b: pointFromDegrees(56, 5),
			intersects: true, interior: true,
		},
		{ // an edge running along the left edge of the rectangle, which is a
			// meridian and therefore a geodesic: it touches the boundary but
			// never enters the interior
			a: pointFromDegrees(52, -10), b: pointFromDegrees(58, -10),
			intersects: true, interior: false,
		},
		{ // an edge that stays west of the rectangle
			a: pointFromDegrees(52, -11), b: pointFromDegrees(58, -11),
			intersects: false, interior: false,
		},
	}

	for i, test := range tests {
		if got := r.IntersectsEdge(test.a, test.b); got != test.intersects {
			t.Errorf("%d: %v.IntersectsEdge(%v, %v) = %t, want %t",
				i, r, test.a, test.b, got, test.intersects)
		}
		if got := r.InteriorIntersectsEdge(test.a, test.b); got != test.interior {
			t.Errorf("%d: %v.InteriorIntersectsEdge(%v, %v) = %t, want %t",
				i, r, test.a, test.b, got, test.interior)
		}
		// an edge is never in the interior without being in the rectangle
		if r.InteriorIntersectsEdge(test.a, test.b) && !r.IntersectsEdge(test.a, test.b) {
			t.Errorf("%d: edge is in the interior of %v but not in it", i, r)
		}
	}
}

// A rectangle at the equator, where both edges of constant latitude bulge away
// from the rectangle rather than into it.
func TestRectIntersectsEdgeAcrossEquator(t *testing.T) {
	r := rectFromDegrees(-5, -10, 5, 10)

	tests := []struct {
		a, b       Point
		intersects bool
		interior   bool
	}{
		{ // an edge along the equator, which is itself a geodesic
			a: pointFromDegrees(0, -20), b: pointFromDegrees(0, 20),
			intersects: true, interior: true,
		},
		{ // just north of the top edge, bulging north, so always outside
			a: pointFromDegrees(5.1, -20), b: pointFromDegrees(5.1, 20),
			intersects: false, interior: false,
		},
		{ // just south of the bottom edge, bulging south, so always outside
			a: pointFromDegrees(-5.1, -20), b: pointFromDegrees(-5.1, 20),
			intersects: false, interior: false,
		},
	}

	for i, test := range tests {
		if got := r.IntersectsEdge(test.a, test.b); got != test.intersects {
			t.Errorf("%d: %v.IntersectsEdge(%v, %v) = %t, want %t",
				i, r, test.a, test.b, got, test.intersects)
		}
		if got := r.InteriorIntersectsEdge(test.a, test.b); got != test.interior {
			t.Errorf("%d: %v.InteriorIntersectsEdge(%v, %v) = %t, want %t",
				i, r, test.a, test.b, got, test.interior)
		}
	}
}

// TestRectIntersectsEdgeMatchesIntersectsCell checks IntersectsEdge against
// Rect.IntersectsCell, which is an existing exact test over the same geometry:
// a Cell is bounded by four geodesic edges, so a Rect intersects a Cell exactly
// when the Cell holds the centre of the Rect or one of those four edges
// intersects the Rect.
func TestRectIntersectsEdgeMatchesIntersectsCell(t *testing.T) {
	for i := 0; i < 5000; i++ {
		r := randomRectForEdgeTest()
		c := CellFromCellID(randomCellID())

		want := r.IntersectsCell(c)

		got := c.ContainsPoint(PointFromLatLng(r.Center()))
		for j := 0; j < 4 && !got; j++ {
			got = r.IntersectsEdge(c.Vertex(j), c.Vertex((j+1)&3))
		}

		if got != want {
			t.Fatalf("%d: rect %v vs cell %v: IntersectsEdge based test = %t, "+
				"IntersectsCell = %t", i, r, c.ID(), got, want)
		}
	}
}

// randomRectForEdgeTest returns a random non degenerate Rect.
func randomRectForEdgeTest() Rect {
	r := EmptyRect()
	r = r.AddPoint(LatLngFromPoint(randomPoint()))
	r = r.AddPoint(LatLngFromPoint(randomPoint()))
	if r.Lat.Lo == r.Lat.Hi || r.Lng.Lo == r.Lng.Hi {
		return rectFromDegrees(-1, -1, 1, 1)
	}
	return r
}
