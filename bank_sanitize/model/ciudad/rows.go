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
	Credito     float64 `excel:"CREDITOS" optional:"true"`
	Saldo       float64 `excel:"SALDO" optional:"true"`
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
	if strings.Contains(row.Fecha, "SALDO FINAL AL DIA") || strings.Contains(row.Descripcion, "SALDO FINAL AL DIA") {
		return true, nil
	}
	if row.Fecha == "" && row.Descripcion == "" && row.Referencia == "" {
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
