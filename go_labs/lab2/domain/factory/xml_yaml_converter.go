// Package factory implements creational design patterns for file format converters.
// It provides Factory Method pattern for converter creation, Object Pool pattern
// for converter reuse, and Builder pattern for pipeline construction.
package factory

import (
	"fmt"
	"io"

	"github.com/clbanning/mxj/v2"
	"gopkg.in/yaml.v3"
	"tmps-go-labs/lab2/domain/models"
)

type XMLToYAMLConverter struct{}

func init() {
	RegisterConverter("xml-yaml", func() models.Converter {
		return &XMLToYAMLConverter{}
	})
}

func (x *XMLToYAMLConverter) Convert(input io.Reader, from, to models.FileFormat) *models.ConversionResult {
	if from != models.FormatXML || to != models.FormatYAML {
		return &models.ConversionResult{Error: fmt.Errorf("unsupported conversion: %s to %s", from, to)}
	}

	// Read XML data
	xmlData, err := io.ReadAll(input)
	if err != nil {
		return &models.ConversionResult{Error: fmt.Errorf("failed to read XML: %w", err)}
	}

	// Parse XML using mxj library
	mv, err := mxj.NewMapXml(xmlData)
	if err != nil {
		return &models.ConversionResult{Error: fmt.Errorf("failed to parse XML: %w", err)}
	}

	// Convert map to YAML using gopkg.in/yaml.v3
	yamlData, err := yaml.Marshal(mv.Old())
	if err != nil {
		return &models.ConversionResult{Error: fmt.Errorf("failed to convert to YAML: %w", err)}
	}

	return &models.ConversionResult{
		Data:   yamlData,
		Format: models.FormatYAML,
	}
}

func (x *XMLToYAMLConverter) SupportsFormat(format models.FileFormat) bool {
	return format == models.FormatXML || format == models.FormatYAML
}
