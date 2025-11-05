// Package main demonstrates creational design patterns through a file conversion pipeline.
// It showcases Factory Method, Object Pool, and Builder patterns working together
// to convert CSV data through JSON and XML formats to YAML output.
package main

import (
	"fmt"
	"log"
	"os"

	"tmps-go-labs/lab2/domain/factory"
)

func main() {
	fmt.Println("Creational Design Patterns Demo: CSV → JSON → XML → YAML")

	converterFactory := factory.NewConverterFactory()
	pool := factory.NewConverterPool(5, converterFactory)

	pipeline, err := factory.NewPipelineBuilder().
		WithInputPath("input_sample.csv").
		WithOutputPath("output_final.yaml").
		WithIndent().
		WithPrettyPrint().
		AddCSVToJSON().
		AddJSONToXML().
		AddXMLToYAML().
		Build()
	if err != nil {
		log.Fatalf("Pipeline build failed: %v", err)
	}

	executor := factory.NewPipelineExecutor(pool)
	result := executor.Execute(pipeline)

	if !result.Success {
		log.Fatalf("Pipeline execution failed: %v", result.Error)
	}

	if _, err := os.Stat(pipeline.OutputPath); err == nil {
		fmt.Printf("Processed %d conversion steps in %d ms\n",
			len(result.Results), result.Duration/1_000_000)

		for i, stepResult := range result.Results {
			if stepResult.Error == nil {
				fmt.Printf("  Step %d: %s → %s (%.1f KB)\n",
					i+1,
					pipeline.Steps[i].From,
					pipeline.Steps[i].To,
					float64(len(stepResult.Data))/1024)
			}
		}
	} else {
		log.Fatalf("Output file not created: %v", err)
	}
}
