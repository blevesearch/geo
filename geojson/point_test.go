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
	"testing"

	index "github.com/blevesearch/bleve_index_api"
	"github.com/blevesearch/geo/s2"
)

func TestPointIntersects(t *testing.T) {
	tests := []struct {
		queryPoint *Point
		other      index.GeoJSON
		output     bool
	}{
		{ // 0 - Same point with 15 decimal places
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1.234567891234567, 1.234567891234567}},
			other:      NewGeoJsonPoint([]float64{1.234567891234567, 1.234567891234567}),
			output:     true,
		},
		{ // 1 - Point with 15th decimal place differing
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1.234567891234567, 1.234567891234567}},
			other:      NewGeoJsonPoint([]float64{1.234567891234568, 1.234567891234567}),
			output:     true,
		},
		{ // 2 - Point with 13th decimal place differing
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1.234567891234567, 1.234567891234567}},
			other:      NewGeoJsonPoint([]float64{1.234567891234667, 1.234567891234567}),
			output:     false,
		},
		{ // 3 - MultiPoint with a match
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1.234567891234567, 1.234567891234567}},
			other:      NewGeoJsonMultiPoint([][]float64{{1.134567891234567, 1.234567891234567}, {1.234567891234567, 1.234567891234567}}),
			output:     true,
		},
		{ // 4 - MultiPoint with no match
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1.234567891234567, 1.234567891234567}},
			other:      NewGeoJsonMultiPoint([][]float64{{1.234567891234567, 1.134567891234567}, {1.134567891234567, 1.234567891234567}}),
			output:     false,
		},
		{ // 5 - Polygon with point on the inside
			queryPoint: &Point{Typ: PointType, Vertices: []float64{0, 0}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}),
			output:     true,
		},
		{ // 6 - Clockwise polygon with point on the outside
			queryPoint: &Point{Typ: PointType, Vertices: []float64{0, 0}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}, {-1, -1}}}),
			output:     false,
		},
		{ // 7 - Polygon with point on the vertex
			queryPoint: &Point{Typ: PointType, Vertices: []float64{-1, -1}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}),
			output:     true,
		},
		{ // 8 - Polygon with point on the edge
			queryPoint: &Point{Typ: PointType, Vertices: []float64{0.5, 1}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}),
			output:     true,
		},
		{ // 9 - Polygon with point in the hole
			queryPoint: &Point{Typ: PointType, Vertices: []float64{0, 0}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}, {{-0.5, -0.5}, {-0.5, 0.5}, {0.5, 0.5}, {0.5, -0.5}, {-0.5, -0.5}}}),
			output:     false,
		},
		{ // 10 - MultiPolygon with point
			queryPoint: &Point{Typ: PointType, Vertices: []float64{2.5, 2.5}},
			other:      NewGeoJsonMultiPolygon([][][][]float64{{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}, {{{2, 2}, {3, 2}, {3, 3}, {2, 3}, {2, 2}}}}),
			output:     true,
		},
		{ // 11 - MultiPolygon without point
			queryPoint: &Point{Typ: PointType, Vertices: []float64{2.5, 2.5}},
			other:      NewGeoJsonMultiPolygon([][][][]float64{{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}, {{{-2, -2}, {-3, -2}, {-3, -3}, {-2, -3}, {-2, -2}}}}),
			output:     false,
		},
		{ // 12 - LineString with point on the line
			queryPoint: &Point{Typ: PointType, Vertices: []float64{0, 0}},
			other:      NewGeoJsonLinestring([][]float64{{-1, 0}, {1, 0}}),
			output:     true,
		},
		{ // 13 - LineString with point on the vertex
			queryPoint: &Point{Typ: PointType, Vertices: []float64{-1, 0}},
			other:      NewGeoJsonLinestring([][]float64{{-1, 0}, {1, 0}}),
			output:     true,
		},
		{ // 14 - LineString with point not on line
			queryPoint: &Point{Typ: PointType, Vertices: []float64{-2, 0}},
			other:      NewGeoJsonLinestring([][]float64{{-1, 0}, {1, 0}}),
			output:     false,
		},
		{ // 15 - MultiLineString with point on the line
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1, 0}},
			other:      NewGeoJsonMultilinestring([][][]float64{{{-5, 0}, {-3, 0}}, {{-2, 0}, {2, 0}}}),
			output:     true,
		},
		{ // 16 - MultiLineString with point on the vertex
			queryPoint: &Point{Typ: PointType, Vertices: []float64{2, 1}},
			other:      NewGeoJsonMultilinestring([][][]float64{{{-1, 0}, {1, 0}}, {{-2, 1}, {2, 1}}}),
			output:     true,
		},
		{ // 17 - MultiLineString with point not on line
			queryPoint: &Point{Typ: PointType, Vertices: []float64{-3, 1}},
			other:      NewGeoJsonMultilinestring([][][]float64{{{-1, 0}, {1, 0}}, {{-2, 1}, {2, 1}}}),
			output:     false,
		},
		{ // 18 - Circle with point not on the inside
			queryPoint: &Point{Typ: PointType, Vertices: []float64{0, 2}},
			other:      NewGeoCircle([]float64{0, 0}, "1km"),
			output:     false,
		},
		{ // 19 - Circle with point on the inside
			queryPoint: &Point{Typ: PointType, Vertices: []float64{0, 0.03}},
			other:      NewGeoCircle([]float64{0, 0}, "10km"),
			output:     true,
		},
		{ // 20 - Envelope with point on the inside
			queryPoint: &Point{Typ: PointType, Vertices: []float64{0, 0}},
			other:      NewGeoEnvelope([][]float64{{-2, 2}, {2, -2}}),
			output:     true,
		},
		{ // 21 - Envelope with point on the outside
			queryPoint: &Point{Typ: PointType, Vertices: []float64{3, 2}},
			other:      NewGeoEnvelope([][]float64{{-2, 2}, {2, -2}}),
			output:     false,
		},
		{ // 22 - Envelope with point on the edge
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1, 2}},
			other:      NewGeoEnvelope([][]float64{{-2, 2}, {2, -2}}),
			output:     true,
		},
	}

	for i, test := range tests {
		result, err := test.queryPoint.Intersects(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestMultiPointIntersects(t *testing.T) {
	tests := []struct {
		queryPoint *MultiPoint
		other      index.GeoJSON
		output     bool
	}{
		{ // 0 - Same point with 15 decimal places
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234567, 1.234567891234567}),
			output:     true,
		},
		{ // 1 - Point with 15th decimal place differing
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234568, 1.234567891234567}),
			output:     true,
		},
		{ // 2 - Point with 13th decimal place differing
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234667, 1.234567891234567}),
			output:     false,
		},
		{ // 3 - MultiPoint with a match
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonMultiPoint([][]float64{{1.134567891234567, 1.234567891234567}, {1.234567891234567, 1.234567891234567}}),
			output:     true,
		},
		{ // 4 - MultiPoint with no match
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonMultiPoint([][]float64{{1.234567891234567, 1.134567891234567}, {1.134567891234567, 1.234567891234567}}),
			output:     false,
		},
		{ // 5 - Polygon with point on the inside
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{0, 0}, {4, 4}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}),
			output:     true,
		},
		{ // 6 - Clockwise polygon with point on the outside
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{0.5, 0.5}, {0, 0}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}, {-1, -1}}}),
			output:     false,
		},
		{ // 7 - Polygon with point on the vertex
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{4, 4}, {-1, -1}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}),
			output:     true,
		},
		{ // 8 - Polygon with point on the edge
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{-0.5, -1}, {4, 4}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}),
			output:     true,
		},
		{ // 9 - Polygon with point in the hole
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{0, 0}, {4, 4}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}, {{-0.5, -0.5}, {-0.5, 0.5}, {0.5, 0.5}, {0.5, -0.5}, {-0.5, -0.5}}}),
			output:     false,
		},
		{ // 10 - MultiPolygon with point
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{4, 4}, {0, 0}}},
			other:      NewGeoJsonMultiPolygon([][][][]float64{{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}, {{{2, 2}, {3, 2}, {3, 3}, {2, 3}, {2, 2}}}}),
			output:     true,
		},
		{ // 11 - MultiPolygon without point
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{4, 4}, {-4, -4}}},
			other:      NewGeoJsonMultiPolygon([][][][]float64{{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}, {{{-2, -2}, {-3, -2}, {-3, -3}, {-2, -3}, {-2, -2}}}}),
			output:     false,
		},
		{ // 12 - LineString with point on the line
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{0, 0}, {-1, -1}}},
			other:      NewGeoJsonLinestring([][]float64{{-1, 0}, {1, 0}}),
			output:     true,
		},
		{ // 13 - LineString with point on the vertex
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1, 0}, {4, 4}}},
			other:      NewGeoJsonLinestring([][]float64{{-1, 0}, {1, 0}}),
			output:     true,
		},
		{ // 14 - LineString with point not on line
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{4, 4}, {2, 3}}},
			other:      NewGeoJsonLinestring([][]float64{{-1, 0}, {1, 0}}),
			output:     false,
		},
		{ // 15 - MultiLineString with point on the line
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{-2, 0}, {4, 4}}},
			other:      NewGeoJsonMultilinestring([][][]float64{{{-5, 0}, {-3, 0}}, {{-2, 0}, {2, 0}}}),
			output:     true,
		},
		{ // 16 - MultiLineString with point on the vertex
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{4, 4}, {-2, 1}}},
			other:      NewGeoJsonMultilinestring([][][]float64{{{-1, 0}, {1, 0}}, {{-2, 1}, {2, 1}}}),
			output:     true,
		},
		{ // 17 - MultiLineString with point not on line
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1, -1}, {4, 4}}},
			other:      NewGeoJsonMultilinestring([][][]float64{{{-1, 0}, {1, 0}}, {{-2, 1}, {2, 1}}}),
			output:     false,
		},
		{ // 18 - Circle with point not on the inside
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{4, 4}, {-1, -3}}},
			other:      NewGeoCircle([]float64{0, 0}, "1km"),
			output:     false,
		},
		{ // 19 - Circle with point on the inside
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{0.024, -0.037}, {4, 4}}},
			other:      NewGeoCircle([]float64{0, 0}, "10km"),
			output:     true,
		},
		{ // 20 - Envelope with point on the inside
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{4, 4}, {0, 0}}},
			other:      NewGeoEnvelope([][]float64{{-2, 2}, {2, -2}}),
			output:     true,
		},
		{ // 21 - Envelope with point on the outside
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{-2, -3}, {4, 4}}},
			other:      NewGeoEnvelope([][]float64{{-2, 2}, {2, -2}}),
			output:     false,
		},
		{ // 22 - Envelope with point on the edge
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{4, 4}, {-1, -2}}},
			other:      NewGeoEnvelope([][]float64{{-2, 2}, {2, -2}}),
			output:     true,
		},
	}

	for i, test := range tests {
		result, err := test.queryPoint.Intersects(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestPointContains(t *testing.T) {
	tests := []struct {
		queryPoint *Point
		other      index.GeoJSON
		output     bool
	}{
		{ // 0 - Same point with 15 decimal places
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1.234567891234567, 1.234567891234567}},
			other:      NewGeoJsonPoint([]float64{1.234567891234567, 1.234567891234567}),
			output:     true,
		},
		{ // 1 - Point with 15th decimal place differing
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1.234567891234567, 1.234567891234567}},
			other:      NewGeoJsonPoint([]float64{1.234567891234568, 1.234567891234567}),
			output:     true,
		},
		{ // 2 - Point with 13th decimal place differing
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1.234567891234567, 1.234567891234567}},
			other:      NewGeoJsonPoint([]float64{1.234567891234667, 1.234567891234567}),
			output:     false,
		},
		{ // 3 - MultiPoint with a match
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1.234567891234567, 1.234567891234567}},
			other:      NewGeoJsonMultiPoint([][]float64{{1.234567891234567, 1.234567891234567}}),
			output:     true,
		},
		{ // 4 - MultiPoint with no match
			queryPoint: &Point{Typ: PointType, Vertices: []float64{1.234567891234567, 1.234567891234567}},
			other:      NewGeoJsonMultiPoint([][]float64{{1.234567891234567, 1.134567891234567}, {1.134567891234567, 1.234567891234567}}),
			output:     false,
		},
	}

	for i, test := range tests {
		result, err := test.queryPoint.Contains(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestMultiPointContains(t *testing.T) {
	tests := []struct {
		queryPoint *MultiPoint
		other      index.GeoJSON
		output     bool
	}{
		{ // 0 - Same point with 15 decimal places
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234567, 1.234567891234567}),
			output:     true,
		},
		{ // 1 - Point with 15th decimal place differing
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234568, 1.234567891234567}),
			output:     true,
		},
		{ // 2 - Point with 13th decimal place differing
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234667, 1.234567891234567}),
			output:     false,
		},
		{ // 3 - MultiPoint with a match
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonMultiPoint([][]float64{{2.234567891234567, 2.234567891234567}, {1.234567891234567, 1.234567891234567}}),
			output:     true,
		},
		{ // 4 - MultiPoint with no match
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonMultiPoint([][]float64{{1.234567891234567, 1.134567891234567}, {1.134567891234567, 1.234567891234567}}),
			output:     false,
		},
	}

	for i, test := range tests {
		result, err := test.queryPoint.Contains(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

// ---------------------------------------------------------------------------
// geo shape v2 cell tests

func TestPointCells(t *testing.T) {
	p := NewGeoJsonPoint([]float64{4.5, 22.5}).(*Point)

	inner, cross := p.IndexCells()
	if len(inner) != 0 {
		t.Fatalf("expected no inner cells for a point, got %v", inner)
	}
	if len(cross) != 1 {
		t.Fatalf("expected exactly one cross cell for a point, got %v", cross)
	}
	if cross[0] != pointCell(*p.s2point) {
		t.Fatalf("expected cross cell %d, got %d", pointCell(*p.s2point), cross[0])
	}
	if !cellsCoverLatLng(cross, 22.5, 4.5) {
		t.Fatal("the point's cell does not cover the point")
	}

	// query cells must be identical to index cells for a point
	qInner, qCross := p.QueryCells()
	if len(qInner) != 0 || len(qCross) != 1 || qCross[0] != cross[0] {
		t.Fatalf("expected query cells to equal index cells, got %v %v",
			qInner, qCross)
	}
}

func TestPointBoundingBox(t *testing.T) {
	p := NewGeoJsonPoint([]float64{4.5, 22.5}).(*Point)

	env, ok := p.BoundingBox().(*Envelope)
	if !ok || env.r == nil {
		t.Fatalf("expected an envelope bounding box, got %v", p.BoundingBox())
	}
	// a point's bounding box is degenerate: lo == hi == the point
	if !env.r.ContainsLatLng(s2.LatLngFromDegrees(22.5, 4.5)) {
		t.Fatal("bounding box does not contain the point")
	}
	if !rectsApproxEqual(*env.r, s2.RectFromLatLng(s2.LatLngFromDegrees(22.5, 4.5))) {
		t.Fatalf("expected a degenerate rect at the point, got %v", env.r)
	}
}

func TestMultiPointCells(t *testing.T) {
	// two distinct points plus an exact duplicate of the first: the
	// duplicate must be deduplicated away
	mp := NewGeoJsonMultiPoint([][]float64{{1, 1}, {50, 50}, {1, 1}}).(*MultiPoint)

	inner, cross := mp.IndexCells()
	if len(inner) != 0 {
		t.Fatalf("expected no inner cells for a multipoint, got %v", inner)
	}
	if len(cross) != 2 {
		t.Fatalf("expected two deduplicated cross cells, got %v", cross)
	}
	for _, ll := range [][2]float64{{1, 1}, {50, 50}} {
		if !cellsCoverLatLng(cross, ll[0], ll[1]) {
			t.Fatalf("cells do not cover point (%v, %v)", ll[0], ll[1])
		}
	}

	// query cells must be identical to index cells for a multipoint
	qInner, qCross := mp.QueryCells()
	if len(qInner) != 0 || len(qCross) != len(cross) {
		t.Fatalf("expected query cells to equal index cells, got %v %v",
			qInner, qCross)
	}
}

func TestMultiPointBoundingBox(t *testing.T) {
	mp := NewGeoJsonMultiPoint([][]float64{{1, 1}, {50, 50}}).(*MultiPoint)

	env, ok := mp.BoundingBox().(*Envelope)
	if !ok || env.r == nil {
		t.Fatalf("expected an envelope bounding box, got %v", mp.BoundingBox())
	}
	for _, ll := range [][2]float64{{1, 1}, {50, 50}} {
		if !rectContainsDegrees(*env.r, ll[0], ll[1]) {
			t.Fatalf("bounding box does not contain point (%v, %v)", ll[0], ll[1])
		}
	}
}
