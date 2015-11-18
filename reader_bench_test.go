package gonx

import (
	"bufio"
	"bytes"
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

func BenchmarkReaderReaderAppend(b *testing.B) {
	file := strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
	for i := 0; i < b.N; i++ {
		reader := bufio.NewReader(file)
		_, err := readLineAppend(reader)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
	}
}

func BenchmarkReaderReaderBuffer(b *testing.B) {
	file := strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
	for i := 0; i < b.N; i++ {
		reader := bufio.NewReader(file)
		_, err := readLineBuffer(reader)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
	}
}

func readLineAppend(reader *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = reader.ReadLine()
		if err == nil {
			ln = append(ln, line...)
		}
	}
	return string(ln), err
}

func readLineBuffer(reader *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line     []byte
		buffer   bytes.Buffer
	)
	for isPrefix && err == nil {
		line, isPrefix, err = reader.ReadLine()
		if err == nil {
			_, err = buffer.Write(line)
		}
	}
	return buffer.String(), err
}
