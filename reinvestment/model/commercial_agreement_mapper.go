package model

import (
	"github.com/gabadi/afip-meli-process/base"
	"github.com/gabadi/afip-meli-process/base/processor"
	"strings"
	"time"
)

var bawCutoffDate = time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC)

func NewCommercialAgreementMapper(p base.ReportRowProcessor[ReportRow]) *processor.MapperProcessor[ReportRow, ReportRow] {
	return processor.NewMapperProcessor[ReportRow, ReportRow](
		func(row *ReportRow, out *ReportRow) {
			out.CopyFrom(row)
			if strings.EqualFold(row.ProductBrand, "Sica") {
				parts, err := out.CostBase.Split(40)
				if err != nil {
					panic(err)
				}
				result, err := out.EarnsBase.Money.Add(parts[0])
				out.EarnsBase.Money = result
				if err != nil {
					panic(err)
				}
				out.EarnsBase.Money = result
			}
			if strings.EqualFold(row.ProductBrand, "Baw") && row.TransactionDate.Before(bawCutoffDate) {
				parts, err := out.CostBase.Split(20)
				if err != nil {
					panic(err)
				}
				result, err := out.EarnsBase.Money.Add(parts[0])
				out.EarnsBase.Money = result
				if err != nil {
					panic(err)
				}
				out.EarnsBase.Money = result
			}
		}, p)
}
