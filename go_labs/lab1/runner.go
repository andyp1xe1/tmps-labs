package main

import (
	"bufio"
	"io"
)

type Runner struct {
	engine SearchEngine
	reader io.Reader
	writer ResultWriter
}

func NewRunner(engine SearchEngine, reader io.Reader, writer ResultWriter) *Runner {
	return &Runner{
		engine: engine,
		reader: reader,
		writer: writer,
	}
}

func (r *Runner) Run(query string) error {
	scanner := bufio.NewScanner(r.reader)
	var results []SearchResult
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		if r.engine.Search(line, query) {
			results = append(results, SearchResult{
				LineNumber: lineNumber,
				Line:       line,
			})
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return r.writer.Write(results)
}
