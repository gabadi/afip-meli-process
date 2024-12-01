package ciudad

import (
	"fmt"
	"github.com/gabadi/afip-meli-process/bank_sanitize/reader"
	"strings"
)

type ExcelRow struct {
	Fecha       string  `excel:"FECHA"`
	Descripcion string  `excel:"DESCRIPCION"`
	Referencia  string  `excel:"REFERENCIA"`
	Debito      float64 `excel:"DEBITOS"`
	Credito     float64 `excel:"CREDITOS"`
	Saldo       float64 `excel:"SALDO"`
}

func NewCiudadSanitizer(processor reader.ReportRowProcessor[ExcelRow]) *Sanitizer {
	return &Sanitizer{
		processor: processor,
	}
}

type Sanitizer struct {
	processor reader.ReportRowProcessor[ExcelRow]
}

func (s *Sanitizer) Process(row *ExcelRow) (bool, error) {
	if strings.HasPrefix(row.Fecha, "SALDO FINAL AL DIA") {
		return false, nil
	}
	if strings.HasPrefix(row.Descripcion, "SALDO FINAL DEL DIA") {
		return true, nil
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
