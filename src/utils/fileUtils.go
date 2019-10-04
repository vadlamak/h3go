package utils

import (
	"encoding/csv"
	"os"
)

func GetFileReader(fileName string) (*csv.Reader,error) {
	//read the file
	f, err := os.Open(fileName)
	return csv.NewReader(f),err
}
