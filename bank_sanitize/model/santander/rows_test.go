package santander

import (
	"github.com/gabadi/afip-meli-process/bank_sanitize/reader"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReportReader_Read_SantanderExcel(t *testing.T) {
	result := reader.ReadTestRows[ExcelRow](t, func(processor reader.ReportRowProcessor[ExcelRow]) reader.ReportRowProcessor[ExcelRow] {
		return NewSantanderSanitizer(processor)
	})

	assert.Equal(t, 23, len(result))
	assert.Contains(t, result, ExcelRow{
		Fecha:       "30/08/22",
		Descripcion: "Impuesto ley 25.413 debito 0,6%",
		Debito:      281.88,
		Saldo:       52046.29,
	})
	assert.Contains(t, result, ExcelRow{
		Fecha:       "30/08/22",
		Descripcion: "Debito automaticoFed patronal sa -2500787136000000150822",
		Referencia:  "73061",
		Debito:      42854.64,
		Saldo:       56453.17,
	})
	assert.Contains(t, result, ExcelRow{
		Fecha:       "01/08/22",
		Descripcion: "Transf recibida cvu mismo titularDe gamalectric s.a. / mercado pago /30715567578",
		Referencia:  "6242321",
		Credito:     350000.00,
		Saldo:       463957.21,
	})
}
