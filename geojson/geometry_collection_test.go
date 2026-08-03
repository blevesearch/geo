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
	"sort"
	"testing"

	index "github.com/blevesearch/bleve_index_api"
)

// stubCellShape is a minimal index.GeoJSON implementation returning fixed
// inner/cross cells, used to test the collection's aggregation logic
// deterministically.
type stubCellShape struct {
	inner []uint64
	cross []uint64
}

func (s *stubCellShape) Type() string { return "stub" }

func (s *stubCellShape) Intersects(other index.GeoJSON) (bool, error) {
	return false, nil
}

func (s *stubCellShape) Contains(other index.GeoJSON) (bool, error) {
	return false, nil
}

func (s *stubCellShape) Value() ([]byte, error) { return nil, nil }

func (s *stubCellShape) IndexCells() (inner, cross []uint64) {
	return s.inner, s.cross
}

func (s *stubCellShape) QueryCells() (inner, cross []uint64) {
	return s.inner, s.cross
}

func (s *stubCellShape) BoundingBox() index.GeoJSON { return nil }

func sortedCopy(cells []uint64) []uint64 {
	rv := append([]uint64{}, cells...)
	sort.Slice(rv, func(i, j int) bool { return rv[i] < rv[j] })
	return rv
}

func TestGeometryCollectionCellAggregation(t *testing.T) {
	// member A reports cell 2 as inner, member B reports the same cell as
	// cross: inner must win. Cells 1 and 3 are duplicated across members
	// and must be deduplicated.
	gc := &GeometryCollection{
		Typ: GeometryCollectionType,
		Shapes: []index.GeoJSON{
			&stubCellShape{inner: []uint64{1, 2}, cross: []uint64{3}},
			&stubCellShape{inner: []uint64{1, 4}, cross: []uint64{2, 3, 5}},
			nil, // nil members must be skipped
		},
	}

	inner, cross := gc.IndexCells()
	if !reflect.DeepEqual(sortedCopy(inner), []uint64{1, 2, 4}) {
		t.Fatalf("expected inner cells [1 2 4], got %v", sortedCopy(inner))
	}
	if !reflect.DeepEqual(sortedCopy(cross), []uint64{3, 5}) {
		t.Fatalf("expected cross cells [3 5], got %v", sortedCopy(cross))
	}

	qInner, qCross := gc.QueryCells()
	if !reflect.DeepEqual(sortedCopy(qInner), sortedCopy(inner)) ||
		!reflect.DeepEqual(sortedCopy(qCross), sortedCopy(cross)) {
		t.Fatalf("expected query cells to match index cells for stub members, "+
			"got %v %v", qInner, qCross)
	}
}

func TestGeometryCollectionCellsEmpty(t *testing.T) {
	for _, gc := range []*GeometryCollection{
		{Typ: GeometryCollectionType},
		{Typ: GeometryCollectionType, Shapes: []index.GeoJSON{nil}},
	} {
		inner, cross := gc.IndexCells()
		if len(inner) != 0 || len(cross) != 0 {
			t.Fatalf("expected no cells for an empty collection, got %v %v",
				inner, cross)
		}

		env, ok := gc.BoundingBox().(*Envelope)
		if !ok || env.r == nil {
			t.Fatalf("expected an envelope bounding box, got %v", gc.BoundingBox())
		}
		if !env.r.IsEmpty() {
			t.Fatalf("expected an empty rect for an empty collection, got %v", env.r)
		}
	}
}

func TestGeometryCollectionCellsRealShapes(t *testing.T) {
	// a polygon and a far-away point: the aggregate covering must cover both
	gc := &GeometryCollection{
		Typ: GeometryCollectionType,
		Shapes: []index.GeoJSON{
			NewGeoJsonPolygon([][][]float64{testSquare(0, 20)}),
			NewGeoJsonPoint([]float64{50, 50}),
		},
	}

	inner, cross := gc.IndexCells()
	all := append(append([]uint64{}, inner...), cross...)
	if !cellsCoverLatLng(all, 10, 10) {
		t.Fatal("covering does not cover the polygon member")
	}
	if !cellsCoverLatLng(cross, 50, 50) {
		t.Fatal("cross cells do not cover the point member")
	}
	if len(inner) == 0 {
		t.Fatal("expected inner cells from the polygon member")
	}
}

func TestGeometryCollectionNestedCells(t *testing.T) {
	innerGC := &GeometryCollection{
		Typ:    GeometryCollectionType,
		Shapes: []index.GeoJSON{NewGeoJsonPoint([]float64{4.5, 22.5})},
	}
	outerGC := &GeometryCollection{
		Typ:    GeometryCollectionType,
		Shapes: []index.GeoJSON{innerGC},
	}

	wantInner, wantCross := innerGC.IndexCells()
	gotInner, gotCross := outerGC.IndexCells()
	if !reflect.DeepEqual(sortedCopy(gotInner), sortedCopy(wantInner)) ||
		!reflect.DeepEqual(sortedCopy(gotCross), sortedCopy(wantCross)) {
		t.Fatalf("expected nested collection cells %v %v, got %v %v",
			wantInner, wantCross, gotInner, gotCross)
	}
}

func TestGeometryCollectionBoundingBox(t *testing.T) {
	gc := &GeometryCollection{
		Typ: GeometryCollectionType,
		Shapes: []index.GeoJSON{
			NewGeoJsonPolygon([][][]float64{testSquare(0, 10)}),
			NewGeoJsonPoint([]float64{50, 50}),
		},
	}

	env, ok := gc.BoundingBox().(*Envelope)
	if !ok || env.r == nil {
		t.Fatalf("expected an envelope bounding box, got %v", gc.BoundingBox())
	}
	for _, ll := range [][2]float64{{0, 0}, {10, 10}, {50, 50}} {
		if !rectContainsDegrees(*env.r, ll[0], ll[1]) {
			t.Fatalf("bounding box does not contain (%v, %v)", ll[0], ll[1])
		}
	}
}

func TestGeometryCollectionIntersectsContains(t *testing.T) {
	gc := &GeometryCollection{
		Typ: GeometryCollectionType,
		Shapes: []index.GeoJSON{
			NewGeoJsonPolygon([][][]float64{testSquare(0, 10)}),
			NewGeoJsonPoint([]float64{50, 50}),
		},
	}

	// a point inside the polygon member intersects the collection
	if ok, err := gc.Intersects(NewGeoJsonPoint([]float64{5, 5})); err != nil || !ok {
		t.Fatalf("expected intersection with a point inside a member, got %v %v", ok, err)
	}
	// a far-away point does not
	if ok, err := gc.Intersects(NewGeoJsonPoint([]float64{-60, -60})); err != nil || ok {
		t.Fatalf("expected no intersection with a far-away point, got %v %v", ok, err)
	}
	// the polygon member contains a point inside it
	if ok, err := gc.Contains(NewGeoJsonPoint([]float64{5, 5})); err != nil || !ok {
		t.Fatalf("expected the collection to contain the point, got %v %v", ok, err)
	}
	// a composite other: both points must be covered by members
	if ok, err := gc.Contains(NewGeoJsonMultiPoint([][]float64{{5, 5}, {50, 50}})); err != nil || !ok {
		t.Fatalf("expected the collection to contain the multipoint, got %v %v", ok, err)
	}
	if ok, err := gc.Contains(NewGeoJsonMultiPoint([][]float64{{5, 5}, {-60, -60}})); err != nil || ok {
		t.Fatalf("expected the collection not to contain the multipoint, got %v %v", ok, err)
	}
}

func TestGeometryCollectionMarshalRoundTrip(t *testing.T) {
	gc := &GeometryCollection{
		Typ: GeometryCollectionType,
		Shapes: []index.GeoJSON{
			NewGeoJsonPoint([]float64{4.5, 22.5}),
			NewGeoJsonPolygon([][][]float64{testSquare(0, 10)}),
			NewGeoCircle([]float64{10, 10}, "100km"),
			NewGeoEnvelope([][]float64{{0, 20}, {20, 0}}),
		},
	}

	data, err := gc.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	var reader *bytes.Reader
	shape, err := ExtractShapesFromBytes(data, &reader, nil)
	if err != nil {
		t.Fatal(err)
	}

	got, ok := shape.(*GeometryCollection)
	if !ok {
		t.Fatalf("expected a geometrycollection, got %T", shape)
	}
	if len(got.Shapes) != len(gc.Shapes) {
		t.Fatalf("expected %d member shapes, got %d", len(gc.Shapes), len(got.Shapes))
	}
	for i, want := range []interface{}{&Point{}, &Polygon{}, &Circle{}, &Envelope{}} {
		if reflect.TypeOf(got.Shapes[i]) != reflect.TypeOf(want) {
			t.Fatalf("expected member %d to be %T, got %T", i, want, got.Shapes[i])
		}
	}
}
