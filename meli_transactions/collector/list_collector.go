package collector

import (
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
)

func NewListCollector[T any](collectors []Collector[T]) *ListCollector[T] {
	return &ListCollector[T]{
		collectors: collectors,
	}
}

type ListCollector[T any] struct {
	collectors []Collector[T]
}

func (c *ListCollector[T]) Collect(key *T, amount *model.MoneyAmount) {
	for _, c := range c.collectors {
		c.Collect(key, amount)
	}
}

func (c *ListCollector[T]) Close() error {
	for _, c := range c.collectors {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}
