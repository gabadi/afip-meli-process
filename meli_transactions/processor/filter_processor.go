package processor

import (
	"fmt"
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
)

type Collector[Key any] interface {
	Collect(key *Key)
	Close() error
}

func NewClassificationFilterProcessor[T any](
	classifications []model.Classification,
	mapper func(classification model.Classification, row *model.ReportRow) *T,
	collector Collector[T]) *ClassificationFilterProcessor[T] {
	m := make(map[model.Classification]bool)
	for _, classification := range classifications {
		m[classification] = true
	}
	return &ClassificationFilterProcessor[T]{
		classifications: m,
		mapper:          mapper,
		collector:       collector,
	}
}

type ClassificationFilterProcessor[T any] struct {
	classifications map[model.Classification]bool
	mapper          func(classification model.Classification, row *model.ReportRow) *T
	collector       Collector[T]
}

func (p *ClassificationFilterProcessor[T]) Process(classification model.Classification, row *model.ReportRow) {
	if _, exists := p.classifications[classification]; exists {
		p.collector.Collect(p.mapper(classification, row))
	}
}

func (p *ClassificationFilterProcessor[T]) Close() error {
	return p.collector.Close()
}

func NewDuplicatesFilterProcessor[T any](
	mapper func(classification model.Classification, row *model.ReportRow) *T,
	collector Collector[T]) *DuplicatesFilterProcessor[T] {
	return &DuplicatesFilterProcessor[T]{
		externalRefsCount: make(map[string]int),
		mapper:            mapper,
		collector:         collector,
	}
}

type DuplicatesFilterProcessor[T any] struct {
	externalRefsCount map[string]int
	mapper            func(classification model.Classification, row *model.ReportRow) *T
	collector         Collector[T]
}

func (p *DuplicatesFilterProcessor[T]) Process(classification model.Classification, row *model.ReportRow) {
	if row.ExternalRef != "" && row.SourceId != "" {
		key := fmt.Sprintf("%s/%s/%s/%s/%s",
			row.SourceId,
			row.Type,
			row.SettlementDate.GoString(),
			row.Amount.Display(),
			row.SettlementNetAmount.Display())
		p.externalRefsCount[key]++
		if p.externalRefsCount[key] == 2 {
			p.collector.Collect(p.mapper(classification, row))
		}
	}
}

func (p *DuplicatesFilterProcessor[T]) Close() error {
	return p.collector.Close()
}
