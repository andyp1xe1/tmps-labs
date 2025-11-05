// Package factory implements creational design patterns for file format converters.
// It provides Factory Method pattern for converter creation, Object Pool pattern
// for converter reuse, and Builder pattern for pipeline construction.
package factory

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"tmps-go-labs/lab2/domain/models"
)

type CSVToJSONConverter struct{}

func init() {
	RegisterConverter("csv-json", func() models.Converter {
		return &CSVToJSONConverter{}
	})
}

func (c *CSVToJSONConverter) Convert(input io.Reader, from, to models.FileFormat) *models.ConversionResult {
	if from != models.FormatCSV || to != models.FormatJSON {
		return &models.ConversionResult{Error: fmt.Errorf("unsupported conversion: %s to %s", from, to)}
	}

	reader := csv.NewReader(input)
	records, err := reader.ReadAll()
	if err != nil {
		return &models.ConversionResult{Error: fmt.Errorf("failed to read CSV: %w", err)}
	}

	if len(records) == 0 {
		return &models.ConversionResult{
			Data:   []byte("[]"),
			Format: models.FormatJSON,
		}
	}

	headers := records[0]
	var jsonData []map[string]string

	for _, record := range records[1:] {
		row := make(map[string]string)
		for i, value := range record {
			if i < len(headers) {
				row[headers[i]] = value
			}
		}
		jsonData = append(jsonData, row)
	}

	data, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return &models.ConversionResult{Error: fmt.Errorf("failed to marshal JSON: %w", err)}
	}

	return &models.ConversionResult{
		Data:   data,
		Format: models.FormatJSON,
	}
}

func (c *CSVToJSONConverter) SupportsFormat(format models.FileFormat) bool {
	return format == models.FormatCSV || format == models.FormatJSON
}
