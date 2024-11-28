package reader

import (
	"github.com/gabadi/afip-meli-process/sanitize_galicia_excel_report/model"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"testing"
)

type Closeable struct {
}

func (c *Closeable) Close() error {
	return nil
}

func currentDir(t *testing.T) string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("cannot get current file")
	}
	return filepath.Dir(currentFile)
}

func TestReportReader_Read_GaliciaExcel(t *testing.T) {
	reportReader := NewReportReader[model.GaliciaExcelRow](&Closeable{})
	err := reportReader.Read(currentDir(t))
	assert.NoError(t, err, "Error al procesar los archivos XLSX")
}
