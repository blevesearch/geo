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
)

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

// ---------------------------------------------------------------------------
// geo shape v2 cell tests

func TestLineStringCells(t *testing.T) {
	ls := NewGeoJsonLinestring([][]float64{{0, 0}, {10, 10}, {20, 10}}).(*LineString)

	inner, cross := ls.IndexCells()
	if len(inner) != 0 {
		t.Fatalf("expected no inner cells for a linestring, got %v", inner)
	}
	if len(cross) == 0 {
		t.Fatal("expected cross cells for a linestring, got none")
	}
	verifyCellPartition(t, ls.pl, inner, cross)

	// every vertex of the line must be covered; vertices are (lng, lat)
	for _, v := range ls.Vertices {
		if !cellsCoverLatLng(cross, v[1], v[0]) {
			t.Fatalf("cells do not cover vertex (%v, %v)", v[0], v[1])
		}
	}

	// the query covering must also be all-cross, and at least as fine
	qInner, qCross := ls.QueryCells()
	if len(qInner) != 0 {
		t.Fatalf("expected no inner query cells for a linestring, got %v", qInner)
	}
	if len(qCross) < len(cross) {
		t.Fatalf("expected the query covering (%d cells) to be at least as "+
			"fine as the index covering (%d cells)", len(qCross), len(cross))
	}
}

func TestLineStringBoundingBox(t *testing.T) {
	ls := NewGeoJsonLinestring([][]float64{{0, 0}, {10, 10}, {20, 10}}).(*LineString)

	env, ok := ls.BoundingBox().(*Envelope)
	if !ok || env.r == nil {
		t.Fatalf("expected an envelope bounding box, got %v", ls.BoundingBox())
	}
	for _, v := range ls.Vertices {
		if !rectContainsDegrees(*env.r, v[1], v[0]) {
			t.Fatalf("bounding box does not contain vertex (%v, %v)", v[0], v[1])
		}
	}
}

func TestMultiLineStringCells(t *testing.T) {
	// two segments far apart: the union covering must cover both
	mls := NewGeoJsonMultilinestring([][][]float64{
		{{0, 0}, {5, 5}},
		{{50, 50}, {55, 55}},
	}).(*MultiLineString)

	inner, cross := mls.IndexCells()
	if len(inner) != 0 {
		t.Fatalf("expected no inner cells for a multilinestring, got %v", inner)
	}
	for _, ll := range [][2]float64{{0, 0}, {5, 5}, {50, 50}, {55, 55}} {
		if !cellsCoverLatLng(cross, ll[0], ll[1]) {
			t.Fatalf("cells do not cover segment endpoint (%v, %v)", ll[0], ll[1])
		}
	}

	qInner, qCross := mls.QueryCells()
	if len(qInner) != 0 {
		t.Fatalf("expected no inner query cells, got %v", qInner)
	}
	if len(qCross) < len(cross) {
		t.Fatalf("expected the query covering (%d cells) to be at least as "+
			"fine as the index covering (%d cells)", len(qCross), len(cross))
	}
}

func TestMultiLineStringBoundingBox(t *testing.T) {
	mls := NewGeoJsonMultilinestring([][][]float64{
		{{0, 0}, {5, 5}},
		{{50, 50}, {55, 55}},
	}).(*MultiLineString)

	env, ok := mls.BoundingBox().(*Envelope)
	if !ok || env.r == nil {
		t.Fatalf("expected an envelope bounding box, got %v", mls.BoundingBox())
	}
	for _, ll := range [][2]float64{{0, 0}, {5, 5}, {50, 50}, {55, 55}} {
		if !rectContainsDegrees(*env.r, ll[0], ll[1]) {
			t.Fatalf("bounding box does not contain endpoint (%v, %v)", ll[0], ll[1])
		}
	}
}

func TestEmptyMultiLineStringCells(t *testing.T) {
	mls := &MultiLineString{Typ: MultiLineStringType}

	inner, cross := mls.IndexCells()
	if len(inner) != 0 || len(cross) != 0 {
		t.Fatalf("expected no cells for an empty multilinestring, got %v %v",
			inner, cross)
	}

	env, ok := mls.BoundingBox().(*Envelope)
	if !ok || env.r == nil {
		t.Fatalf("expected an envelope bounding box, got %v", mls.BoundingBox())
	}
	if !env.r.IsEmpty() {
		t.Fatalf("expected an empty rect for an empty multilinestring, got %v", env.r)
	}
}
