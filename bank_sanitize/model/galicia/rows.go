package galicia

import (
	"fmt"
	"github.com/gabadi/afip-meli-process/base"
	"strings"
)

type ExcelRow struct {
	Fecha       string  `excel:"Fecha"`
	Descripcion string  `excel:"Descripción"`
	Origen      string  `excel:"Origen"`
	Debito      float64 `excel:"Débito(-)"`
	Credito     float64 `excel:"Crédito(+)"`
	Saldo       float64 `excel:"Saldo" optional:"true"`
}

func NewGaliciaSanitizer(processor base.ReportRowProcessor[ExcelRow]) *GaliciaSanitizer {
	return &GaliciaSanitizer{
		processor: processor,
		open:      false,
		row:       nil,
	}
}

type GaliciaSanitizer struct {
	processor base.ReportRowProcessor[ExcelRow]
	open      bool
	row       *ExcelRow
}

func (s *GaliciaSanitizer) Process(row *ExcelRow) (bool, error) {
	if !s.open {
		if row.Descripcion != "SALDO INICIAL" {
			return false, fmt.Errorf("expected first row to be 'SALDO INICIAL', got %s", row.Descripcion)
		}
		s.open = true
		return true, nil
	}
	if s.row == nil {
		s.row = row
		return true, nil
	}
	if row.Fecha == "" {
		s.row.Descripcion += ", " + row.Descripcion
		return true, nil
	}
	if strings.HasPrefix(row.Fecha, "TOTAL RETENCION IMPUESTO") || strings.HasPrefix(row.Fecha, "Los depósitos en pesos y en moneda extranjera cuentan con") {
		return s.closeRow()
	}
	goOn, err := s.processor.Process(s.row)
	if err != nil {
		return false, fmt.Errorf("error processing row %v: %v", s.row, err)
	}
	s.row = row
	return goOn, nil
}

func (s *GaliciaSanitizer) closeRow() (bool, error) {
	if s.row != nil {
		return s.processor.Process(s.row)
	}
	return true, nil
}

func (s *GaliciaSanitizer) Close() error {
	if _, err := s.closeRow(); err != nil {
		return fmt.Errorf("error closing row: %v", err)
	}
	return s.processor.Close()
}
