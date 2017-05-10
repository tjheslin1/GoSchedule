package database

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func insertStatement(jobEntry TableEntry) string {
	var keys []string
	for key := range jobEntry.Data {
		keys = append(keys, key)
	}

	fmt.Printf("STORED KEYS: '%v'", keys)

	columnNames := make([]string, len(jobEntry.Data))
	values := make([]string, len(jobEntry.Data))

	i := 0
	for _, key := range keys {
		columnNames[i] = key

		val := jobEntry.Data[key]
		switch val.(type) {
		case string:
			values[i] = "'" + val.(string) + "'"
		case bool:
			values[i] = strconv.FormatBool(val.(bool))
		case int:
			values[i] = strconv.Itoa(val.(int))
		case float32:
			values[i] = fmt.Sprintf("%.4f", val.(float32))
		case float64:
			values[i] = fmt.Sprintf("%.4f", val.(float64))
		case time.Time:
			values[i] = fmt.Sprintf("%v", val.(time.Time))
		}
		i++
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", jobEntry.Name, strings.Join(columnNames, ", "), strings.Join(values, ", "))
}
