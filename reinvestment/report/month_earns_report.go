package report

import (
	"github.com/gabadi/afip-meli-process/base/collector"
	"github.com/gabadi/afip-meli-process/base/processor"
	"github.com/gabadi/afip-meli-process/base/values"
	"github.com/gabadi/afip-meli-process/reinvestment/model"
	"path/filepath"
	"strings"
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

func NewMonthReport(outputDir string, withMelech bool) *processor.SummarizationByKeyProcessor[model.ReportRow, monthEarnsKey, model.EarnCost] {
	reportName := "year-month-aggregations.csv"
	if !withMelech {
		reportName = "year-month-no-melech-aggregations.csv"
	}
	return processor.NewSummarizationByKeyProcessor[model.ReportRow, monthEarnsKey, model.EarnCost](
		func(row *model.ReportRow, key *monthEarnsKey) {
			key.Year = row.TransactionDate.Year()
			key.Month = int(row.TransactionDate.Month())
		}, func(row *model.ReportRow) model.EarnCost {
			if !withMelech && strings.EqualFold(row.ProductBrand, "Melech") {
				return model.EarnCost{
					Cost:  values.NewZeroMoneyAmount(),
					Earns: values.NewZeroMoneyAmount(),
				}
			}
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
				filepath.Join(outputDir, reportName),
			)))
}
