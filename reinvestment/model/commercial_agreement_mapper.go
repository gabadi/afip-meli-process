package model

import (
	"github.com/gabadi/afip-meli-process/base"
	"github.com/gabadi/afip-meli-process/base/processor"
	"strings"
)

func NewCommercialAgreementMapper(p base.ReportRowProcessor[ReportRow]) *processor.MapperProcessor[ReportRow, ReportRow] {
	return processor.NewMapperProcessor[ReportRow, ReportRow](
		func(row *ReportRow, out *ReportRow) {
			out.CopyFrom(row)
			if strings.EqualFold(row.ProductBrand, "Sica") || strings.EqualFold(row.ProductBrand, "Baw") {
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
