package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestCsv(t *testing.T) {
	result, _ := os.Create("/Users/chenjia/Documents/code/kupu/kupu_go/gemini/local/call_gemini/result.csv")
	defer result.Close()
	var records = [][]string{{"a", "b", "c"}, {"aa", "bb", "cc"}, {"a1", "b1", "c1"}}
	for _, record := range records {
		writer := csv.NewWriter(result)
		err := writer.Write(record)
		defer writer.Flush()
		if err != nil {
			return
		}
	}
}

func TestRead(t *testing.T) {
	f, _ := os.Open("/Users/chenjia/Documents/code/kupu/kupu_go/gemini/local/call_gemini/test.csv")
	reader := csv.NewReader(f)
	records, _ := reader.ReadAll()
	for _, record := range records {
		fmt.Printf("%s,", strings.TrimSpace(record[0]))
	}
}
