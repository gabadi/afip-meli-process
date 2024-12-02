package main

import (
	"github.com/gabadi/afip-meli-process/bank_sanitize/collector"
	"github.com/gabadi/afip-meli-process/bank_sanitize/model/ciudad"
	"github.com/gabadi/afip-meli-process/bank_sanitize/model/credicop"
	"github.com/gabadi/afip-meli-process/bank_sanitize/model/galicia"
	"github.com/gabadi/afip-meli-process/bank_sanitize/model/santander"
	"github.com/gabadi/afip-meli-process/bank_sanitize/reader"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Por favor, especifica el directorio de archivos CSV")
	}
	inputDir := os.Args[1]

	fileCollector := collector.NewMergeCollector(collector.NewCSVCollector(filepath.Join(inputDir, "output.csv")))

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

func collectGalicia(collect reader.ReportRowProcessor[collector.ExcelRow], inputDir string) error {
	reportReader := reader.NewReportReader(
		galicia.NewGaliciaSanitizer(
			reader.NewMapperReader[galicia.ExcelRow, collector.ExcelRow](
				collect,
				func(row *galicia.ExcelRow) *collector.ExcelRow {
					return &collector.ExcelRow{
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

func collectSantander(collect reader.ReportRowProcessor[collector.ExcelRow], inputDir string) error {
	reportReader := reader.NewReportReader(
		santander.NewSantanderSanitizer(
			reader.NewMapperReader[santander.ExcelRow, collector.ExcelRow](
				collect,
				func(row *santander.ExcelRow) *collector.ExcelRow {
					return &collector.ExcelRow{
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

func collectCiudad(collect reader.ReportRowProcessor[collector.ExcelRow], inputDir string) error {
	reportReader := reader.NewReportReader(
		ciudad.NewCiudadSanitizer(
			reader.NewMapperReader[ciudad.ExcelRow, collector.ExcelRow](
				collect,
				func(row *ciudad.ExcelRow) *collector.ExcelRow {
					return &collector.ExcelRow{
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

func collectCredicop(collect reader.ReportRowProcessor[collector.ExcelRow], inputDir string) error {
	reportReader := reader.NewReportReader(
		credicop.NewCredicopSanitizer(
			reader.NewMapperReader[credicop.ExcelRow, collector.ExcelRow](
				collect,
				func(row *credicop.ExcelRow) *collector.ExcelRow {
					return &collector.ExcelRow{
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
