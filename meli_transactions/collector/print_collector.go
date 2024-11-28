package collector

import (
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"log"
)

type PrintCollector[T any] struct {
}

func (c *PrintCollector[T]) Collect(key *T, amount *model.MoneyAmount) {
	log.Println(
		"Key: ",
		key,
		", Amount:",
		amount.Display(),
	)
}

func (c *PrintCollector[T]) Close() error {
	log.Println("Close")
	return nil
}
