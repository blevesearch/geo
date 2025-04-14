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
		{ // 8 - Polygon with point on the vertex
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{-0.5, -1}, {4, 4}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}}),
			output:     true,
		},
		{ // 9 - Polygon with point in the hole
			queryPoint: &MultiPoint{Typ: MultiPointType, Vertices: [][]float64{{0, 0}, {4, 4}}},
			other:      NewGeoJsonPolygon([][][]float64{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}, {{-0.5, -0.5}, {-0.5, 0.5}, {0.5, 0.5}, {0.5, -0.5}, {-0.5, -0.5}}}),
			output:     false,
		},
		{ // 10 - MulitiPolygon with point
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

func TestLineStringIntersects(t *testing.T) {
	tests := []struct {
		query  *LineString
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not on the line
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPoint([]float64{1, 1}),
			output: false,
		},
		{ // 1 - Point on edge
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPoint([]float64{0, 0}),
			output: true,
		},
		{ // 2 - Point on inner vertex
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPoint([]float64{2, 3}),
			output: true,
		},
		{ // 3 - Point on outer vertex
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPoint([]float64{0, 3}),
			output: true,
		},
		{ // 4 - Multipoint with one intersection
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultiPoint([][]float64{{1, 0}, {1, 1}}),
			output: true,
		},
		{ // 5 - Multipoint with no intersection
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultiPoint([][]float64{{2, 2}, {1, 1}}),
			output: false,
		},
		{ // 6 - Polygon with one vertex overlap
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPolygon([][][]float64{{{1, 0}, {1, -1}, {2, -1}, {2, 0}, {1, 0}}}),
			output: true,
		},
		{ // 7 - Polygon with one edge overlap
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 0}, {1, -1}, {2, -1}, {2, 0}, {-1, 0}}}),
			output: true,
		},
		{ // 8 - Polygon with no vertex overlap, but crossing edge
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 1}, {-5, 5}, {-5, -5}, {5, -5}, {-1, 1}}}),
			output: true,
		},
		{ // 9 - Polygon containing linestring
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-5, 5}, {-5, -5}, {5, -5}, {5, 5}, {-5, 5}}}),
			output: true,
		},
		{ // 10 - Polygon with no intersection
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-5, 5}, {5, 5}, {5, -5}, {-5, -5}, {-5, 5}}}),
			output: false,
		},
		{ // 11 - Multipolygon with one vertex overlap
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{1, 0}, {1, -1}, {2, -1}, {2, 0}, {1, 0}}}, {{{5, 5}, {4, 5}, {4, 4}, {5, 4}, {5, 5}}}}),
			output: true,
		},
		{ // 12 - Multipolygon with one edge overlap
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{5, 5}, {4, 5}, {4, 4}, {5, 4}, {5, 5}}}, {{{-1, 0}, {1, -1}, {2, -1}, {2, 0}, {-1, 0}}}}),
			output: true,
		},
		{ // 13 - Multipolygon with no vertex overlap, but crossing edge
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-1, 1}, {-5, 5}, {-5, -5}, {5, -5}, {-1, 1}}}, {{{5, 5}, {4, 5}, {4, 4}, {5, 4}, {5, 5}}}}),
			output: true,
		},
		{ // 14 - Multipolygon containing linestring
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{5, 5}, {4, 5}, {4, 4}, {5, 4}, {5, 5}}}, {{{-5, 5}, {-5, -5}, {5, -5}, {5, 5}, {-5, 5}}}}),
			output: true,
		},
		{ // 15 - Multipolygon with no intersection
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-5, 5}, {5, 5}, {5, -5}, {-5, -5}, {-5, 5}}}, {{{5, 5}, {4, 5}, {4, 4}, {5, 4}, {5, 5}}}}),
			output: false,
		},
		{ // 16 - Linestring with one vertex overlap
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonLinestring([][]float64{{2, 3}, {3, 3}, {4, 3}}),
			output: true,
		},
		{ // 17 - Linestring with one edge overlap
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonLinestring([][]float64{{2, 3}, {1, 0}, {1, -1}}),
			output: true,
		},
		{ // 18 - Linestring overlapping but no vertex overlap
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonLinestring([][]float64{{-2, 0}, {2, 0}, {2, 2}}),
			output: true,
		},
		{ // 19 - Linestring with intersection
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {2, 0}, {2, 2}}),
			output: true,
		},
		{ // 20 - Linestring with no intersection
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {0, 5}, {5, 5}}),
			output: false,
		},
		{ // 21 - Multilinestring with one vertex overlap
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{2, 3}, {3, 3}, {4, 3}}}),
			output: true,
		},
		{ // 22 - Multilinestring with one edge overlap
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{2, 3}, {1, 0}, {1, -1}}, {{5, 5}, {6, 6}, {5, 6}}}),
			output: true,
		},
		{ // 23 - Multilinestring with intersection
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 4}, {2, 0}, {2, 2}}}),
			output: true,
		},
		{ // 24 - Multilinestring with no intersection
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{0, 4}, {0, 5}, {5, 5}}, {{5, 5}, {6, 6}, {5, 6}}}),
			output: false,
		},
		{ // 25 - Circle with intersection
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoCircle([]float64{1, 1}, "100km"),
			output: true,
		},
		{ // 26 - Circle with no intersection
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoCircle([]float64{0, 1}, "10km"),
			output: false,
		},
		{ // 27 - Envelope with one vertex overlap
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoEnvelope([][]float64{{1, 0}, {2, -2}}),
			output: true,
		},
		{ // 28 - Envelope with one edge overlap
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoEnvelope([][]float64{{-2, 0}, {2, -2}}),
			output: true,
		},
		{ // 29 - Envelope containing linestring
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoEnvelope([][]float64{{-5, 5}, {5, -5}}),
			output: true,
		},
		{ // 30 - Envelope with no intersection
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoEnvelope([][]float64{{-5, 5}, {-4, 4}}),
			output: false,
		},
	}

	for i, test := range tests {
		result, err := test.query.Intersects(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestMultiLineStringIntersects(t *testing.T) {
	tests := []struct {
		query  *MultiLineString
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not on the line
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonPoint([]float64{1, 1}),
			output: false,
		},
		{ // 1 - Point on edge
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonPoint([]float64{0, 0}),
			output: true,
		},
		{ // 2 - Point on inner vertex
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonPoint([]float64{2, 3}),
			output: true,
		},
		{ // 3 - Point on outer vertex
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonPoint([]float64{0, 3}),
			output: true,
		},
		{ // 4 - Multipoint with one intersection
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{1, 0}, {1, 1}}),
			output: true,
		},
		{ // 5 - Multipoint with no intersection
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{2, 2}, {1, 1}}),
			output: false,
		},
		{ // 6 - Polygon with one vertex overlap
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{1, 0}, {1, -1}, {2, -1}, {2, 0}, {1, 0}}}),
			output: true,
		},
		{ // 7 - Polygon with one edge overlap
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 0}, {1, -1}, {2, -1}, {2, 0}, {-1, 0}}}),
			output: true,
		},
		{ // 8 - Polygon with no vertex overlap, but crossing edge
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 1}, {-5, 5}, {-5, -5}, {5, -5}, {-1, 1}}}),
			output: true,
		},
		{ // 9 - Polygon containing linestring
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-5, 5}, {-5, -5}, {5, -5}, {5, 5}, {-5, 5}}}),
			output: true,
		},
		{ // 10 - Polygon with no intersection
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{5, 5}, {4, 5}, {4, 4}, {5, 4}, {5, 5}}}),
			output: false,
		},
		{ // 11 - Multipolygon with one vertex overlap
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{1, 0}, {1, -1}, {2, -1}, {2, 0}, {1, 0}}}, {{{5, 5}, {4, 5}, {4, 4}, {5, 4}, {5, 5}}}}),
			output: true,
		},
		{ // 12 - Multipolygon with one edge overlap
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{5, 5}, {4, 5}, {4, 4}, {5, 4}, {5, 5}}}, {{{-1, 0}, {1, -1}, {2, -1}, {2, 0}, {-1, 0}}}}),
			output: true,
		},
		{ // 13 - Multipolygon with no vertex overlap, but crossing edge
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-1, 1}, {-5, 5}, {-5, -5}, {5, -5}, {-1, 1}}}, {{{5, 5}, {4, 5}, {4, 4}, {5, 4}, {5, 5}}}}),
			output: true,
		},
		{ // 14 - Multipolygon containing linestring
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{5, 5}, {4, 5}, {4, 4}, {5, 4}, {5, 5}}}, {{{-5, 5}, {-5, -5}, {5, -5}, {5, 5}, {-5, 5}}}}),
			output: true,
		},
		{ // 15 - Multipolygon with no intersection
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{6, 6}, {5, 6}, {5, 5}, {6, 5}, {6, 6}}}, {{{5, 5}, {4, 5}, {4, 4}, {5, 4}, {5, 5}}}}),
			output: false,
		},
		{ // 16 - Linestring with one vertex overlap
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonLinestring([][]float64{{2, 3}, {3, 3}, {4, 3}}),
			output: true,
		},
		{ // 17 - Linestring with one edge overlap
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonLinestring([][]float64{{2, 3}, {1, 0}, {1, -1}}),
			output: true,
		},
		{ // 18 - Linestring overlapping but no vertex overlap
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonLinestring([][]float64{{-2, 0}, {2, 0}, {2, 2}}),
			output: true,
		},
		{ // 19 - Linestring with intersection
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {2, 0}, {2, 2}}),
			output: true,
		},
		{ // 20 - Linestring with no intersection
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {0, 5}, {5, 5}}),
			output: false,
		},
		{ // 21 - Multilinestring with one vertex overlap
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{2, 3}, {3, 3}, {4, 3}}}),
			output: true,
		},
		{ // 22 - Multilinestring with one edge overlap
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{2, 3}, {1, 0}, {1, -1}}, {{5, 5}, {6, 6}, {5, 6}}}),
			output: true,
		},
		{ // 23 - Multilinestring with intersection
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 4}, {2, 0}, {2, 2}}}),
			output: true,
		},
		{ // 24 - Multilinestring with no intersection
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{0, 4}, {0, 5}, {5, 5}}, {{5, 5}, {6, 6}, {5, 6}}}),
			output: false,
		},
		{ // 25 - Circle with intersection
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoCircle([]float64{1, 1}, "100km"),
			output: true,
		},
		{ // 26 - Circle with no intersection
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoCircle([]float64{0, 1}, "10km"),
			output: false,
		},
		{ // 27 - Envelope with one vertex overlap
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoEnvelope([][]float64{{1, 0}, {2, -2}}),
			output: true,
		},
		{ // 28 - Envelope with one edge overlap
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoEnvelope([][]float64{{-2, 0}, {2, -2}}),
			output: true,
		},
		{ // 29 - Envelope containing linestring
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoEnvelope([][]float64{{-5, 5}, {5, -5}}),
			output: true,
		},
		{ // 30 - Envelope with no intersection
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoEnvelope([][]float64{{-5, 5}, {-4, 4}}),
			output: false,
		},
	}

	for i, test := range tests {
		result, err := test.query.Intersects(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestLineStringContains(t *testing.T) {
	tests := []struct {
		query  *LineString
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not on the line
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPoint([]float64{1, 1}),
			output: false,
		},
		{ // 1 - Point on edge
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPoint([]float64{0, 0}),
			output: true,
		},
		{ // 2 - Point on inner vertex
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPoint([]float64{2, 3}),
			output: true,
		},
		{ // 3 - Point on outer vertex
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonPoint([]float64{0, 3}),
			output: true,
		},
		{ // 4 - Multipoint with two intersecting points
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultiPoint([][]float64{{0, 0}, {0, 3}}),
			output: true,
		},
		{ // 5 - Multipoint with one intersecting point
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultiPoint([][]float64{{0, 0}, {1, 1}}),
			output: false,
		},
		{ // 6 - Multipoint with no intersecting point
			query:  &LineString{Typ: LineStringType, Vertices: [][]float64{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}},
			other:  NewGeoJsonMultiPoint([][]float64{{2, 2}, {1, 1}}),
			output: false,
		},
	}

	for i, test := range tests {
		result, err := test.query.Contains(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestMultiLineStringContains(t *testing.T) {
	tests := []struct {
		query  *MultiLineString
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not on the line
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonPoint([]float64{1, 1}),
			output: false,
		},
		{ // 1 - Point on edge
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonPoint([]float64{0, 0}),
			output: true,
		},
		{ // 2 - Point on inner vertex
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonPoint([]float64{2, 3}),
			output: true,
		},
		{ // 3 - Point on outer vertex
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonPoint([]float64{0, 3}),
			output: true,
		},
		{ // 4 - Multipoint with two intersecting points
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{0, 0}, {0, 3}}),
			output: true,
		},
		{ // 5 - Multipoint with one intersecting point
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{100, 101}, {102, 103}, {104, 105}}, {{-1, 0}, {1, 0}, {2, 3}, {0, 3}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{0, 0}, {1, 1}}),
			output: false,
		},
		{ // 6 - Multipoint with no intersecting point
			query:  &MultiLineString{Typ: MultiLineStringType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}}, {{100, 101}, {102, 103}, {104, 105}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{2, 2}, {1, 1}}),
			output: false,
		},
	}

	for i, test := range tests {
		result, err := test.query.Contains(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

// Erratic edge cases not covered due to performance concerns
// 2, 16, 17, 22, 23, 32
func TestPolygonIntersects(t *testing.T) {
	tests := []struct {
		query  *Polygon
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not in polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPoint([]float64{5, 5}),
			output: false,
		},
		{ // 1 - Point on vertex
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPoint([]float64{-1, 0}),
			output: true,
		},
		{ // 2 - Point on edge
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPoint([]float64{0, 0}),
			output: false,
		},
		{ // 3 - Point inside polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPoint([]float64{1, 1}),
			output: true,
		},
		{ // 4 - Multipoint with one point inside polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {1, 2}}),
			output: true,
		},
		{ // 5 - Multipoint with no points inside polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {4, 2}}),
			output: false,
		},
		{ // 6 - Linestring with one vertex overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 3}, {3, 3}, {4, 3}}),
			output: true,
		},
		{ // 7 - Linestring with one edge overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 3}, {2, 3}, {4, 3}}),
			output: true,
		},
		{ // 8 - Linestring with intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {2, 0}, {4, 3}}),
			output: true,
		},
		{ // 9 - Linestring contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonLinestring([][]float64{{2, 2}, {1, 1}, {0, 1}}),
			output: true,
		},
		{ // 10 - Linestring with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {4, 4}, {4, 3}}),
			output: false,
		},
		{ // 11 - Multilinestring with one vertex overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 3}, {3, 3}, {4, 3}}}),
			output: true,
		},
		{ // 12 - Multilinestring with one edge overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{0, 3}, {2, 3}, {4, 3}}, {{5, 5}, {6, 6}, {5, 6}}}),
			output: true,
		},
		{ // 13 - Multilinestring with intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 4}, {2, 0}, {4, 3}}}),
			output: true,
		},
		{ // 14 - Multilinestring contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{2, 2}, {1, 1}, {0, 1}}, {{5, 5}, {6, 6}, {5, 6}}}),
			output: true,
		},
		{ // 15 - Multilinestring with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 4}, {4, 4}, {4, 3}}}),
			output: false,
		},
		{ // 16 - Polygon with one vertex overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 0}, {-1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}),
			output: false,
		},
		{ // 17 - Polygon with one edge overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 0}, {-1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}),
			output: false,
		},
		{ // 18 - Polygon with intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 0}, {1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}),
			output: true,
		},
		{ // 19 - Polygon containing polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-5, -5}, {5, -5}, {5, 5}, {-5, 5}, {-5, -5}}}),
			output: true,
		},
		{ // 20 - Polygon with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}),
			output: false,
		},
		{ // 21 - Polygon contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{0, 1}, {0.5, 1}, {0.5, 1.5}, {0, 1.5}, {0, 1}}}),
			output: true,
		},
		{ // 22 - MultiPolygon with one vertex overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}, {{{-1, 0}, {-1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}}),
			output: false,
		},
		{ // 23 - MultiPolygon with one edge overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-1, 0}, {-1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}, {{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}}),
			output: false,
		},
		{ // 24 - MultiPolygon with intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}, {{{-1, 0}, {1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}}),
			output: true,
		},
		{ // 25 - MultiPolygon contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-5, -5}, {5, -5}, {5, 5}, {-5, 5}, {-5, -5}}}, {{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}}),
			output: true,
		},
		{ // 26 - MultiPolygon with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}, {{{5, 6}, {6, 5}, {6, 6}, {5, 6}, {5, 5}}}}),
			output: false,
		},
		{ // 27 - MultiPolygon containing polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{0, 1}, {0.5, 1}, {0.5, 1.5}, {0, 1.5}, {0, 1}}}, {{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}}),
			output: true,
		},
		{ // 28 - Circle with overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoCircle([]float64{1, 0}, "100km"),
			output: true,
		},
		{ // 29 - Circle with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoCircle([]float64{5, 0}, "100km"),
			output: false,
		},
		{ // 30 - Circle containing polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoCircle([]float64{1, 1}, "100000km"),
			output: true,
		},
		{ // 31 - Circle contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoCircle([]float64{0.5, 1}, "1km"),
			output: true,
		},
		{ // 32 - Envelope with one vertex overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoEnvelope([][]float64{{1, 0}, {2, -2}}),
			output: false,
		},
		{ // 33 - Envelope with one edge overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoEnvelope([][]float64{{-1, 0}, {2, -2}}),
			output: true,
		},
		{ // 34 - Envelope with intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoEnvelope([][]float64{{-1, 1}, {2, -2}}),
			output: true,
		},
		{ // 35 - Envelope contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoEnvelope([][]float64{{0.5, 1}, {0.75, 0.5}}),
			output: true,
		},
		{ // 36 - Envelope with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoEnvelope([][]float64{{5, 5}, {6, 4}}),
			output: false,
		},
		{ // 37 - Envelope containing polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoEnvelope([][]float64{{-5, 5}, {5, -5}}),
			output: true,
		},
	}

	for i, test := range tests {
		result, err := test.query.Intersects(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestMultiPolygonIntersects(t *testing.T) {
	tests := []struct {
		query  *MultiPolygon
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not in multipolygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonPoint([]float64{5, 5}),
			output: false,
		},
		{ // 1 - Point on vertex
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonPoint([]float64{-1, 0}),
			output: true,
		},
		{ // 2 - Point on edge
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonPoint([]float64{0, 0}),
			output: false,
		},
		{ // 3 - Point inside multipolygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonPoint([]float64{1, 1}),
			output: true,
		},
		{ // 4 - Multipoint with one point inside multipolygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {1, 2}}),
			output: true,
		},
		{ // 5 - Multipoint with no points inside multipolygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {4, 2}}),
			output: false,
		},
		{ // 6 - Linestring with one vertex overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 3}, {3, 3}, {4, 3}}),
			output: true,
		},
		{ // 7 - Linestring with one edge overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 3}, {2, 3}, {4, 3}}),
			output: true,
		},
		{ // 8 - Linestring with intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {2, 0}, {4, 3}}),
			output: true,
		},
		{ // 9 - Linestring contained by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonLinestring([][]float64{{2, 2}, {1, 1}, {0, 1}}),
			output: true,
		},
		{ // 10 - Linestring with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {4, 4}, {4, 3}}),
			output: false,
		},
		{ // 11 - Multilinestring with one vertex overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 3}, {3, 3}, {4, 3}}}),
			output: true,
		},
		{ // 12 - Multilinestring with one edge overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{0, 3}, {2, 3}, {4, 3}}, {{5, 5}, {6, 6}, {5, 6}}}),
			output: true,
		},
		{ // 13 - Multilinestring with intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 4}, {2, 0}, {4, 3}}}),
			output: true,
		},
		{ // 14 - Multilinestring contained by multipolygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{2, 2}, {1, 1}, {0, 1}}, {{5, 5}, {6, 6}, {5, 6}}}),
			output: true,
		},
		{ // 15 - Multilinestring with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 4}, {4, 4}, {4, 3}}}),
			output: false,
		},
		{ // 16 - Polygon with one vertex overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 0}, {-1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}),
			output: false,
		},
		{ // 17 - Polygon with one edge overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 0}, {-1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}),
			output: false,
		},
		{ // 18 - Polygon with intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 0}, {1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}),
			output: true,
		},
		{ // 19 - Polygon containing multipolygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-5, -5}, {5, -5}, {5, 5}, {-5, 5}, {-5, -5}}}),
			output: true,
		},
		{ // 20 - Polygon with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}),
			output: false,
		},
		{ // 21 - Polygon contained by multipolygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{0, 1}, {0.5, 1}, {0.5, 1.5}, {0, 1.5}, {0, 1}}}),
			output: true,
		},
		{ // 22 - MultiPolygon with one vertex overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}, {{{-1, 0}, {-1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}}),
			output: false,
		},
		{ // 23 - MultiPolygon with one edge overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-1, 0}, {-1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}, {{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}}),
			output: false,
		},
		{ // 24 - MultiPolygon with intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}, {{{-1, 0}, {1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}}),
			output: true,
		},
		{ // 25 - MultiPolygon containing multipolygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-5, -5}, {5, -5}, {5, 5}, {-5, 5}, {-5, -5}}}, {{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}}),
			output: true,
		},
		{ // 26 - MultiPolygon with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}, {{{5, 6}, {6, 5}, {6, 6}, {5, 6}, {5, 5}}}}),
			output: false,
		},
		{ // 27 - MultiPolygon contained by multipolygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{0, 1}, {0.5, 1}, {0.5, 1.5}, {0, 1.5}, {0, 1}}}, {{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}}),
			output: true,
		},
		{ // 28 - Circle with overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoCircle([]float64{1, 0}, "100km"),
			output: true,
		},
		{ // 29 - Circle with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoCircle([]float64{5, 0}, "100km"),
			output: false,
		},
		{ // 30 - Circle containing polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoCircle([]float64{1, 1}, "100000km"),
			output: true,
		},
		{ // 31 - Circle contained by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoCircle([]float64{0.5, 1}, "1km"),
			output: true,
		},
		{ // 32 - Envelope with one vertex overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoEnvelope([][]float64{{1, 0}, {2, -2}}),
			output: false,
		},
		{ // 33 - Envelope with one edge overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoEnvelope([][]float64{{-1, 0}, {2, -2}}),
			output: true,
		},
		{ // 34 - Envelope with intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoEnvelope([][]float64{{-1, 1}, {2, -2}}),
			output: true,
		},
		{ // 35 - Envelope contained by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoEnvelope([][]float64{{0.5, 1}, {0.75, 0.5}}),
			output: true,
		},
		{ // 36 - Envelope with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoEnvelope([][]float64{{5, 5}, {6, 4}}),
			output: false,
		},
		{ // 37 - Envelope containing polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoEnvelope([][]float64{{-5, 5}, {5, -5}}),
			output: true,
		},
	}

	for i, test := range tests {
		result, err := test.query.Intersects(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestPolygonContains(t *testing.T) {
	tests := []struct {
		query  *Polygon
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not in polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPoint([]float64{5, 5}),
			output: false,
		},
		{ // 1 - Point inside polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPoint([]float64{1, 1}),
			output: true,
		},
		{ // 2 - Multipoint with one point inside polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {1, 2}}),
			output: false,
		},
		{ // 3 - Multipoint with no points inside polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {4, 2}}),
			output: false,
		},
		{ // 4 - Multipoint with all points inside polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{1, 1}, {1, 2}}),
			output: true,
		},
		{ // 5 - Linestring with intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {2, 0}, {4, 3}}),
			output: false,
		},
		{ // 6 - Linestring contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonLinestring([][]float64{{2, 2}, {1, 1}, {0, 1}}),
			output: true,
		},
		{ // 7 - Linestring with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {4, 4}, {4, 3}}),
			output: false,
		},
		{ // 8 - Multilinestring with intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 4}, {2, 0}, {4, 3}}}),
			output: false,
		},
		{ // 9 - Multilinestring with one linestring contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{2, 2}, {1, 1}, {0, 1}}, {{5, 5}, {6, 6}, {5, 6}}}),
			output: false,
		},
		{ // 10 - Multilinestring with both linestrings contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{2, 2}, {1, 1}, {0, 1}}, {{0.5, 0.5}, {0, 1}, {0.5, 1.5}}}),
			output: true,
		},
		{ // 11 - Multilinestring with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 4}, {4, 4}, {4, 3}}}),
			output: false,
		},
		{ // 12 - Polygon with intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 0}, {1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}),
			output: false,
		},
		{ // 13 - Polygon containing by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-5, -5}, {5, -5}, {5, 5}, {-5, 5}, {-5, -5}}}),
			output: false,
		},
		{ // 14 - Polygon with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}),
			output: false,
		},
		{ // 15 - Polygon contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{0, 1}, {0.5, 1}, {0.5, 1.5}, {0, 1.5}, {0, 1}}}),
			output: true,
		},
		{ // 16 - Polygon with same polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}),
			output: true,
		},
		{ // 17 - MultiPolygon with intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}, {{{-1, 0}, {1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}}),
			output: false,
		},
		{ // 18 - MultiPolygon containing by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-5, -5}, {5, -5}, {5, 5}, {-5, 5}, {-5, -5}}}, {{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}}),
			output: false,
		},
		{ // 19 - MultiPolygon with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}, {{{5, 6}, {6, 5}, {6, 6}, {5, 6}, {5, 5}}}}),
			output: false,
		},
		{ // 20 - MultiPolygon contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{0, 1}, {0.5, 1}, {0.5, 1.5}, {0, 1.5}, {0, 1}}}, {{{1, 1}, {1.1, 1}, {1.1, 1.1}, {1, 1.1}, {1, 1}}}}),
			output: true,
		},
		{ // 21 - Circle with overlap
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoCircle([]float64{1, 0}, "100km"),
			output: false,
		},
		{ // 22 - Circle with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoCircle([]float64{5, 0}, "100km"),
			output: false,
		},
		{ // 23 - Circle containing polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoCircle([]float64{1, 1}, "100000km"),
			output: false,
		},
		{ // 24 - Circle contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoCircle([]float64{0.5, 1}, "1km"),
			output: true,
		},
		{ // 25 - Envelope with intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoEnvelope([][]float64{{-1, 1}, {2, -2}}),
			output: false,
		},
		{ // 26 - Envelope contained by polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoEnvelope([][]float64{{0.5, 1}, {0.75, 0.5}}),
			output: true,
		},
		{ // 27 - Envelope with no intersection
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoEnvelope([][]float64{{5, 5}, {6, 4}}),
			output: false,
		},
		{ // 28 - Envelope containing polygon
			query:  &Polygon{Typ: PolygonType, Vertices: [][][]float64{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}},
			other:  NewGeoEnvelope([][]float64{{-5, 5}, {5, -5}}),
			output: false,
		},
	}

	for i, test := range tests {
		result, err := test.query.Contains(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestMultiPolygonContains(t *testing.T) {
	tests := []struct {
		query  *MultiPolygon
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not in polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonPoint([]float64{5, 5}),
			output: false,
		},
		{ // 1 - Point inside polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonPoint([]float64{1, 1}),
			output: true,
		},
		{ // 2 - Multipoint with one point inside polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {1, 2}}),
			output: false,
		},
		{ // 3 - Multipoint with no points inside polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {4, 2}}),
			output: false,
		},
		{ // 4 - Multipoint with all points inside polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultiPoint([][]float64{{1, 1}, {1, 2}}),
			output: true,
		},
		{ // 5 - Linestring with intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {2, 0}, {4, 3}}),
			output: false,
		},
		{ // 6 - Linestring contained by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonLinestring([][]float64{{2, 2}, {1, 1}, {0, 1}}),
			output: true,
		},
		{ // 7 - Linestring with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonLinestring([][]float64{{0, 4}, {4, 4}, {4, 3}}),
			output: false,
		},
		{ // 8 - Multilinestring with intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 4}, {2, 0}, {4, 3}}}),
			output: false,
		},
		{ // 9 - Multilinestring with one linestring contained by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{2, 2}, {1, 1}, {0, 1}}, {{5, 5}, {6, 6}, {5, 6}}}),
			output: false,
		},
		{ // 10 - Multilinestring with both linestrings contained by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{2, 2}, {1, 1}, {0, 1}}, {{0.5, 0.5}, {0, 1}, {0.5, 1.5}}}),
			output: true,
		},
		{ // 11 - Multilinestring with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {6, 6}, {5, 6}}, {{0, 4}, {4, 4}, {4, 3}}}),
			output: false,
		},
		{ // 12 - Polygon with intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-1, 0}, {1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}),
			output: false,
		},
		{ // 13 - Polygon containing by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-5, -5}, {5, -5}, {5, 5}, {-5, 5}, {-5, -5}}}),
			output: false,
		},
		{ // 14 - Polygon with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}),
			output: false,
		},
		{ // 15 - Polygon contained by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonPolygon([][][]float64{{{0, 1}, {0.5, 1}, {0.5, 1.5}, {0, 1.5}, {0, 1}}}),
			output: true,
		},
		{ // 16 - MultiPolygon with intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}, {{{-1, 0}, {1, 1}, {-2, -1}, {-2, 0}, {-1, 0}}}}),
			output: false,
		},
		{ // 17 - MultiPolygon containing by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-5, -5}, {5, -5}, {5, 5}, {-5, 5}, {-5, -5}}}, {{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}}),
			output: false,
		},
		{ // 18 - MultiPolygon with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}}, {{{5, 6}, {6, 5}, {6, 6}, {5, 6}, {5, 5}}}}),
			output: false,
		},
		{ // 19 - MultiPolygon contained by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{0, 1}, {0.5, 1}, {0.5, 1.5}, {0, 1.5}, {0, 1}}}, {{{1, 1}, {1.1, 1}, {1.1, 1.1}, {1, 1.1}, {1, 1}}}}),
			output: true,
		},
		{ // 20 - MultiPolygon with exact same multipolygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}),
			output: true,
		},
		{ // 21 - Circle with overlap
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoCircle([]float64{1, 0}, "100km"),
			output: false,
		},
		{ // 22 - Circle with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoCircle([]float64{5, 0}, "100km"),
			output: false,
		},
		{ // 23 - Circle containing polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoCircle([]float64{1, 1}, "100000km"),
			output: false,
		},
		{ // 24 - Circle contained by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoCircle([]float64{0.5, 1}, "1km"),
			output: true,
		},
		{ // 25 - Envelope with intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoEnvelope([][]float64{{-1, 1}, {2, -2}}),
			output: false,
		},
		{ // 26 - Envelope contained by polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoEnvelope([][]float64{{0.5, 1}, {0.75, 0.5}}),
			output: true,
		},
		{ // 27 - Envelope with no intersection
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}, {{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}}},
			other:  NewGeoEnvelope([][]float64{{5, 5}, {6, 4}}),
			output: false,
		},
		{ // 28 - Envelope containing polygon
			query:  &MultiPolygon{Typ: MultiPolygonType, Vertices: [][][][]float64{{{{-1, 0}, {1, 0}, {2, 3}, {0, 3}, {-1, 0}}}, {{{100, 100}, {100, 101}, {101, 101}, {101, 100}, {100, 100}}}}},
			other:  NewGeoEnvelope([][]float64{{-5, 5}, {5, -5}}),
			output: false,
		},
	}

	for i, test := range tests {
		result, err := test.query.Contains(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestCircleIntersects(t *testing.T) {
	tests := []struct {
		query  *Circle
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not in circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPoint([]float64{5, 5}),
			output: false,
		},
		{ // 1 - Point inside circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPoint([]float64{1.2, 1.2}),
			output: true,
		},
		{ // 2 - Multipoint with one point inside circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {0.8, 0.8}}),
			output: true,
		},
		{ // 3 - Multipoint with no points inside circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {8, 8}}),
			output: false,
		},
		{ // 4 - Multipoint with all points inside circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultiPoint([][]float64{{1.1, 1.1}, {0.8, 0.8}}),
			output: true,
		},
		{ // 5 - Linestring with intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonLinestring([][]float64{{5, 5}, {1.2, 0.8}}),
			output: true,
		},
		{ // 6 - Linestring contained by circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonLinestring([][]float64{{0.8, 0.8}, {1.2, 1.2}}),
			output: true,
		},
		{ // 7 - Linestring with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonLinestring([][]float64{{5, 5}, {8, 8}}),
			output: false,
		},
		{ // 8 - Multilinestring with intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {0.8, 0.8}}, {{-5, -5}, {-2, -4}}}),
			output: true,
		},
		{ // 9 - Multilinestring with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultilinestring([][][]float64{{{-5, -5}, {-2, -4}}, {{5, 5}, {8, 7}}}),
			output: false,
		},
		{ // 10 - Polygon with intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPolygon([][][]float64{{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}}}),
			output: true,
		},
		{ // 11 - Polygon contained by circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPolygon([][][]float64{{{0.9, 0.9}, {1.1, 0.9}, {1.1, 1.1}, {0.9, 1.1}, {0.9, 0.9}}}),
			output: true,
		},
		{ // 12 - Polygon containing circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPolygon([][][]float64{{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}}}),
			output: true,
		},
		{ // 13 - Polygon with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPolygon([][][]float64{{{-5, -5}, {-4, -5}, {-4, -4}, {-5, -4}, {-5, -5}}}),
			output: false,
		},
		{ // 14 - MultiPolygon with intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}}}, {{{-5, -5}, {-4, -5}, {-4, -4}, {-5, -4}, {-5, -5}}}}),
			output: true,
		},
		{ // 15 - MultiPolygon with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-4, -4}, {-3, -4}, {-3, -3}, {-4, -3}, {-4, -4}}}, {{{-5, -5}, {-4, -5}, {-4, -4}, {-5, -4}, {-5, -5}}}}),
			output: false,
		},
		{ // 16 - Circle with intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoCircle([]float64{1.5, 1.5}, "100km"),
			output: true,
		},
		{ // 17 - Circle contained by circle
			query:  NewGeoCircle([]float64{1, 1}, "100000km").(*Circle),
			other:  NewGeoCircle([]float64{1.5, 1.5}, "100km"),
			output: true,
		},
		{ // 18 - Circle containing circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoCircle([]float64{1.5, 1.5}, "100000km"),
			output: true,
		},
		{ // 19 - Circle with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "1km").(*Circle),
			other:  NewGeoCircle([]float64{1.5, 1.5}, "1km"),
			output: false,
		},
		{ // 20 - Envelope with intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoEnvelope([][]float64{{0, 2}, {2, 0}}),
			output: true,
		},
		{ // 21 - Envelope with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoEnvelope([][]float64{{4, 6}, {6, 4}}),
			output: false,
		},
	}

	for i, test := range tests {
		result, err := test.query.Intersects(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

// Erratic edge cases not covered due to performance concerns
// 13
func TestCircleContains(t *testing.T) {
	tests := []struct {
		query  *Circle
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not in circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPoint([]float64{5, 5}),
			output: false,
		},
		{ // 1 - Point inside circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPoint([]float64{1.2, 1.2}),
			output: true,
		},
		{ // 2 - Multipoint with one point inside circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {0.8, 0.8}}),
			output: false,
		},
		{ // 3 - Multipoint with no points inside circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {8, 8}}),
			output: false,
		},
		{ // 4 - Multipoint with all points inside circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultiPoint([][]float64{{1.1, 1.1}, {0.8, 0.8}}),
			output: true,
		},
		{ // 5 - Linestring with intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonLinestring([][]float64{{5, 5}, {1.2, 0.8}}),
			output: false,
		},
		{ // 6 - Linestring contained by circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonLinestring([][]float64{{0.8, 0.8}, {1.2, 1.2}}),
			output: true,
		},
		{ // 7 - Linestring with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonLinestring([][]float64{{5, 5}, {8, 8}}),
			output: false,
		},
		{ // 8 - Multilinestring contained by circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultilinestring([][][]float64{{{0.8, 0.8}, {1.2, 1.2}}, {{0.8, 1.2}, {1.2, 0.8}}}),
			output: true,
		},
		{ // 9 - Multilinestring with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultilinestring([][][]float64{{{-5, -5}, {-2, -4}}, {{5, 5}, {8, 7}}}),
			output: false,
		},
		{ // 10 - Polygon contained by circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPolygon([][][]float64{{{0.9, 0.9}, {1.1, 0.9}, {1.1, 1.1}, {0.9, 1.1}, {0.9, 0.9}}}),
			output: true,
		},
		{ // 11 - Polygon containing circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPolygon([][][]float64{{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}}}),
			output: false,
		},
		{ // 12 - Polygon with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPolygon([][][]float64{{{-5, -5}, {-4, -5}, {-4, -4}, {-5, -4}, {-5, -5}}}),
			output: false,
		},
		{ // 13 - Clockwise Polygon within circle but not contained by it
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonPolygon([][][]float64{{{0.9, 0.9}, {0.9, 1.1}, {1.1, 1.1}, {1.1, 0.9}, {0.9, 0.9}}}),
			output: true,
		},
		{ // 14 - MultiPolygon contained by circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{0.9, 0.9}, {1.1, 0.9}, {1.1, 1.1}, {0.9, 1.1}, {0.9, 0.9}}}, {{{0.8, 0.8}, {0.9, 0.8}, {0.9, 0.9}, {0.9, 0.8}, {0.8, 0.8}}}}),
			output: true,
		},
		{ // 15 - MultiPolygon with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-4, -4}, {-3, -4}, {-3, -3}, {-4, -3}, {-4, -4}}}, {{{-5, -5}, {-4, -5}, {-4, -4}, {-5, -4}, {-5, -5}}}}),
			output: false,
		},
		{ // 16 - Circle contained by circle
			query:  NewGeoCircle([]float64{1, 1}, "100000km").(*Circle),
			other:  NewGeoCircle([]float64{1.5, 1.5}, "100km"),
			output: true,
		},
		{ // 17 - Circle containing circle
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoCircle([]float64{1.5, 1.5}, "100000km"),
			output: false,
		},
		{ // 18 - Circle with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "1km").(*Circle),
			other:  NewGeoCircle([]float64{1.5, 1.5}, "1km"),
			output: false,
		},
		{ // 19 - Envelope contained by circle
			query:  NewGeoCircle([]float64{1, 1}, "100000km").(*Circle),
			other:  NewGeoEnvelope([][]float64{{0, 2}, {2, 0}}),
			output: true,
		},
		{ // 20 - Envelope with no intersection
			query:  NewGeoCircle([]float64{1, 1}, "100km").(*Circle),
			other:  NewGeoEnvelope([][]float64{{4, 6}, {6, 4}}),
			output: false,
		},
	}

	for i, test := range tests {
		result, err := test.query.Contains(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestEnvelopeIntersects(t *testing.T) {
	tests := []struct {
		query  *Envelope
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not in envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonPoint([]float64{5, 5}),
			output: false,
		},
		{ // 1 - Point inside envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonPoint([]float64{1.2, 1.2}),
			output: true,
		},
		{ // 2 - Multipoint with one point inside envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {1.8, 1.8}}),
			output: true,
		},
		{ // 3 - Multipoint with no points inside envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {8, 8}}),
			output: false,
		},
		{ // 4 - Multipoint with all points inside envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultiPoint([][]float64{{1.1, 1.1}, {1.8, 1.8}}),
			output: true,
		},
		{ // 5 - Linestring with intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonLinestring([][]float64{{5, 5}, {1.2, 1.8}}),
			output: true,
		},
		{ // 6 - Linestring contained by envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonLinestring([][]float64{{1.8, 1.8}, {1.2, 1.2}}),
			output: true,
		},
		{ // 7 - Linestring with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonLinestring([][]float64{{5, 5}, {8, 8}}),
			output: false,
		},
		{ // 8 - Multilinestring with intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{5, 5}, {1.8, 1.8}}, {{-5, -5}, {-2, -4}}}),
			output: true,
		},
		{ // 9 - Multilinestring with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{-5, -5}, {-2, -4}}, {{5, 5}, {8, 7}}}),
			output: false,
		},
		{ // 10 - Polygon with intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonPolygon([][][]float64{{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}}}),
			output: true,
		},
		{ // 11 - Polygon contained by envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonPolygon([][][]float64{{{1.1, 1.1}, {1.2, 1.1}, {1.2, 1.2}, {1.1, 1.2}, {1.1, 1.1}}}),
			output: true,
		},
		{ // 12 - Polygon containing envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonPolygon([][][]float64{{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}}}),
			output: true,
		},
		{ // 13 - Polygon with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-5, -5}, {-4, -5}, {-4, -4}, {-5, -4}, {-5, -5}}}),
			output: false,
		},
		{ // 14 - MultiPolygon with intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}}}, {{{-5, -5}, {-4, -5}, {-4, -4}, {-5, -4}, {-5, -5}}}}),
			output: true,
		},
		{ // 15 - MultiPolygon with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-4, -4}, {-3, -4}, {-3, -3}, {-4, -3}, {-4, -4}}}, {{{-5, -5}, {-4, -5}, {-4, -4}, {-5, -4}, {-5, -5}}}}),
			output: false,
		},
		{ // 16 - Circle with intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoCircle([]float64{1.5, 1.5}, "100km"),
			output: true,
		},
		{ // 17 - Circle with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoCircle([]float64{2.5, 2.5}, "1km"),
			output: false,
		},
		{ // 18 - Envelope with intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoEnvelope([][]float64{{0, 2}, {2, 0}}),
			output: true,
		},
		{ //  - Envelope with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoEnvelope([][]float64{{4, 6}, {6, 4}}),
			output: false,
		},
	}

	for i, test := range tests {
		result, err := test.query.Intersects(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}

func TestEnvelopeContains(t *testing.T) {
	tests := []struct {
		query  *Envelope
		other  index.GeoJSON
		output bool
	}{
		{ // 0 - Point not in envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonPoint([]float64{5, 5}),
			output: false,
		},
		{ // 1 - Point inside envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonPoint([]float64{1.2, 1.2}),
			output: true,
		},
		{ // 2 - Multipoint with one point inside envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {1.8, 1.8}}),
			output: false,
		},
		{ // 3 - Multipoint with no points inside envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultiPoint([][]float64{{5, 5}, {8, 8}}),
			output: false,
		},
		{ // 4 - Multipoint with all points inside envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultiPoint([][]float64{{1.1, 1.1}, {1.8, 1.8}}),
			output: true,
		},
		{ // 5 - Linestring with intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonLinestring([][]float64{{5, 5}, {1.2, 1.8}}),
			output: false,
		},
		{ // 6 - Linestring contained by envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonLinestring([][]float64{{1.8, 1.8}, {1.2, 1.2}}),
			output: true,
		},
		{ // 7 - Linestring with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonLinestring([][]float64{{5, 5}, {8, 8}}),
			output: false,
		},
		{ // 8 - Multilinestring contained by envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{1.8, 1.8}, {1.2, 1.2}}, {{1.8, 1.2}, {1.2, 1.8}}}),
			output: true,
		},
		{ // 9 - Multilinestring with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultilinestring([][][]float64{{{-5, -5}, {-2, -4}}, {{5, 5}, {8, 7}}}),
			output: false,
		},
		{ // 10 - Polygon contained by envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonPolygon([][][]float64{{{1.1, 1.1}, {1.2, 1.1}, {1.2, 1.2}, {1.1, 1.2}, {1.1, 1.1}}}),
			output: true,
		},
		{ // 11 - Polygon with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonPolygon([][][]float64{{{-5, -5}, {-4, -5}, {-4, -4}, {-5, -4}, {-5, -5}}}),
			output: false,
		},
		{ // 12 - MultiPolygon contained by envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{1.1, 1.1}, {1.2, 1.1}, {1.2, 1.2}, {1.1, 1.2}, {1.1, 1.1}}}, {{{1.2, 1.2}, {1.3, 1.2}, {1.3, 1.3}, {1.2, 1.3}, {1.2, 1.2}}}}),
			output: true,
		},
		{ // 13 - MultiPolygon with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoJsonMultiPolygon([][][][]float64{{{{-4, -4}, {-3, -4}, {-3, -3}, {-4, -3}, {-4, -4}}}, {{{-5, -5}, {-4, -5}, {-4, -4}, {-5, -4}, {-5, -5}}}}),
			output: false,
		},
		{ // 14 - Circle contained by envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoCircle([]float64{1.5, 1.5}, "1km"),
			output: true,
		},
		{ // 15 - Circle with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoCircle([]float64{2.5, 2.5}, "1km"),
			output: false,
		},
		{ // 16 - Envelope contained by envelope
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoEnvelope([][]float64{{1.5, 1.25}, {1.25, 1.5}}),
			output: true,
		},
		{ // 17 - Envelope with no intersection
			query:  &Envelope{Typ: EnvelopeType, Vertices: [][]float64{{2, 1}, {1, 2}}},
			other:  NewGeoEnvelope([][]float64{{4, 6}, {6, 4}}),
			output: false,
		},
	}

	for i, test := range tests {
		result, err := test.query.Contains(test.other)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if result != test.output {
			t.Errorf("Test - %d, expected %v, got %v", i, test.output, result)
		}
	}
}
