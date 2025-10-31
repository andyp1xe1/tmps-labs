// Package factory implements a pool for managing converter instances.
package factory

import (
	"sync"

	"tmps-go-labs/lab2/domain/models"
)

type ConverterPool struct {
	pool    chan models.Converter
	factory ConverterFactory
	mu      sync.Mutex
	created int
	maxSize int
}

func NewConverterPool(maxSize int, factory ConverterFactory) *ConverterPool {
	return &ConverterPool{
		pool:    make(chan models.Converter, maxSize),
		factory: factory,
		maxSize: maxSize,
	}
}

func (p *ConverterPool) Get(converterType string) (models.Converter, error) {
	select {
	case converter := <-p.pool:
		return converter, nil
	default:
		p.mu.Lock()
		defer p.mu.Unlock()

		if p.created < p.maxSize {
			converter, err := p.factory.CreateConverter(converterType)
			if err != nil {
				return nil, err
			}
			p.created++
			return converter, nil
		}

		select {
		case converter := <-p.pool:
			return converter, nil
		default:
			return p.factory.CreateConverter(converterType)
		}
	}
}

func (p *ConverterPool) Put(converter models.Converter) {
	select {
	case p.pool <- converter:
	default:
	}
}

func (p *ConverterPool) Size() int {
	return len(p.pool)
}

func (p *ConverterPool) Created() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.created
}
