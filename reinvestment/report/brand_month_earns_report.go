package report

import (
	"github.com/gabadi/afip-meli-process/base/collector"
	"github.com/gabadi/afip-meli-process/base/processor"
	"github.com/gabadi/afip-meli-process/base/values"
	"github.com/gabadi/afip-meli-process/reinvestment/model"
	"path/filepath"
)

type brandMonthEarnsKey struct {
	Year  int
	Month int
	Brand string
}

type brandMonthEarnsSummarizationResult struct {
	Brand string             `csv:"marca"`
	Year  int                `csv:"Ano"`
	Month int                `csv:"Mes"`
	Cost  values.MoneyAmount `csv:"Costo"`
	Earns values.MoneyAmount `csv:"Ganancia"`
	ROI   float64            `csv:"ROI"`
}

func NewBrandMonthReport(outputDir string) *processor.SummarizationByKeyProcessor[model.ReportRow, brandMonthEarnsKey, model.EarnCost] {
	return processor.NewSummarizationByKeyProcessor[model.ReportRow, brandMonthEarnsKey, model.EarnCost](
		func(row *model.ReportRow, key *brandMonthEarnsKey) {
			key.Brand = row.ProductBrand
			key.Year = row.TransactionDate.Year()
			key.Month = int(row.TransactionDate.Month())
		}, func(row *model.ReportRow) model.EarnCost {
			return model.EarnCost{
				Cost:  row.CostBase,
				Earns: row.EarnsBase,
			}
		}, processor.NewMapperProcessor[processor.Summarization[brandMonthEarnsKey, model.EarnCost], brandMonthEarnsSummarizationResult](
			func(row *processor.Summarization[brandMonthEarnsKey, model.EarnCost], out *brandMonthEarnsSummarizationResult) {
				out.Brand = row.Key.Brand
				out.Month = row.Key.Month
				out.Year = row.Key.Year
				out.Earns = row.Aggregation.Earns
				out.Cost = row.Aggregation.Cost
				out.ROI = row.Aggregation.Roi()
			}, collector.NewCSVCollector[brandMonthEarnsSummarizationResult](
				filepath.Join(outputDir, "brand-month-account-aggregations.csv"),
			)))
}
