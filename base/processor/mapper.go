package processor

import "github.com/gabadi/afip-meli-process/base"

func NewMapperProcessor[IN any, OUT any](
	out base.ReportRowProcessor[OUT],
	mapper func(row *IN) *OUT,
) *MapperProcessor[IN, OUT] {
	return &MapperProcessor[IN, OUT]{
		out:    out,
		mapper: mapper,
	}
}

type MapperProcessor[IN any, OUT any] struct {
	out    base.ReportRowProcessor[OUT]
	mapper func(row *IN) *OUT
}

func (m *MapperProcessor[IN, OUT]) Process(row *IN) (bool, error) {
	return m.out.Process(m.mapper(row))
}

func (m *MapperProcessor[IN, OUT]) Close() error {
	return m.out.Close()
}
