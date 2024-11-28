package processor

import (
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"reflect"
)

type MapKeyCollector[Key comparable] interface {
	Collect(key *Key, amount *model.MoneyAmount)
	Close() error
}

type MapKeySummarizationProcessor[Key comparable] struct {
	keyFactory    func(classification model.Classification, row *model.ReportRow, outputKey *Key)
	summarization map[Key]model.MoneyAmount
	collector     MapKeyCollector[Key]
}

func NewMapKeySummarizationProcessor[Key comparable](
	keyFactory func(classification model.Classification, row *model.ReportRow, outputKey *Key),
	collector MapKeyCollector[Key]) *MapKeySummarizationProcessor[Key] {
	return &MapKeySummarizationProcessor[Key]{
		keyFactory:    keyFactory,
		summarization: make(map[Key]model.MoneyAmount),
		collector:     collector,
	}
}

func (p *MapKeySummarizationProcessor[Key]) Process(classification model.Classification, row *model.ReportRow) {
	key := p.new()
	p.keyFactory(classification, row, &key)
	amount, exists := p.summarization[key]
	if !exists {
		amount = model.NewMoneyAmount()
	}

	newAmount, err := amount.Add(&row.Amount)
	if err != nil {
		panic(err)
	}
	p.summarization[key] = *newAmount
}

func (p *MapKeySummarizationProcessor[Key]) Close() error {
	for key, amount := range p.summarization {
		p.collector.Collect(&key, &amount)
	}
	return p.collector.Close()
}

func (p *MapKeySummarizationProcessor[Key]) new() Key {
	valueType := reflect.TypeOf(*new(Key))
	value := reflect.New(valueType).Elem()
	return value.Interface().(Key)
}
