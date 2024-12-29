package main

import (
	"github.com/gabadi/afip-meli-process/bank_sanitize/model/ciudad"
	"github.com/gabadi/afip-meli-process/bank_sanitize/model/credicop"
	"github.com/gabadi/afip-meli-process/bank_sanitize/model/galicia"
	"github.com/gabadi/afip-meli-process/bank_sanitize/model/santander"
	"github.com/gabadi/afip-meli-process/base"
	"github.com/gabadi/afip-meli-process/base/collector"
	"github.com/gabadi/afip-meli-process/base/processor"
	reader2 "github.com/gabadi/afip-meli-process/base/reader"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ExcelRow struct {
	Fuente      string  `excel:"FUENTE"`
	Fecha       string  `excel:"FECHA"`
	Descripcion string  `excel:"DESCRIPCION"`
	Referencia  string  `excel:"REFERENCIA"`
	Debito      float64 `excel:"DEBITO"`
	Credito     float64 `excel:"CREDITO"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Por favor, especifica el directorio de archivos CSV")
	}
	inputDir := os.Args[1]

	fileCollector := collector.NewMergeCollector(collector.NewCSVCollector[ExcelRow](filepath.Join(inputDir, "output.csv")))

	if err := collectGalicia(fileCollector.NewInstance(), inputDir); err != nil {
		log.Fatal("Error al procesar el archivo de Galicia:", err)
	}

	if err := collectSantander(fileCollector.NewInstance(), inputDir); err != nil {
		log.Fatal("Error al procesar el archivo de Santander:", err)
	}

	if err := collectCiudad(fileCollector.NewInstance(), inputDir); err != nil {
		log.Fatal("Error al procesar el archivo de Ciudad:", err)
	}

	if err := collectCredicop(fileCollector.NewInstance(), inputDir); err != nil {
		log.Fatal("Error al procesar el archivo de Credicop:", err)
	}

	log.Println("Procesado correctamente")
}

func collectGalicia(collect base.ReportRowProcessor[ExcelRow], inputDir string) error {
	reportReader := reader2.NewExcelReader(
		galicia.NewGaliciaSanitizer(
			processor.NewMapperProcessor[galicia.ExcelRow, ExcelRow](
				collect,
				func(row *galicia.ExcelRow) *ExcelRow {
					return &ExcelRow{
						Fecha:       row.Fecha,
						Descripcion: row.Descripcion,
						Referencia:  row.Origen,
						Debito:      row.Debito,
						Credito:     row.Credito,
						Fuente:      "GALICIA",
					}
				},
			),
		))

	return reportReader.Read(filepath.Join(inputDir, "galicia", "inputs"))
}

func collectSantander(collect base.ReportRowProcessor[ExcelRow], inputDir string) error {
	reportReader := reader2.NewExcelReader(
		santander.NewSantanderSanitizer(
			processor.NewMapperProcessor[santander.ExcelRow, ExcelRow](
				collect,
				func(row *santander.ExcelRow) *ExcelRow {
					return &ExcelRow{
						Fecha:       row.Fecha,
						Descripcion: row.Descripcion,
						Referencia:  row.Referencia,
						Debito:      row.Debito,
						Credito:     row.Credito,
						Fuente:      "SANTANDER",
					}
				},
			),
		))

	return reportReader.Read(filepath.Join(inputDir, "santander", "inputs"))
}

func collectCiudad(collect base.ReportRowProcessor[ExcelRow], inputDir string) error {
	reportReader := reader2.NewExcelReader(
		ciudad.NewCiudadSanitizer(
			processor.NewMapperProcessor[ciudad.ExcelRow, ExcelRow](
				collect,
				func(row *ciudad.ExcelRow) *ExcelRow {
					return &ExcelRow{
						Fecha:       strings.Trim(row.Fecha, " "),
						Descripcion: strings.Trim(row.Descripcion, " "),
						Referencia:  strings.Trim(row.Referencia, " "),
						Debito:      row.Debito,
						Credito:     row.Credito,
						Fuente:      "CIUDAD",
					}
				},
			),
		))

	return reportReader.Read(filepath.Join(inputDir, "ciudad", "inputs"))
}

func collectCredicop(collect base.ReportRowProcessor[ExcelRow], inputDir string) error {
	reportReader := reader2.NewExcelReader(
		credicop.NewCredicopSanitizer(
			processor.NewMapperProcessor[credicop.ExcelRow, ExcelRow](
				collect,
				func(row *credicop.ExcelRow) *ExcelRow {
					return &ExcelRow{
						Fecha:       row.Fecha,
						Descripcion: row.Descripcion,
						Referencia:  row.Referencia,
						Debito:      row.Debito,
						Credito:     row.Credito,
						Fuente:      "CREDICOP",
					}
				},
			),
		))

	return reportReader.Read(filepath.Join(inputDir, "credicop", "inputs"))
}
