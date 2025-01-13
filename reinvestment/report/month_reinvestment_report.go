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

type monthReinvestmentResult struct {
	Year                       int                `csv:"Ano"`
	Month                      int                `csv:"Mes"`
	AccreditedIva21            values.MoneyAmount `csv:"AcreditadoIva21"`
	PendingAccreditationIva21  values.MoneyAmount `csv:"PendienteAcreditacionIva21"`
	AccreditedIva105           values.MoneyAmount `csv:"AcreditadoIva105"`
	PendingAccreditationIva105 values.MoneyAmount `csv:"PendienteAcreditacionIva105"`
}

type reinvestmentAggregations struct {
	AccreditedIva21            values.MoneyAmount
	PendingAccreditationIva21  values.MoneyAmount
	AccreditedIva105           values.MoneyAmount
	PendingAccreditationIva105 values.MoneyAmount
}

func (ra reinvestmentAggregations) Add(other *reinvestmentAggregations) *reinvestmentAggregations {
	if ra == (reinvestmentAggregations{}) {
		if other == nil {
			panic("both nil")
		}
		return other
	}
	if other == nil {
		return &ra
	}
	PendingAccreditationIva21, err := ra.PendingAccreditationIva21.Money.Add(other.PendingAccreditationIva21.Money)
	if err != nil {
		panic(err)
	}
	PendingAccreditationIva105, err := ra.PendingAccreditationIva105.Money.Add(other.PendingAccreditationIva105.Money)
	if err != nil {
		panic(err)
	}
	AccreditedIva21, err := ra.AccreditedIva21.Money.Add(other.AccreditedIva21.Money)
	if err != nil {
		panic(err)
	}
	AccreditedIva105, err := ra.AccreditedIva105.Money.Add(other.AccreditedIva105.Money)
	if err != nil {
		panic(err)
	}
	return &reinvestmentAggregations{
		PendingAccreditationIva21:  values.NewMoneyAmount(PendingAccreditationIva21),
		PendingAccreditationIva105: values.NewMoneyAmount(PendingAccreditationIva105),
		AccreditedIva21:            values.NewMoneyAmount(AccreditedIva21),
		AccreditedIva105:           values.NewMoneyAmount(AccreditedIva105),
	}
}

var tenDaysAgo = time.Now().AddDate(0, 0, -10)

func NewMonthMelechReinvestmentReport(outputDir string) *processor.SummarizationByKeyProcessor[model.ReportRow, monthEarnsKey, reinvestmentAggregations] {
	return processor.NewSummarizationByKeyProcessor[model.ReportRow, monthEarnsKey, reinvestmentAggregations](
		func(row *model.ReportRow, key *monthEarnsKey) {
			key.Year = row.TransactionDate.Year()
			key.Month = int(row.TransactionDate.Month())
		}, func(row *model.ReportRow) reinvestmentAggregations {
			if !strings.EqualFold(row.ProductBrand, "Melech") {
				return reinvestmentAggregations{
					AccreditedIva21:            values.NewZeroMoneyAmount(),
					PendingAccreditationIva21:  values.NewZeroMoneyAmount(),
					AccreditedIva105:           values.NewZeroMoneyAmount(),
					PendingAccreditationIva105: values.NewZeroMoneyAmount(),
				}
			}
			Reinvestment21 := row.GrossReinvestmentIva21
			Reinvestment21Part, err := Reinvestment21.Split(20)
			if err != nil {
				panic(err)
			}
			Reinvestment105 := row.GrossReinvestmentIva105
			Reinvestment105Part, err := Reinvestment105.Split(20)
			if err != nil {
				panic(err)
			}

			//parts := int64(2)
			//if row.TransactionDate.Year() > 2024 || row.TransactionDate.Month() > 9 {
			//	parts = 3
			//}
			parts := int64(4)

			R105, err := Reinvestment105.Subtract(Reinvestment105Part[0].Multiply(parts))
			if err != nil {
				panic(err)
			}
			R21, err := Reinvestment21.Subtract(Reinvestment21Part[0].Multiply(parts))
			if err != nil {
				panic(err)
			}

			if row.TransactionDate.After(tenDaysAgo) {
				return reinvestmentAggregations{
					AccreditedIva21:            values.NewZeroMoneyAmount(),
					PendingAccreditationIva21:  values.NewMoneyAmount(R21),
					AccreditedIva105:           values.NewZeroMoneyAmount(),
					PendingAccreditationIva105: values.NewMoneyAmount(R105),
				}
			} else {
				return reinvestmentAggregations{
					AccreditedIva21:            values.NewMoneyAmount(R21),
					PendingAccreditationIva21:  values.NewZeroMoneyAmount(),
					AccreditedIva105:           values.NewMoneyAmount(R105),
					PendingAccreditationIva105: values.NewZeroMoneyAmount(),
				}
			}
		}, processor.NewMapperProcessor[processor.Summarization[monthEarnsKey, reinvestmentAggregations], monthReinvestmentResult](
			func(row *processor.Summarization[monthEarnsKey, reinvestmentAggregations], out *monthReinvestmentResult) {
				out.Month = row.Key.Month
				out.Year = row.Key.Year
				out.AccreditedIva105 = row.Aggregation.AccreditedIva105
				out.AccreditedIva21 = row.Aggregation.AccreditedIva21
				out.PendingAccreditationIva105 = row.Aggregation.PendingAccreditationIva105
				out.PendingAccreditationIva21 = row.Aggregation.PendingAccreditationIva21
			}, collector.NewCSVCollector[monthReinvestmentResult](
				filepath.Join(outputDir, "melech-reinvestment-report.csv"),
			)))
}
