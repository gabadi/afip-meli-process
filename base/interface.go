package base

type ReportRowProcessor[T any] interface {
	Process(row *T) (bool, error)
	Close() error
}
