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

func init() {
	flag.StringVar(&conf, "conf", "/etc/nginx/nginx.conf", "Nginx config file")
	flag.StringVar(&format, "format", "main", "Nginx log_format name")
}

func main() {
	flag.Parse()

	// Use nginx config file to extract format by the name
	nginxConfig, err := os.Open(conf)
	if err != nil {
		panic(err)
	}
	defer nginxConfig.Close()

	// Read from STDIN and use log_format to parse log records
	reader := gonx.NewNginxReader(os.Stdin, nginxConfig, format)
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
