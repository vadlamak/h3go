package ingest

import (
	"bytes"
	"fmt"
	"github.com/gistao/RedisGo-Async/redis"
	"github.com/tidwall/gjson"
	"github.com/uber/h3-go"
	"github.com/vadlamak/h3go/utils"
	"io"
	"log"
	"time"
)

//initialize logger
var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func getGeoJson(record []string) string {
	return record[6]
}

func IngestPolygons(fileName string) {

	// debug end
	start := time.Now()
	count := 0
	redisConn, _ := utils.GetConn()

	r, _ := utils.GetFileReader(fileName)

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

		handlePolygons(record, redisConn)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	logger.Println("elapsed time: ", elapsed)
	logger.Println("total records written: %d", count)
	utils.CloseConn(redisConn)

	fmt.Println(&buf)
}

// query KV store of choice
func getAssetsForH3Index(h3Indices []h3.H3Index, redisConn redis.AsynConn) {
	for _, h3Index := range h3Indices {
		fmt.Println("h3Index", h3.ToString(h3Index))

		booleanExists, e := redisConn.Do("SMEMBERS", h3.ToString(h3Index))
		exists, e := redis.Int(booleanExists, e)
		println(exists)
		if exists == 1 {
			fmt.Println("found a member set:", h3Index)
		}

	}
}

func handlePolygons(record []string, conn redis.AsynConn) {
	jsonStr := getGeoJson(record)
	geoJsonType := gjson.Get(jsonStr, "type").String()
	switch geoJsonType {
	case "Polygon":
		coordinates := gjson.Get(jsonStr, "coordinates")
		h3polygon := utils.CoordinatesToH3Polygon(coordinates)
		h3Indices := h3.Polyfill(h3polygon, 9) //polygons matching
		getAssetsForH3Index(h3Indices, conn)   //reads

	case "MultiPolygon":

	default:
		panic("unsupported type:" + geoJsonType)

	}
}
