package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type SearchResult struct {
	LineNumber int    `json:"line_number"`
	Line       string `json:"line"`
}

type ResultWriter interface {
	Write(results []SearchResult) error
}

type PlainWriter struct {
	output io.Writer
}

func (p *PlainWriter) Write(results []SearchResult) error {
	for _, result := range results {
		_, err := fmt.Fprintf(p.output, "%d: %s\n", result.LineNumber, result.Line)
		if err != nil {
			return err
		}
	}
	return nil
}

type JSONWriter struct {
	output io.Writer
}

func (j *JSONWriter) Write(results []SearchResult) error {
	encoder := json.NewEncoder(j.output)
	return encoder.Encode(results)
}
