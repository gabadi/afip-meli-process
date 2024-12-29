package report

import (
	"github.com/gabadi/afip-meli-process/base/collector"
	"github.com/gabadi/afip-meli-process/base/processor"
	"github.com/gabadi/afip-meli-process/base/values"
	"github.com/gabadi/afip-meli-process/reinvestment/model"
	"path/filepath"
)

type yearMonthSellCostKey struct {
	Year  int
	Month int
}

type yearMonthSellCostSummarizationResult struct {
	Year   int                `csv:"Ano"`
	Month  int                `csv:"Mes"`
	Amount values.MoneyAmount `csv:"Costo"`
}

func NewYearMonthSellCostReport(outputDir string) *processor.SummarizationByKeyProcessor[model.ReportRow, yearMonthSellCostKey] {
	return processor.NewSummarizationByKeyProcessor[model.ReportRow, yearMonthSellCostKey](
		func(row *model.ReportRow, key *yearMonthSellCostKey) {
			key.Year = row.TransactionDate.Year()
			key.Month = int(row.TransactionDate.Month())
		}, func(row *model.ReportRow) values.MoneyAmount {
			return row.CostBase
		}, processor.NewMapperProcessor[processor.Summarization[yearMonthSellCostKey], yearMonthSellCostSummarizationResult](
			func(row *processor.Summarization[yearMonthSellCostKey], out *yearMonthSellCostSummarizationResult) {
				out.Month = row.Key.Month
				out.Year = row.Key.Year
				out.Amount = row.Amount
			}, collector.NewCSVCollector[yearMonthSellCostSummarizationResult](
				filepath.Join(outputDir, "year-month-cost.csv"),
			)))
}
