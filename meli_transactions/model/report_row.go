package model

import (
	"github.com/Rhymond/go-money"
	"strconv"
	"strings"
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
	ExternalRef         string      `csv:"EXTERNAL_REFERENCE"`
	SourceId            string      `csv:"SOURCE_ID"`
	Type                string      `csv:"TRANSACTION_TYPE"`
	PaymentMethod       string      `csv:"PAYMENT_METHOD"`
	UserId              int         `csv:"USER_ID"`
	Amount              MoneyAmount `csv:"TRANSACTION_AMOUNT"`
	SettlementNetAmount MoneyAmount `csv:"SETTLEMENT_NET_AMOUNT"`
	TransactionDate     Date        `csv:"TRANSACTION_DATE"`
	SettlementDate      Date        `csv:"SETTLEMENT_DATE"`
}

type Classification int

const (
	Unclassified Classification = iota
	SubventionFreeShipCancelled
	SubventionFreeShip
	ShipmentPaymentReceivedCancel
	ShipmentPaymentReceived
	ShipmentSettlementPaymentReceivedCancel
	ShipmentSettlementPaymentReceived
	PaymentReceivedCancel
	PaymentReceived
	PaymentMade
	TransferReceived
	TransferMade
	Tax
)

var classificationToString = map[Classification]string{
	Unclassified:                            "Unclassified",
	SubventionFreeShipCancelled:             "Envio subvencionado cancelado",
	SubventionFreeShip:                      "Envio subvencionado",
	ShipmentPaymentReceivedCancel:           "Envio pago recibido cancelado",
	ShipmentPaymentReceived:                 "Envio pago recibido",
	ShipmentSettlementPaymentReceivedCancel: "Envio settlement pago recibido cancelado",
	ShipmentSettlementPaymentReceived:       "Envio settlement pago recibido",
	PaymentReceivedCancel:                   "Pago recibido cancelado",
	PaymentReceived:                         "Pago recibido",
	PaymentMade:                             "Pago hecho",
	TransferReceived:                        "Transferencia recibida",
	TransferMade:                            "Transferencia hecha",
	Tax:                                     "Impuesto",
}

var hasRetentions = map[Classification]bool{
	Unclassified:                            false,
	SubventionFreeShipCancelled:             false,
	SubventionFreeShip:                      false,
	ShipmentSettlementPaymentReceivedCancel: true,
	ShipmentSettlementPaymentReceived:       true,
	ShipmentPaymentReceivedCancel:           true,
	ShipmentPaymentReceived:                 true,
	PaymentReceivedCancel:                   true,
	PaymentReceived:                         true,
	PaymentMade:                             false,
	TransferReceived:                        false,
	TransferMade:                            false,
	Tax:                                     false,
}

var paysEarns = map[Classification]bool{
	Unclassified:                            false,
	SubventionFreeShipCancelled:             false,
	SubventionFreeShip:                      false,
	ShipmentSettlementPaymentReceivedCancel: true,
	ShipmentSettlementPaymentReceived:       true,
	ShipmentPaymentReceivedCancel:           true,
	ShipmentPaymentReceived:                 true,
	PaymentReceivedCancel:                   true,
	PaymentReceived:                         true,
	PaymentMade:                             false,
	TransferReceived:                        true,
	TransferMade:                            false,
	Tax:                                     false,
}

func (c Classification) String() string {
	return classificationToString[c]
}

func (c Classification) HasRetentions() bool {
	return hasRetentions[c]
}

func (c Classification) PaysEarns() bool {
	return paysEarns[c]
}

func (row *ReportRow) Classify() Classification {
	if strings.HasPrefix(row.Type, "TAX_") {
		return Tax
	}
	switch expression := row.Type; expression {
	case "CASHBACK_CANCEL":
		if row.Amount.IsNegative() {
			return SubventionFreeShipCancelled
		} else {
			return Unclassified
		}
	case "CASHBACK":
		if row.Amount.IsPositive() {
			return SubventionFreeShip
		} else {
			return Unclassified
		}
	case "DISPUTE":
		return PaymentReceivedCancel
	case "PAYOUTS":
		if row.Amount.IsNegative() {
			return TransferMade
		} else {
			return Unclassified
		}
	case "REFUND":
		return PaymentReceivedCancel
	case "REFUND_SHIPPING":
		if row.Amount.IsNegative() {
			return ShipmentPaymentReceivedCancel
		} else {
			return Unclassified
		}
	case "SHIPPING":
		if row.Amount.IsPositive() {
			return ShipmentPaymentReceived
		} else {
			return Unclassified
		}
	case "SETTLEMENT_SHIPPING":
		if row.Amount.IsPositive() {
			return ShipmentSettlementPaymentReceived
		} else {
			return Unclassified
		}
	case "SETTLEMENT":
		if specialSettlement := specialSettlementClassify(row); specialSettlement != 0 {
			return specialSettlement
		}
		if row.Amount.IsPositive() {
			return PaymentReceived
		} else {
			return PaymentMade
		}
	case "TRANSFER":
		if row.Amount.IsPositive() {
			return TransferReceived
		} else {
			return Unclassified
		}
	case "WITHDRAWAL_CANCEL":
		if row.Amount.IsPositive() {
			return TransferMade
		} else {
			return Unclassified
		}
	case "WITHDRAWAL":
		if row.Amount.IsNegative() {
			return TransferMade
		} else {
			return Unclassified
		}
	default:
		return Unclassified
	}
}

func specialSettlementClassify(row *ReportRow) Classification {
	if row.Amount.IsNegative() {
		return 0
	}

	externalRef := row.ExternalRef
	if strings.HasPrefix(externalRef, "money_transfer") {
		return TransferReceived
	}
	// estos quiza no retienen
	if strings.HasPrefix(externalRef, "MP-QR") {
		return PaymentReceived
	}
	if (strings.HasPrefix(externalRef, "2") && !strings.HasPrefix(externalRef, "2000")) || externalRef == "" {
		if strings.ToLower(row.PaymentMethod) == "cvu" {
			return TransferReceived
		} else {
			return PaymentReceived
		}
	}

	return 0
}
