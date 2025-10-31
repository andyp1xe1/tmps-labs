package main

import (
	"fmt"
	"log"
	"strings"

	"tmps-go-labs/lab2/domain/factory"
	_ "tmps-go-labs/lab2/domain/factory" // Import to trigger init functions
	"tmps-go-labs/lab2/domain/models"
)

func main() {
	fmt.Println("File Converter - Creational Design Patterns Demo")
	fmt.Println("================================================")

	demonstrateFactoryMethod()
	fmt.Println()
	demonstrateBuilder()
	fmt.Println()
	demonstrateObjectPool()
}

func demonstrateFactoryMethod() {
	fmt.Println("1. Factory Method Pattern Demo")
	fmt.Println("------------------------------")

	converterFactory := factory.NewConverterFactory()

	converterTypes := []string{"csv-json", "json-xml", "xml-yaml", "json-csv"}

	for _, converterType := range converterTypes {
		converter, err := converterFactory.CreateConverter(converterType)
		if err != nil {
			log.Printf("Error creating %s converter: %v", converterType, err)
			continue
		}

		fmt.Printf("✓ Created %s converter\n", converterType)

		testData := strings.NewReader("name,age\nJohn,25\nJane,30")
		result := converter.Convert(testData, models.FormatCSV, models.FormatJSON)

		if result.Error != nil {
			log.Printf("Conversion error: %v", result.Error)
		} else {
			fmt.Printf("  Sample output: %s\n", string(result.Data[:50])+"...")
		}
	}
}

func demonstrateBuilder() {
	fmt.Println("2. Builder Pattern Demo")
	fmt.Println("-----------------------")

	pipeline, err := factory.NewPipelineBuilder().
		WithInputPath("/tmp/input.csv").
		WithOutputPath("/tmp/output.yaml").
		WithIndent().
		WithPrettyPrint().
		WithHeaders([]string{"name", "age", "city"}).
		AddCSVToJSON().
		AddJSONToXML().
		AddXMLToYAML().
		Build()

	if err != nil {
		log.Printf("Error building pipeline: %v", err)
		return
	}

	fmt.Printf("✓ Built conversion pipeline with %d steps\n", len(pipeline.Steps))
	fmt.Printf("  Input: %s\n", pipeline.InputPath)
	fmt.Printf("  Output: %s\n", pipeline.OutputPath)
	fmt.Printf("  Options: Indent=%v, PrettyPrint=%v\n",
		pipeline.Options.Indent, pipeline.Options.PrettyPrint)

	for i, step := range pipeline.Steps {
		fmt.Printf("  Step %d: %s → %s\n", i+1, step.From, step.To)
	}
}

func demonstrateObjectPool() {
	fmt.Println("3. Object Pool Pattern Demo")
	fmt.Println("---------------------------")

	converterFactory := factory.NewConverterFactory()
	pool := factory.NewConverterPool(3, converterFactory)

	fmt.Printf("✓ Created converter pool with max size: 3\n")

	converterTypes := []string{"csv-json", "json-xml", "xml-yaml", "csv-json"}

	var converters []models.Converter

	for i, converterType := range converterTypes {
		converter, err := pool.Get(converterType)
		if err != nil {
			log.Printf("Error getting converter from pool: %v", err)
			continue
		}

		converters = append(converters, converter)
		fmt.Printf("  Got converter %d (%s), pool size: %d, created: %d\n",
			i+1, converterType, pool.Size(), pool.Created())
	}

	for i, converter := range converters {
		pool.Put(converter)
		fmt.Printf("  Returned converter %d, pool size: %d\n", i+1, pool.Size())
	}

	fmt.Printf("✓ Pool demonstration complete, final pool size: %d\n", pool.Size())
}
