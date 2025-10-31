package factory

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"sync"

	"tmps-go-labs/lab2/domain/models"
)

type ConverterCreator func() models.Converter

var (
	converterRegistry = make(map[string]ConverterCreator)
	registryMutex     sync.RWMutex
)

func RegisterConverter(formatType string, creator ConverterCreator) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	converterRegistry[formatType] = creator
}

type ConverterFactory interface {
	CreateConverter(formatType string) (models.Converter, error)
}

type DefaultConverterFactory struct{}

func NewConverterFactory() ConverterFactory {
	return &DefaultConverterFactory{}
}

func (f *DefaultConverterFactory) CreateConverter(formatType string) (models.Converter, error) {
	registryMutex.RLock()
	creator, exists := converterRegistry[formatType]
	registryMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unsupported converter type: %s", formatType)
	}

	return creator(), nil
}

type CSVToJSONConverter struct{}

func init() {
	RegisterConverter("csv-json", func() models.Converter {
		return &CSVToJSONConverter{}
	})
}

func (c *CSVToJSONConverter) Convert(input io.Reader, from, to models.FileFormat) *models.ConversionResult {
	if from != models.FormatCSV || to != models.FormatJSON {
		return &models.ConversionResult{Error: fmt.Errorf("unsupported conversion")}
	}

	reader := csv.NewReader(input)
	records, err := reader.ReadAll()
	if err != nil {
		return &models.ConversionResult{Error: err}
	}

	var jsonData []map[string]string
	if len(records) > 0 {
		headers := records[0]
		for _, record := range records[1:] {
			row := make(map[string]string)
			for i, value := range record {
				if i < len(headers) {
					row[headers[i]] = value
				}
			}
			jsonData = append(jsonData, row)
		}
	}

	data, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return &models.ConversionResult{Error: err}
	}

	return &models.ConversionResult{
		Data:   data,
		Format: models.FormatJSON,
	}
}

func (c *CSVToJSONConverter) SupportsFormat(format models.FileFormat) bool {
	return format == models.FormatCSV || format == models.FormatJSON
}

type JSONToXMLConverter struct{}

func init() {
	RegisterConverter("json-xml", func() models.Converter {
		return &JSONToXMLConverter{}
	})
}

func (j *JSONToXMLConverter) Convert(input io.Reader, from, to models.FileFormat) *models.ConversionResult {
	if from != models.FormatJSON || to != models.FormatXML {
		return &models.ConversionResult{Error: fmt.Errorf("unsupported conversion")}
	}

	var data any
	decoder := json.NewDecoder(input)
	if err := decoder.Decode(&data); err != nil {
		return &models.ConversionResult{Error: err}
	}

	wrapped := map[string]any{"root": data}
	xmlData, err := xml.MarshalIndent(wrapped, "", "  ")
	if err != nil {
		return &models.ConversionResult{Error: err}
	}

	return &models.ConversionResult{
		Data:   xmlData,
		Format: models.FormatXML,
	}
}

func (j *JSONToXMLConverter) SupportsFormat(format models.FileFormat) bool {
	return format == models.FormatJSON || format == models.FormatXML
}

type XMLToYAMLConverter struct{}

func init() {
	RegisterConverter("xml-yaml", func() models.Converter {
		return &XMLToYAMLConverter{}
	})
}

func (x *XMLToYAMLConverter) Convert(input io.Reader, from, to models.FileFormat) *models.ConversionResult {
	if from != models.FormatXML || to != models.FormatYAML {
		return &models.ConversionResult{Error: fmt.Errorf("unsupported conversion")}
	}

	content, err := io.ReadAll(input)
	if err != nil {
		return &models.ConversionResult{Error: err}
	}

	yamlData := "# Converted from XML\n" + strings.ReplaceAll(string(content), "<", "")
	yamlData = strings.ReplaceAll(yamlData, ">", ":")

	return &models.ConversionResult{
		Data:   []byte(yamlData),
		Format: models.FormatYAML,
	}
}

func (x *XMLToYAMLConverter) SupportsFormat(format models.FileFormat) bool {
	return format == models.FormatXML || format == models.FormatYAML
}
