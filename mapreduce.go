package gonx

import (
	"bufio"
	"io"
)

func (r *Reader) handleError(err error) {
	//fmt.Fprintln(os.Stderr, err)
}

func newMap(files chan io.Reader, parser *Parser) *Reader {
	reader := &Reader{
		parser:  parser,
		files:   files,
		entries: make(chan Entry, 10),
	}

	for file := range files {
		reader.wg.Add(1)
		go reader.readFile(file)
	}

	go func() {
		reader.wg.Wait()
		close(reader.entries)
	}()

	return reader
}

func oneFileChannel(file io.Reader) chan io.Reader {
	ch := make(chan io.Reader, 1)
	ch <- file
	close(ch)
	return ch
}

func (r *Reader) readFile(file io.Reader) {
	// Iterate over log file lines and spawn new mapper goroutine
	// to parse it into given format
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r.wg.Add(1)
		go func(line string) {
			entry, err := r.parser.ParseString(line)
			if err == nil {
				r.entries <- entry
			} else {
				r.handleError(err)
			}
			r.wg.Done()
		}(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		r.handleError(err)
	}
	r.wg.Done()
}
