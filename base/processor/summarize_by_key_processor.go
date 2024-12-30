package processor

import (
	"fmt"
	"github.com/gabadi/afip-meli-process/base"
	"reflect"
)

type Aggregation[T any] interface {
	Add(another *T) *T
}

type Summarization[Key comparable, T Aggregation[T]] struct {
	Key         *Key
	Aggregation *T
}

type SummarizationByKeyProcessor[T any, Key comparable, S Aggregation[S]] struct {
	keyFactory    func(row *T, key *Key)
	summarization map[Key]S
	amountFactory func(row *T) S
	processor     base.ReportRowProcessor[Summarization[Key, S]]
}

func NewSummarizationByKeyProcessor[T any, Key comparable, S Aggregation[S]](
	keyFactory func(row *T, key *Key),
	amountFactory func(row *T) S,
	processor base.ReportRowProcessor[Summarization[Key, S]]) *SummarizationByKeyProcessor[T, Key, S] {
	return &SummarizationByKeyProcessor[T, Key, S]{
		keyFactory:    keyFactory,
		summarization: make(map[Key]S),
		processor:     processor,
		amountFactory: amountFactory,
	}
}

func (p *SummarizationByKeyProcessor[T, Key, S]) Process(row *T) (bool, error) {
	key := p.new()
	p.keyFactory(row, &key)
	summarization, exists := p.summarization[key]
	if !exists {
		summarization = p.newSummarization()
	}

	toAdd := p.amountFactory(row)
	newAmount := summarization.Add(&toAdd)
	p.summarization[key] = *newAmount
	return true, nil
}

func (p *SummarizationByKeyProcessor[T, Key, S]) Close() error {
	for key, aggregation := range p.summarization {
		processed, err := p.processor.Process(&Summarization[Key, S]{
			Key:         &key,
			Aggregation: &aggregation,
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

func (p *SummarizationByKeyProcessor[T, Key, S]) new() Key {
	valueType := reflect.TypeOf(*new(Key))
	value := reflect.New(valueType).Elem()
	return value.Interface().(Key)
}

func (p *SummarizationByKeyProcessor[T, Key, S]) newSummarization() S {
	valueType := reflect.TypeOf(*new(S))
	value := reflect.New(valueType).Elem()
	return value.Interface().(S)
}
