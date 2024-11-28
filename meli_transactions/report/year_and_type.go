package report

import (
	collector2 "github.com/gabadi/afip-meli-process/meli_transactions/collector"
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"github.com/gabadi/afip-meli-process/meli_transactions/processor"
	"path/filepath"
)

type yearAndTypeReport struct {
	Year            int                  `csv:"Ano"`
	Classification  model.Classification `csv:"Clasificacion"`
	Amount          model.MoneyAmount    `csv:"Monto"`
	RetentionAmount model.MoneyAmount    `csv:"Monto Retenido"`
	EarnsAmount     model.MoneyAmount    `csv:"Monto Ganancias"`
}

type yearAndTypeKey struct {
	Year           int
	Classification model.Classification
}

func YearAndTypeReport(outputDir string) *processor.MapKeySummarizationProcessor[yearAndTypeKey] {
	return processor.NewMapKeySummarizationProcessor[yearAndTypeKey](
		func(classification model.Classification, row *model.ReportRow, outputKey *yearAndTypeKey) {
			outputKey.Year = row.SettlementDate.Year()
			outputKey.Classification = classification
		},
		collector2.NewListCollector[yearAndTypeKey](
			[]collector2.Collector[yearAndTypeKey]{
				&collector2.PrintCollector[yearAndTypeKey]{},
				collector2.NewSummarizedCSVCollector[yearAndTypeKey, yearAndTypeReport](
					func(p *yearAndTypeKey, amount *model.MoneyAmount) *yearAndTypeReport {
						RetentionAmount := *amount
						if !p.Classification.HasRetentions() {
							RetentionAmount = model.NewMoneyAmount()
						}

						EarnsAmount := *amount
						if !p.Classification.PaysEarns() {
							EarnsAmount = model.NewMoneyAmount()
						}
						return &yearAndTypeReport{
							Year:            p.Year,
							Classification:  p.Classification,
							Amount:          *amount,
							RetentionAmount: RetentionAmount,
							EarnsAmount:     EarnsAmount,
						}
					},
					filepath.Join(outputDir, "year-classification-summarization.csv"),
				),
			}),
	)
}
