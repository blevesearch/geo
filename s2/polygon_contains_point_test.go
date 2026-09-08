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
	"math/rand"
	"testing"

	"github.com/blevesearch/geo/s1"
)

// ringPolygon builds a polygon out of concentric rings around the given centre:
// a shell, a hole inside it, and an island inside the hole. Every ring is given
// enough vertices to put the polygon past the brute force cutoff used elsewhere,
// so that the code paths under test are the ones a stored shape of any real size
// would take.
func ringPolygon(centre Point, radius s1.Angle, numVertices int) *Polygon {
	shell := RegularLoop(centre, radius, numVertices)
	hole := RegularLoop(centre, radius/3, numVertices)
	hole.Invert()
	island := RegularLoop(centre, radius/9, numVertices)

	return PolygonFromOrientedLoops([]*Loop{shell, hole, island})
}

func TestPolygonContainsPointClosed(t *testing.T) {
	centre := pointFromDegrees(40, -73)
	pgn := ringPolygon(centre, s1.Angle(0.1), 40)
	if got := pgn.NumLoops(); got != 3 {
		t.Fatalf("ringPolygon has %d loops, want 3", got)
	}

	shell, hole, island := pgn.Loop(0), pgn.Loop(1), pgn.Loop(2)

	tests := []struct {
		msg  string
		p    Point
		want bool
	}{
		{"between the shell and the hole", pointFromDegrees(43, -73), true},
		{"inside the hole", pointFromDegrees(41.2, -73), false},
		{"inside the island within the hole", centre, true},
		{"outside the shell", pointFromDegrees(48, -73), false},
		{"far outside the bound of the polygon", pointFromDegrees(0, 0), false},
		{"on a vertex of the shell", shell.Vertex(0), true},
		{"on a vertex of the hole", hole.Vertex(0), true},
		{"on a vertex of the island", island.Vertex(0), true},
	}

	for _, test := range tests {
		if got := pgn.ContainsPointClosed(test.p); got != test.want {
			t.Errorf("%v: ContainsPointClosed(%v) = %v, want %v",
				test.msg, LatLngFromPoint(test.p), got, test.want)
		}
	}
}

func TestPolygonContainsPointClosedDegenerate(t *testing.T) {
	p := pointFromDegrees(40, -73)

	if got := FullPolygon().ContainsPointClosed(p); !got {
		t.Errorf("the full polygon contains every point, got %v", got)
	}

	if got := (&Polygon{}).ContainsPointClosed(p); got {
		t.Errorf("the empty polygon contains no point, got %v", got)
	}
}

// TestPolygonContainsPointClosedMatchesQuery checks the answers against a
// ContainsPointQuery over a ShapeIndex under the same vertex model, which is the
// definition this method is a cheaper stand in for.
func TestPolygonContainsPointClosedMatchesQuery(t *testing.T) {
	rnd := rand.New(rand.NewSource(4))

	for trial := 0; trial < 30; trial++ {
		centre := randomPoint(rnd)
		radius := s1.Angle(randomUniformFloat64(0.001, 0.3, rnd))
		pgn := ringPolygon(centre, radius, 33+randomUniformInt(200, rnd))

		idx := NewShapeIndex()
		idx.Add(pgn)
		query := NewContainsPointQuery(idx, VertexModelClosed)

		var points []Point
		// every vertex of every loop, where the two disagree unless the
		// vertices are answered before any crossing count is taken
		for i := 0; i < pgn.NumLoops(); i++ {
			l := pgn.Loop(i)
			for j := 0; j < l.NumVertices(); j++ {
				points = append(points, l.Vertex(j))
			}
		}
		// the midpoint of every eighth edge, which lands on the boundary
		for i := 0; i < pgn.NumLoops(); i++ {
			l := pgn.Loop(i)
			for j := 0; j < l.NumVertices(); j += 8 {
				e := l.Edge(j)
				points = append(points, Point{e.V0.Add(e.V1.Vector).Normalize()})
			}
		}
		// points spread from the centre out past the shell, covering the
		// island, the hole, the ring between the hole and the shell, and the
		// outside
		frame := randomFrameAtPoint(centre, rnd)
		for i := 0; i < 60; i++ {
			d := s1.Angle(randomUniformFloat64(0, 1.4, rnd)) * radius
			points = append(points, InterpolateAtDistance(d, centre,
				fromFrame(*frame, PointFromCoords(
					randomUniformFloat64(-1, 1, rnd),
					randomUniformFloat64(-1, 1, rnd), 0))))
		}

		for _, p := range points {
			want := query.Contains(p)
			if got := pgn.ContainsPointClosed(p); got != want {
				t.Fatalf("trial %d: ContainsPointClosed(%v) = %v, ContainsPointQuery = %v",
					trial, LatLngFromPoint(p), got, want)
			}
		}
	}
}
