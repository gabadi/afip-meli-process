package collector

import (
	"log"
)

func NewPrintCollector[T any]() *PrintCollector[T] {
	return &PrintCollector[T]{}
}

type PrintCollector[T any] struct {
}

func (p *PrintCollector[T]) Process(row *T) {
	log.Println(row)
}

func (p *PrintCollector[T]) Close() error {
	return nil
}
