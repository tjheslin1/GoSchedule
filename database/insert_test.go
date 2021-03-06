package database

import (
	"testing"
	"time"
)

func TestInsertStatement(t *testing.T) {
	tableEntry := TableEntry{
		Name: "TestTableName",
		Data: map[string]TableCell{
			"column1": IntCell{1},
			"column2": Int64Cell{123456789123345},
			"column3": Float32Cell{50.0},
			"column4": BoolCell{true},
			"column5": TimeCell{time.Date(1970, 2, 3, 4, 5, 6, 0, time.UTC)},
			"column6": StringCell{"stringstring"},
		},
	}

	expectedStatement := "INSERT INTO TestTableName (column1, column2, column3, column4, column5, column6) VALUES (1, 123456789123345, 50.0000, true, 1970-02-03 04:05:06 +0000 UTC, 'stringstring')"
	actualStatement := insertStatement(tableEntry)

	if actualStatement != expectedStatement {
		t.Errorf("Expected generated statement:\n%s\nto equal expected:\n%s\n", actualStatement, expectedStatement)
	}
}
