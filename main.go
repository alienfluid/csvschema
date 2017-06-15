package main

import "encoding/csv"
import "flag"
import "fmt"
import "os"

var _minLines int
var _delimiter string

func determineType (s string) string {
	if s == "" {
		return "unknown"
	} else {
		return "string"
	}
}

func main() {
	flag.IntVar(&_minLines, "lines", 1000, "Minimum number of lines to sample (unless file is smaller)")
	flag.StringVar(&_delimiter, "delimiter", ",", "Column delimiter")
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

	var nlines int = 0
	var lines []([]string)
	for {
		record, err := reader.Read()
		if err != nil {
			continue
		}

		if record == nil {
			break
		}

		lines = append(lines, record)
		nlines += 1

		if nlines == _minLines {
			break
		}
	}

}
