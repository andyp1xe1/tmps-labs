package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var engine = flag.String("e", "literal", "search engine: literal, regex, fuzzy")
	var query = flag.String("q", "", "search query")
	var format = flag.String("f", "plain", "output format: plain, json")
	var path = flag.String("p", "", "file path to search in")

	flag.Parse()

	if *query == "" || *path == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -e <engine> -q <query> -f <format> -p <path>\n", os.Args[0])
		os.Exit(1)
	}

	file, err := os.Open(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	searchEngine := createSearchEngine(*engine)
	writer := createWriter(*format, os.Stdout)

	runner := NewRunner(searchEngine, file, writer)

	if err := runner.Run(*query); err != nil {
		fmt.Fprintf(os.Stderr, "Error running search: %v\n", err)
		os.Exit(1)
	}
}

func createSearchEngine(engineType string) SearchEngine {
	switch engineType {
	case "literal":
		return &LiteralSearch{}
	case "regex":
		return &RegexSearch{}
	case "fuzzy":
		return &FuzzySearch{}
	default:
		fmt.Fprintf(os.Stderr, "Unknown engine type: %s\n", engineType)
		os.Exit(1)
		return nil
	}
}

func createWriter(format string, output io.Writer) ResultWriter {
	switch format {
	case "plain":
		return &PlainWriter{output: output}
	case "json":
		return &JSONWriter{output: output}
	default:
		fmt.Fprintf(os.Stderr, "Unknown format: %s\n", format)
		os.Exit(1)
		return nil
	}
}
