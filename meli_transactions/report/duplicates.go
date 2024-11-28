package report

import (
	"github.com/gabadi/afip-meli-process/meli_transactions/collector"
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"github.com/gabadi/afip-meli-process/meli_transactions/processor"
	"path/filepath"
)

func DuplicatesReport(outputDir string) *processor.DuplicatesFilterProcessor[model.ReportRow] {
	return processor.NewDuplicatesFilterProcessor[model.ReportRow](
		func(classification model.Classification, row *model.ReportRow) *model.ReportRow {
			return row
		},
		collector.NewCSVCollector[model.ReportRow](
			filepath.Join(outputDir, "duplicates.csv"),
		))
}

func ShippingSettlementReport(outputDir string) *processor.ClassificationFilterProcessor[model.ReportRow] {
	return processor.NewClassificationFilterProcessor[model.ReportRow](
		[]model.Classification{model.ShipmentSettlementPaymentReceived, model.ShipmentSettlementPaymentReceivedCancel},
		func(classification model.Classification, row *model.ReportRow) *model.ReportRow {
			return row
		},
		collector.NewCSVCollector[model.ReportRow](
			filepath.Join(outputDir, "shipping-settlement.csv"),
		))
}
