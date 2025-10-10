package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	input := "hello world\ntest line\nworld again"
	reader := strings.NewReader(input)

	var output bytes.Buffer
	engine := &LiteralSearch{}
	writer := &PlainWriter{output: &output}

	runner := NewRunner(engine, reader, writer)
	err := runner.Run("world")

	assert.NoError(t, err)
	assert.Contains(t, output.String(), "1: hello world")
	assert.Contains(t, output.String(), "3: world again")
}
