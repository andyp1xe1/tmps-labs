// Package models defines the core interfaces and data structures for file format
// conversion operations. It provides the foundation types used by the creational
// design patterns implemented in the factory package.
package models

type Pipeline struct {
	Steps      []ConversionStep
	Options    ConversionOptions
	InputPath  string
	OutputPath string
}

type ConversionStep struct {
	From FileFormat
	To   FileFormat
}

type PipelineResult struct {
	Success  bool
	Results  []*ConversionResult
	Error    error
	Duration int64
}
