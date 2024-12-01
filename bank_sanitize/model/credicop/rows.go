package credicop

import (
	"fmt"
	"github.com/gabadi/afip-meli-process/bank_sanitize/reader"
)

type ExcelRow struct {
	Fecha       string  `excel:"FECHA"`
	Descripcion string  `excel:"DESCRIPCION"`
	Referencia  string  `excel:"COMBTE"`
	Debito      float64 `excel:"DEBITO"`
	Credito     float64 `excel:"CREDITO"`
	Saldo       float64 `excel:"SALDO"`
}

func NewCredicopSanitizer(processor reader.ReportRowProcessor[ExcelRow]) *Sanitizer {
	return &Sanitizer{
		processor: processor,
	}
}

type Sanitizer struct {
	processor reader.ReportRowProcessor[ExcelRow]
}

func (s *Sanitizer) Process(row *ExcelRow) (bool, error) {
	if row.Fecha == "" && row.Descripcion == "" && row.Referencia == "" {
		return true, nil
	}
	if row.Referencia == "SALDO" {
		return false, nil
	}
	goOn, err := s.processor.Process(row)
	if err != nil {
		return false, fmt.Errorf("error processing row %v: %v", row, err)
	}
	return goOn, nil
}

func (s *Sanitizer) Close() error {
	return s.processor.Close()
}
