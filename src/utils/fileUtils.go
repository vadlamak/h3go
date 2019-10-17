package utils

import (
	"encoding/csv"
	"os"
)

func GetFileReader(fileName string) (*csv.Reader,error) {
	//read the file
	f, err := os.Open(fileName)
	//defer f.Close()
	return csv.NewReader(f),err
}

func GetFileWriter(fileName string) (*csv.Writer,error) {

	f,err := os.Open(fileName)
	if err !=nil {
		nf,nerr := os.Create(fileName)
		if nerr !=nil {
			println(nerr)
		}else {
			//defer nf.Close()
			return csv.NewWriter(nf),nil
		}
	}
	//defer f.Close()
	return csv.NewWriter(f),err

}
