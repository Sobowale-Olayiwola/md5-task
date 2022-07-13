package helpers

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// LoadCSV is used to load a csv file into an array of array of strings, which index in the
// array can be referred to as csvline
func LoadCSV(path string) (csvLine [][]string, err error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	firstRow, err := bufio.NewReader(csvFile).ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	_, err = csvFile.Seek(int64(len(firstRow)), io.SeekStart)
	if err != nil {
		return nil, err
	}
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}
	return csvLines, err
}
