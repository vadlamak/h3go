package main

import (
	"github.com/gistao/RedisGo-Async/redis"
	"os"
	"strconv"
)

////initialize logger
//var (
//	buf    bytes.Buffer
//	logger = log.New(&buf, "logger: ", log.Lshortfile)
//)
//
//func logError(e error) {
//	if e != nil {
//		logger.Fatal(e)
//	}
//}
func getConn() redis.AsynConn {
	//create conn to standalone redis
	c, err := redis.AsyncDial("tcp", ":6379")
	logError(err)
	return c
}
func closeConn(c redis.AsynConn) {
	defer c.Close()
}

//func getFileName() string {
//	return os.Args[1]
//}

func getLatIndex() (int, error) {
	return strconv.Atoi(os.Args[2])
}

func getLonIndex() (int, error) {
	return strconv.Atoi(os.Args[3])
}

func getGroundResolution() (int, error) {
	return strconv.Atoi(os.Args[4])
}

//func getFileReader(fileName string) *csv.Reader {
//	//read the file
//	f, err := os.Open("../samples/" + fileName)
//	logError(err)
//	return csv.NewReader(f)
//}

/*
func main() {

	start := time.Now()

	fileName := getFileName()
	latIndex, _ := getLatIndex()
	lonIndex, _ := getLonIndex()
	ground_resolution, _ := getGroundResolution()

	logger.Print("passed args: ")
	logger.Println(os.Args)

	r := getFileReader(fileName)

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
		for resolution := 0; resolution <= ground_resolution; resolution++ {
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

	fmt.Println(&buf)
}

*/
