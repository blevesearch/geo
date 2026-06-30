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
