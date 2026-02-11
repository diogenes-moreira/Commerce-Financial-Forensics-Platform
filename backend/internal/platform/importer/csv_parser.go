package importer

import (
	"encoding/csv"
	"fmt"
	"io"
)

type CSVParser struct{}

func NewCSVParser() *CSVParser { return &CSVParser{} }

// ParseCSV reads a CSV file and returns rows as []map[string]string.
// The first row is treated as headers.
func (p *CSVParser) ParseCSV(reader io.Reader) ([]map[string]string, error) {
	r := csv.NewReader(reader)
	r.TrimLeadingSpace = true

	headers, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	var rows []map[string]string
	lineNum := 1
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row %d: %w", lineNum, err)
		}
		lineNum++

		row := make(map[string]string, len(headers))
		for i, header := range headers {
			if i < len(record) {
				row[header] = record[i]
			}
		}
		rows = append(rows, row)
	}

	return rows, nil
}
