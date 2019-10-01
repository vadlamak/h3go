package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/gistao/RedisGo-Async/redis"
	"github.com/tidwall/gjson"
	"github.com/uber/h3-go"
	"io"
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

func main() {
	c:=getConn()
	c.Do("SADD","894cc5b2523ffff","fo1","ba1")
	re,_ := redis.Int(c.Do("EXISTS","894cc5b2523ffff"))
	asses,_ := redis.Strings(c.Do("SMEMBERS","894cc5b2523ffff"))
	for _,ass := range asses {
		println("ass:",ass)
	}
	println("reply",re)
	start := time.Now()
	fileName := getFileName()
	count := 0

	logger.Print("passed args: ")
	logger.Println(os.Args)

	r := getFileReader(fileName)

	//read line # 0 - header
	header, _ := r.Read()
	fmt.Println(header)
	for i := 0; i < 2; i++ {
	//for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			logger.Println("finished reading the file")
			break
		}

		handlePolygons(record)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	logger.Println("elapsed time: ", elapsed)
	logger.Println("total records written: %d", count)

	fmt.Println(&buf)
}

// query KV store of choice
func getAssetsForH3Index(h3Indices []h3.H3Index) {
	c := getConn()
	for _, h3Index := range h3Indices {
		fmt.Println("h3Index",h3.ToString(h3Index))

		booleanExists, e := c.Do("EXISTS", h3.ToString(h3Index))
		exists,e:=redis.Int(booleanExists,e)
		println(exists)
		if exists ==1 {
			fmt.Println("found a member set:",h3Index)
		}

	}
}

func handlePolygons(record []string) {
	jsonStr := getGeoJson(record)
	geoJsonType := gjson.Get(jsonStr, "type").String()
	switch geoJsonType {
	case "Polygon":
		coordinates := gjson.Get(jsonStr, "coordinates")

		h3polygon := CoordinatesToH3Polygon(coordinates)
		h3Indices := h3.Polyfill(h3polygon, 9)
		//if len(h3Indices) == 0 {
		//	println("no match")
		//} else {println("matched")}
		//println("len:", len(h3Indices))
		getAssetsForH3Index(h3Indices)
	case "MultiPolygon":
	default:
		panic("unsupported type:" + geoJsonType)

	}
}
