package reader

import (
	"github.com/gabadi/afip-meli-process/base"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"testing"
)

type testProcessor[T any] struct {
	content []T
}

func (c *testProcessor[T]) Close() error {
	return nil
}
func (c *testProcessor[T]) Process(row *T) (bool, error) {
	c.content = append(c.content, *row)
	return true, nil
}

func currentDir(t *testing.T) string {
	_, currentFile, _, ok := runtime.Caller(2)
	if !ok {
		t.Fatalf("cannot get current file")
	}
	return filepath.Dir(currentFile)
}

func ReadTestRows[T any](t *testing.T, processorFactory func(processor base.ReportRowProcessor[T]) base.ReportRowProcessor[T]) []T {
	collector := testProcessor[T]{
		content: make([]T, 0),
	}
	processor := processorFactory(&collector)
	reportReader := NewExcelReader[T](processor)
	err := reportReader.Read(currentDir(t))
	assert.NoError(t, err, "Error al procesar los archivos XLSX")
	return collector.content
}
