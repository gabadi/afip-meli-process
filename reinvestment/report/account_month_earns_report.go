package report

import (
	"github.com/gabadi/afip-meli-process/base/collector"
	"github.com/gabadi/afip-meli-process/base/processor"
	"github.com/gabadi/afip-meli-process/base/values"
	"github.com/gabadi/afip-meli-process/reinvestment/model"
	"path/filepath"
)

type accountMonthEarnsKey struct {
	Year      int
	Month     int
	AccountId string
}

type accountMonthEarnsSummarizationResult struct {
	AccountId string             `csv:"cuenta"`
	Year      int                `csv:"Ano"`
	Month     int                `csv:"Mes"`
	Amount    values.MoneyAmount `csv:"Ganancia"`
}

func NewAccountMonthReport(outputDir string) *processor.SummarizationByKeyProcessor[model.ReportRow, accountMonthEarnsKey] {
	return processor.NewSummarizationByKeyProcessor[model.ReportRow, accountMonthEarnsKey](
		func(row *model.ReportRow, key *accountMonthEarnsKey) {
			key.AccountId = row.SellerId
			key.Year = row.TransactionDate.Year()
			key.Month = int(row.TransactionDate.Month())
		}, func(row *model.ReportRow) values.MoneyAmount {
			return row.EarnsBase
		}, processor.NewMapperProcessor[processor.Summarization[accountMonthEarnsKey], accountMonthEarnsSummarizationResult](
			func(row *processor.Summarization[accountMonthEarnsKey], out *accountMonthEarnsSummarizationResult) {
				out.AccountId = row.Key.AccountId
				out.Month = row.Key.Month
				out.Year = row.Key.Year
				out.Amount = row.Amount
			}, collector.NewCSVCollector[accountMonthEarnsSummarizationResult](
				filepath.Join(outputDir, "year-month-account-earns.csv"),
			)))
}
