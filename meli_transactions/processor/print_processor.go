package processor

import (
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"log"
)

func NewUnclassifiedPrintProcessor() *UnclassifiedPrintProcessor {
	return &UnclassifiedPrintProcessor{}
}

type UnclassifiedPrintProcessor struct {
}

func (p *UnclassifiedPrintProcessor) Process(classification model.Classification, row *model.ReportRow) {
	if classification == model.Unclassified {
		log.Println(classification, row.Amount.Display(), row)
	}
}

func (p *UnclassifiedPrintProcessor) Close() error {
	return nil
}
