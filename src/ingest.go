package main

import (
	"fmt"
	"github.com/vadlamak/h3go/ingest"
	"os"
	"strconv"
)

func getPointUsageText() string {
	return "for `point` sample it is possible to just specify the file name or specify the filename along with the lat,lon and ground resolution as args.\n" +
		"For ex: ./main point lightning_2016.csv\n" +
		" or \n" +
		"./main point lightning_2016.csv 3 2 9" +
		" or \n" +
		"./main point -- will by default ingest lightning_2016.csv from samples/point folder"
}

func getPolygonUsageText() string {
	return "similarly we can ingest polygon data by specifying the file name\n" +
		"./main polygon zipcodes.csv\n" +
		"or \n" +
		"./main polygon -- will be default ingest zipcodes.csv from samples/polygon folder"
}

func getUsageHelp() string {
	resp := "The first argument is to specify the type of sample to be ingested. " +
		"accepted values are `point` and `polygon`.\n" + getPointUsageText() + getPolygonUsageText()
	return resp
}

func main() {
	argsLen := len(os.Args)

	if argsLen == 1 {
		fmt.Println(getUsageHelp())
		return
	}

	switch os.Args[1] {
	case "point":
		switch argsLen {
		case 2:
			ingest.IngestSamples("../samples/point/lightning_2016.csv", 3, 2, 9)
		case 3:
			fname := os.Args[2]
			ingest.IngestSamples("../samples/point/"+fname, 3, 2, 9)
		case 5:
			fname := os.Args[2]
			lat, _ := strconv.Atoi(os.Args[3])
			lon, _ := strconv.Atoi(os.Args[4])
			groundResolution, _ := strconv.Atoi(os.Args[5])

			ingest.IngestSamples("../samples/point/"+fname, lat, lon, groundResolution)
		default:
			fmt.Println(getPointUsageText())
		}
	case "polygon":
		switch argsLen {
		case 3:
			fname := os.Args[2]
			ingest.IngestPolygons("../samples/polygon" + fname)
		case 2:
			ingest.IngestPolygons("../samples/polygon/zipcodes.csv")
		default:
			fmt.Println(getPolygonUsageText())
		}
	case "help":
	default:
		fmt.Println(getUsageHelp())
	}

}
