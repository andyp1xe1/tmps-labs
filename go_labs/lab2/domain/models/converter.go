package models

import "io"

type FileFormat string

const (
	FormatCSV  FileFormat = "csv"
	FormatJSON FileFormat = "json"
	FormatXML  FileFormat = "xml"
	FormatYAML FileFormat = "yaml"
)

type ConversionResult struct {
	Data   []byte
	Format FileFormat
	Error  error
}

type Converter interface {
	Convert(input io.Reader, from, to FileFormat) *ConversionResult
	SupportsFormat(format FileFormat) bool
}

type ConversionOptions struct {
	Indent      bool
	PrettyPrint bool
	Headers     []string
}
