package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlainWriter(t *testing.T) {
	var buf bytes.Buffer
	writer := &PlainWriter{output: &buf}

	results := []SearchResult{
		{LineNumber: 1, Line: "hello"},
		{LineNumber: 3, Line: "world"},
	}

	err := writer.Write(results)
	assert.NoError(t, err)
	assert.Equal(t, "1: hello\n3: world\n", buf.String())
}

func TestJSONWriter(t *testing.T) {
	var buf bytes.Buffer
	writer := &JSONWriter{output: &buf}

	results := []SearchResult{
		{LineNumber: 1, Line: "hello"},
	}

	err := writer.Write(results)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), `"line_number":1`)
	assert.Contains(t, buf.String(), `"line":"hello"`)
}
