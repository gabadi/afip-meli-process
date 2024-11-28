package collector

import (
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
)

type Collector[T any] interface {
	Collect(key *T, amount *model.MoneyAmount)
	Close() error
}
