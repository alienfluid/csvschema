package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Global variables set by the flag parsing logic
var _minLines int
var _delimiter string
var _noheader bool

// Given a string, tries to convert it to various types to figure out
// the actual data of the underlying data
func determineType(s string) string {
	// Data is missing (NULL case)
	if s == "" {
		return "unknown"
	}

	_, err := strconv.ParseInt(s, 10, 32)
	if err == nil {
		return "int32"
	}

	_, err = strconv.ParseInt(s, 10, 64)
	if err == nil {
		return "int64"
	}

	_, err = strconv.ParseFloat(s, 32)
	if err == nil {
		return "float32"
	}

	_, err = strconv.ParseFloat(s, 64)
	if err == nil {
		return "float64"
	}

	// Define various format dates and times can be specified in
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05.000000",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.Kitchen,
		time.Stamp,
		time.StampMicro,
		time.StampMilli,
		time.StampMilli,
		time.StampNano}

	for _, f := range formats {
		_, err = time.Parse(f, s)
		if err == nil {
			return "timestamp"
		}
	}

	// If we can't convert it to anything, it's probably string
	return "string"
}

// Main entry point to the tool
func main() {
	// Define and parse the command line flags
	flag.IntVar(&_minLines, "lines", 1000, "Minimum number of lines to sample (unless file is smaller)")
	flag.StringVar(&_delimiter, "delimiter", ",", "Column delimiter")
	flag.BoolVar(&_noheader, "noheader", false, "Don't consider the first line to be the header")
	flag.Parse()

	filename := flag.Args()
	if len(filename) == 0 {
		fmt.Printf("Please provide the path to the file to be parsed.")
		return
	}

	file, err := os.Open(filename[0])
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = []rune(_delimiter)[0]

	fmt.Printf("Sampling %v records from the file\n", _minLines)

	var nlines int
	var totlines int
	var header []string

	// Reservoir sampling algorithm implementation
	reservoir := make([][]string, 0, _minLines)
	for {
		record, err := reader.Read()
		if err == io.EOF || record == nil {
			break
		}

		if !_noheader && nlines == 0 && len(header) == 0 {
			header = record
			continue
		}

		if nlines < _minLines {
			// Reservoir is not yet full, just add it
			reservoir = append(reservoir, record)
			nlines++
		} else {
			// Reservoir is full, sample
			index := rand.Intn(nlines)
			if index < _minLines {
				reservoir[index] = record
			}
		}

		totlines++
	}

	fmt.Printf("Sampled %v records from the file (out of total %v)\n", nlines, totlines)

	// Iterate over the columns and determine the type of the data within. If all the values
	// in a given column are of the same type, we can say that the column type is of that type.
	// Unknown types are ignored since they are considered to be missing values.
	ncols := len(reservoir[0])
	for i := 0; i < ncols; i++ {
		types := make([]string, 0, len(reservoir))
		for _, r := range reservoir {
			types = append(types, determineType(r[i]))
		}
		lastType := types[0]
		same := true
		for idx, t := range types {
			if idx == 0 {
				continue
			}
			if lastType != t && lastType != "unknown" && t != "unknown" {
				same = false
				break
			}
			lastType = t
		}

		if !same {
			lastType = "unknown"
		}

		if !_noheader {
			fmt.Printf("Column %v: %v\n", header[i], lastType)
		} else {
			fmt.Printf("Column %v: %v\n", i, lastType)
		}
	}
}
