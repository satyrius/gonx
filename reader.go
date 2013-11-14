package gonx

import (
	"io"
	"sync"
)

type Reader struct {
	parser  *Parser
	files   chan io.Reader
	entries chan Entry
	wg      sync.WaitGroup
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

// Read next line from entries channel, and return parsed record. If channel is closed
// then method returns io.EOF error
func (r *Reader) Read() (entry Entry, err error) {
	entry, ok := <-r.entries
	if !ok {
		err = io.EOF
	}
	return
}
