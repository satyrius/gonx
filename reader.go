package gonx

import (
	"io"
)

type Reader struct {
	entryMap *Map
}

func NewReader(logFile io.Reader, format string) *Reader {
	m := NewMap(oneFileChannel(logFile), NewParser(format))
	return &Reader{entryMap: m}
}

func NewNginxReader(logFile io.Reader, nginxConf io.Reader, formatName string) (reader *Reader, err error) {
	parser, err := NewNginxParser(nginxConf, formatName)
	if err != nil {
		return nil, err
	}
	m := NewMap(oneFileChannel(logFile), parser)
	reader = &Reader{entryMap: m}
	return
}

// Read next the map. Return EOF if there is no Entries to read
func (r *Reader) Read() (Entry, error) {
	// TODO return Entry reference instead of instance
	entry := r.entryMap.GetEntry()
	if entry == nil {
		// Have to return emtry entry for backward capability
		return Entry{}, io.EOF
	}
	return *entry, nil
}
