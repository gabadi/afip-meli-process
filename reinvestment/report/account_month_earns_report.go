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
	Cost      values.MoneyAmount `csv:"Costo"`
	Earns     values.MoneyAmount `csv:"Ganancia"`
	ROI       float64            `csv:"ROI"`
	Orders    int                `csv:"Cantidad"`
}

func NewAccountMonthReport(outputDir string) *processor.SummarizationByKeyProcessor[model.ReportRow, accountMonthEarnsKey, model.EarnCost] {
	return processor.NewSummarizationByKeyProcessor[model.ReportRow, accountMonthEarnsKey, model.EarnCost](
		func(row *model.ReportRow, key *accountMonthEarnsKey) {
			key.AccountId = row.SellerId
			key.Year = row.TransactionDate.Year()
			key.Month = int(row.TransactionDate.Month())
		}, func(row *model.ReportRow) model.EarnCost {
			return model.EarnCost{
				Cost:   row.CostBase,
				Earns:  row.EarnsBase,
				Orders: 1,
			}
		}, processor.NewMapperProcessor[processor.Summarization[accountMonthEarnsKey, model.EarnCost], accountMonthEarnsSummarizationResult](
			func(row *processor.Summarization[accountMonthEarnsKey, model.EarnCost], out *accountMonthEarnsSummarizationResult) {
				out.AccountId = row.Key.AccountId
				out.Month = row.Key.Month
				out.Year = row.Key.Year
				out.Earns = row.Aggregation.Earns
				out.Cost = row.Aggregation.Cost
				out.Orders = row.Aggregation.Orders
				out.ROI = row.Aggregation.Roi()
			}, collector.NewCSVCollector[accountMonthEarnsSummarizationResult](
				filepath.Join(outputDir, "year-month-account-aggregations.csv"),
			)))
}
