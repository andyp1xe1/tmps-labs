// Package factory implements creational design patterns for file format converters.
// It provides Factory Method pattern for converter creation, Object Pool pattern
// for converter reuse, and Builder pattern for pipeline construction.
package factory

import (
	"sync"

	"tmps-go-labs/lab2/domain/models"
)

type ConverterPool struct {
	pools   map[string]chan models.Converter
	factory ConverterFactory
	mu      sync.Mutex
	created map[string]int
	maxSize int
}

func NewConverterPool(maxSize int, factory ConverterFactory) *ConverterPool {
	return &ConverterPool{
		pools:   make(map[string]chan models.Converter),
		factory: factory,
		created: make(map[string]int),
		maxSize: maxSize,
	}
}

func (p *ConverterPool) Get(converterType string) (models.Converter, error) {
	p.mu.Lock()

	if _, exists := p.pools[converterType]; !exists {
		p.pools[converterType] = make(chan models.Converter, p.maxSize)
		p.created[converterType] = 0
	}

	pool := p.pools[converterType]
	p.mu.Unlock()

	select {
	case converter := <-pool:
		return converter, nil
	default:
		p.mu.Lock()
		if p.created[converterType] < p.maxSize {
			converter, err := p.factory.CreateConverter(converterType)
			if err != nil {
				p.mu.Unlock()
				return nil, err
			}
			p.created[converterType]++
			p.mu.Unlock()
			return converter, nil
		}
		p.mu.Unlock()

		select {
		case converter := <-pool:
			return converter, nil
		default:
			return p.factory.CreateConverter(converterType)
		}
	}
}

func (p *ConverterPool) Put(converter models.Converter) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, pool := range p.pools {
		select {
		case pool <- converter:
			return
		default:
			continue
		}
	}
}

func (p *ConverterPool) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	total := 0
	for _, pool := range p.pools {
		total += len(pool)
	}
	return total
}

func (p *ConverterPool) Created() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	total := 0
	for _, count := range p.created {
		total += count
	}
	return total
}
