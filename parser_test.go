package gonx

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type ParserTestSuite struct {
	suite.Suite
	format string
	parser *Parser
}

func (suite *ParserTestSuite) SetupTest() {
	suite.format = "$remote_addr [$time_local] \"$request\" $status"
	suite.parser = NewParser(suite.format)
}

func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}

func (suite *ParserTestSuite) TestFormatSaved() {
	assert.Equal(suite.T(), suite.parser.format, suite.format)
}

func (suite *ParserTestSuite) TestRegexp() {
	assert.Equal(suite.T(),
		suite.parser.regexp.String(),
		`^(?P<remote_addr>[^ ]*) \[(?P<time_local>[^]]*)\] "(?P<request>[^"]*)" (?P<status>[^ ]*)$`)
}

func (suite *ParserTestSuite) TestParseString() {
	line := `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1" 200`
	expected := NewEntry(Fields{
		"remote_addr": "89.234.89.123",
		"time_local":  "08/Nov/2013:13:39:18 +0000",
		"request":     "GET /api/foo/bar HTTP/1.1",
		"status":      "200",
		"raw_line":    `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1" 200`,
	})
	entry, err := suite.parser.ParseString(line)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entry, expected)
}

func (suite *ParserTestSuite) TestParseInvalidString() {
	line := `GET /api/foo/bar HTTP/1.1`
	_, err := suite.parser.ParseString(line)
	assert.Error(suite.T(), err)
}

func (suite *ParserTestSuite) TestEmptyValue() {
	line := `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "" 200`
	expected := NewEntry(Fields{
		"remote_addr": "89.234.89.123",
		"time_local":  "08/Nov/2013:13:39:18 +0000",
		"request":     "",
		"status":      "200",
		"raw_line":    `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "" 200`,
	})
	entry, err := suite.parser.ParseString(line)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entry, expected)
}

func TestNginxParser(t *testing.T) {
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
	parser, err := NewNginxParser(conf, "main")
	assert.NoError(t, err)
	assert.Equal(t, parser.format, expected)
}
