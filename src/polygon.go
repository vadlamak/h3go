package main

import (
	"github.com/paulmach/orb/geojson"
	"github.com/skelterjohn/geom"
	"github.com/tidwall/gjson"
	h3 "github.com/uber/h3-go"
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

func GjsonToH3Polygon(r gjson.Result)(h3.GeoPolygon){

	rings := r.Array()

	count_rings := len(rings)
	count_interior := count_rings - 1

	exterior, err := gjsonLinearRingToH3Geocoord(rings[0])

	if err != nil {
		println(err)
	}

	interior := make([][]h3.GeoCoord, count_interior)

	for i := 1; i <= count_interior; i++ {

		poly, err := gjsonLinearRingToH3Geocoord(rings[i])

		if err != nil {
			println(err)
		}

		interior = append(interior, poly)
	}

	polygon := h3.GeoPolygon{
		Geofence: exterior,
		Holes: interior,
	}

	return polygon

}

/*
type GeoCoord struct {
	Latitude, Longitude float64
}

func (g GeoCoord) toCPtr() *C.GeoCoord {
	return &C.GeoCoord{
		lat: C.double(deg2rad * g.Latitude),
		lon: C.double(deg2rad * g.Longitude),
	}
}

func (g GeoCoord) toC() C.GeoCoord {
	return *g.toCPtr()
}

// GeoPolygon is a geofence with 0 or more geofence holes
type GeoPolygon struct {
	// Geofence is the exterior boundary of the polygon
	Geofence []GeoCoord

	// Holes is a slice of interior boundary (holes) in the polygon
	Holes [][]GeoCoord
}
**/

func gjsonLinearRingToH3Geocoord(r gjson.Result) ([]h3.GeoCoord, error) {

	coords := make([]h3.GeoCoord, 0)

	for _, pt := range r.Array() {

		lonlat := pt.Array()

		lat := lonlat[1].Float()
		lon := lonlat[0].Float()
		//println("lat:", lat)
		//println("lon:", lon)

		coord, _ := H3CoordFromLatLons(lat, lon)
		coords = append(coords, coord)
	}

	return coords,nil
}

func GjsonCoordsToPolygon(r gjson.Result) (Polygon, error) {
	//println("result",r.String())
	//println("isArray", r.IsArray())
	//println("r.Type", r.Type)
	rings := r.Array()

	count_rings := len(rings)
	count_interior := count_rings - 1

	//println("rings:", rings)
	//println("ring[0]:",rings[0])
	exterior, err := gjsonLinearRingTogeomPolygon(rings[0])

	if err != nil {
		println(err)
	}

	interior := make([]geom.Polygon, count_interior)

	for i := 1; i <= count_interior; i++ {

		poly, err := gjsonLinearRingTogeomPolygon(rings[i])

		if err != nil {
			println(err)
		}

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

func H3CoordFromLatLons(lat float64, lon float64) (h3.GeoCoord, error) {

	coord := new(h3.GeoCoord)

	coord.Latitude = lat
	coord.Longitude = lon

	return *coord, nil
}

func NewPolygonFromCoords(coords []geom.Coord) (geom.Polygon, error) {

	path := geom.Path{}

	for _, c := range coords {
		path.AddVertex(c)
	}

	poly := new(geom.Polygon)
	//println("new polygon:", poly)
	poly.Path = path

	return *poly, nil
}

func gjsonLinearRingTogeomPolygon(r gjson.Result) (geom.Polygon, error) {

	coords := make([]geom.Coord, 0)

	for _, pt := range r.Array() {

		lonlat := pt.Array()

		lat := lonlat[1].Float()
		lon := lonlat[0].Float()
		//println("lat:", lat)
		//println("lon:", lon)

		coord, _ := NewCoordinateFromLatLons(lat, lon)
		coords = append(coords, coord)
	}

	return NewPolygonFromCoords(coords)
}
