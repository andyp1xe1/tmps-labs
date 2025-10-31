package factory

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"tmps-go-labs/lab2/domain/models"
)

type JSONToCSVConverter struct{}

func init() {
	RegisterConverter("json-csv", func() models.Converter {
		return &JSONToCSVConverter{}
	})
}

func (j *JSONToCSVConverter) Convert(input io.Reader, from, to models.FileFormat) *models.ConversionResult {
	if from != models.FormatJSON || to != models.FormatCSV {
		return &models.ConversionResult{Error: fmt.Errorf("unsupported conversion")}
	}

	var data []map[string]interface{}
	decoder := json.NewDecoder(input)
	if err := decoder.Decode(&data); err != nil {
		return &models.ConversionResult{Error: err}
	}

	if len(data) == 0 {
		return &models.ConversionResult{
			Data:   []byte(""),
			Format: models.FormatCSV,
		}
	}

	var headers []string
	for key := range data[0] {
		headers = append(headers, key)
	}

	var csvLines []string
	csvLines = append(csvLines, strings.Join(headers, ","))

	for _, row := range data {
		var values []string
		for _, header := range headers {
			if val, exists := row[header]; exists {
				values = append(values, fmt.Sprintf("%v", val))
			} else {
				values = append(values, "")
			}
		}
		csvLines = append(csvLines, strings.Join(values, ","))
	}

	csvData := strings.Join(csvLines, "\n")

	return &models.ConversionResult{
		Data:   []byte(csvData),
		Format: models.FormatCSV,
	}
}

func (j *JSONToCSVConverter) SupportsFormat(format models.FileFormat) bool {
	return format == models.FormatJSON || format == models.FormatCSV
}
