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
		{ // 19 - Envelope with no intersection
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

// ---------------------------------------------------------------------------
// geo shape v2 cell tests

func TestEnvelopeCells(t *testing.T) {
	// envelope vertices are [[minLng, maxLat], [maxLng, minLat]]
	e := NewGeoEnvelope([][]float64{{0, 20}, {20, 0}}).(*Envelope)

	inner, cross := e.IndexCells()
	if len(inner) == 0 {
		t.Fatal("expected inner cells for a large envelope, got none")
	}
	if len(cross) == 0 {
		t.Fatal("expected cross cells along the envelope boundary, got none")
	}
	verifyCellPartition(t, *e.r, inner, cross)

	if !cellsCoverLatLng(append(inner, cross...), 10, 10) {
		t.Fatal("covering does not cover the envelope's center")
	}

	qInner, qCross := e.QueryCells()
	verifyCellPartition(t, *e.r, qInner, qCross)
	if len(qInner)+len(qCross) < len(inner)+len(cross) {
		t.Fatalf("expected the query covering (%d cells) to be at least as "+
			"fine as the index covering (%d cells)",
			len(qInner)+len(qCross), len(inner)+len(cross))
	}
}

func TestEnvelopeBoundingBox(t *testing.T) {
	e := NewGeoEnvelope([][]float64{{0, 20}, {20, 0}}).(*Envelope)

	// an envelope's bounding box is an equivalent envelope
	env, ok := e.BoundingBox().(*Envelope)
	if !ok || env.r == nil {
		t.Fatalf("expected an envelope bounding box, got %v", e.BoundingBox())
	}
	if !rectsApproxEqual(*env.r, *e.r) {
		t.Fatalf("expected bounding box rect %v to equal the envelope's own "+
			"rect %v", env.r, e.r)
	}
}
