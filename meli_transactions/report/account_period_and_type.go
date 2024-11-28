package report

import (
	collector2 "github.com/gabadi/afip-meli-process/meli_transactions/collector"
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"github.com/gabadi/afip-meli-process/meli_transactions/processor"
	"path/filepath"
	"time"
)

type accountPeriodAndTypeReport struct {
	Year            int                  `csv:"Ano"`
	Month           int                  `csv:"Mes"`
	Account         int                  `csv:"Cuenta"`
	Classification  model.Classification `csv:"Clasificacion"`
	Amount          model.MoneyAmount    `csv:"Monto"`
	RetentionAmount model.MoneyAmount    `csv:"Monto Retenido"`
	EarnsAmount     model.MoneyAmount    `csv:"Monto Ganancias"`
}

type accountPeriodAndTypeKey struct {
	Year           int
	Month          time.Month
	Account        int
	Classification model.Classification
}

func AccountPeriodAndTypeReport(outputDir string) *processor.MapKeySummarizationProcessor[accountPeriodAndTypeKey] {
	return processor.NewMapKeySummarizationProcessor[accountPeriodAndTypeKey](
		func(classification model.Classification, row *model.ReportRow, outputKey *accountPeriodAndTypeKey) {
			outputKey.Year = row.SettlementDate.Year()
			outputKey.Month = row.SettlementDate.Month()
			outputKey.Classification = classification
			outputKey.Account = row.UserId
		},
		collector2.NewListCollector[accountPeriodAndTypeKey](
			[]collector2.Collector[accountPeriodAndTypeKey]{
				&collector2.PrintCollector[accountPeriodAndTypeKey]{},
				collector2.NewSummarizedCSVCollector[accountPeriodAndTypeKey, accountPeriodAndTypeReport](
					func(p *accountPeriodAndTypeKey, amount *model.MoneyAmount) *accountPeriodAndTypeReport {
						RetentionAmount := *amount
						if !p.Classification.HasRetentions() {
							RetentionAmount = model.NewMoneyAmount()
						}

						EarnsAmount := *amount
						if !p.Classification.PaysEarns() {
							EarnsAmount = model.NewMoneyAmount()
						}
						return &accountPeriodAndTypeReport{
							Year:            p.Year,
							Month:           int(p.Month),
							Classification:  p.Classification,
							Amount:          *amount,
							Account:         p.Account,
							RetentionAmount: RetentionAmount,
							EarnsAmount:     EarnsAmount,
						}
					},
					filepath.Join(outputDir, "account-period-classification-summarization.csv"),
				),
			}),
	)
}
