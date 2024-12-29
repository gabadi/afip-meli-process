package report

import (
	"github.com/gabadi/afip-meli-process/base/collector"
	"github.com/gabadi/afip-meli-process/base/processor"
	"github.com/gabadi/afip-meli-process/base/values"
	"github.com/gabadi/afip-meli-process/reinvestment/model"
	"path/filepath"
)

type yearMonthEarnsKey struct {
	Year  int
	Month int
}

type yearMonthEarnsSummarizationResult struct {
	Year   int                `csv:"Ano"`
	Month  int                `csv:"Mes"`
	Amount values.MoneyAmount `csv:"Ganancia"`
}

func NewYearMonthEarnsReport(outputDir string) *processor.SummarizationByKeyProcessor[model.ReportRow, yearMonthEarnsKey] {
	return processor.NewSummarizationByKeyProcessor[model.ReportRow, yearMonthEarnsKey](
		func(row *model.ReportRow, key *yearMonthEarnsKey) {
			key.Year = row.TransactionDate.Year()
			key.Month = int(row.TransactionDate.Month())
		}, func(row *model.ReportRow) values.MoneyAmount {
			return row.EarnsBase
		}, processor.NewMapperProcessor[processor.Summarization[yearMonthEarnsKey], yearMonthEarnsSummarizationResult](
			func(row *processor.Summarization[yearMonthEarnsKey], out *yearMonthEarnsSummarizationResult) {
				out.Month = row.Key.Month
				out.Year = row.Key.Year
				out.Amount = row.Amount
			}, collector.NewCSVCollector[yearMonthEarnsSummarizationResult](
				filepath.Join(outputDir, "year-month-earns.csv"),
			)))
}
