package gonx

import (
	"bufio"
	"io"
	//"os"
	"sync"
)

type Reader struct {
	parser  *Parser
	files   chan io.Reader
	entries chan Entry
	wg      sync.WaitGroup
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

func (r *Reader) handleError(err error) {
	//fmt.Fprintln(os.Stderr, err)
}

func oneFileChannel(file io.Reader) chan io.Reader {
	ch := make(chan io.Reader, 1)
	ch <- file
	close(ch)
	return ch
}

func NewReader(logFile io.Reader, format string) *Reader {
	return newMap(oneFileChannel(logFile), NewParser(format))
}

func NewNginxReader(logFile io.Reader, nginxConf io.Reader, formatName string) (reader *Reader, err error) {
	parser, err := NewNginxParser(nginxConf, formatName)
	if err != nil {
		return nil, err
	}
	reader = newMap(oneFileChannel(logFile), parser)
	return
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

// Read next line from entries channel, and return parsed record. If channel is closed
// then method returns io.EOF error
func (r *Reader) Read() (entry Entry, err error) {
	entry, ok := <-r.entries
	if !ok {
		err = io.EOF
	}
	return
}
