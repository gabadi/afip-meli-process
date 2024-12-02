package reader

func NewMapperReader[IN any, OUT any](
	out ReportRowProcessor[OUT],
	mapper func(row *IN) *OUT,
) *MapperReader[IN, OUT] {
	return &MapperReader[IN, OUT]{
		out:    out,
		mapper: mapper,
	}
}

type MapperReader[IN any, OUT any] struct {
	out    ReportRowProcessor[OUT]
	mapper func(row *IN) *OUT
}

func (m *MapperReader[IN, OUT]) Process(row *IN) (bool, error) {
	return m.out.Process(m.mapper(row))
}

func (m *MapperReader[IN, OUT]) Close() error {
	return m.out.Close()
}
