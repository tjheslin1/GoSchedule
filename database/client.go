package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

// TableEntry represents data to be inserted into the database,
// the Data field is a map of column name -> value for a new row.
type TableEntry struct {
	Name string
	Data map[string]TableCell
}

// DBClient defines the interactions with the database.
type DBClient interface {
	Connection() *sql.DB
	SubmitJob(TableEntry) error
}

// TableCell represents any entry into the database.
//
// The `asString()` function defines how the cell's value is represented in SQL.
type TableCell interface {
	asString() string
}

// IntCell represents an integer to be stored in the database.
type IntCell struct {
	Value int
}

func (cell IntCell) asString() string {
	return strconv.Itoa(cell.Value)
}

// StringCell represnts a string to be stored in the database.
// Its `asString()` function wraps the string in single quotes ('')
// 	to be written in SQL.
type StringCell struct {
	Value string
}

func (cell StringCell) asString() string {
	return "'" + cell.Value + "'"
}

// BoolCell represents a bool value to be stored in the database.
type BoolCell struct {
	Value bool
}

func (cell BoolCell) asString() string {
	return strconv.FormatBool(cell.Value)
}

// TimeCell represents a time.Time value to be stored in the database.
//
// Its `asString()` function uses the fmt package to format the time.
type TimeCell struct {
	Value time.Time
}

func (cell TimeCell) asString() string {
	return fmt.Sprintf("%v", cell.Value)
}

// Float32Cell presents a float32 value to be stored in the database.
//
// Its `asString()` function represents the float32 as a string to four
// 	decimal places (%.4f).
type Float32Cell struct {
	Value float32
}

func (cell Float32Cell) asString() string {
	return fmt.Sprintf("%.4f", cell.Value)
}

// Float64Cell presents a float64 value to be stored in the database.
//
// Its `asString()` function represents the float64 as a string to four
// 	decimal places (%.4f).
type Float64Cell struct {
	Value float32
}

func (cell Float64Cell) asString() string {
	return fmt.Sprintf("%.4f", cell.Value)
}

func check(err error, logger *log.Logger) {
	if err != nil {
		logger.Fatal(err)
	}
}
