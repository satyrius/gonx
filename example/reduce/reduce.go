package main

import (
	gonx "../.."
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var format string
var logFile string

func init() {
	flag.StringVar(&format, "format", `$remote_addr [$time_local] "$request" $status $request_length $body_bytes_sent $request_time "$t_size" $read_time $gen_time`, "Log format")
	flag.StringVar(&logFile, "log", "dummy", "Log file name to read. Read from STDIN if file name is '-'")
}

func main() {
	flag.Parse()

	// Create a parser based on given format
	parser := gonx.NewParser(format)

	// Read given file or from STDIN
	var file io.Reader
	if logFile == "dummy" {
		file = strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /t/100x100/foo/bar.jpeg HTTP/1.1" 200 1027 2430 0.014 "100x100" 10 1`)
	} else if logFile == "-" {
		file = os.Stdin
	} else {
		file, err := os.Open(logFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	// Create reader and call Read method until EOF
	reducer := gonx.NewChain(&gonx.Avg{[]string{"request_time"}}, &gonx.Count{})
	output := gonx.MapReduce(file, parser, reducer)
	for res := range output {
		// Process the record... e.g.
		fmt.Printf("Parsed entry: %+v\n", res)
	}
}
