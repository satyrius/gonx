package gonx

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	format := "$remote_addr [$time_local] \"$request\""
	file := strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
	reader := NewReader(file, format)
	assert.Nil(t, reader.entries)

	expected := NewEntry(Fields{
		"remote_addr": "89.234.89.123",
		"time_local":  "08/Nov/2013:13:39:18 +0000",
		"request":     "GET /api/foo/bar HTTP/1.1",
		"raw_line":    `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`,
	})

	// Read entry from incoming channel
	entry, err := reader.Read()
	assert.NoError(t, err)
	assert.Equal(t, entry, expected)

	// It was only one line, nothing to read
	_, err = reader.Read()
	assert.Equal(t, err, io.EOF)
}

func TestInvalidLineFormat(t *testing.T) {
	t.Skip("Read method does not return errors anymore, because of asynchronios algorithm")
	format := "$remote_addr [$time_local] \"$request\""
	file := strings.NewReader(`89.234.89.123 - - [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
	reader := NewReader(file, format)

	// Invalid entries do not go to the entries channel, so nothing to read
	_, err := reader.Read()
	assert.Equal(t, err, io.EOF)

	// TODO test Reader internal error handling
}
