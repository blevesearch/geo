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
