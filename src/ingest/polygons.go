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
	"sync"
	"time"
)

//global vars
var (
	buf      bytes.Buffer
	logger   = log.New(&buf, "logger: ", log.Lshortfile)
	match    = 0
	polygons = 0
	geojsonMap = sync.Map{}
	h3IndexMap = sync.Map{}
	assetsMap = sync.Map{}

)

func getGeoJson(record []string) string {
	return record[6]
}

func getZipcode(record[] string) string {
	return record[0];
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

		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			logger.Println("finished reading the file")
			break
		}
		go HandlePolygons(record, redisConn)

	}

	writeToCSV("../out/geojsons.csv",geojsonMap)
	writeToCSV("../out/h3index.csv",h3IndexMap)
	writeToCSV("../out/assets.csv",assetsMap)

	t := time.Now()
	elapsed := t.Sub(start)
	logger.Println("elapsed time: ", elapsed)
	logger.Println("total records written: %d", count)
	utils.CloseConn(redisConn)
	logger.Println("total matches", match)
	logger.Println("total polygons", polygons)

	fmt.Println(&buf)
}


func mapToSlice(data sync.Map) [][]string {
	resp := [][]string{}

	data.Range(func(k, v interface{}) bool{
		record :=[]string{}
		record = append(record,fmt.Sprintf("%v", k))
		record = append(record,fmt.Sprintf("%v", v))
		resp = append(resp,record)
		return  true
	})
	//
	//for _,stra := range resp {
	//	for _, rec := range stra {
	//		println("rec",rec)
	//	}
	//}

	return resp

}
func writeToCSV(fileName string,data sync.Map){
	println("in writeToCSV:",fileName)
	convertedData := mapToSlice(data)
	println("converted data",convertedData)
	geoJsonWriter,e := utils.GetFileWriter(fileName)
	if e!=nil {
		log.Fatal("err getting reader",e)
	}
	err := geoJsonWriter.WriteAll(convertedData)
	if err !=nil {
		log.Fatal("err",err)
	}


}

// query KV store of choice
func getAssetsForH3Index(h3Indicies []h3.H3Index, redisConn redis.AsynConn,zipcode string)  {
	for _, h3Index := range h3Indicies {
		smembers, e := redisConn.Do("SMEMBERS", h3.ToString(h3Index))
		members, e := redis.Strings(smembers, e)
		for _, member := range members {
			record := strings.Split(member, " ")
			assetsMap.Store(zipcode,record)
			//println(len(record))
			println("record", record)
		}

		if len(members) > 0 {
			match = match + 1
		}
	}

}

func H3IndiciesToCSV(h3indicies []h3.H3Index) string{
	resp := strings.Builder{}
	for _,h3Index := range h3indicies{
		resp.WriteString(h3.ToString(h3Index))
		resp.WriteString(",")
	}
	return resp.String()
}
func HandlePolygons(record []string, conn redis.AsynConn) {
	jsonStr := getGeoJson(record)
	zipcode := getZipcode(record)
	geojsonMap.Store(zipcode,jsonStr) //add to map zipcode -> geojson

	geoJsonType := gjson.Get(jsonStr, "type").String()
	switch geoJsonType {
	case "Polygon":
		polygons = polygons + 1
		coordinates := gjson.Get(jsonStr, "coordinates")
		h3polygon := utils.CoordinatesToH3Polygon(coordinates)
		//st := time.Now()
		h3Indices := h3.Polyfill(h3polygon, 9) //polygons matching
		h3IndexMap.Store(zipcode,H3IndiciesToCSV(h3Indices))
		//elapsed := time.Now().Sub(st)
		//println("time for polyfill",elapsed )
		getAssetsForH3Index(h3Indices, conn,zipcode)

	case "MultiPolygon":
		//TODO support multi polygons

	default:
		//TODO figure better way to handle this case
		//panic("unsupported type:" + geoJsonType)
		println("geoJsonType", geoJsonType)
		println("jsonStr", jsonStr)

	}

}
