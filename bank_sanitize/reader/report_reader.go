package reader

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type UnmarshalCSV interface {
	UnmarshalCSV(string) error
}

type ReportRowProcessor[T any] interface {
	Process(row *T) (bool, error)
	Close() error
}

func NewReportReader[T any](processor ReportRowProcessor[T]) *ReportReader[T] {
	return &ReportReader[T]{
		processor: processor,
	}
}

type ReportReader[T any] struct {
	processor ReportRowProcessor[T]
}

func (rr *ReportReader[T]) Read(dir string) error {
	if err := rr.readDir(dir); err != nil {
		return err
	}
	return rr.processor.Close()
}

func (rr *ReportReader[T]) readDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking dir: %v", err)
		}
		if !info.IsDir() && filepath.Ext(path) == ".xlsx" {
			err := rr.processXlsx(path)
			if err != nil {
				return fmt.Errorf("error processing file %s: %v", path, err)
			}
		}
		return nil
	})
}

func (rr *ReportReader[T]) processXlsx(path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("was not able to open file %s: %v", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("error closing file %s: %v", path, err)
		}
	}()

	isEmpty := true
	for _, sheet := range f.GetSheetList() {
		isEmptySheet, err := rr.processSheet(f, sheet)
		if err != nil {
			return fmt.Errorf("error processing sheet %s: %v", sheet, err)
		}
		if !isEmptySheet {
			isEmpty = false
		}
	}
	if isEmpty {
		return fmt.Errorf("no rows were processed from file %s", path)
	}
	return nil
}

func (rr *ReportReader[T]) processSheet(f *excelize.File, sheet string) (bool, error) {
	rows, err := f.Rows(sheet)
	if err != nil {
		return true, fmt.Errorf("error getting rows: %v", err)
	}

	sheetReader := newExcelSheetReader[T](rr.processor)
	defer func(sheetReader *excelSheetReader[T]) {
		err := sheetReader.Close()
		if err != nil {
			log.Fatalf("error closing sheet %s from reader %s: %v", sheet, f.Path, err)
		}
	}(sheetReader)

	i := 0
	for rows.Next() {
		i++
		column, err := rows.Columns()
		if err != nil {
			return sheetReader.WasEmpty(), fmt.Errorf("error getting columns from row %d: %v", i, err)
		}
		goOn, err := sheetReader.Process(column)
		if err != nil {
			return sheetReader.WasEmpty(), fmt.Errorf("error processing row number %d: %v", i, err)
		}
		if !goOn {
			return sheetReader.WasEmpty(), nil
		}
	}
	return sheetReader.WasEmpty(), nil
}

func newExcelSheetReader[T any](
	processor ReportRowProcessor[T],
) *excelSheetReader[T] {
	elemType := reflect.TypeOf(new(T)).Elem()
	elem := reflect.New(elemType).Elem()
	headers := make([]string, elem.NumField())
	for i := 0; i < elem.NumField(); i++ {
		fieldType := elemType.Field(i)
		tag := fieldType.Tag.Get("excel")
		headers[i] = tag
	}
	return &excelSheetReader[T]{
		headersMap: make(map[string]int),
		headers:    headers,
		processor:  processor,
	}
}

type excelSheetReader[T any] struct {
	headersMap map[string]int
	headers    []string
	processor  ReportRowProcessor[T]
}

func (esr *excelSheetReader[T]) Close() error {
	return nil
}

func (esr *excelSheetReader[T]) WasEmpty() bool {
	return len(esr.headersMap) == 0
}

func (esr *excelSheetReader[T]) new() T {
	valueType := reflect.TypeOf(*new(T))
	value := reflect.New(valueType).Elem()
	return value.Interface().(T)
}

func (esr *excelSheetReader[T]) Process(row []string) (bool, error) {
	if len(esr.headersMap) != 0 {
		return esr.processReportRow(row)
	}
	esr.tryStart(row)
	return true, nil
}

func (esr *excelSheetReader[T]) processReportRow(row []string) (bool, error) {
	elemType := reflect.TypeOf(new(T)).Elem()

	elem := reflect.New(elemType).Elem()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := elemType.Field(i)
		tag := fieldType.Tag.Get("excel")

		if idx, ok := esr.headersMap[tag]; ok && idx < len(row) {
			value := row[idx]

			switch field.Kind() {
			case reflect.String:
				field.SetString(strings.Replace(value, "\n", "", 10000))
			case reflect.Int:
				if value == "" {
					field.SetInt(0)
					continue
				}
				intVal, err := strconv.Atoi(value)
				if err != nil {
					return false, fmt.Errorf("error parsing int value %s: %v", value, err)
				}
				field.SetInt(int64(intVal))
			case reflect.Float64:
				if value == "" {
					field.SetFloat(0)
					continue
				}
				value = strings.Replace(value, ",", "", 10000)
				value = strings.Replace(value, "$", "", 10000)
				value = strings.Trim(value, " ")
				floatVal, err := strconv.ParseFloat(strings.Replace(value, ",", "", 10000), 64)
				if err != nil {
					if fieldType.Tag.Get("optional") != "true" {
						return false, fmt.Errorf("error parsing float value %s: %v", value, err)
					}
					log.Printf("ignoring error parsing float value %s: %v", value, err)
				}
				field.SetFloat(floatVal)
			default:
				if unmarshaler, ok := interface{}(field.Addr().Interface()).(UnmarshalCSV); ok {
					if err := unmarshaler.UnmarshalCSV(value); err != nil {
						return false, fmt.Errorf("error unmarshalling CSV value for %s: %v", fieldType.Name, err)
					}
					continue
				}
				return false, fmt.Errorf("unknown field type %s in %s", fieldType.Type, fieldType.Name)
			}
		}
	}
	goOn, err := esr.processor.Process(elem.Addr().Interface().(*T))
	if err != nil {
		return false, fmt.Errorf("error processing row: %v", err)
	}
	return goOn, nil
}

func (esr *excelSheetReader[T]) tryStart(row []string) bool {
	rowMapping := make(map[string]int)
	for i, header := range row {
		rowMapping[header] = i
	}
	for _, header := range esr.headers {
		if idx, ok := rowMapping[header]; ok {
			esr.headersMap[header] = idx
		} else {
			return false
		}
	}
	return true
}
