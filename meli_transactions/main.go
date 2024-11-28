package main

import (
	processor2 "github.com/gabadi/afip-meli-process/meli_transactions/processor"
	reader2 "github.com/gabadi/afip-meli-process/meli_transactions/reader"
	report2 "github.com/gabadi/afip-meli-process/meli_transactions/report"
	"log"
	"os"
)

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

	p := processor2.NewListProcessor(
		[]processor2.ReportRowProcessor{
			report2.YearAccountAndTypeReport(outputDir),
			report2.AccountPeriodAndTypeReport(outputDir),
			report2.YearAndTypeReport(outputDir),
			report2.TransferReceivedAsSettlementReport(outputDir),
			report2.PeriodAndTypeReport(outputDir),
			report2.DuplicatesReport(outputDir),
			report2.ShippingSettlementReport(outputDir),
			processor2.NewUnclassifiedPrintProcessor(),
		})

	reportReader := reader2.ReportReader{
		LinePreprocessors: []reader2.LinePreprocessor{
			reader2.NewOfficialStoreLinePreProcessor(),
			reader2.NewMeliPaymentsRemoveLinePreProcessor(),
		},
		Processor: p,
	}

	err := reportReader.Read(inputDir)

	if err != nil {
		log.Fatal("Error al procesar los archivos CSV:", err)
	}
	log.Println("Procesado correctamente")
}
