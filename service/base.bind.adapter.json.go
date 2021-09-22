package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
)

func csvToJson(data []byte) []byte {
	c := csv.NewReader(bytes.NewReader(data))
	lines, err := c.ReadAll()
	if err != nil {
		panic(err)
	}

	var header []string
	records := make([]map[string]string, 0, len(lines) - 1)
	for _, line := range lines {
		if header == nil || len(header) <= 0 {
			header = line
			continue
		}
		record := map[string]string{}
		for k, h := range header {
			record[h] = line[k]
		}
		records = append(records, record)
	}

	out, err := json.Marshal(records)
	if err != nil {
		panic(err)
	}

	return out
}