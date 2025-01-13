package values

import (
	"github.com/Rhymond/go-money"
	"strconv"
)

func NewZeroMoneyAmount() MoneyAmount {
	return MoneyAmount{
		Money: money.New(0, "ARS"),
	}
}

func NewMoneyAmount(mon *money.Money) MoneyAmount {
	return MoneyAmount{
		Money: mon,
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
