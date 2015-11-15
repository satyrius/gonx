package gonx

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func BenchmarkScannerReader(b *testing.B) {
	file, err := os.Open("example/access.log")
	if err != nil {
		b.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_ = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkReaderReader(b *testing.B) {
	file, err := os.Open("example/access.log")
	if err != nil {
		b.Fatal(err)
	}
	reader := bufio.NewReader(file)
	_, err = readLine(reader)
	for err == nil {
		_, err = readLine(reader)
	}
	if err != nil && err != io.EOF {
		b.Fatal(err)
	}
}
