package database

import (
	"fmt"
	"sort"
	"strings"
)

func insertStatement(jobEntry TableEntry) string {
	var keys []string
	for key := range jobEntry.Data {
		keys = append(keys, key)
	}

	columnNames := make([]string, len(jobEntry.Data))
	values := make([]string, len(jobEntry.Data))

	// sorting is necesary to produce the same SQL string every time.
	sort.Strings(keys)

	i := 0
	for _, key := range keys {
		columnNames[i] = key
		values[i] = jobEntry.Data[key].asString()
		i++
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", jobEntry.Name,
		strings.Join(columnNames, ", "), strings.Join(values, ", "))
}
