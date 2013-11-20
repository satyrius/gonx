package gonx

import (
	"io"
)

type Reader struct {
	file    io.Reader
	parser  *Parser
	entries chan Entry
}

func NewReader(logFile io.Reader, format string) *Reader {
	return &Reader{
		file:   logFile,
		parser: NewParser(format),
	}
}

func NewNginxReader(logFile io.Reader, nginxConf io.Reader, formatName string) (reader *Reader, err error) {
	parser, err := NewNginxParser(nginxConf, formatName)
	if err != nil {
		return nil, err
	}
	reader = &Reader{
		file:   logFile,
		parser: parser,
	}
	return
}

// Read next the map. Return EOF if there is no Entries to read
func (r *Reader) Read() (entry Entry, err error) {
	if r.entries == nil {
		r.entries = make(chan Entry, 10)
		go func() {
			EntryMap(r.file, r.parser, r.entries)
			close(r.entries)
		}()
	}
	// TODO return Entry reference instead of instance
	entry, ok := <-r.entries
	if !ok {
		err = io.EOF
	}
	return
}
