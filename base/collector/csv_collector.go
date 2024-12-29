package collector

import (
	"fmt"
	"github.com/gabadi/afip-meli-process/base"
	"github.com/gocarina/gocsv"
	"os"
)

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

func (c *CSVCollector[T]) Process(row *T) (bool, error) {
	c.rows = append(c.rows, *row)
	return true, nil
}

func (c *CSVCollector[T]) Close() error {
	targetFile, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer targetFile.Close()
	return gocsv.MarshalFile(c.rows, targetFile)
}

func NewMergeCollector[T any](target base.ReportRowProcessor[T]) *MergeCollector[T] {
	return &MergeCollector[T]{
		target:     target,
		collectors: make([]base.ReportRowProcessor[T], 0),
		index:      0,
		active:     nil,
	}
}

type MergeCollector[T any] struct {
	target     base.ReportRowProcessor[T]
	collectors []base.ReportRowProcessor[T]
	index      int
	active     base.ReportRowProcessor[T]
}

func (c *MergeCollector[T]) Process(row *T) (bool, error) {
	if c.active == nil {
		return false, fmt.Errorf("no collector in the collector merger")
	}
	return c.active.Process(row)
}

func (c *MergeCollector[T]) Close() error {
	if len(c.collectors) == c.index {
		return c.target.Close()
	}
	c.index++
	c.active = c.collectors[c.index]
	return nil
}

func (c *MergeCollector[T]) NewInstance() *MergeCollectorInstance[T] {
	collector := &MergeCollectorInstance[T]{
		target: c.target,
	}
	c.collectors = append(c.collectors, collector)
	if c.active == nil {
		c.active = collector
		c.index = 1
	}
	return collector
}

type MergeCollectorInstance[T any] struct {
	target base.ReportRowProcessor[T]
}

func (c *MergeCollectorInstance[T]) Process(row *T) (bool, error) {
	return c.target.Process(row)
}

func (c *MergeCollectorInstance[T]) Close() error {
	return c.target.Close()
}
