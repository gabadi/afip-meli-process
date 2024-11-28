package processor

import (
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"log"
)

func NewDuplicatedLogProcessor() *DuplicatedLogProcessor {
	return &DuplicatedLogProcessor{
		externalRefsCount: make(map[string]int),
	}
}

type DuplicatedLogProcessor struct {
	externalRefsCount map[string]int
}

func (p *DuplicatedLogProcessor) Process(_ model.Classification, row *model.ReportRow) {
	if row.ExternalRef != "" && row.SourceId != "" {
		p.externalRefsCount[row.SourceId+"/"+row.Type+"/"+row.SettlementDate.GoString()+"/"+row.Amount.Display()]++
	}
}

func (p *DuplicatedLogProcessor) Close() error {
	for ref, count := range p.externalRefsCount {
		if count > 1 {
			log.Printf("Duplicated SourceId: %s, Count: %d\n", ref, count)
		}
	}
	return nil
}
