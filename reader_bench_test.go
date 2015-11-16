package gonx

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func BenchmarkScannerReader(b *testing.B) {
	file := strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
	for i := 0; i < b.N; i++ {
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		_ = scanner.Text()
		if err := scanner.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReaderReader(b *testing.B) {
	file := strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
	for i := 0; i < b.N; i++ {
		reader := bufio.NewReader(file)
		_, err := readLine(reader)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
	}
}
