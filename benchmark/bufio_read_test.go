package benchmark_test

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func BenchmarkBufioRead(b *testing.B) {
	start := time.Now()
	count := 0
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parent := filepath.Dir(wd)
	filePath := filepath.Join(parent, `example`, `access.log`)
	for i := 0; i < b.N; i++ {
		count = count + bufioRead(filePath)
	}
	fmt.Printf("%v lines readed, it takes %v\n", count, time.Since(start))
}

func bufioRead(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	var count int
	for scanner.Scan() {
		// A dummy action, jest read line by line
		scanner.Text()
		count++
	}
	return count
}
