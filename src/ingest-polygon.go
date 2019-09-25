package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gistao/RedisGo-Async/redis"
	"github.com/paulmach/orb/geojson"
	h3 "github.com/uber/h3-go"
)

//initialize logger
var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func logError(e error) {
	if e != nil {
		logger.Fatal(e)
	}
}
func getConn() redis.AsynConn {
	//create conn to standalone redis
	c, err := redis.AsyncDial("tcp", ":6379")
	logError(err)
	return c
}
func closeConn(c redis.AsynConn) {
	defer c.Close()
}
func getFileName() string {
	return os.Args[1]
}

func getFileReader(fileName string) *csv.Reader {
	//read the file
	f, err := os.Open("../samples/polygon/" + fileName)
	logError(err)
	return csv.NewReader(f)
}

func getGeoJson(record []string) string {
	return record[6]
}

func getGeocoord(arr []float64) h3.GeoCoord {
	return h3.GeoCoord{
		Latitude:  arr[1],
		Longitude: arr[0],
	}
}

func processPolygon(polygon [][][]float64) {
	println(len(polygon[0]))
	if len(polygon) > 0 && len(polygon) == 1 {
		for i := 0; i < len(polygon[0]); i++ {
			if len(polygon[0][i]) == 2 {
				// print(polygon[0][i][0])
				// print(":", polygon[0][i][1])
				// println()
				geo := getGeocoord(polygon[0][i])
				println("lat: %f", geo.Latitude)
				println("lon: %f", geo.Longitude)

			} else {
				println("should never happen")
			}

		}
	} else {
		//throw error
		println("handle me")
	}

	//h3.GeoPolygon

}
func processMultiPolygon(mPolygon [][][][]float64) {
	//println(mPolygon[0][0][0][1])
	for i := 0; i < len(mPolygon); i++ {
		processPolygon(mPolygon[i])
	}
}

func main() {
	start := time.Now()
	fileName := getFileName()
	count := 0

	logger.Print("passed args: ")
	logger.Println(os.Args)

	r := getFileReader(fileName)

	//skip header
	header, _ := r.Read()
	fmt.Println(header)
	rec, _ := r.Read()
	gjson := getGeoJson(rec)
	rawData := []byte(gjson)
	geoJsonGeom := &geojson.Geometry{}
	err := geoJsonGeom.UnmarshalJSON(rawData)
	logError(err)
	orbGeom := geoJsonGeom.Geometry() //convert to diff geometry type
	println("GeoJSONType: ", orbGeom.GeoJSONType())
	println("Dimensions: ", orbGeom.Dimensions())
	println("closed?: ", orbGeom.Bound().ToRing().Closed())

	t := time.Now()
	elapsed := t.Sub(start)
	logger.Println("elapsed time: ", elapsed)
	logger.Println("total records written: %d", count)

	fmt.Println(&buf)
}
