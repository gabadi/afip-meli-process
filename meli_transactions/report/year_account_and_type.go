package report

import (
	collector2 "github.com/gabadi/afip-meli-process/meli_transactions/collector"
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"github.com/gabadi/afip-meli-process/meli_transactions/processor"
	"path/filepath"
)

type yearAccountAndTypeReport struct {
	Year            int                  `csv:"Ano"`
	Cuenta          int                  `csv:"Cuenta"`
	Classification  model.Classification `csv:"Clasificacion"`
	Amount          model.MoneyAmount    `csv:"Monto"`
	RetentionAmount model.MoneyAmount    `csv:"Monto Retenido"`
	EarnsAmount     model.MoneyAmount    `csv:"Monto Ganancias"`
}

type yearAccountAndTypeKey struct {
	Year           int
	Account        int
	Classification model.Classification
}

func YearAccountAndTypeReport(outputDir string) *processor.MapKeySummarizationProcessor[yearAccountAndTypeKey] {
	return processor.NewMapKeySummarizationProcessor[yearAccountAndTypeKey](
		func(classification model.Classification, row *model.ReportRow, outputKey *yearAccountAndTypeKey) {
			outputKey.Year = row.SettlementDate.Year()
			outputKey.Account = row.UserId
			outputKey.Classification = classification
		},
		collector2.NewListCollector[yearAccountAndTypeKey](
			[]collector2.Collector[yearAccountAndTypeKey]{
				&collector2.PrintCollector[yearAccountAndTypeKey]{},
				collector2.NewSummarizedCSVCollector[yearAccountAndTypeKey, yearAccountAndTypeReport](
					func(p *yearAccountAndTypeKey, amount *model.MoneyAmount) *yearAccountAndTypeReport {
						RetentionAmount := *amount
						if !p.Classification.HasRetentions() {
							RetentionAmount = model.NewMoneyAmount()
						}

						EarnsAmount := *amount
						if !p.Classification.PaysEarns() {
							EarnsAmount = model.NewMoneyAmount()
						}
						return &yearAccountAndTypeReport{
							Year:            p.Year,
							Cuenta:          p.Account,
							Classification:  p.Classification,
							Amount:          *amount,
							RetentionAmount: RetentionAmount,
							EarnsAmount:     EarnsAmount,
						}
					},
					filepath.Join(outputDir, "year-account-classification-summarization.csv"),
				),
			}),
	)
}
