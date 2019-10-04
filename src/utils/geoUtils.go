package utils

import (
	"github.com/tidwall/gjson"
	"github.com/uber/h3-go"
	_ "log"
)

func CoordinatesToH3Polygon(r gjson.Result) h3.GeoPolygon {
	rings := r.Array()
	count_rings := len(rings)
	count_interior := count_rings - 1
	if count_interior > 0 {
		println("count interior:", count_interior)
	}

	exterior, err := linearRingToH3Geocoord(rings[0])

	if err != nil {
		println(err)
	}

	interior := make([][]h3.GeoCoord, count_interior)

	for i := 1; i <= count_interior; i++ {

		poly, err := linearRingToH3Geocoord(rings[i])

		if err != nil {
			println(err)
		}

		interior = append(interior, poly)
	}

	polygon := h3.GeoPolygon{
		Geofence: exterior,
		Holes:    interior,
	}

	return polygon

}

func linearRingToH3Geocoord(r gjson.Result) ([]h3.GeoCoord, error) {

	coords := make([]h3.GeoCoord, 0)

	for _, pt := range r.Array() {

		lonlat := pt.Array()

		lat := lonlat[1].Float()
		lon := lonlat[0].Float()
		//fmt.Print("\nlat:lon: ",lat,lon)

		coord, _ := H3CoordFromLatLons(lat, lon)
		coords = append(coords, coord)
	}

	return coords, nil
}

func H3CoordFromLatLons(lat float64, lon float64) (h3.GeoCoord, error) {

	coord := new(h3.GeoCoord)

	coord.Latitude = lat
	coord.Longitude = lon

	return *coord, nil
}
