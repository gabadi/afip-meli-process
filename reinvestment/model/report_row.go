package model

import (
	"github.com/gabadi/afip-meli-process/base/values"
)

type ReportRow struct {
	SellerId          string             `excel:"vendedor_id"`
	TransactionId     string             `excel:"transaccion_id"`
	TransactionAmount values.MoneyAmount `excel:"transaccion_monto"`
	TransactionType   string             `excel:"transaccion_tipo"`
	TransactionDate   values.Date        `excel:"transaccion_date"`
	ShippingType      string             `excel:"orden_envio"`
	Product           string             `excel:"orden_product"`
	ProductBrand      string             `excel:"orden_product_marca"`
	ProductID         int                `excel:"orden_product_id"`
	ReinvestmentBase  values.MoneyAmount `excel:"reinversion_base"`
	CostBase          values.MoneyAmount `excel:"costo_base"`
	EarnsBase         values.MoneyAmount `excel:"ganancia_base"`
}

func (r *ReportRow) CopyFrom(from *ReportRow) {
	if from == nil {
		return
	}
	r.SellerId = from.SellerId
	r.TransactionId = from.TransactionId
	r.TransactionAmount = from.TransactionAmount
	r.TransactionType = from.TransactionType
	r.TransactionDate = from.TransactionDate
	r.ShippingType = from.ShippingType
	r.Product = from.Product
	r.ProductBrand = from.ProductBrand
	r.ProductID = from.ProductID
	r.ReinvestmentBase = from.ReinvestmentBase
	r.CostBase = from.CostBase
	r.EarnsBase = from.EarnsBase
}
