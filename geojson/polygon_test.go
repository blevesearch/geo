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
			other:  NewGeoJsonLinestring([][]float64{{1, 2}, {1, 1}, {0, 1}}),
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
			other:  NewGeoJsonMultilinestring([][][]float64{{{1, 2}, {1, 1}, {0, 1}}, {{0.5, 0.5}, {0, 1}, {0.5, 1.5}}}),
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
			other:  NewGeoJsonLinestring([][]float64{{1, 2}, {1, 1}, {0, 1}}),
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
			other:  NewGeoJsonMultilinestring([][][]float64{{{1, 2}, {1, 1}, {0, 1}}, {{0.5, 0.5}, {0, 1}, {0.5, 1.5}}}),
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
