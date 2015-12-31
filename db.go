package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/serve/libs/metal"
	_ "github.com/serve/libs/pq"
)

func getData(query string, m *metal.Metal) {
	db, _ := sql.Open("postgres", "user=postgres password=root dbname=test sslmode=disable")
	rows, _ := db.Query(query)
	defer rows.Close()

	columns, _ := rows.Columns()

	rowIdx := 0

	for rows.Next() {

		row := make([]interface{}, len(columns))
		for idx := range columns {
			row[idx] = new(MetalScanner)
		}

		err := rows.Scan(row...)
		if err != nil {
			fmt.Println(err)
		}
		var mrow = metal.NewMetal()
		for idx, column := range columns {
			var scanner = row[idx].(*MetalScanner)
			mrow.Set(column, scanner.value)
		}
		m.Set("@"+strconv.Itoa(rowIdx), mrow)
		rowIdx++
	}
}

type MetalScanner struct {
	valid bool
	value interface{}
}

func (scanner *MetalScanner) getBytes(src interface{}) []byte {
	if a, ok := src.([]uint8); ok {
		return a
	}
	return nil
}

func (scanner *MetalScanner) Scan(src interface{}) error {
	switch src.(type) {
	case int64:
		if value, ok := src.(int64); ok {
			scanner.value = value
			scanner.valid = true
		}
	case float64:
		if value, ok := src.(float64); ok {
			scanner.value = value
			scanner.valid = true
		}
	case bool:
		if value, ok := src.(bool); ok {
			scanner.value = value
			scanner.valid = true
		}
	case string:
		value := scanner.getBytes(src)
		scanner.value = string(value)
		scanner.valid = true
	case []byte:
		value := scanner.getBytes(src)
		scanner.value = string(value)
		scanner.valid = true
	case time.Time:
		if value, ok := src.(time.Time); ok {
			scanner.value = value
			scanner.valid = true
		}
	case nil:
		scanner.value = nil
		scanner.valid = true
	}
	return nil
}
