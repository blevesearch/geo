//  Copyright (c) 2025 Couchbase, Inc.
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
		{ // 10 - MulitiPolygon with point
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
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234567, 1.234567891234567}),
			output:     true,
		},
		{ // 1 - Point with 15th decimal place differing
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234568, 1.234567891234567}),
			output:     true,
		},
		{ // 2 - Point with 13th decimal place differing
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234667, 1.234567891234567}),
			output:     false,
		},
		{ // 3 - MultiPoint with a match
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonMultiPoint([][]float64{{1.134567891234567, 1.234567891234567}, {1.234567891234567, 1.234567891234567}}),
			output:     true,
		},
		{ // 4 - MultiPoint with no match
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonMultiPoint([][]float64{{1.234567891234567, 1.134567891234567}, {1.134567891234567, 1.234567891234567}}),
			output:     false,
		},
		{ // 5 - Polygon with point on the inside
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{0, 0}, {4, 4}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}),
			output:     true,
		},
		{ // 6 - Clockwise polygon with point on the outside
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{0.5, 0.5}, {0, 0}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}, {-1, -1}}}),
			output:     false,
		},
		{ // 7 - Polygon with point on the vertex
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{4, 4}, {-1, -1}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}),
			output:     true,
		},
		{ // 8 - Polygon with point on the vertex
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{-0.5, -1}, {4, 4}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}),
			output:     true,
		},
		{ // 9 - Polygon with point in the hole
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{0, 0}, {4, 4}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}, {{-0.5, -0.5}, {-0.5, 0.5}, {0.5, 0.5}, {0.5, -0.5}, {-0.5, -0.5}}}),
			output:     false,
		},
		{ // 10 - MulitiPolygon with point
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{4, 4}, {0, 0}}},
			other:      NewGeoJsonMultiPolygon([][][][]float64{{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}, {{{2, 2}, {3, 2}, {3, 3}, {2, 3}, {2, 2}}}}),
			output:     true,
		},
		{ // 11 - MultiPolygon without point
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{4, 4}, {-4, -4}}},
			other:      NewGeoJsonMultiPolygon([][][][]float64{{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}, {{{-2, -2}, {-3, -2}, {-3, -3}, {-2, -3}, {-2, -2}}}}),
			output:     false,
		},
		{ // 12 - LineString with point on the line
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{0, 0}, {-1, -1}}},
			other:      NewGeoJsonLinestring([][]float64{{-1, 0}, {1, 0}}),
			output:     true,
		},
		{ // 13 - LineString with point on the vertex
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1, 0}, {4, 4}}},
			other:      NewGeoJsonLinestring([][]float64{{-1, 0}, {1, 0}}),
			output:     true,
		},
		{ // 14 - LineString with point not on line
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{4, 4}, {2, 3}}},
			other:      NewGeoJsonLinestring([][]float64{{-1, 0}, {1, 0}}),
			output:     false,
		},
		{ // 15 - MultiLineString with point on the line
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{-2, 0}, {4, 4}}},
			other:      NewGeoJsonMultilinestring([][][]float64{{{-5, 0}, {-3, 0}}, {{-2, 0}, {2, 0}}}),
			output:     true,
		},
		{ // 16 - MultiLineString with point on the vertex
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{4, 4}, {-2, 1}}},
			other:      NewGeoJsonMultilinestring([][][]float64{{{-1, 0}, {1, 0}}, {{-2, 1}, {2, 1}}}),
			output:     true,
		},
		{ // 17 - MultiLineString with point not on line
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1, -1}, {4, 4}}},
			other:      NewGeoJsonMultilinestring([][][]float64{{{-1, 0}, {1, 0}}, {{-2, 1}, {2, 1}}}),
			output:     false,
		},
		{ // 18 - Circle with point not on the inside
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{4, 4}, {-1, -3}}},
			other:      NewGeoCircle([]float64{0, 0}, "1km"),
			output:     false,
		},
		{ // 19 - Circle with point on the inside
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{0.024, -0.037}, {4, 4}}},
			other:      NewGeoCircle([]float64{0, 0}, "10km"),
			output:     true,
		},
		{ // 20 - Envelope with point on the inside
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{4, 4}, {0, 0}}},
			other:      NewGeoEnvelope([][]float64{{-2, 2}, {2, -2}}),
			output:     true,
		},
		{ // 21 - Envelope with point on the outside
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{-2, -3}, {4, 4}}},
			other:      NewGeoEnvelope([][]float64{{-2, 2}, {2, -2}}),
			output:     false,
		},
		{ // 22 - Envelope with point on the edge
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{4, 4}, {-1, -2}}},
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
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234567, 1.234567891234567}),
			output:     true,
		},
		{ // 1 - Point with 15th decimal place differing
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234568, 1.234567891234567}),
			output:     true,
		},
		{ // 2 - Point with 13th decimal place differing
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonPoint([]float64{1.234567891234667, 1.234567891234567}),
			output:     false,
		},
		{ // 3 - MultiPoint with a match
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
			other:      NewGeoJsonMultiPoint([][]float64{{2.234567891234567, 2.234567891234567}, {1.234567891234567, 1.234567891234567}}),
			output:     true,
		},
		{ // 4 - MultiPoint with no match
			queryPoint: &MultiPoint{Typ: PointType, Vertices: [][]float64{{1.234567891234567, 1.234567891234567}, {2.234567891234567, 2.234567891234567}}},
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
