package ingest

import (
	"bytes"
	"fmt"
	"github.com/uber/h3-go"
	"github.com/vadlamak/h3go/utils"
	"io"
	"log"
	"strconv"
	"time"
)

//initialize logger
var (
	buffer bytes.Buffer
	logg   = log.New(&buf, "logger: ", log.Lshortfile)
)

func logError(e error) {
	if e != nil {
		logger.Fatal(e)
	}
}

func IngestSamples(fileName string, latIndex int, lonIndex int, groundResolution int) {

	start := time.Now()

	r, _ := utils.GetFileReader(fileName)
	c, _ := utils.GetConn()

	//skip header
	r.Read()
	count := 0
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			logger.Println("finished reading the file")
			break
		}
		logError(err)
		lat, _ := strconv.ParseFloat(record[latIndex], 64) //index 3
		lon, _ := strconv.ParseFloat(record[lonIndex], 64) //index 2

		geo := h3.GeoCoord{
			Latitude:  lat,
			Longitude: lon,
		}

		//store asset location in hierarchy
		for resolution := 0; resolution <= groundResolution; resolution++ {
			h3Index := h3.ToString(h3.FromGeo(geo, resolution))
			//add in redis sets as there can be many assets with the same h3 index
			c.AsyncDo("SADD", h3Index, record)
			count++
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	logger.Println("elapsed time: ", elapsed)
	logger.Println("total records written: %d", count)
	c.Close()
	fmt.Println(&buf)
}
