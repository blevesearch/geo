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
	"math"
	"testing"

	"github.com/blevesearch/geo/s2"
)

// ---------------------------------------------------------------------------
// helpers shared by the v2 cell tests across the *_test.go files

func cellSet(cells []uint64) map[uint64]struct{} {
	rv := make(map[uint64]struct{}, len(cells))
	for _, c := range cells {
		rv[c] = struct{}{}
	}
	return rv
}

// verifyCellPartition asserts the inner/cross invariants of a covering:
// inner cells are fully contained in the region, cross cells are not, the
// two sets are disjoint and hold no duplicates.
func verifyCellPartition(t *testing.T, region s2.Region, inner, cross []uint64) {
	t.Helper()

	innerSet := cellSet(inner)
	crossSet := cellSet(cross)
	if len(innerSet) != len(inner) {
		t.Fatalf("inner cells contain duplicates: %v", inner)
	}
	if len(crossSet) != len(cross) {
		t.Fatalf("cross cells contain duplicates: %v", cross)
	}
	for c := range innerSet {
		if _, ok := crossSet[c]; ok {
			t.Fatalf("cell %d present in both inner and cross", c)
		}
	}

	for _, c := range inner {
		if !region.ContainsCell(s2.CellFromCellID(s2.CellID(c))) {
			t.Fatalf("inner cell %d is not contained in the region", c)
		}
	}
	for _, c := range cross {
		if region.ContainsCell(s2.CellFromCellID(s2.CellID(c))) {
			t.Fatalf("cross cell %d is fully contained in the region, "+
				"should be inner", c)
		}
	}
}

// cellsCoverLatLng reports whether any of the cells contains the given
// (lat, lng) location.
func cellsCoverLatLng(cells []uint64, lat, lng float64) bool {
	leaf := s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng))
	for _, c := range cells {
		if s2.CellID(c).Contains(leaf) {
			return true
		}
	}
	return false
}

// testSquare returns a closed CCW square ring from (lo, lo) to (hi, hi),
// in (lng, lat) vertex order.
func testSquare(lo, hi float64) [][]float64 {
	return [][]float64{{lo, lo}, {hi, lo}, {hi, hi}, {lo, hi}, {lo, lo}}
}

// ---------------------------------------------------------------------------

func TestCellsFromRegion(t *testing.T) {
	// an areal region: a polygon spanning (0,0)..(20,20)
	pgn := s2PolygonFromCoordinates([][][]float64{testSquare(0, 20)})

	inner, cross := cellsFromRegion(pgn, regionCovererIndexV2)
	if len(inner) == 0 {
		t.Fatal("expected inner cells for a large polygon, got none")
	}
	if len(cross) == 0 {
		t.Fatal("expected cross cells along the polygon boundary, got none")
	}
	if len(inner)+len(cross) > maxIndexCells {
		t.Fatalf("covering has %d cells, exceeding the coverer's max of %d",
			len(inner)+len(cross), maxIndexCells)
	}
	verifyCellPartition(t, pgn, inner, cross)

	// a point inside the polygon must be covered, one far outside must not be
	if !cellsCoverLatLng(append(inner, cross...), 10, 10) {
		t.Fatal("covering does not cover a point inside the polygon")
	}
	if cellsCoverLatLng(append(inner, cross...), 60, 60) {
		t.Fatal("covering covers a point far outside the polygon")
	}
}

func TestCellsFromRegionArealess(t *testing.T) {
	// a polyline has no area: every covering cell must be a cross cell
	pl := s2.PolylineFromLatLngs([]s2.LatLng{
		s2.LatLngFromDegrees(0, 0),
		s2.LatLngFromDegrees(10, 10),
		s2.LatLngFromDegrees(10, 20),
	})

	inner, cross := cellsFromRegion(pl, regionCovererIndexV2)
	if len(inner) != 0 {
		t.Fatalf("expected no inner cells for a polyline, got %v", inner)
	}
	if len(cross) == 0 {
		t.Fatal("expected cross cells for a polyline, got none")
	}
	verifyCellPartition(t, pl, inner, cross)

	// the polyline's endpoints must be covered
	for _, ll := range [][2]float64{{0, 0}, {10, 10}, {10, 20}} {
		if !cellsCoverLatLng(cross, ll[0], ll[1]) {
			t.Fatalf("covering does not cover polyline point (%v, %v)",
				ll[0], ll[1])
		}
	}
}

func TestIndexVsQueryCells(t *testing.T) {
	// for a large region the query-time coverer (maxQueryCells) must
	// produce a finer covering than the index-time coverer (maxIndexCells)
	pgn := s2PolygonFromCoordinates([][][]float64{testSquare(0, 60)})

	indexInner, indexCross := indexCellsFromRegion(pgn)
	queryInner, queryCross := queryCellsFromRegion(pgn)

	indexTotal := len(indexInner) + len(indexCross)
	queryTotal := len(queryInner) + len(queryCross)

	if indexTotal > maxIndexCells {
		t.Fatalf("index covering has %d cells, max is %d", indexTotal, maxIndexCells)
	}
	if queryTotal > maxQueryCells {
		t.Fatalf("query covering has %d cells, max is %d", queryTotal, maxQueryCells)
	}
	if queryTotal <= indexTotal {
		t.Fatalf("expected the query covering (%d cells) to be finer than "+
			"the index covering (%d cells)", queryTotal, indexTotal)
	}
	verifyCellPartition(t, pgn, indexInner, indexCross)
	verifyCellPartition(t, pgn, queryInner, queryCross)
}

func TestPointCell(t *testing.T) {
	lat, lng := 22.5, 4.5
	pt := s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng))

	cell := s2.CellID(pointCell(pt))
	if !cell.IsValid() {
		t.Fatalf("pointCell returned an invalid cell: %d", cell)
	}
	if cell.Level() != maxCellLevel {
		t.Fatalf("expected a level %d cell, got level %d", maxCellLevel, cell.Level())
	}
	if !cell.Contains(s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng))) {
		t.Fatal("pointCell does not contain the point it was derived from")
	}
}

func TestEnvelopeFromRect(t *testing.T) {
	minLng, minLat, maxLng, maxLat := -10.0, 5.0, 15.0, 25.0
	rect := s2RectFromBounds([]float64{minLng, maxLat}, []float64{maxLng, minLat})

	env := envelopeFromRect(*rect)
	if env.Typ != EnvelopeType {
		t.Fatalf("expected type %q, got %q", EnvelopeType, env.Typ)
	}

	// vertices must follow the [[minLng, maxLat], [maxLng, minLat]] convention
	want := [][]float64{{minLng, maxLat}, {maxLng, minLat}}
	for i := range want {
		for j := range want[i] {
			if math.Abs(env.Vertices[i][j]-want[i][j]) > 1e-9 {
				t.Fatalf("expected vertices %v, got %v", want, env.Vertices)
			}
		}
	}

	// the envelope's rect must round-trip through its own vertices
	env2 := NewGeoEnvelope(env.Vertices).(*Envelope)
	if !rectsApproxEqual(*env2.r, *rect) {
		t.Fatalf("expected rect %v to round-trip, got %v", rect, env2.r)
	}
}

// rectContainsDegrees checks rect containment of a (lat, lng) location,
// round-tripping the location through an s2.Point the same way the shapes'
// RectBound implementations do, to avoid boundary float mismatches.
func rectContainsDegrees(r s2.Rect, lat, lng float64) bool {
	return r.ContainsLatLng(s2.LatLngFromPoint(
		s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng))))
}

func rectsApproxEqual(a, b s2.Rect) bool {
	const eps = 1e-9
	return math.Abs(a.Lo().Lat.Degrees()-b.Lo().Lat.Degrees()) < eps &&
		math.Abs(a.Lo().Lng.Degrees()-b.Lo().Lng.Degrees()) < eps &&
		math.Abs(a.Hi().Lat.Degrees()-b.Hi().Lat.Degrees()) < eps &&
		math.Abs(a.Hi().Lng.Degrees()-b.Hi().Lng.Degrees()) < eps
}
