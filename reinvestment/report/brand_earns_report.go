package report

import (
	"github.com/gabadi/afip-meli-process/base/collector"
	"github.com/gabadi/afip-meli-process/base/processor"
	"github.com/gabadi/afip-meli-process/base/values"
	"github.com/gabadi/afip-meli-process/reinvestment/model"
	"path/filepath"
)

type brandEarnsKey struct {
	Brand string
}

type brandEarnsSummarizationResult struct {
	Brand string             `csv:"marca"`
	Cost  values.MoneyAmount `csv:"Costo"`
	Earns values.MoneyAmount `csv:"Ganancia"`
	ROI   float64            `csv:"ROI"`
}

func NewBrandReport(outputDir string) *processor.SummarizationByKeyProcessor[model.ReportRow, brandEarnsKey, model.EarnCost] {
	return processor.NewSummarizationByKeyProcessor[model.ReportRow, brandEarnsKey, model.EarnCost](
		func(row *model.ReportRow, key *brandEarnsKey) {
			key.Brand = row.ProductBrand
		}, func(row *model.ReportRow) model.EarnCost {
			return model.EarnCost{
				Cost:  row.CostBase,
				Earns: row.EarnsBase,
			}
		}, processor.NewMapperProcessor[processor.Summarization[brandEarnsKey, model.EarnCost], brandEarnsSummarizationResult](
			func(row *processor.Summarization[brandEarnsKey, model.EarnCost], out *brandEarnsSummarizationResult) {
				out.Brand = row.Key.Brand
				out.Earns = row.Aggregation.Earns
				out.Cost = row.Aggregation.Cost
				out.ROI = row.Aggregation.Roi()
			}, collector.NewCSVCollector[brandEarnsSummarizationResult](
				filepath.Join(outputDir, "brand-aggregations.csv"),
			)))
}
