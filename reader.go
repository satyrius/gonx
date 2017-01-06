package gonx

import (
	"io"
)

// Reader is a log file reader. Use specific constructors to create it.
type Reader struct {
	file    io.Reader
	parser  StringParser
	entries chan *Entry
}

// NewReader creates a reader for a custom log format.
func NewReader(logFile io.Reader, format string) *Reader {
	return NewParserReader(logFile, NewParser(format))
}

// NewParserReader creates a reader with the given parser
func NewParserReader(logFile io.Reader, parser StringParser) *Reader {
	return &Reader{
		file:   logFile,
		parser: parser,
	}
}

// NewNginxReader creates a reader for the nginx log format. Nginx config parser will be used
// to get particular format from the conf file.
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

// Read next parsed Entry from the log file. Return EOF if there are no Entries to read.
func (r *Reader) Read() (entry *Entry, err error) {
	if r.entries == nil {
		r.entries = MapReduce(r.file, r.parser, new(ReadAll))
	}
	entry, ok := <-r.entries
	if !ok {
		err = io.EOF
	}
	return
}
