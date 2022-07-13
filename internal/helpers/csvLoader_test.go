package helpers

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"testing"
)

type Employee struct {
	ID  string
	Age int
}

func TestLoadCSC(t *testing.T) {
	records := []Employee{
		{"E01", 25},
		{"E02", 26},
		{"E03", 24},
		{"E04", 26},
	}
	file, err := os.Create("records.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer file.Close()
	defer func() { os.Remove("records.csv") }()
	w := csv.NewWriter(file)
	defer w.Flush()

	// Using WriteAll
	var data [][]string
	for _, record := range records {
		row := []string{record.ID, strconv.Itoa(record.Age)}
		data = append(data, row)
	}
	w.WriteAll(data)
	csvLine, err := LoadCSV("records.csv")
	if err != nil {
		t.Error("should not return an error")
	}
	if len(csvLine) != 3 {
		t.Errorf("expected %d got %d", 3, len(csvLine))
	}

}
