package main

import (
	"github.com/gabadi/afip-meli-process/sanitize_galicia_excel_report/model"
	"github.com/gabadi/afip-meli-process/sanitize_galicia_excel_report/reader"
	"log"
	"os"
)

type Closeable struct {
}

func (c *Closeable) Close() error {
	return nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Por favor, especifica el directorio de archivos CSV")
	}
	inputDir := os.Args[1]
	outputDir := os.Args[2]
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			log.Fatal("Error al crear el directorio de salida:", err)
		}
	}

	reportReader := reader.NewReportReader[model.GaliciaExcelRow](&Closeable{})

	err := reportReader.Read(inputDir)

	if err != nil {
		log.Fatal("Error al procesar los archivos XLSX:", err)
	}
	log.Println("Procesado correctamente")
}
