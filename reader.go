package gonx

import (
	"io"
)

type Reader struct {
	entries chan Entry
	file    io.Reader
	parser  *Parser
}

func NewEntryReader(logFile io.Reader, parser *Parser) *Reader {
	return &Reader{
		file:   logFile,
		parser: parser,
	}
}

func NewReader(logFile io.Reader, format string) *Reader {
	return NewEntryReader(logFile, NewParser(format))
}

func NewNginxReader(logFile io.Reader, nginxConf io.Reader, formatName string) (reader *Reader, err error) {
	parser, err := NewNginxParser(nginxConf, formatName)
	if err != nil {
		return nil, err
	}
	reader = NewEntryReader(logFile, parser)
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
