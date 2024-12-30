package report

import (
	"github.com/gabadi/afip-meli-process/base/collector"
	"github.com/gabadi/afip-meli-process/base/processor"
	"github.com/gabadi/afip-meli-process/base/values"
	"github.com/gabadi/afip-meli-process/reinvestment/model"
	"path/filepath"
)

type monthEarnsKey struct {
	Year  int
	Month int
}

type monthEarnsSummarizationResult struct {
	Year  int                `csv:"Ano"`
	Month int                `csv:"Mes"`
	Cost  values.MoneyAmount `csv:"Costo"`
	Earns values.MoneyAmount `csv:"Ganancia"`
	ROI   float64            `csv:"ROI"`
}

func NewMonthReport(outputDir string) *processor.SummarizationByKeyProcessor[model.ReportRow, monthEarnsKey, model.EarnCost] {
	return processor.NewSummarizationByKeyProcessor[model.ReportRow, monthEarnsKey, model.EarnCost](
		func(row *model.ReportRow, key *monthEarnsKey) {
			key.Year = row.TransactionDate.Year()
			key.Month = int(row.TransactionDate.Month())
		}, func(row *model.ReportRow) model.EarnCost {
			return model.EarnCost{
				Cost:  row.CostBase,
				Earns: row.EarnsBase,
			}
		}, processor.NewMapperProcessor[processor.Summarization[monthEarnsKey, model.EarnCost], monthEarnsSummarizationResult](
			func(row *processor.Summarization[monthEarnsKey, model.EarnCost], out *monthEarnsSummarizationResult) {
				out.Month = row.Key.Month
				out.Year = row.Key.Year
				out.Earns = row.Aggregation.Earns
				out.Cost = row.Aggregation.Cost
				out.ROI = row.Aggregation.Roi()
			}, collector.NewCSVCollector[monthEarnsSummarizationResult](
				filepath.Join(outputDir, "year-month-aggregations.csv"),
			)))
}
