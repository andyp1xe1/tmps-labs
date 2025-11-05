// Package factory implements creational design patterns for file format converters.
// It provides Factory Method pattern for converter creation, Object Pool pattern
// for converter reuse, and Builder pattern for pipeline construction.
package factory

import (
	"fmt"
	"sync"

	"tmps-go-labs/lab2/domain/models"
)

type ConverterCreator func() models.Converter

var (
	converterRegistry = make(map[string]ConverterCreator)
	registryMutex     sync.RWMutex
)

func RegisterConverter(formatType string, creator ConverterCreator) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	converterRegistry[formatType] = creator
}

type ConverterFactory interface {
	CreateConverter(formatType string) (models.Converter, error)
}

type DefaultConverterFactory struct{}

func NewConverterFactory() ConverterFactory {
	return &DefaultConverterFactory{}
}

func (f *DefaultConverterFactory) CreateConverter(formatType string) (models.Converter, error) {
	registryMutex.RLock()
	creator, exists := converterRegistry[formatType]
	registryMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unsupported converter type: %s", formatType)
	}

	return creator(), nil
}
