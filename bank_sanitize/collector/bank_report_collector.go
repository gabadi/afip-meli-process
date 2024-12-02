package collector

import (
	"fmt"
	"github.com/gabadi/afip-meli-process/bank_sanitize/reader"
	"github.com/gocarina/gocsv"
	"os"
)

func NewCSVCollector(path string) *CSVCollector {
	return &CSVCollector{
		rows: make([]ExcelRow, 0),
		path: path,
	}
}

type ExcelRow struct {
	Fuente      string  `excel:"FUENTE"`
	Fecha       string  `excel:"FECHA"`
	Descripcion string  `excel:"DESCRIPCION"`
	Referencia  string  `excel:"REFERENCIA"`
	Debito      float64 `excel:"DEBITO"`
	Credito     float64 `excel:"CREDITO"`
}

type CSVCollector struct {
	rows []ExcelRow
	path string
}

func (c *CSVCollector) Process(row *ExcelRow) (bool, error) {
	c.rows = append(c.rows, *row)
	return true, nil
}

func (c *CSVCollector) Close() error {
	targetFile, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer targetFile.Close()
	return gocsv.MarshalFile(c.rows, targetFile)
}

func NewMergeCollector(target reader.ReportRowProcessor[ExcelRow]) *MergeCollector {
	return &MergeCollector{
		target:     target,
		collectors: make([]reader.ReportRowProcessor[ExcelRow], 0),
		index:      0,
		active:     nil,
	}
}

type MergeCollector struct {
	target     reader.ReportRowProcessor[ExcelRow]
	collectors []reader.ReportRowProcessor[ExcelRow]
	index      int
	active     reader.ReportRowProcessor[ExcelRow]
}

func (c *MergeCollector) Process(row *ExcelRow) (bool, error) {
	if c.active == nil {
		return false, fmt.Errorf("no collector in the collector merger")
	}
	return c.active.Process(row)
}

func (c *MergeCollector) Close() error {
	if len(c.collectors) == c.index {
		return c.target.Close()
	}
	c.index++
	c.active = c.collectors[c.index]
	return nil
}

func (c *MergeCollector) NewInstance() *MergeCollectorInstance {
	collector := &MergeCollectorInstance{
		target: c.target,
	}
	c.collectors = append(c.collectors, collector)
	if c.active == nil {
		c.active = collector
		c.index = 1
	}
	return collector
}

type MergeCollectorInstance struct {
	target reader.ReportRowProcessor[ExcelRow]
}

func (c *MergeCollectorInstance) Process(row *ExcelRow) (bool, error) {
	return c.target.Process(row)
}

func (c *MergeCollectorInstance) Close() error {
	return c.target.Close()
}
