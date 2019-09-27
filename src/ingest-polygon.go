package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/uber/h3-go"
	"log"
	"os"
	"time"
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
	jsonStr := getGeoJson(rec)
	result := gjson.Parse(jsonStr)
	polygon,_:= GjsonCoordsToPolygon(result)
	extCoord := polygon.Exterior.Vertices()
	println("extCoords")
	println(extCoord)

	t := time.Now()
	elapsed := t.Sub(start)
	logger.Println("elapsed time: ", elapsed)
	logger.Println("total records written: %d", count)

	fmt.Println(&buf)
}
