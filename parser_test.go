package gonx

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	Convey("Test Parser", t, func() {
		Convey("Parse custom format", func() {
			format := "$remote_addr [$time_local] \"$request\" $status"
			parser := NewParser(format)

			Convey("Ensure parser format is ok", func() {
				So(parser.format, ShouldEqual, format)
			})

			Convey("Test format to regexp translation", func() {
				So(parser.regexp.String(), ShouldEqual,
					`^(?P<remote_addr>[^ ]*) \[(?P<time_local>[^]]*)\] "(?P<request>[^"]*)" (?P<status>[^ ]*)`)
			})

			Convey("ParseString", func() {
				line := `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1" 200`
				expected := NewEntry(Fields{
					"remote_addr": "89.234.89.123",
					"time_local":  "08/Nov/2013:13:39:18 +0000",
					"request":     "GET /api/foo/bar HTTP/1.1",
					"status":      "200",
				})
				entry, err := parser.ParseString(line)
				So(err, ShouldBeNil)
				So(entry, ShouldResemble, expected)
			})

			Convey("Handle empty values", func() {
				line := `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "" 200`
				expected := NewEntry(Fields{
					"remote_addr": "89.234.89.123",
					"time_local":  "08/Nov/2013:13:39:18 +0000",
					"request":     "",
					"status":      "200",
				})
				entry, err := parser.ParseString(line)
				So(err, ShouldBeNil)
				So(entry, ShouldResemble, expected)
			})

			Convey("Parse invalid string", func() {
				line := `GET /api/foo/bar HTTP/1.1`
				_, err := parser.ParseString(line)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Test multiple fields concatenated toggle", func() {
			format := `$remote_addr [$time_local] "$host$request_uri" $status`
			parser := NewParser(format)

			Convey("Ensure two fields concatenated toggle regexp ok", func() {
				So(parser.format, ShouldEqual, format)
				So(
					parser.regexp.String(),
					ShouldEqual,
					`^(?P<remote_addr>[^ ]*) \[(?P<time_local>[^]]*)\] "(?P<host>[^"]*)(?P<request_uri>[^"]*)" (?P<status>[^ ]*)`,
				)
			})

			format = `$remote_addr [$time_local] "$host$request_uri$demo" $status`
			parser = NewParser(format)
			Convey("Ensure three fields concatenated toggle regexp ok", func() {
				So(parser.format, ShouldEqual, format)
				So(
					parser.regexp.String(),
					ShouldEqual,
					`^(?P<remote_addr>[^ ]*) \[(?P<time_local>[^]]*)\] "(?P<host>[^"]*)(?P<request_uri>[^"]*)(?P<demo>[^"]*)" (?P<status>[^ ]*)`,
				)
			})
		})

		Convey("Nginx format parser", func() {
			expected := "$remote_addr - $remote_user [$time_local] \"$request\" $status \"$http_referer\" \"$http_user_agent\""

			testCases := []string{
				`http {
					include      conf/mime.types;
					log_format   minimal  '$remote_addr [$time_local] "$request"';
					log_format   main     '$remote_addr - $remote_user [$time_local] '
										'"$request" $status '
										'"$http_referer" "$http_user_agent"';
					log_format   download '$remote_addr - $remote_user [$time_local] '
										'"$request" $status $bytes_sent '
										'"$http_referer" "$http_user_agent" '
										'"$http_range" "$sent_http_content_range"';
				}`,

				// Note that the line containing the log_format directive ("log_format main") has NO trailing whitespace.
				// This is perfectly valid nginx config.
				`http {
					include      conf/mime.types;
					log_format   minimal  '$remote_addr [$time_local] "$request"';
					log_format   main 
'$remote_addr - $remote_user [$time_local] '
'"$request" $status '
'"$http_referer" "$http_user_agent"';
				}`,

				// Similar to the above, but contains both trailing whitespace and empty lines in the log_format def
				// This is perfectly valid nginx config.
				`
				http {
					include      conf/mime.types;
					log_format   minimal  '$remote_addr [$time_local] "$request"';
					log_format   main   

'$remote_addr - $remote_user [$time_local] '

'"$request" $status '

'"$http_referer" "$http_user_agent"';
				}`,
			}

			for _, conf := range testCases {
				parser, err := NewNginxParser(strings.NewReader(conf), "main")
				So(err, ShouldBeNil)
				So(parser.format, ShouldEqual, expected)
			}
		})
	})
}
