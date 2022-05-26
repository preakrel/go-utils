package files

import (
	"bufio"
	"encoding/csv"
	"github.com/tealeg/xlsx"
	"io"
)

// ExtractTXT 提取txt文件
func ExtractTXT(r io.Reader) ([]string, error) {
	es := make([]string, 0, 10)
	br := bufio.NewReader(r)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		es = append(es, string(line))
	}
	return es, nil
}

// ExtractCSV 提取csv文件
func ExtractCSV(r io.Reader) ([]string, error) {
	es := make([]string, 0, 10)
	r1 := csv.NewReader(r)
	content, err := r1.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, row := range content {
		if len(row) > 0 {
			es = append(es, row[0])
		}
	}
	return es, nil
}

// ExtractXLSX 提取xlsx文件
func ExtractXLSX(filePath string) ([]string, error) {
	es := make([]string, 0, 10)
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			if len(row.Cells) > 0 {
				es = append(es, row.Cells[0].String())
			}

		}
	}
	return es, nil
}
