package credicop

import (
	"github.com/gabadi/afip-meli-process/base"
	reader2 "github.com/gabadi/afip-meli-process/base/reader"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReportReader_Read_CredicopExcel(t *testing.T) {
	result := reader2.ReadTestRows[ExcelRow](t, func(processor base.ReportRowProcessor[ExcelRow]) base.ReportRowProcessor[ExcelRow] {
		return NewCredicopSanitizer(processor)
	})
	assert.Equal(t, 18, len(result))
	assert.Contains(t, result, ExcelRow{
		Fecha:       "11/08/22",
		Descripcion: "Devolucion Comision de Chequeras",
		Referencia:  "418484",
		Credito:     260.0,
	})
	assert.Contains(t, result, ExcelRow{
		Fecha:       "26/08/22",
		Descripcion: "Com. mantenimiento cuenta",
		Referencia:  "222383",
		Debito:      1650.00,
	})
	assert.Contains(t, result, ExcelRow{
		Fecha:       "31/08/22",
		Descripcion: "Impuesto Ley 25.413 Alic Gral s/Debitos",
		Debito:      0.36,
		Saldo:       -2089.8,
	})
}
