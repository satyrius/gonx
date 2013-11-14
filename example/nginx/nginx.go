package main

import (
	"flag"
	"fmt"
	"github.com/satyrius/gonx"
	"io"
	"os"
	"strings"
)

var conf string
var format string
var logFile string

func init() {
	flag.StringVar(&conf, "conf", "dummy", "Nginx config file (e.g. /etc/nginx/nginx.conf)")
	flag.StringVar(&format, "format", "main", "Nginx log_format name")
	flag.StringVar(&logFile, "log", "dummy", "Log file name to read. Read from STDIN if file name is '-'")
}

func main() {
	flag.Parse()

	// Read given file or from STDIN
	var file io.Reader
	if logFile == "dummy" {
		file = strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
	} else if logFile == "-" {
		file = os.Stdin
	} else {
		file, err := os.Open(logFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	// Use nginx config file to extract format by the name
	var nginxConfig io.Reader
	if conf == "dummy" {
		nginxConfig = strings.NewReader(`
            http {
                log_format   main  '$remote_addr [$time_local] "$request"';
            }
        `)
	} else {
		nginxConfig, err := os.Open(conf)
		if err != nil {
			panic(err)
		}
		defer nginxConfig.Close()
	}

	// Read from STDIN and use log_format to parse log records
	reader, err := gonx.NewNginxReader(file, nginxConfig, format)
	if err != nil {
		panic(err)
	}
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		// Process the record... e.g.
		fmt.Printf("%+v\n", rec)
	}
}
