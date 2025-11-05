// Package factory implements creational design patterns for file format converters.
// It provides Factory Method pattern for converter creation, Object Pool pattern
// for converter reuse, and Builder pattern for pipeline construction.
package factory

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/clbanning/mxj/v2"
	"tmps-go-labs/lab2/domain/models"
)

type JSONToXMLConverter struct{}

func init() {
	RegisterConverter("json-xml", func() models.Converter {
		return &JSONToXMLConverter{}
	})
}

func (j *JSONToXMLConverter) Convert(input io.Reader, from, to models.FileFormat) *models.ConversionResult {
	if from != models.FormatJSON || to != models.FormatXML {
		return &models.ConversionResult{Error: fmt.Errorf("unsupported conversion: %s to %s", from, to)}
	}

	// Read JSON data
	jsonData, err := io.ReadAll(input)
	if err != nil {
		return &models.ConversionResult{Error: fmt.Errorf("failed to read JSON: %w", err)}
	}

	// Parse JSON into generic interface
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return &models.ConversionResult{Error: fmt.Errorf("failed to parse JSON: %w", err)}
	}

	// Convert to XML using mxj library
	mv := mxj.Map{"root": data}
	xmlData, err := mv.XmlIndent("", "  ")
	if err != nil {
		return &models.ConversionResult{Error: fmt.Errorf("failed to convert to XML: %w", err)}
	}

	return &models.ConversionResult{
		Data:   xmlData,
		Format: models.FormatXML,
	}
}

func (j *JSONToXMLConverter) SupportsFormat(format models.FileFormat) bool {
	return format == models.FormatJSON || format == models.FormatXML
}
