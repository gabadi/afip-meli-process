package processor

import (
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
)

type ReportRowProcessor interface {
	Process(classification model.Classification, row *model.ReportRow)
	Close() error
}

func NewListProcessor[T ReportRowProcessor](processors []T) *ListProcessor[T] {
	return &ListProcessor[T]{
		processors: processors,
	}
}

type ListProcessor[T ReportRowProcessor] struct {
	processors []T
}

func (p *ListProcessor[T]) Process(classification model.Classification, row *model.ReportRow) {
	for _, processor := range p.processors {
		processor.Process(classification, row)
	}
}

func (p *ListProcessor[T]) Close() error {
	for _, processor := range p.processors {
		err := processor.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
