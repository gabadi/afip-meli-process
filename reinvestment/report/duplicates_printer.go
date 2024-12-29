package report

import (
	"github.com/gabadi/afip-meli-process/base/collector"
	"github.com/gabadi/afip-meli-process/base/processor"
	"github.com/gabadi/afip-meli-process/reinvestment/model"
)

type UniqueKey struct {
	TransactionId     string
	TransactionType   string
	AccountId         string
	Date              string
	ProductID         int
	TransactionAmount string
}

func NewDuplicatesPrinterReport() *processor.DuplicatesFilterProcessor[model.ReportRow, UniqueKey] {
	return processor.NewDuplicatesFilterProcessor[model.ReportRow, UniqueKey](
		func(row *model.ReportRow, key *UniqueKey) bool {
			if row.TransactionId == "" || row.TransactionId == "Venta presencial" {
				return false
			}
			key.TransactionId = row.TransactionId
			key.TransactionType = row.TransactionType
			key.ProductID = row.ProductID
			key.TransactionAmount = row.TransactionAmount.Display()
			key.AccountId = row.SellerId
			key.Date = row.TransactionDate.Format("2006-01-02T01:11:22")
			return true
		},
		collector.NewPrintCollector[UniqueKey]("Duplicated transaction"),
	)
}
