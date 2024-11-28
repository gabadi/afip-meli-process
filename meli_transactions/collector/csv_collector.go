package collector

import (
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"github.com/gocarina/gocsv"
	"os"
)

func NewSummarizedCSVCollector[I any, O any](
	mapper func(*I, *model.MoneyAmount) *O,
	path string,
) *CSVMapper[I, O] {
	collector := CSVCollector[O]{
		rows: make([]O, 0),
		path: path,
	}
	return &CSVMapper[I, O]{
		mapper:       mapper,
		csvCollector: collector,
	}
}

func NewCSVCollector[T any](path string) *CSVCollector[T] {
	return &CSVCollector[T]{
		rows: make([]T, 0),
		path: path,
	}
}

type CSVCollector[T any] struct {
	rows []T
	path string
}

func (c *CSVCollector[T]) Collect(row *T) {
	c.rows = append(c.rows, *row)
}

func (c *CSVCollector[T]) Close() error {
	targetFile, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer targetFile.Close()
	return gocsv.MarshalFile(c.rows, targetFile)
}

type CSVMapper[I any, O any] struct {
	mapper       func(*I, *model.MoneyAmount) *O
	csvCollector CSVCollector[O]
}

func (c *CSVMapper[I, O]) Collect(key *I, amount *model.MoneyAmount) {
	row := c.mapper(key, amount)
	c.csvCollector.Collect(row)
}

func (c *CSVMapper[I, O]) Close() error {
	return c.csvCollector.Close()
}
