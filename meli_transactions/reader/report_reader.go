package reader

import (
	"bufio"
	"fmt"
	"github.com/gabadi/afip-meli-process/meli_transactions/model"
	"github.com/gocarina/gocsv"
	"os"
	"path/filepath"
)

type LinePreprocessor interface {
	PreProcess(line string) string
}

type ReportRowProcessor interface {
	Process(classification model.Classification, row *model.ReportRow)
	Close() error
}

type ReportReader struct {
	LinePreprocessors []LinePreprocessor
	Processor         ReportRowProcessor
}

func (rr *ReportReader) Read(dir string) error {
	if err := rr.readDir(dir); err != nil {
		return err
	}
	return rr.Processor.Close()
}

func (rr *ReportReader) readDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking dir", err)
		}
		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			err := rr.processCSV(path)
			if err != nil {
				return fmt.Errorf("error processing file %s: %v", path, err)
			}
		}
		return nil
	})
}

func (rr *ReportReader) processCSV(filename string) error {
	tempFile, err := os.CreateTemp("", "processed_*.csv")
	if err != nil {
		return fmt.Errorf("was not able to create temporary file: %v", err)
	}
	defer tempFile.Close()

	if err := rr.preProcessCSV(filename, tempFile); err != nil {
		return err
	}

	if _, err := tempFile.Seek(0, 0); err != nil {
		return fmt.Errorf("was not able to reset file %s: %v", filename, err)
	}

	tempReader := bufio.NewReader(tempFile)
	return gocsv.UnmarshalToCallback(tempReader, func(line *model.ReportRow) {
		rr.Processor.Process(line.Classify(), line)
	})
}

func (rr *ReportReader) preProcessCSV(inputFileName string, outputFileName *os.File) error {
	originalFile, err := os.Open(inputFileName)
	if err != nil {
		return fmt.Errorf("was not able to open file %s: %v", inputFileName, err)
	}
	defer originalFile.Close()

	scanner := bufio.NewScanner(originalFile)
	writer := bufio.NewWriter(outputFileName)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		for _, preprocessor := range rr.LinePreprocessors {
			line = preprocessor.PreProcess(line)
		}
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing to temporary file: %v", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing to temporary file: %v", err)
	}

	return nil
}
