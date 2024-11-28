package model

type GaliciaExcelRow struct {
	Fecha       string  `excel:"Fecha"`
	Descripcion string  `excel:"Descripción"`
	Origen      string  `excel:"Origen"`
	Debito      float64 `excel:"Débito(-)"`
	Credito     float64 `excel:"Crédito(+)"`
	Saldo       float64 `excel:"Saldo" optional:"true"`
}
