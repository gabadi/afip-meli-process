package reader

import "regexp"

type RegexRemoveLinePreProcessor struct {
	pattern *regexp.Regexp
}

func NewOfficialStoreLinePreProcessor() *RegexRemoveLinePreProcessor {
	re := regexp.MustCompile(`\{"official_store_id":(null|\d+)\}`)
	return &RegexRemoveLinePreProcessor{
		pattern: re,
	}
}

func NewMeliPaymentsRemoveLinePreProcessor() *RegexRemoveLinePreProcessor {
	re := regexp.MustCompile(`"\[\{.*?\}\]"`)
	return &RegexRemoveLinePreProcessor{
		pattern: re,
	}
}

func (pp *RegexRemoveLinePreProcessor) PreProcess(line string) string {
	return pp.pattern.ReplaceAllString(line, "")
}
