package processor

import (
	"fmt"
	"github.com/gabadi/afip-meli-process/base"
	"github.com/gabadi/afip-meli-process/base/values"
	"reflect"
)

type Summarization[Key comparable] struct {
	Key    *Key
	Amount values.MoneyAmount
}

type SummarizationByKeyProcessor[T any, Key comparable] struct {
	keyFactory    func(row *T, key *Key)
	summarization map[Key]values.MoneyAmount
	amountFactory func(row *T) values.MoneyAmount
	processor     base.ReportRowProcessor[Summarization[Key]]
}

func NewSummarizationByKeyProcessor[T any, Key comparable](
	keyFactory func(row *T, key *Key),
	amountFactory func(row *T) values.MoneyAmount,
	processor base.ReportRowProcessor[Summarization[Key]]) *SummarizationByKeyProcessor[T, Key] {
	return &SummarizationByKeyProcessor[T, Key]{
		keyFactory:    keyFactory,
		summarization: make(map[Key]values.MoneyAmount),
		processor:     processor,
		amountFactory: amountFactory,
	}
}

func (p *SummarizationByKeyProcessor[T, Key]) Process(row *T) (bool, error) {
	key := p.new()
	p.keyFactory(row, &key)
	amount, exists := p.summarization[key]
	if !exists {
		amount = values.NewMoneyAmount()
	}

	toAdd := p.amountFactory(row)
	newAmount, err := amount.Add(&toAdd)
	if err != nil {
		return false, fmt.Errorf(
			"error adding amount %v to %v: %v",
			toAdd,
			amount,
			err)
	}
	p.summarization[key] = *newAmount
	return true, nil
}

func (p *SummarizationByKeyProcessor[T, Key]) Close() error {
	for key, amount := range p.summarization {
		processed, err := p.processor.Process(&Summarization[Key]{
			Key:    &key,
			Amount: amount,
		})
		if err != nil {
			return fmt.Errorf(
				"error processing row %v: %v",
				key,
				err)
		}
		if !processed {
			return p.processor.Close()
		}
	}
	return p.processor.Close()
}

func (p *SummarizationByKeyProcessor[T, Key]) new() Key {
	valueType := reflect.TypeOf(*new(Key))
	value := reflect.New(valueType).Elem()
	return value.Interface().(Key)
}
