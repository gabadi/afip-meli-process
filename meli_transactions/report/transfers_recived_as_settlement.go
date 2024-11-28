package report

import (
	"github.com/gabadi/afip-meli-process/meli_transactions/collector"
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"github.com/gabadi/afip-meli-process/meli_transactions/processor"
	"path/filepath"
)

func TransferReceivedAsSettlementReport(outputDir string) *processor.ClassificationFilterProcessor[model.ReportRow] {
	return processor.NewClassificationFilterProcessor[model.ReportRow](
		[]model.Classification{model.TransferReceived},
		func(classification model.Classification, row *model.ReportRow) *model.ReportRow {
			return row
		},
		collector.NewCSVCollector[model.ReportRow](
			filepath.Join(outputDir, "transfer-received-as-settlement.csv"),
		))
}
