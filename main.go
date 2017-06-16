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

var _minLines int
var _delimiter string
var _noheader bool

func determineType(s string) string {
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

	formats := []string{time.ANSIC,
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
		_, err = time.Parse(s, f)
		if err == nil {
			return "timestamp"
		}
	}

	return "string"
}

func main() {
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
			reservoir = append(reservoir, record)
			nlines++
		} else {
			index := rand.Intn(nlines)
			if index < _minLines {
				reservoir[index] = record
			}
		}

		totlines++
	}

	fmt.Printf("Sampled %v records from the file (out of total %v)\n", nlines, totlines)

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
