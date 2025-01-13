package model

import (
	"github.com/gabadi/afip-meli-process/base/values"
)

type ReportRow struct {
	SellerId                string             `excel:"vendedor_id"`
	TransactionId           string             `excel:"transaccion_id"`
	TransactionAmount       values.MoneyAmount `excel:"transaccion_monto"`
	TransactionType         string             `excel:"transaccion_tipo"`
	TransactionDate         values.Date        `excel:"transaccion_date"`
	ShippingType            string             `excel:"orden_envio"`
	Product                 string             `excel:"orden_product"`
	ProductBrand            string             `excel:"orden_product_marca"`
	ProductID               int                `excel:"orden_product_id"`
	ReinvestmentBase        values.MoneyAmount `excel:"reinversion_base"`
	CostBase                values.MoneyAmount `excel:"costo_base"`
	EarnsBase               values.MoneyAmount `excel:"ganancia_base"`
	GrossReinvestmentIva105 values.MoneyAmount `excel:"reinversion_iva105_base"`
	GrossReinvestmentIva21  values.MoneyAmount `excel:"reinversion_iva21_base"`
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
	r.GrossReinvestmentIva21 = from.GrossReinvestmentIva21
	r.GrossReinvestmentIva105 = from.GrossReinvestmentIva105
}

type EarnCost struct {
	Earns  values.MoneyAmount
	Cost   values.MoneyAmount
	Orders int
}

func (ec EarnCost) Roi() float64 {
	if ec == (EarnCost{}) || ec.Cost.Money.IsZero() {
		return 1
	}
	return (float64(ec.Earns.Amount()) / float64(ec.Cost.Amount())) + 1.0
}

func (ec EarnCost) Add(other *EarnCost) *EarnCost {
	if ec == (EarnCost{}) {
		if other == nil {
			panic("both nil")
		}
		return other
	}
	if other == nil {
		return &ec
	}
	resultEarn, err := ec.Earns.Money.Add(other.Earns.Money)
	if err != nil {
		panic(err)
	}
	resultCost, err := ec.Cost.Money.Add(other.Cost.Money)
	if err != nil {
		panic(err)
	}
	result := &EarnCost{
		Earns:  values.NewZeroMoneyAmount(),
		Cost:   values.NewZeroMoneyAmount(),
		Orders: ec.Orders + other.Orders,
	}
	result.Earns.Money = resultEarn
	result.Cost.Money = resultCost
	return result
}
