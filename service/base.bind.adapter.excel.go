package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

func csvToExcel(data []byte) []byte {
	c := csv.NewReader(bytes.NewReader(data))
	lines, err := c.ReadAll()
	if err != nil {
		panic(err)
	}

	f := excelize.NewFile()
	for row, line := range lines {
		for i, cell := range line {
			column, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				panic(err)
			}
			if err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", column, row+1), cell); err != nil {
				panic(err)
			}
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func excelToCsv(data []byte) ([]byte, error) {
	buf := bytes.NewBufferString("")
	c := csv.NewWriter(buf)

	f, err := excelize.OpenReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	sheetName := "Sheet1"
	if f.GetSheetIndex(sheetName) == -1 {
		sheetName = f.GetSheetName(0)
	}

	headers := -1
	row := 1
	for {
		column := 1
		record := make([]string, 0, 10)
		for {
			letter, err := excelize.ColumnNumberToName(column)
			if err != nil {
				panic(err)
			}
			cell, err := f.GetCellValue(sheetName, fmt.Sprintf("%s%d", letter, row))
			if err != nil {
				panic(err)
			}
			if headers < 0 && strings.TrimSpace(cell) == "" {
				break
			} else if headers >= 0 && column >= headers {
				break
			} else if headers >= 0 && column == 1 && strings.TrimSpace(cell) == "" {
				break
			}
			record = append(record, strings.TrimSpace(strings.TrimPrefix(cell, "'")))
			column++
		}
		if headers < 0 {
			headers = column
		}
		if column == 1 {
			break
		}
		if err := c.Write(record); err != nil {
			panic(err)
		}
		row++
	}
	c.Flush()

	return buf.Bytes(), nil
}