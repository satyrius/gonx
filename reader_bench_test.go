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
	file, err := os.Open("example/access.log")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		reader := bufio.NewReader(file)
		_, err = readLine(reader)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
	}
}
