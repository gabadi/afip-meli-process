package processor

import (
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
)

type Processor interface {
	Process(classification model.Classification, row *model.ReportRow)
	Close() error
}
