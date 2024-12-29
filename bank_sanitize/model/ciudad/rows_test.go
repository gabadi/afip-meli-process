package ciudad

import (
	"github.com/gabadi/afip-meli-process/base"
	reader2 "github.com/gabadi/afip-meli-process/base/reader"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReportReader_Read_CiudadExcel(t *testing.T) {
	result := reader2.ReadTestRows[ExcelRow](t, func(processor base.ReportRowProcessor[ExcelRow]) base.ReportRowProcessor[ExcelRow] {
		return NewCiudadSanitizer(processor)
	})

	assert.Equal(t, 14, len(result))
	assert.Contains(t, result, ExcelRow{
		Fecha:       "01/04/2022",
		Descripcion: "DEBITO FISCAL IVA BASICO",
		Referencia:  "0",
		Debito:      0.28,
	})
	assert.Contains(t, result, ExcelRow{
		Fecha:       "25/04/2022",
		Descripcion: "DEBIN: PDX4OGNYGOZOG74N0L6EY5           - VAR",
		Referencia:  "270382",
		Credito:     10000.00,
	})
	assert.Contains(t, result, ExcelRow{
		Fecha:       "25/04/2022",
		Descripcion: "N/D DEBITO PRESTAMOS",
		Referencia:  "520",
		Debito:      84004.12,
	})

}
