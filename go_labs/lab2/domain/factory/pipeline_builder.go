// Package factory implements creational design patterns for file format converters.
// It provides Factory Method pattern for converter creation, Object Pool pattern
// for converter reuse, and Builder pattern for pipeline construction.
package factory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tmps-go-labs/lab2/domain/models"
)

type PipelineBuilder struct {
	pipeline *models.Pipeline
	factory  ConverterFactory
}

func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{
		pipeline: &models.Pipeline{
			Steps:   make([]models.ConversionStep, 0),
			Options: models.ConversionOptions{},
		},
		factory: NewConverterFactory(),
	}
}

func (b *PipelineBuilder) WithInputPath(path string) *PipelineBuilder {
	b.pipeline.InputPath = path
	return b
}

func (b *PipelineBuilder) WithOutputPath(path string) *PipelineBuilder {
	b.pipeline.OutputPath = path
	return b
}

func (b *PipelineBuilder) WithOptions(options models.ConversionOptions) *PipelineBuilder {
	b.pipeline.Options = options
	return b
}

func (b *PipelineBuilder) WithIndent() *PipelineBuilder {
	b.pipeline.Options.Indent = true
	return b
}

func (b *PipelineBuilder) WithPrettyPrint() *PipelineBuilder {
	b.pipeline.Options.PrettyPrint = true
	return b
}

func (b *PipelineBuilder) WithHeaders(headers []string) *PipelineBuilder {
	b.pipeline.Options.Headers = headers
	return b
}

func (b *PipelineBuilder) WithSaveIntermediarySteps() *PipelineBuilder {
	b.pipeline.Options.SaveIntermediarySteps = true
	return b
}

func (b *PipelineBuilder) AddConversionStep(from, to models.FileFormat) *PipelineBuilder {
	step := models.ConversionStep{
		From: from,
		To:   to,
	}

	b.pipeline.Steps = append(b.pipeline.Steps, step)
	return b
}

func (b *PipelineBuilder) AddCSVToJSON() *PipelineBuilder {
	return b.AddConversionStep(models.FormatCSV, models.FormatJSON)
}

func (b *PipelineBuilder) AddJSONToXML() *PipelineBuilder {
	return b.AddConversionStep(models.FormatJSON, models.FormatXML)
}

func (b *PipelineBuilder) AddXMLToYAML() *PipelineBuilder {
	return b.AddConversionStep(models.FormatXML, models.FormatYAML)
}

func (b *PipelineBuilder) Build() (*models.Pipeline, error) {
	if len(b.pipeline.Steps) == 0 {
		return nil, fmt.Errorf("pipeline must have at least one conversion step")
	}

	if b.pipeline.InputPath == "" {
		return nil, fmt.Errorf("input path is required")
	}

	if b.pipeline.OutputPath == "" {
		return nil, fmt.Errorf("output path is required")
	}

	return b.pipeline, nil
}

type PipelineExecutor struct {
	pool *ConverterPool
}

func NewPipelineExecutor(pool *ConverterPool) *PipelineExecutor {
	return &PipelineExecutor{pool: pool}
}

func (e *PipelineExecutor) Execute(pipeline *models.Pipeline) *models.PipelineResult {
	start := time.Now()
	result := &models.PipelineResult{
		Success: true,
		Results: make([]*models.ConversionResult, 0),
	}

	if len(pipeline.Steps) == 0 {
		result.Success = false
		result.Error = fmt.Errorf("no conversion steps in pipeline")
		return result
	}

	inputData, err := os.ReadFile(pipeline.InputPath)
	if err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to read input file: %w", err)
		return result
	}

	if pipeline.Options.SaveIntermediarySteps {
		if err := os.MkdirAll("steps", 0755); err != nil {
			result.Success = false
			result.Error = fmt.Errorf("failed to create steps directory: %w", err)
			return result
		}
	}

	currentData := inputData
	for i, step := range pipeline.Steps {
		converterType := string(step.From) + "-" + string(step.To)
		converter, err := e.pool.Get(converterType)
		if err != nil {
			result.Success = false
			result.Error = fmt.Errorf("failed to get converter from pool for step %d: %w", i+1, err)
			return result
		}

		conversionResult := converter.Convert(
			strings.NewReader(string(currentData)),
			step.From,
			step.To,
		)

		e.pool.Put(converter)

		result.Results = append(result.Results, conversionResult)

		if conversionResult.Error != nil {
			result.Success = false
			result.Error = fmt.Errorf("step %d failed (%s→%s): %w",
				i+1, step.From, step.To, conversionResult.Error)
			return result
		}

		currentData = conversionResult.Data

		if pipeline.Options.SaveIntermediarySteps {
			stepFileName := filepath.Join("steps", fmt.Sprintf("step_%d_%s_to_%s.%s",
				i+1, step.From, step.To, step.To))
			if err := os.WriteFile(stepFileName, currentData, 0644); err != nil {
				result.Success = false
				result.Error = fmt.Errorf("failed to save intermediary step %d to file: %w", i+1, err)
				return result
			}
		}
	}

	if err := os.WriteFile(pipeline.OutputPath, currentData, 0644); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to write output file: %w", err)
		return result
	}

	result.Duration = time.Since(start).Nanoseconds()
	return result
}
