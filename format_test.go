package gonx

import (
	"reflect"
	"strings"
	"testing"
)

func TestFormatRegexp(t *testing.T) {
	format := "$remote_addr [$time_local] \"$request\""
	expected := `^(?P<remote_addr>[^ ]+) \[(?P<time_local>[^]]+)\] "(?P<request>[^"]+)"$`
	if re := FormatRegexp(format); re.String() != expected {
		t.Errorf("Wrong RE '%v'", re)
	}
}

func TestGetRecord(t *testing.T) {
	format := "$remote_addr [$time_local] \"$request\""
	line := `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`
	expected := map[string]string{
		"remote_addr": "89.234.89.123",
		"time_local":  "08/Nov/2013:13:39:18 +0000",
		"request":     "GET /api/foo/bar HTTP/1.1",
	}
	rec, err := GetRecord(line, FormatRegexp(format))
	if err != nil {
		t.Errorf("Got unexpected error %v", err)
	}
	if !reflect.DeepEqual(rec, expected) {
		t.Errorf("Get invalid record %v", rec)
	}
}

func TestInvalidLineFormat(t *testing.T) {
	format := "$remote_addr [$time_local] \"$request\""
	line := `89.234.89.123 - - [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`
	if rec, err := GetRecord(line, FormatRegexp(format)); err == nil {
		t.Errorf("Invalid record error expected, but get the record %+v", rec)
	}
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
	format, err := NginxLogFormat(conf, "main")
	if err != nil {
		t.Error(err)
	}
	if format != expected {
		t.Errorf("Wrong format was read from conf file \n%v\nExpected\n%v", format, expected)
	}
}
