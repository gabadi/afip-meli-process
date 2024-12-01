package santander

import (
	"fmt"
	"github.com/gabadi/afip-meli-process/bank_sanitize/reader"
	"strings"
)

type ExcelRow struct {
	Fecha       string  `excel:"Fecha"`
	Descripcion string  `excel:"Movimiento"`
	Referencia  string  `excel:"Comprobante"`
	Debito      float64 `excel:"Débito"`
	Credito     float64 `excel:"Crédito"`
	Saldo       float64 `excel:"Saldo en cuenta"`
}

func NewSantanderSanitizer(processor reader.ReportRowProcessor[ExcelRow]) *Sanitizer {
	return &Sanitizer{
		processor: processor,
	}
}

type Sanitizer struct {
	processor reader.ReportRowProcessor[ExcelRow]
	date      string
}

func (s *Sanitizer) Process(row *ExcelRow) (bool, error) {
	if strings.HasPrefix(row.Fecha, "Saldo total") {
		return false, nil
	}
	if row.Descripcion == "Saldo Inicial" {
		return true, nil
	}
	if row.Fecha != "" {
		s.date = row.Fecha
	}
	row.Fecha = s.date
	goOn, err := s.processor.Process(row)
	if err != nil {
		return false, fmt.Errorf("error processing row %v: %v", row, err)
	}
	return goOn, nil
}

func (s *Sanitizer) Close() error {
	return s.processor.Close()
}
