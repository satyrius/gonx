package gonx

import (
	"bufio"
	"io"
	"sync"
)

func handleError(err error) {
	//fmt.Fprintln(os.Stderr, err)
}

// Iterate over given file and map each it's line into Entry record using parser.
// Results will be written into output Entries channel.
func EntryMap(file io.Reader, parser *Parser, output chan Entry) {
	var wg sync.WaitGroup

	// Iterate over log file lines and spawn new mapper goroutine
	// to parse it into given format
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			entry, err := parser.ParseString(line)
			if err == nil {
				output <- entry
			} else {
				handleError(err)
			}
		}(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		handleError(err)
	}
	// Wait until all files will be read and all lines will be
	// parsed and wrote to the Entries channel
	wg.Wait()
}
