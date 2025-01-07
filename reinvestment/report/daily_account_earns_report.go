package report

import (
	"github.com/gabadi/afip-meli-process/base/collector"
	"github.com/gabadi/afip-meli-process/base/processor"
	"github.com/gabadi/afip-meli-process/base/values"
	"github.com/gabadi/afip-meli-process/reinvestment/model"
	"path/filepath"
	"time"
)

type dailyAccountEarnsKey struct {
	Year      int
	Month     int
	Day       int
	AccountId string
}

type dailyAccountEarnsSummarizationResult struct {
	Date      string             `csv:"Date"`
	AccountId string             `csv:"Cuenta"`
	Cost      values.MoneyAmount `csv:"Costo"`
	Earns     values.MoneyAmount `csv:"Ganancia"`
	ROI       float64            `csv:"ROI"`
}

func NewDailyAccountReport(outputDir string) *processor.SummarizationByKeyProcessor[model.ReportRow, dailyAccountEarnsKey, model.EarnCost] {
	return processor.NewSummarizationByKeyProcessor[model.ReportRow, dailyAccountEarnsKey, model.EarnCost](
		func(row *model.ReportRow, key *dailyAccountEarnsKey) {
			key.Year = row.TransactionDate.Year()
			key.Month = int(row.TransactionDate.Month())
			key.Day = row.TransactionDate.Day()
			key.AccountId = row.SellerId
		}, func(row *model.ReportRow) model.EarnCost {
			return model.EarnCost{
				Cost:  row.CostBase,
				Earns: row.EarnsBase,
			}
		}, processor.NewMapperProcessor[processor.Summarization[dailyAccountEarnsKey, model.EarnCost], dailyAccountEarnsSummarizationResult](
			func(row *processor.Summarization[dailyAccountEarnsKey, model.EarnCost], out *dailyAccountEarnsSummarizationResult) {
				date := time.Date(row.Key.Year, time.Month(row.Key.Month), row.Key.Day, 0, 0, 0, 0, time.UTC)

				out.Date = date.Format("2006-01-02")
				out.Earns = row.Aggregation.Earns
				out.Cost = row.Aggregation.Cost
				out.ROI = row.Aggregation.Roi()
			}, collector.NewCSVCollector[dailyAccountEarnsSummarizationResult](
				filepath.Join(outputDir, "daily-account-aggregations.csv"),
			)))
}
