package report

import (
	"github.com/gabadi/afip-meli-process/base/collector"
	"github.com/gabadi/afip-meli-process/base/processor"
	"github.com/gabadi/afip-meli-process/base/values"
	"github.com/gabadi/afip-meli-process/reinvestment/model"
	"path/filepath"
	"strings"
	"time"
)

type dailyEarnsKey struct {
	Year  int
	Month int
	Day   int
}

type dailyEarnsSummarizationResult struct {
	Date  string             `csv:"Date"`
	Cost  values.MoneyAmount `csv:"Costo"`
	Earns values.MoneyAmount `csv:"Ganancia"`
	ROI   float64            `csv:"ROI"`
}

func NewDailyReport(outputDir string, withMelech bool) *processor.SummarizationByKeyProcessor[model.ReportRow, dailyEarnsKey, model.EarnCost] {
	reportName := "daily-aggregations.csv"
	if !withMelech {
		reportName = "daily-aggregations-no-melech.csv"
	}
	return processor.NewSummarizationByKeyProcessor[model.ReportRow, dailyEarnsKey, model.EarnCost](
		func(row *model.ReportRow, key *dailyEarnsKey) {
			key.Year = row.TransactionDate.Year()
			key.Month = int(row.TransactionDate.Month())
			key.Day = row.TransactionDate.Day()
		}, func(row *model.ReportRow) model.EarnCost {
			if !withMelech && strings.EqualFold(row.ProductBrand, "Melech") {
				return model.EarnCost{
					Cost:  values.NewMoneyAmount(),
					Earns: values.NewMoneyAmount(),
				}
			}
			return model.EarnCost{
				Cost:  row.CostBase,
				Earns: row.EarnsBase,
			}
		}, processor.NewMapperProcessor[processor.Summarization[dailyEarnsKey, model.EarnCost], dailyEarnsSummarizationResult](
			func(row *processor.Summarization[dailyEarnsKey, model.EarnCost], out *dailyEarnsSummarizationResult) {
				date := time.Date(row.Key.Year, time.Month(row.Key.Month), row.Key.Day, 0, 0, 0, 0, time.UTC)

				out.Date = date.Format("2006-01-02")
				out.Earns = row.Aggregation.Earns
				out.Cost = row.Aggregation.Cost
				out.ROI = row.Aggregation.Roi()
			}, collector.NewCSVCollector[dailyEarnsSummarizationResult](
				filepath.Join(outputDir, reportName),
			)))
}
