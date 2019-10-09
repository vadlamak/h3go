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
	"strings"
	"time"
)

//initialize logger
var (
	buf      bytes.Buffer
	logger   = log.New(&buf, "logger: ", log.Lshortfile)
	match    = 0
	polygons = 0
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
	for i := 0; i < 250; i++ {
		//for {
		polygons = polygons + 1
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			logger.Println("finished reading the file")
			break
		}
		go handlePolygons(record, redisConn)

	}

	t := time.Now()
	elapsed := t.Sub(start)
	logger.Println("elapsed time: ", elapsed)
	logger.Println("total records written: %d", count)
	utils.CloseConn(redisConn)
	logger.Println("total matches", match)
	logger.Println("total polygons", polygons)

	fmt.Println(&buf)
}

// query KV store of choice
func getAssetsForH3Index(h3Indicies []h3.H3Index, redisConn redis.AsynConn) {

	for _, h3Index := range h3Indicies {
		smembers, e := redisConn.Do("SMEMBERS", h3.ToString(h3Index))
		members, e := redis.Strings(smembers, e)

		for _, member := range members {
			record := strings.Split(member, " ")
			println(len(record))
			println("member", member)
		}

		if len(members) > 0 {
			match = match + 1
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
		//st := time.Now()
		h3Indices := h3.Polyfill(h3polygon, 9) //polygons matching
		//elapsed := time.Now().Sub(st)
		//println("time for polyfill",elapsed )
		getAssetsForH3Index(h3Indices, conn)

	case "MultiPolygon":
		//TODO support multi polygons

	default:
		//TODO figure better way to handle this case
		//panic("unsupported type:" + geoJsonType)
		println("geoJsonType", geoJsonType)
		println("jsonStr", jsonStr)

	}

}
