package collector

import (
	"log"
)

func NewPrintCollector[T any](
	prefix string,
) *PrintCollector[T] {
	return &PrintCollector[T]{
		prefix: prefix,
	}
}

type PrintCollector[T any] struct {
	prefix string
}

func (p *PrintCollector[T]) Process(row *T) (bool, error) {
	log.Println(p.prefix, ": ", row)
	return true, nil
}

func (p *PrintCollector[T]) Close() error {
	return nil
}
