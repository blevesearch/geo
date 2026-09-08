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
	"bytes"
	"math"
	"testing"

	index "github.com/blevesearch/bleve_index_api"
	"github.com/blevesearch/geo/s2"
)

// A shape that a search reads out of a segment is decoded, not constructed, and
// a decoded shape takes a different path through s2 than a constructed one: it
// carries no ShapeIndex until something asks for one. These tests pin the two
// paths together, so that a relation cannot come out differently, or fail
// outright, only because the shape arrived through a decode.

// testRing returns a closed ring of n points around (cx, cy).
func testRing(n int, cx, cy, r float64) [][]float64 {
	ring := make([][]float64, 0, n+1)
	for i := 0; i < n; i++ {
		theta := 2 * math.Pi * float64(i) / float64(n)
		ring = append(ring, []float64{cx + r*math.Cos(theta), cy + r*math.Sin(theta)})
	}

	return append(ring, ring[0])
}

// roundTrip marshals a shape and decodes it back, the way a shape reaches a
// search from the index.
func roundTrip(t *testing.T, shape index.GeoJSON) index.GeoJSON {
	t.Helper()

	marshaller, ok := shape.(interface{ Marshal() ([]byte, error) })
	if !ok {
		t.Fatalf("shape of type %T cannot be marshalled", shape)
	}

	encoded, err := marshaller.Marshal()
	if err != nil {
		t.Fatalf("marshalling %T: %v", shape, err)
	}

	var reader *bytes.Reader
	decoded, err := ExtractShapesFromBytes(encoded, &reader,
		s2.NewGeoBufferPool(24*1024, 24))
	if err != nil {
		t.Fatalf("decoding %T: %v", shape, err)
	}

	return decoded
}

func TestRelationsOnDecodedShapes(t *testing.T) {
	// A 4-point ring stays under the vertex count at which s2 switches from
	// checking every edge to consulting an index, and a 64-point ring goes over
	// it, so both sides of that switch are covered.
	smallSquare := [][][]float64{{{0, 0}, {4, 0}, {4, 4}, {0, 4}, {0, 0}}}
	overlapping := [][][]float64{{{2, 2}, {6, 2}, {6, 6}, {2, 6}, {2, 2}}}
	disjointSquare := [][][]float64{{{20, 20}, {24, 20}, {24, 24}, {20, 24}, {20, 20}}}
	bigCircle := [][][]float64{testRing(64, 2, 2, 3)}
	withHole := [][][]float64{
		testRing(64, 2, 2, 5),
		testRing(16, 2, 2, 1),
	}

	tests := []struct {
		name  string
		doc   index.GeoJSON
		query index.GeoJSON
	}{
		{"polygon/polygon overlapping", NewGeoJsonPolygon(smallSquare), NewGeoJsonPolygon(overlapping)},
		{"polygon/polygon disjoint", NewGeoJsonPolygon(smallSquare), NewGeoJsonPolygon(disjointSquare)},
		{"polygon/polygon indexed", NewGeoJsonPolygon(bigCircle), NewGeoJsonPolygon(overlapping)},
		{"polygon with hole/polygon", NewGeoJsonPolygon(withHole), NewGeoJsonPolygon(overlapping)},
		{"polygon/point inside", NewGeoJsonPolygon(smallSquare), NewGeoJsonPoint([]float64{2, 2})},
		{"polygon/point outside", NewGeoJsonPolygon(smallSquare), NewGeoJsonPoint([]float64{40, 40})},
		{"polygon/point on vertex", NewGeoJsonPolygon(smallSquare), NewGeoJsonPoint([]float64{0, 0})},
		{"polygon/multipoint", NewGeoJsonPolygon(bigCircle),
			NewGeoJsonMultiPoint([][]float64{{2, 2}, {40, 40}})},
		{"polygon/linestring crossing", NewGeoJsonPolygon(smallSquare),
			NewGeoJsonLinestring([][]float64{{-2, 2}, {6, 2}})},
		{"polygon/linestring inside", NewGeoJsonPolygon(bigCircle),
			NewGeoJsonLinestring([][]float64{{1.5, 2}, {2.5, 2}})},
		{"polygon/multilinestring", NewGeoJsonPolygon(bigCircle),
			NewGeoJsonMultilinestring([][][]float64{{{1.5, 2}, {2.5, 2}}, {{40, 40}, {41, 41}}})},
		{"multipolygon/polygon", NewGeoJsonMultiPolygon([][][][]float64{smallSquare, disjointSquare}),
			NewGeoJsonPolygon(overlapping)},
		{"multipolygon/point", NewGeoJsonMultiPolygon([][][][]float64{smallSquare, disjointSquare}),
			NewGeoJsonPoint([]float64{22, 22})},
		{"linestring/polygon", NewGeoJsonLinestring([][]float64{{-2, 2}, {6, 2}}),
			NewGeoJsonPolygon(smallSquare)},
		{"point/polygon", NewGeoJsonPoint([]float64{2, 2}), NewGeoJsonPolygon(smallSquare)},
		{"envelope/polygon", NewGeoEnvelope([][]float64{{0, 4}, {4, 0}}), NewGeoJsonPolygon(overlapping)},
	}

	relations := []string{"intersects", "contains", "within", "disjoint"}

	for _, test := range tests {
		for _, relation := range relations {
			t.Run(test.name+"/"+relation, func(t *testing.T) {
				want, err := filterShapes(test.query, test.doc, relation)
				if err != nil {
					t.Fatalf("relation %q on constructed shapes: %v", relation, err)
				}

				// The decode is what a search actually feeds to the relation.
				got, err := filterShapes(test.query, roundTrip(t, test.doc), relation)
				if err != nil {
					t.Fatalf("relation %q on decoded shape: %v", relation, err)
				}

				if got != want {
					t.Errorf("relation %q: decoded shape gave %v, constructed shape gave %v",
						relation, got, want)
				}
			})
		}
	}
}

// A decoded polygon has to answer questions that would normally be served by a
// ShapeIndex, so check that asking for one, and asking repeatedly, is consistent
// with a polygon that was constructed.
func TestDecodedPolygonIndexedQueries(t *testing.T) {
	coords := [][][]float64{testRing(128, 0, 0, 10)}

	constructed := NewGeoJsonPolygon(coords)
	if _, err := constructed.Intersects(NewGeoJsonPoint([]float64{0, 0})); err != nil {
		t.Fatal(err)
	}

	decoded := roundTrip(t, NewGeoJsonPolygon(coords))

	probes := []index.GeoJSON{
		NewGeoJsonPoint([]float64{0, 0}),
		NewGeoJsonPoint([]float64{50, 50}),
		NewGeoJsonPolygon([][][]float64{testRing(64, 1, 1, 2)}),
		NewGeoJsonLinestring([][]float64{{-20, 0}, {20, 0}}),
	}

	// Repeat the sweep so that a polygon which has already built its index on
	// demand is exercised as well as one that has not.
	for round := 0; round < 2; round++ {
		for i, probe := range probes {
			want, err := constructed.Intersects(probe)
			if err != nil {
				t.Fatalf("probe %d on constructed polygon: %v", i, err)
			}

			got, err := decoded.Intersects(probe)
			if err != nil {
				t.Fatalf("probe %d on decoded polygon: %v", i, err)
			}

			if got != want {
				t.Errorf("round %d probe %d: decoded gave %v, constructed gave %v",
					round, i, got, want)
			}
		}
	}
}
