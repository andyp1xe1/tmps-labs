package main

import (
	"regexp"
	"strings"
)

type SearchEngine interface {
	Search(text, query string) bool
}

type LiteralSearch struct{}

func (l *LiteralSearch) Search(text, query string) bool {
	return strings.Contains(text, query)
}

type RegexSearch struct{}

func (r *RegexSearch) Search(text, query string) bool {
	matched, err := regexp.MatchString(query, text)
	if err != nil {
		return false
	}
	return matched
}

type FuzzySearch struct{}

func (f *FuzzySearch) Search(text, query string) bool {
	textLower := strings.ToLower(text)
	queryLower := strings.ToLower(query)

	if len(queryLower) == 0 {
		return true
	}

	textIdx := 0
	queryIdx := 0

	for textIdx < len(textLower) && queryIdx < len(queryLower) {
		if textLower[textIdx] == queryLower[queryIdx] {
			queryIdx++
		}
		textIdx++
	}

	return queryIdx == len(queryLower)
}
