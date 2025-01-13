package main

import (
	"github.com/gabadi/afip-meli-process/base"
	"github.com/gabadi/afip-meli-process/base/processor"
	"github.com/gabadi/afip-meli-process/base/reader"
	"github.com/gabadi/afip-meli-process/reinvestment/model"
	"github.com/gabadi/afip-meli-process/reinvestment/report"
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

	p := processor.NewListProcessor[base.ReportRowProcessor[model.ReportRow], model.ReportRow](
		[]base.ReportRowProcessor[model.ReportRow]{
			report.NewDuplicatesPrinterReport(),
			report.NewAccountMonthReport(outputDir),
			report.NewBrandMonthReport(outputDir),
			report.NewBrandReport(outputDir),
			report.NewMonthReport(outputDir),
			report.NewDailyAccountReport(outputDir),
			report.NewDailyReport(outputDir, true),
			report.NewDailyReport(outputDir, false),
			report.NewMonthMelechReinvestmentReport(outputDir),
		})

	agreementProcessor := model.NewCommercialAgreementMapper(p)

	reportReader := reader.NewExcelReader[model.ReportRow](
		agreementProcessor,
	)

	err := reportReader.Read(inputDir)

	if err != nil {
		log.Fatal("Error al procesar los archivos CSV:", err)
	}
	log.Println("Procesado correctamente")
}
