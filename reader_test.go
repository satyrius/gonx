package gonx

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	Convey("Test Reader", t, func() {
		format := "$remote_addr [$time_local] \"$request\""

		Convey("Test valid file", func() {
			file := strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
			reader := NewReader(file, format)
			So(reader.entries, ShouldBeNil)

			expected := NewEntry(Fields{
				"remote_addr": "89.234.89.123",
				"time_local":  "08/Nov/2013:13:39:18 +0000",
				"request":     "GET /api/foo/bar HTTP/1.1",
			})

			// Read entry from incoming channel
			entry, err := reader.Read()
			So(err, ShouldBeNil)
			So(entry, ShouldResemble, expected)

			// It was only one line, nothing to read
			_, err = reader.Read()
			So(err, ShouldEqual, io.EOF)
		})
	})
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
