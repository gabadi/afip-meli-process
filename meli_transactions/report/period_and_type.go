package report

import (
	collector2 "github.com/gabadi/afip-meli-process/meli_transactions/collector"
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"github.com/gabadi/afip-meli-process/meli_transactions/processor"
	"path/filepath"
	"time"
)

type periodAndTypeReport struct {
	Year            int                  `csv:"Ano"`
	Month           int                  `csv:"Mes"`
	Classification  model.Classification `csv:"Clasificacion"`
	Amount          model.MoneyAmount    `csv:"Monto"`
	RetentionAmount model.MoneyAmount    `csv:"Monto Retenido"`
	EarnsAmount     model.MoneyAmount    `csv:"Monto Ganancias"`
}

type periodAndTypeKey struct {
	Year           int
	Month          time.Month
	Classification model.Classification
}

func PeriodAndTypeReport(outputDir string) *processor.MapKeySummarizationProcessor[periodAndTypeKey] {
	return processor.NewMapKeySummarizationProcessor[periodAndTypeKey](
		func(classification model.Classification, row *model.ReportRow, outputKey *periodAndTypeKey) {
			outputKey.Year = row.SettlementDate.Year()
			outputKey.Month = row.SettlementDate.Month()
			outputKey.Classification = classification
		},
		collector2.NewListCollector[periodAndTypeKey](
			[]collector2.Collector[periodAndTypeKey]{
				&collector2.PrintCollector[periodAndTypeKey]{},
				collector2.NewSummarizedCSVCollector[periodAndTypeKey, periodAndTypeReport](
					func(p *periodAndTypeKey, amount *model.MoneyAmount) *periodAndTypeReport {
						RetentionAmount := *amount
						if !p.Classification.HasRetentions() {
							RetentionAmount = model.NewMoneyAmount()
						}

						EarnsAmount := *amount
						if !p.Classification.PaysEarns() {
							EarnsAmount = model.NewMoneyAmount()
						}
						return &periodAndTypeReport{
							Year:            p.Year,
							Month:           int(p.Month),
							Classification:  p.Classification,
							Amount:          *amount,
							RetentionAmount: RetentionAmount,
							EarnsAmount:     EarnsAmount,
						}
					},
					filepath.Join(outputDir, "period-classification-summarization.csv"),
				),
			}),
	)
}
