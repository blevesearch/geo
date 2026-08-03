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
	"reflect"
	"testing"

	index "github.com/blevesearch/bleve_index_api"
)

// TestMarshalExtractRoundTrip marshals every shape type and decodes it back
// through ExtractShapesFromBytes, verifying the type prefix dispatch.
func TestMarshalExtractRoundTrip(t *testing.T) {
	shapes := []index.GeoJSON{
		NewGeoJsonPoint([]float64{4.5, 22.5}),
		NewGeoJsonMultiPoint([][]float64{{1, 1}, {50, 50}}),
		NewGeoJsonLinestring([][]float64{{0, 0}, {10, 10}}),
		NewGeoJsonMultilinestring([][][]float64{{{0, 0}, {5, 5}}, {{50, 50}, {55, 55}}}),
		NewGeoJsonPolygon([][][]float64{testSquare(0, 10)}),
		NewGeoJsonMultiPolygon([][][][]float64{{testSquare(0, 10)}, {testSquare(40, 50)}}),
		NewGeoCircle([]float64{10, 10}, "100km"),
		NewGeoEnvelope([][]float64{{0, 20}, {20, 0}}),
	}

	for _, shape := range shapes {
		s, ok := shape.(s2Serializable)
		if !ok {
			t.Fatalf("%T does not implement s2Serializable", shape)
		}
		data, err := s.Marshal()
		if err != nil {
			t.Fatalf("%T: marshal failed: %v", shape, err)
		}

		var reader *bytes.Reader
		got, err := ExtractShapesFromBytes(data, &reader, nil)
		if err != nil {
			t.Fatalf("%T: extract failed: %v", shape, err)
		}
		if reflect.TypeOf(got) != reflect.TypeOf(shape) {
			t.Fatalf("expected %T after the round trip, got %T", shape, got)
		}

		// the decoded shape carries only s2 state; it must still intersect
		// the original shape it was derived from
		if ok, err := got.Intersects(shape); err != nil || !ok {
			t.Fatalf("%T: expected the decoded shape to intersect the "+
				"original, got %v %v", shape, ok, err)
		}
	}
}

func TestExtractShapesFromBytesUnknownPrefix(t *testing.T) {
	var reader *bytes.Reader
	if _, err := ExtractShapesFromBytes([]byte{0xFF, 0x01}, &reader, nil); err == nil {
		t.Fatal("expected an error for an unknown shape prefix")
	}
}

func TestParseGeoJSONShape(t *testing.T) {
	tests := []struct {
		input    string
		wantType string
	}{
		{`{"type": "point", "coordinates": [1, 1]}`, PointType},
		{`{"type": "MultiPoint", "coordinates": [[1, 1], [2, 2]]}`, MultiPointType},
		{`{"type": "LineString", "coordinates": [[0, 0], [1, 1]]}`, LineStringType},
		{`{"type": "multilinestring", "coordinates": [[[0, 0], [1, 1]]]}`, MultiLineStringType},
		{`{"type": "Polygon", "coordinates": [[[0, 0], [1, 0], [1, 1], [0, 1], [0, 0]]]}`, PolygonType},
		{`{"type": "MULTIPOLYGON", "coordinates": [[[[0, 0], [1, 0], [1, 1], [0, 1], [0, 0]]]]}`, MultiPolygonType},
		{`{"type": "circle", "coordinates": [1, 1], "radius": "10km"}`, CircleType},
		{`{"type": "envelope", "coordinates": [[0, 1], [1, 0]]}`, EnvelopeType},
		{`{"type": "geometrycollection", "geometries": [{"type": "point", "coordinates": [1, 1]}]}`, GeometryCollectionType},
	}

	for i, test := range tests {
		shape, err := ParseGeoJSONShape([]byte(test.input))
		if err != nil {
			t.Fatalf("case %d: unexpected error: %v", i, err)
		}
		if shape.Type() != test.wantType {
			t.Fatalf("case %d: expected type %q, got %q", i, test.wantType, shape.Type())
		}
	}

	if _, err := ParseGeoJSONShape([]byte(`{"type": "hexagon", "coordinates": []}`)); err == nil {
		t.Fatal("expected an error for an unknown shape type")
	}
}

func TestFilterGeoShapesOnRelation(t *testing.T) {
	// the document holds a point inside the query polygon
	docShape := NewGeoJsonPoint([]float64{5, 5}).(*Point)
	docBytes, err := docShape.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	queryShape := NewGeoJsonPolygon([][][]float64{testSquare(0, 10)})

	tests := []struct {
		relation string
		output   bool
	}{
		{"intersects", true},
		// the doc's point cannot contain the query polygon
		{"contains", false},
		// the query polygon contains the doc's point
		{"within", true},
		{"disjoint", false},
	}

	for _, test := range tests {
		var reader *bytes.Reader
		got, err := FilterGeoShapesOnRelation(queryShape, docBytes, test.relation, &reader, nil)
		if err != nil {
			t.Fatalf("relation %q: unexpected error: %v", test.relation, err)
		}
		if got != test.output {
			t.Fatalf("relation %q: expected %v, got %v", test.relation, test.output, got)
		}
	}

	var reader *bytes.Reader
	if _, err := FilterGeoShapesOnRelation(queryShape, docBytes, "overlapping", &reader, nil); err == nil {
		t.Fatal("expected an error for an unknown relation")
	}
}
