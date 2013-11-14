package gonx

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestGetFormatRegexp(t *testing.T) {
	format := "$remote_addr [$time_local] \"$request\""
	reader := NewReader(strings.NewReader(""), format)
	assert.Equal(t,
		reader.GetFormatRegexp().String(),
		`^(?P<remote_addr>[^ ]+) \[(?P<time_local>[^]]+)\] "(?P<request>[^"]+)"$`)
}

func TestGetRecord(t *testing.T) {
	format := "$remote_addr [$time_local] \"$request\""
	file := strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
	reader := NewReader(file, format)
	expected := Entry{
		"remote_addr": "89.234.89.123",
		"time_local":  "08/Nov/2013:13:39:18 +0000",
		"request":     "GET /api/foo/bar HTTP/1.1",
	}
	rec, err := reader.Read()
	assert.NoError(t, err)
	assert.Equal(t, rec, expected)

	_, err = reader.Read()
	assert.Equal(t, err, io.EOF)
}

func TestInvalidLineFormat(t *testing.T) {
	format := "$remote_addr [$time_local] \"$request\""
	file := strings.NewReader(`89.234.89.123 - - [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
	reader := NewReader(file, format)
	rec, err := reader.Read()
	assert.Error(t, err)
	assert.Empty(t, rec)
}

func TestReadLogFormatFromFile(t *testing.T) {
	expected := "$remote_addr - $remote_user [$time_local] \"$request\" $status \"$http_referer\" \"$http_user_agent\""
	conf := strings.NewReader(`
        http {
            include      conf/mime.types;
            log_format   minimal  '$remote_addr [$time_local] "$request"';
            log_format   main     '$remote_addr - $remote_user [$time_local] '
                                  '"$request" $status '
                                  '"$http_referer" "$http_user_agent"';
            log_format   download '$remote_addr - $remote_user [$time_local] '
                                  '"$request" $status $bytes_sent '
                                  '"$http_referer" "$http_user_agent" '
                                  '"$http_range" "$sent_http_content_range"';
        }
    `)
	file := strings.NewReader("")
	reader, err := NewNginxReader(file, conf, "main")
	if err != nil {
		t.Error(err)
	}
	if format := reader.GetFormat(); format != expected {
		t.Errorf("Wrong format was read from conf file \n%v\nExpected\n%v", format, expected)
	}
}

func TestEntry(t *testing.T) {
	entry := Entry{"foo": "1"}

	// Get existings field
	val, err := entry.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, val, "1")

	// Get field that does not exist
	val, err = entry.Get("bar")
	assert.Error(t, err)
	assert.Equal(t, val, "")
}
