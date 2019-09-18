package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gistao/RedisGo-Async/redis"
)

//initialize logger
var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

type GeoCoord struct {
	Latitude, Longitude float64
}

type InnerLoop struct {
	Temp []GeoCoord
}

type OuterLoop struct {
	Temp []InnerLoop
}

type Polygon struct {
	Temp []OuterLoop
}

type geojson struct {
	Type        string  `json:"type"`
	Coordinates Polygon `json:"coordinates"`
}

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
	var geojson1 geojson
	j := getGeoJson(rec)
	//fmt.Println(j)
	err := json.Unmarshal([]byte(j), &geojson1)
	fmt.Println("error", err)
	logError(err)
	fmt.Println("unmarshalled")
	fmt.Println(geojson1)
	// fmt.Println(rec)
	// for {
	// 	// Read each record from csv
	// 	record, err := r.Read()
	// 	if err == io.EOF {
	// 		logger.Println("finished reading the file")
	// 		break
	// 	}
	// 	logError(err)
	// 	count++
	// 	fmt.Println(record)
	// }

	t := time.Now()
	elapsed := t.Sub(start)
	logger.Println("elapsed time: ", elapsed)
	logger.Println("total records written: %d", count)

	fmt.Println(&buf)
}
