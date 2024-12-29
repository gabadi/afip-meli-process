package processor

import (
	"fmt"
	"github.com/gabadi/afip-meli-process/base"
	"reflect"
)

func NewDuplicatesFilterProcessor[IN any, KEY comparable](
	keyFactory func(row *IN, key *KEY) bool,
	processor base.ReportRowProcessor[KEY]) *DuplicatesFilterProcessor[IN, KEY] {
	return &DuplicatesFilterProcessor[IN, KEY]{
		keysCount:  make(map[KEY]int),
		keyFactory: keyFactory,
		processor:  processor,
	}
}

type DuplicatesFilterProcessor[IN any, KEY comparable] struct {
	keysCount  map[KEY]int
	keyFactory func(row *IN, KEY *KEY) bool
	processor  base.ReportRowProcessor[KEY]
}

func (p *DuplicatesFilterProcessor[IN, KEY]) Process(row *IN) (bool, error) {
	key := p.new()
	if p.keyFactory(row, &key) {
		p.keysCount[key]++
		if p.keysCount[key] == 2 {
			processed, err := p.processor.Process(&key)
			if err != nil {
				return false, fmt.Errorf(
					"error processing row %v: %v",
					row,
					err,
				)
			}
			return processed, nil
		}
	}
	return true, nil
}

func (p *DuplicatesFilterProcessor[IN, KEY]) Close() error {
	return p.processor.Close()
}

func (p *DuplicatesFilterProcessor[IN, KEY]) new() KEY {
	valueType := reflect.TypeOf(*new(KEY))
	value := reflect.New(valueType).Elem()
	return value.Interface().(KEY)
}
