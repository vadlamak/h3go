package main

import (

	"github.com/paulmach/orb/geojson"
	"github.com/skelterjohn/geom"
	"github.com/tidwall/gjson"
	_ "log"
)

type Polygon struct {
	geojson.Polygon `json:",omitempty"`
	Exterior        geom.Polygon   `json:"exterior"`
	Interior        []geom.Polygon `json:"interior"`
}

func (p Polygon) ExteriorRing() geom.Polygon {
	return p.Exterior
}

func (p Polygon) InteriorRings() []geom.Polygon {
	return p.Interior
}

func (p Polygon) ContainsCoord(c geom.Coord) bool {

	ext := p.ExteriorRing()

	contains := false

	if ext.ContainsCoord(c) {

		contains = true

		for _, int := range p.InteriorRings() {

			if int.ContainsCoord(c) {
				contains = false
				break
			}
		}
	}

	return contains
}


func GjsonCoordsToPolygon(r gjson.Result) (Polygon, error) {

	rings := r.Array()

	count_rings := len(rings)
	count_interior := count_rings - 1

	exterior, err := gjsonLinearringtogeompolygon(rings[0])

	panic(err)

	interior := make([]geom.Polygon, count_interior)

	for i := 1; i <= count_interior; i++ {

		poly, err := gjsonLinearringtogeompolygon(rings[i])

		panic(err)

		interior = append(interior, poly)
	}

	polygon := Polygon{
		Exterior: exterior,
		Interior: interior,
	}

	return polygon, nil
}

func NewCoordinateFromLatLons(lat float64, lon float64) (geom.Coord, error) {

	coord := new(geom.Coord)

	coord.Y = lat
	coord.X = lon

	return *coord, nil
}

func NewPolygonFromCoords(coords []geom.Coord) (geom.Polygon, error) {

	path := geom.Path{}

	for _, c := range coords {
		path.AddVertex(c)
	}

	poly := new(geom.Polygon)
	poly.Path = path

	return *poly, nil
}

func gjsonLinearringtogeompolygon(r gjson.Result) (geom.Polygon, error) {

	coords := make([]geom.Coord, 0)

	for _, pt := range r.Array() {

		lonlat := pt.Array()

		lat := lonlat[1].Float()
		lon := lonlat[0].Float()

		coord, _ := NewCoordinateFromLatLons(lat, lon)
		coords = append(coords, coord)
	}

	return NewPolygonFromCoords(coords)
}