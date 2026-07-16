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
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	index "github.com/blevesearch/bleve_index_api"
	"github.com/blevesearch/geo/s2"
)

// Polygon represents the geoJSON polygon type
// and it implements the index.GeoJSON interface.
type Polygon struct {
	Typ      string        `json:"type"`
	Vertices [][][]float64 `json:"coordinates"`
	s2pgn    *s2.Polygon
}

func NewGeoJsonPolygon(points [][][]float64) index.GeoJSON {
	rv := &Polygon{Typ: PolygonType, Vertices: points}
	rv.init()
	return rv
}

func (p *Polygon) init() {
	if p.s2pgn == nil {
		p.s2pgn = s2PolygonFromCoordinates(p.Vertices)
	}
}

func (p *Polygon) Type() string {
	return strings.ToLower(p.Typ)
}

func (p *Polygon) Value() ([]byte, error) {
	return jsoniter.Marshal(p)
}

func (p *Polygon) Marshal() ([]byte, error) {
	p.init()

	var b bytes.Buffer
	b.Grow(128)
	w := bufio.NewWriter(&b)
	err := p.s2pgn.Encode(w)
	if err != nil {
		return nil, err
	}

	w.Flush()
	return append([]byte{PolygonTypePrefix}, b.Bytes()...), nil
}

func (p *Polygon) Intersects(other index.GeoJSON) (bool, error) {
	// make an s2polygon for reuse.
	p.init()

	return checkPolygonIntersectsShape(p.s2pgn, p, other)
}

func (p *Polygon) Contains(other index.GeoJSON) (bool, error) {
	// make an s2polygon for reuse.
	p.init()

	return checkMultiPolygonContainsShape([]*s2.Polygon{p.s2pgn}, p, other)
}

func (p *Polygon) Coordinates() [][][]float64 {
	return p.Vertices
}

func (pg *Polygon) Cells() (inner, cross []uint64) {
	if pg.s2pgn == nil {
		return nil, nil
	}
	return indexCellsFromRegion(pg.s2pgn)
}

func (pg *Polygon) QueryCells() ([]uint64, []uint64) {
	if pg.s2pgn == nil {
		return nil, nil
	}
	return queryCellsFromRegion(pg.s2pgn)
}

func (pg *Polygon) BoundingBox() index.GeoJSON {
	if pg.s2pgn == nil {
		return nil
	}
	return envelopeFromRect(pg.s2pgn.RectBound())
}

// --------------------------------------------------------
// MultiPolygon represents the geoJSON multipolygon type
// and it implements the index.GeoJSON interface as well as the
// compositeShap interface.
type MultiPolygon struct {
	Typ      string          `json:"type"`
	Vertices [][][][]float64 `json:"coordinates"`
	s2pgns   []*s2.Polygon
}

func NewGeoJsonMultiPolygon(points [][][][]float64) index.GeoJSON {
	rv := &MultiPolygon{Typ: MultiPolygonType, Vertices: points}
	rv.init()
	return rv
}

func (p *MultiPolygon) init() {
	if p.s2pgns == nil {
		p.s2pgns = make([]*s2.Polygon, len(p.Vertices))
		for i, vertices := range p.Vertices {
			pgn := s2PolygonFromCoordinates(vertices)
			p.s2pgns[i] = pgn
		}
	}
}

func (p *MultiPolygon) Type() string {
	return strings.ToLower(p.Typ)
}

func (p *MultiPolygon) Value() ([]byte, error) {
	return jsoniter.Marshal(p)
}

func (p *MultiPolygon) Marshal() ([]byte, error) {
	p.init()

	var b bytes.Buffer
	b.Grow(512)
	w := bufio.NewWriter(&b)

	// first write the number of polygons.
	count := int32(len(p.s2pgns))
	err := binary.Write(w, binary.BigEndian, count)
	if err != nil {
		return nil, err
	}
	// write the polygons.
	for _, pgn := range p.s2pgns {
		err := pgn.Encode(w)
		if err != nil {
			return nil, err
		}
	}

	w.Flush()
	return append([]byte{MultiPolygonTypePrefix}, b.Bytes()...), nil
}

func (p *MultiPolygon) Intersects(other index.GeoJSON) (bool, error) {
	p.init()

	for _, pgn := range p.s2pgns {
		rv, err := checkPolygonIntersectsShape(pgn, p, other)
		if rv && err == nil {
			return true, nil
		}
	}

	return false, nil
}

func (p *MultiPolygon) Contains(other index.GeoJSON) (bool, error) {
	p.init()

	return checkMultiPolygonContainsShape(p.s2pgns, p, other)
}

func (p *MultiPolygon) Coordinates() [][][][]float64 {
	return p.Vertices
}

func (p *MultiPolygon) Members() []index.GeoJSON {
	if len(p.Vertices) > 0 && len(p.s2pgns) == 0 {
		polygons := make([]index.GeoJSON, len(p.Vertices))
		for pos, vertices := range p.Vertices {
			polygons[pos] = NewGeoJsonPolygon(vertices)
		}
		return polygons
	}

	polygons := make([]index.GeoJSON, len(p.s2pgns))
	for pos, pgn := range p.s2pgns {
		polygons[pos] = &Polygon{s2pgn: pgn}
	}
	return polygons
}

func (mp *MultiPolygon) Cells() (inner, cross []uint64) {
	ru := make(s2.RegionUnion, 0, len(mp.s2pgns))
	for _, pg := range mp.s2pgns {
		if pg != nil {
			ru = append(ru, pg)
		}
	}
	if len(ru) == 0 {
		return nil, nil
	}
	return indexCellsFromRegion(ru)
}

func (mp *MultiPolygon) QueryCells() ([]uint64, []uint64) {
	ru := make(s2.RegionUnion, 0, len(mp.s2pgns))
	for _, pg := range mp.s2pgns {
		if pg != nil {
			ru = append(ru, pg)
		}
	}
	if len(ru) == 0 {
		return nil, nil
	}
	return queryCellsFromRegion(ru)
}

func (mp *MultiPolygon) BoundingBox() index.GeoJSON {
	r := s2.EmptyRect()
	for _, pg := range mp.s2pgns {
		if pg != nil {
			r = r.Union(pg.RectBound())
		}
	}
	return envelopeFromRect(r)
}

// checkPolygonIntersectsShape checks the intersection between the
// s2 polygon and the other shapes in the documents.
func checkPolygonIntersectsShape(s2pgn *s2.Polygon, shapeIn,
	other index.GeoJSON) (bool, error) {
	// check if the other shape is a point.
	if p2, ok := other.(*Point); ok {
		if polygonsIntersectsPoint([]*s2.Polygon{s2pgn}, p2.s2point) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a multipoint.
	if p2, ok := other.(*MultiPoint); ok {
		for _, s2point := range p2.s2points {
			if polygonsIntersectsPoint([]*s2.Polygon{s2pgn}, s2point) {
				return true, nil
			}
		}

		return false, nil
	}

	// check if the other shape is a polygon.
	if p2, ok := other.(*Polygon); ok {
		if s2pgn.Intersects(p2.s2pgn) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a multipolygon.
	if p2, ok := other.(*MultiPolygon); ok {
		// check the intersection for any polygon in the collection.
		for _, s2pgn1 := range p2.s2pgns {
			if s2pgn.Intersects(s2pgn1) {
				return true, nil
			}
		}

		return false, nil
	}

	// check if the other shape is a linestring.
	if ls, ok := other.(*LineString); ok {
		if polylineIntersectsPolygons([]*s2.Polyline{ls.pl},
			[]*s2.Polygon{s2pgn}) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a multilinestring.
	if mls, ok := other.(*MultiLineString); ok {
		if polylineIntersectsPolygons(mls.pls, []*s2.Polygon{s2pgn}) {
			return true, nil
		}

		return false, nil
	}

	if gc, ok := other.(*GeometryCollection); ok {
		// check whether the polygon intersects with any of the
		// member shapes of the geometry collection.
		if geometryCollectionIntersectsShape(gc, shapeIn) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a circle.
	if c, ok := other.(*Circle); ok {
		cp := c.s2cap.Center()
		radius := c.s2cap.Radius()

		projected := s2pgn.Project(&cp)
		distance := projected.Distance(cp)

		return distance <= radius, nil
	}

	// check if the other shape is a envelope.
	if e, ok := other.(*Envelope); ok {
		s2pgnInDoc := s2PolygonFromS2Rectangle(e.r)
		if s2pgn.Intersects(s2pgnInDoc) {
			return true, nil
		}
		return false, nil
	}

	return false, fmt.Errorf("unknown geojson type: %s "+
		" found in document", other.Type())
}

// checkMultiPolygonContainsShape checks whether the given polygons
// collectively contains the shape in the document.
func checkMultiPolygonContainsShape(s2pgns []*s2.Polygon,
	shapeIn, other index.GeoJSON) (bool, error) {
	// check if the other shape is a point.
	if p2, ok := other.(*Point); ok {
		if polygonsIntersectsPoint(s2pgns, p2.s2point) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a multipoint.
	if p2, ok := other.(*MultiPoint); ok {
		// check the containment for every point in the collection.
		idx := s2.NewShapeIndex()
		for _, s2pgn := range s2pgns {
			idx.Add(s2pgn)
		}

		for _, point := range p2.s2points {
			if !s2.NewContainsPointQuery(idx, s2.VertexModelClosed).Contains(*point) {
				return false, nil
			}
		}

		return true, nil
	}

	// check if the other shape is a polygon.
	if p2, ok := other.(*Polygon); ok {
		for _, s2pgn := range s2pgns {
			if s2pgn.Contains(p2.s2pgn) {
				return true, nil
			}
		}

		return false, nil
	}

	// check if the other shape is a multipolygon.
	if p2, ok := other.(*MultiPolygon); ok {
		// check the intersection for every polygon in the collection.
		polygonsWithIn := make(map[int]struct{})
	nextPolygon:
		for pgnIndex, pgn := range p2.s2pgns {
			for _, s2pgn := range s2pgns {
				if s2pgn.Contains(pgn) {
					polygonsWithIn[pgnIndex] = struct{}{}
					continue nextPolygon
				}
			}
		}

		return len(p2.s2pgns) == len(polygonsWithIn), nil
	}

	// check if the other shape is a linestring.
	if ls, ok := other.(*LineString); ok {
		if polygonsContainsLineStrings(s2pgns,
			[]*s2.Polyline{ls.pl}) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a multilinestring.
	if mls, ok := other.(*MultiLineString); ok {
		// check whether any of the linestring is inside the polygon.
		if polygonsContainsLineStrings(s2pgns, mls.pls) {
			return true, nil
		}

		return false, nil
	}

	if gc, ok := other.(*GeometryCollection); ok {
		shapesWithIn := make(map[int]struct{})
	nextShape:
		for pos, shape := range gc.Members() {
			for _, s2pgn := range s2pgns {
				contains, err := checkMultiPolygonContainsShape(
					[]*s2.Polygon{s2pgn}, shapeIn, shape)
				if err == nil && contains {
					shapesWithIn[pos] = struct{}{}
					continue nextShape
				}
			}
		}
		return len(shapesWithIn) == len(gc.Members()), nil
	}

	// check if the other shape is a circle.
	if c, ok := other.(*Circle); ok {
		cp := c.s2cap.Center()
		radius := c.s2cap.Radius()

		for _, s2pgn := range s2pgns {
			if s2pgn.ContainsPoint(cp) {
				projected := s2pgn.ProjectToBoundary(&cp)
				distance := projected.Distance(cp)
				if distance >= radius {
					return true, nil
				}
			}
		}

		return false, nil
	}

	// check if the other shape is a envelope.
	if e, ok := other.(*Envelope); ok {
		// create a polygon from the rectangle and checks the containment.
		s2pgnInDoc := s2PolygonFromS2Rectangle(e.r)
		for _, s2pgn := range s2pgns {
			if s2pgn.Contains(s2pgnInDoc) {
				return true, nil
			}
		}

		return false, nil
	}

	return false, fmt.Errorf("unknown geojson type: %s"+
		" found in document", other.Type())
}
