package model

import (
	"github.com/Rhymond/go-money"
	"strconv"
	"time"
)

const dateFormat = "2006-01-02T15:04:05.000-03:00"

type Date struct {
	time.Time
}

func (date *Date) UnmarshalCSV(csv string) (err error) {
	parsedTime, err := time.Parse(dateFormat, csv)
	if err != nil {
		return err
	}

	date.Time = parsedTime.In(buenosAiresLocation)
	return nil
}

func getBuenosAiresLocation() *time.Location {
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		panic(err)
	}
	return loc
}

var buenosAiresLocation = getBuenosAiresLocation()

func NewMoneyAmount() MoneyAmount {
	return MoneyAmount{
		Money: money.New(0, "ARS"),
	}
}

type MoneyAmount struct {
	*money.Money
}

func (ma *MoneyAmount) Add(other *MoneyAmount) (*MoneyAmount, error) {
	result, err := ma.Money.Add(other.Money)
	if err != nil {
		return nil, err
	} else {
		return &MoneyAmount{result}, nil
	}
}

func (ma *MoneyAmount) IsPositive() bool {
	return ma.Money.IsPositive()
}

func (ma *MoneyAmount) UnmarshalCSV(csv string) error {
	amount, err := strconv.ParseFloat(csv, 64)
	if err != nil {
		return err
	}

	ma.Money = money.NewFromFloat(amount, "ARS")
	return nil
}

func (ma *MoneyAmount) MarshalCSV() (string, error) {
	return strconv.FormatFloat(ma.Money.AsMajorUnits(), 'f', 2, 64), nil
}

type ReportRow struct {
	SellerId         string      `excel:"vendedor_id"`
	TransactionId    string      `excel:"transaccion_id"`
	TransactionDate  Date        `excel:"transaccion_date"`
	ShippingType     string      `excel:"orden_envio"`
	Product          string      `excel:"orden_product"`
	ProductBrand     string      `excel:"orden_product_marca"`
	ProductID        int         `excel:"orden_product_id"`
	ReinvestmentBase MoneyAmount `excel:"reinversion_base"`
	CostBase         MoneyAmount `excel:"costo_base"`
	EarnsBase        MoneyAmount `excel:"ganancia_base"`
}
