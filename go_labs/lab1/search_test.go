package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiteralSearch(t *testing.T) {
	engine := &LiteralSearch{}

	assert.True(t, engine.Search("hello world", "world"))
	assert.False(t, engine.Search("hello world", "xyz"))
	assert.True(t, engine.Search("test", ""))
}

func TestRegexSearch(t *testing.T) {
	engine := &RegexSearch{}

	assert.True(t, engine.Search("hello123", "\\d+"))
	assert.False(t, engine.Search("hello", "\\d+"))
	assert.False(t, engine.Search("hello", "["))
}

func TestFuzzySearch(t *testing.T) {
	engine := &FuzzySearch{}

	assert.True(t, engine.Search("hello world", "hlowrd"))
	assert.False(t, engine.Search("hello", "xyz"))
	assert.True(t, engine.Search("test", ""))
}
