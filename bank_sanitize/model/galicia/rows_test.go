package galicia

import (
	"github.com/gabadi/afip-meli-process/bank_sanitize/reader"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReportReader_Read_GaliciaExcel(t *testing.T) {
	result := reader.ReadTestRows[ExcelRow](t, func(processor reader.ReportRowProcessor[ExcelRow]) reader.ReportRowProcessor[ExcelRow] {
		return NewGaliciaSanitizer(processor)
	})

	assert.Equal(t, 17, len(result))
	assert.Contains(t, result, ExcelRow{
		Fecha:       "31-03",
		Descripcion: "ECHEQ GALICIA NRO:     659",
		Origen:      "",
		Debito:      55369.8,
		Saldo:       2383866.94,
	})
	assert.Contains(t, result, ExcelRow{
		Fecha:       "01-04",
		Descripcion: "TRANSFERENCIA DE TERCEROS, ALDERETE SERGIO ENRI, 20239067221, VARIOS, OPERACION 1468547290, SALDO ANTERIOR",
		Origen:      "",
		Credito:     2914.6,
		Saldo:       1503086.93,
	})
	assert.Contains(t, result, ExcelRow{
		Fecha:       "12-04",
		Descripcion: "ECHEQ 48 HS. NRO.",
		Origen:      "",
		Debito:      185120.62,
		Saldo:       406397.32,
	})
}
