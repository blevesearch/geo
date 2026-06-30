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

// --------------------------------------------------------
// Point represents the geoJSON point type and it
// implements the index.GeoJSON interface.
type Point struct {
	Typ      string    `json:"type"`
	Vertices []float64 `json:"coordinates"`
	s2point  *s2.Point
}

func NewGeoJsonPoint(v []float64) index.GeoJSON {
	rv := &Point{Typ: PointType, Vertices: v}
	rv.init()
	return rv
}

func (p *Point) Type() string {
	return strings.ToLower(p.Typ)
}

func (p *Point) Value() ([]byte, error) {
	return jsoniter.Marshal(p)
}

func (p *Point) init() {
	if p.s2point == nil {
		s2point := s2.PointFromLatLng(s2.LatLngFromDegrees(
			p.Vertices[1], p.Vertices[0]))
		p.s2point = &s2point
	}
}

func (p *Point) Marshal() ([]byte, error) {
	p.init()

	var b bytes.Buffer
	b.Grow(32)
	w := bufio.NewWriter(&b)
	err := p.s2point.Encode(w)
	if err != nil {
		return nil, err
	}

	w.Flush()
	return append([]byte{PointTypePrefix}, b.Bytes()...), nil
}

func (p *Point) Intersects(other index.GeoJSON) (bool, error) {
	p.init()

	return checkPointIntersectsShape(p.s2point, p, other)
}

func (p *Point) Contains(other index.GeoJSON) (bool, error) {
	p.init()

	return checkPointContainsShape([]*s2.Point{p.s2point}, other)
}

func (p *Point) Coordinates() []float64 {
	return p.Vertices
}

// Point can only have a single cross cell
func (p *Point) Cells() (inner, cross []uint64) {
	if p.s2point == nil {
		return nil, nil
	}
	return nil, []uint64{pointCell(*p.s2point)}
}

func (p *Point) BoundingBox() index.GeoJSON {
	if p.s2point == nil {
		return nil
	}
	return envelopeFromRect(p.s2point.RectBound())
}

// --------------------------------------------------------
// MultiPoint represents the geoJSON multipoint type and it
// implements the index.GeoJSON interface as well as the
// compositeShap interface.
type MultiPoint struct {
	Typ      string      `json:"type"`
	Vertices [][]float64 `json:"coordinates"`
	s2points []*s2.Point
}

func NewGeoJsonMultiPoint(v [][]float64) index.GeoJSON {
	rv := &MultiPoint{Typ: MultiPointType, Vertices: v}
	rv.init()
	return rv
}

func (mp *MultiPoint) init() {
	if mp.s2points == nil {
		mp.s2points = make([]*s2.Point, len(mp.Vertices))
		for i, point := range mp.Vertices {
			s2point := s2.PointFromLatLng(s2.LatLngFromDegrees(
				point[1], point[0]))
			mp.s2points[i] = &s2point
		}
	}
}

func (p *MultiPoint) Marshal() ([]byte, error) {
	p.init()

	var b bytes.Buffer
	b.Grow(64)
	w := bufio.NewWriter(&b)

	// first write the number of points.
	count := int32(len(p.s2points))
	err := binary.Write(w, binary.BigEndian, count)
	if err != nil {
		return nil, err
	}
	// write the points.
	for _, s2point := range p.s2points {
		err := s2point.Encode(w)
		if err != nil {
			return nil, err
		}
	}

	w.Flush()
	return append([]byte{MultiPointTypePrefix}, b.Bytes()...), nil
}

func (p *MultiPoint) Type() string {
	return strings.ToLower(p.Typ)
}

func (mp *MultiPoint) Value() ([]byte, error) {
	return jsoniter.Marshal(mp)
}

func (p *MultiPoint) Intersects(other index.GeoJSON) (bool, error) {
	p.init()

	for _, s2point := range p.s2points {
		rv, err := checkPointIntersectsShape(s2point, p, other)
		if rv && err == nil {
			return rv, nil
		}
	}

	return false, nil
}

func (p *MultiPoint) Contains(other index.GeoJSON) (bool, error) {
	p.init()

	rv, err := checkPointContainsShape(p.s2points, other)
	if rv && err == nil {
		return rv, nil
	}

	return false, nil
}

func (p *MultiPoint) Coordinates() [][]float64 {
	return p.Vertices
}

func (p *MultiPoint) Members() []index.GeoJSON {
	if len(p.Vertices) > 0 && len(p.s2points) == 0 {
		points := make([]index.GeoJSON, len(p.Vertices))
		for pos, vertices := range p.Vertices {
			points[pos] = NewGeoJsonPoint(vertices)
		}
		return points
	}

	points := make([]index.GeoJSON, len(p.s2points))
	for pos, point := range p.s2points {
		points[pos] = &Point{s2point: point}
	}
	return points
}

// MultiPoints can only have cross cells
func (mp *MultiPoint) Cells() (inner, cross []uint64) {
	cross = make([]uint64, 0, len(mp.s2points))
	for _, pt := range mp.s2points {
		if pt == nil {
			continue
		}
		cross = append(cross, pointCell(*pt))
	}
	return nil, cross
}

func (mp *MultiPoint) BoundingBox() index.GeoJSON {
	r := s2.EmptyRect()
	for _, pt := range mp.s2points {
		if pt == nil {
			continue
		}
		r = r.Union(pt.RectBound())
	}
	return envelopeFromRect(r)
}

// checkPointIntersectsShape checks for intersection between
// the point and the shape in the document.
func checkPointIntersectsShape(point *s2.Point, shapeIn, other index.GeoJSON) (bool, error) {
	// check if the other shape is a point.
	if p2, ok := other.(*Point); ok {
		// check if the points are equal
		if point.ApproxEqual(*p2.s2point) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a multipoint.
	if p2, ok := other.(*MultiPoint); ok {
		// check if any of the points are equal
		for _, p := range p2.s2points {
			if point.ApproxEqual(*p) {
				return true, nil
			}
		}

		return false, nil
	}

	// check if the other shape is a polygon.
	if p2, ok := other.(*Polygon); ok {
		// check if the point is contained within the polygon.
		if polygonsIntersectsPoint([]*s2.Polygon{p2.s2pgn}, point) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a multipolygon.
	if p2, ok := other.(*MultiPolygon); ok {
		// check if the point is contained within any of the polygons
		if polygonsIntersectsPoint(p2.s2pgns, point) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a linestring.
	if p2, ok := other.(*LineString); ok {
		// project the point to the linestring and check if
		// the projection is equal to the point.
		if polylineIntersectsPoint([]*s2.Polyline{p2.pl}, point) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a multilinestring.
	if p2, ok := other.(*MultiLineString); ok {
		// check the intersection for any linestring in the array.
		if polylineIntersectsPoint(p2.pls, point) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a geometrycollection.
	if gc, ok := other.(*GeometryCollection); ok {
		// check for intersection across every member shape.
		if geometryCollectionIntersectsShape(gc, shapeIn) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is a circle.
	if c, ok := other.(*Circle); ok {
		// check if the point is contained within the circle
		// by calculating the distance between the point and the
		// center of the circle.
		if c.s2cap.ContainsPoint(*point) {
			return true, nil
		}

		return false, nil
	}

	// check if the other shape is an envelope.
	if e, ok := other.(*Envelope); ok {
		// check if the point is contained by the envelope
		// by checking if the point is within its bounds
		if e.r.ContainsPoint(*point) {
			return true, nil
		}

		return false, nil
	}

	return false, fmt.Errorf("unknown geojson type: %s "+
		" found in document", other.Type())
}

// checkPointContainsShape checks whether the given shape in
// in the document is approximately contained with the point.
func checkPointContainsShape(points []*s2.Point,
	other index.GeoJSON) (bool, error) {
	// check if the other shape is a point.
	if p2, ok := other.(*Point); ok {
		for _, point := range points {
			if point.ApproxEqual(*p2.s2point) {
				return true, nil
			}
		}

		return false, nil
	}

	// check if the other shape is a multipoint, if so containment is
	// checked for every point in the multipoint with every given point.
	if p2, ok := other.(*MultiPoint); ok {
		// check the containment for every point in the collection.
		lookup := make(map[int]struct{})
		for _, qpoint := range points {
			for pos, dpoint := range p2.s2points {
				if _, done := lookup[pos]; done {
					continue
				}
				// already processed all the points in the multipoint.
				if len(lookup) == len(p2.s2points) {
					return true, nil
				}

				if qpoint.ApproxEqual(*dpoint) {
					lookup[pos] = struct{}{}
				}
			}
		}

		return len(lookup) == len(p2.s2points), nil
	}

	// as point is a non closed shape, containment isn't feasible
	// for other higher dimensions.
	return false, nil
}
