package processor

import (
	"fmt"
	"github.com/gabadi/afip-meli-process/base"
)

func NewListProcessor[P base.ReportRowProcessor[T], T any](processors []P) *ListProcessor[P, T] {
	return &ListProcessor[P, T]{
		processors: processors,
	}
}

type ListProcessor[P base.ReportRowProcessor[T], T any] struct {
	processors []P
}

func (p *ListProcessor[P, T]) Process(row *T) (bool, error) {
	for _, processor := range p.processors {
		if _, err := processor.Process(row); err != nil {
			return false, fmt.Errorf("error processing row: %w", err)
		}
	}
	return true, nil
}

func (p *ListProcessor[P, T]) Close() error {
	for _, processor := range p.processors {
		err := processor.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
