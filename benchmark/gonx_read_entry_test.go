package benchmark_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/satyrius/gonx"
)

func BenchmarkGonxReadEntry(b *testing.B) {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu + 1)
	start := time.Now()
	count := 0
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parent := filepath.Dir(wd)
	filePath := filepath.Join(parent, `example`, `access.log`)
	for i := 0; i < b.N; i++ {
		count = count + gonxReadEntry(filePath)
	}
	fmt.Printf("%v lines readed, it takes %v\n", count, time.Since(start))
}

func gonxReadEntry(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	var count int
	format := `$remote_addr [$time_local] "$request"`
	reader := gonx.NewReader(file, format)
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		}
		count++
	}
	return count
}
