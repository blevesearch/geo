//  Copyright (c) 2026 Couchbase, Inc.
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

package geojson

import (
	"math"
	"math/rand"
	"testing"

	"github.com/blevesearch/geo/s1"
	"github.com/blevesearch/geo/s2"
)

// ringWithHolePolygon builds a polygon out of a ring of the given number of
// vertices around (40, -73) with a concentric hole a fifth of its radius, so
// that a segment across the middle of the ring leaves the polygon and comes
// back into it.
func ringWithHolePolygon(numVertices int) *s2.Polygon {
	ring := make([]s2.Point, 0, numVertices)
	hole := make([]s2.Point, 0, numVertices)
	for i := 0; i < numVertices; i++ {
		th := 2 * math.Pi * float64(i) / float64(numVertices)
		ring = append(ring, s2.PointFromLatLng(s2.LatLngFromDegrees(
			40+math.Sin(th), -73+math.Cos(th))))
		// the hole runs the other way around, as a hole must
		hole = append(hole, s2.PointFromLatLng(s2.LatLngFromDegrees(
			40+0.2*math.Sin(-th), -73+0.2*math.Cos(-th))))
	}

	return s2.PolygonFromOrientedLoops([]*s2.Loop{
		s2.LoopFromPoints(ring), s2.LoopFromPoints(hole)})
}

func polylineFromDegrees(latLngs ...float64) *s2.Polyline {
	pl := make(s2.Polyline, 0, len(latLngs)/2)
	for i := 0; i < len(latLngs); i += 2 {
		pl = append(pl, s2.PointFromLatLng(
			s2.LatLngFromDegrees(latLngs[i], latLngs[i+1])))
	}

	return &pl
}

// A segment whose endpoints both lie in the polygon but whose interior passes
// through a hole is not contained by the polygon. The crossing that rules it out
// used to be looked for with a CrossingEdgeQuery over an empty ShapeIndex, which
// reported no crossings at all for any polygon of more than 27 edges, so the
// vertex counts here straddle that threshold.
func TestPolygonsContainsLineStringsThroughHole(t *testing.T) {
	for _, numVertices := range []int{8, 14, 30, 200} {
		pgn := ringWithHolePolygon(numVertices)

		// from one side of the ring to the other, straight through the hole
		across := polylineFromDegrees(40, -73.75, 40, -72.25)
		if polygonsContainsLineStrings([]*s2.Polygon{pgn},
			[]*s2.Polyline{across}) {
			t.Errorf("polygon of %d edges contains a segment that runs "+
				"through its hole", pgn.NumEdges())
		}

		// wholly within the ring, on one side of the hole, and so contained
		beside := polylineFromDegrees(40.4, -73.4, 40.4, -72.6)
		if !polygonsContainsLineStrings([]*s2.Polygon{pgn},
			[]*s2.Polyline{beside}) {
			t.Errorf("polygon of %d edges does not contain a segment that "+
				"stays within it", pgn.NumEdges())
		}

		// leaving the ring entirely is not contained either
		out := polylineFromDegrees(40, -73.5, 40, -70)
		if polygonsContainsLineStrings([]*s2.Polygon{pgn},
			[]*s2.Polyline{out}) {
			t.Errorf("polygon of %d edges contains a segment that leaves it",
				pgn.NumEdges())
		}
	}
}

// The distances the circle relations are decided by are taken over the edges of
// the polygon rather than through a ClosestEdgeQuery over its ShapeIndex. The two
// answer the same question, so they are compared here against each other.
func TestDistanceFromPointToPolygonMatchesProjection(t *testing.T) {
	rnd := rand.New(rand.NewSource(11))

	// the largest difference seen between the two, in radians
	var worstBoundary, worstPolygon s1.Angle

	for trial := 0; trial < 40; trial++ {
		pgn := ringWithHolePolygon(4 + rnd.Intn(300))

		for i := 0; i < 40; i++ {
			// points spread over the hole, the ring and the outside of it
			p := s2.PointFromLatLng(s2.LatLngFromDegrees(
				40+3*(rnd.Float64()*2-1), -73+3*(rnd.Float64()*2-1)))

			gotBoundary := distanceFromPointToPolygonBoundary(p, pgn)
			wantBoundary := pgn.ProjectToBoundary(&p).Distance(p)
			if d := (gotBoundary - wantBoundary).Abs(); d > worstBoundary {
				worstBoundary = d
			}

			gotPolygon := distanceFromPointToPolygon(p, pgn)
			wantPolygon := pgn.Project(&p).Distance(p)
			if d := (gotPolygon - wantPolygon).Abs(); d > worstPolygon {
				worstPolygon = d
			}
		}
	}

	// a radian is a little over six thousand kilometres of the earth's surface,
	// so this tolerance is well under a micrometre of it
	const tolerance = s1.Angle(1e-13)
	if worstBoundary > tolerance {
		t.Errorf("distanceFromPointToPolygonBoundary differs from "+
			"ProjectToBoundary by %v radians, tolerance is %v",
			float64(worstBoundary), float64(tolerance))
	}
	if worstPolygon > tolerance {
		t.Errorf("distanceFromPointToPolygon differs from Project by %v "+
			"radians, tolerance is %v",
			float64(worstPolygon), float64(tolerance))
	}
}

// segmentCrossesPolygonBoundary replaces a CrossingEdgeQuery under
// CrossingTypeInterior, which is the same test taken over an index.
func TestSegmentCrossesPolygonBoundaryMatchesQuery(t *testing.T) {
	rnd := rand.New(rand.NewSource(12))

	for trial := 0; trial < 40; trial++ {
		pgn := ringWithHolePolygon(4 + rnd.Intn(300))

		// a query over an index the polygon is actually in, which is what the
		// method under test stands in for
		idx := s2.NewShapeIndex()
		idx.Add(pgn)
		query := s2.NewCrossingEdgeQuery(idx)

		for i := 0; i < 40; i++ {
			a := s2.PointFromLatLng(s2.LatLngFromDegrees(
				40+3*(rnd.Float64()*2-1), -73+3*(rnd.Float64()*2-1)))
			b := s2.PointFromLatLng(s2.LatLngFromDegrees(
				40+3*(rnd.Float64()*2-1), -73+3*(rnd.Float64()*2-1)))

			want := len(query.Crossings(a, b, pgn, s2.CrossingTypeInterior)) > 0
			if got := segmentCrossesPolygonBoundary(a, b, pgn); got != want {
				t.Fatalf("trial %d: segmentCrossesPolygonBoundary(%v, %v) = %v, "+
					"CrossingEdgeQuery = %v", trial,
					s2.LatLngFromPoint(a), s2.LatLngFromPoint(b), got, want)
			}
		}
	}
}
