package processor

import (
	"github.com/gabadi/afip-meli-process/base"
	"reflect"
)

func NewMapperProcessor[IN any, OUT any](
	mapper func(row *IN, out *OUT),
	out base.ReportRowProcessor[OUT],
) *MapperProcessor[IN, OUT] {
	return &MapperProcessor[IN, OUT]{
		out:    out,
		mapper: mapper,
	}
}

type MapperProcessor[IN any, OUT any] struct {
	out    base.ReportRowProcessor[OUT]
	mapper func(row *IN, out *OUT)
}

func (m *MapperProcessor[IN, OUT]) Process(row *IN) (bool, error) {
	out := m.new()
	m.mapper(row, &out)
	return m.out.Process(&out)
}

func (m *MapperProcessor[IN, OUT]) Close() error {
	return m.out.Close()
}

func (m *MapperProcessor[IN, OUT]) new() OUT {
	valueType := reflect.TypeOf(*new(OUT))
	value := reflect.New(valueType).Elem()
	return value.Interface().(OUT)
}
