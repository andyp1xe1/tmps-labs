package models

type Pipeline struct {
	Steps      []ConversionStep
	Options    ConversionOptions
	InputPath  string
	OutputPath string
}

type ConversionStep struct {
	From      FileFormat
	To        FileFormat
	Converter Converter
}

type PipelineResult struct {
	Success  bool
	Results  []*ConversionResult
	Error    error
	Duration int64
}
