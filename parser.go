package gonx

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// StringParser is the interface that wraps the ParseString method.
type StringParser interface {
	ParseString(line string) (entry *Entry, err error)
}

// Parser is a log record parser. Use specific constructors to initialize it.
type Parser struct {
	format string
	regexp *regexp.Regexp
}

func getSpecialNginxRegexes() map[string]string {
	return map[string]string{
		"http_x_forwarded_for": `[^, ]*(?:, ?[^, ]+)*`}
}

// NewParser returns a new Parser, use given log format to create its internal
// strings parsing regexp.
func NewParser(format string) *Parser {
	// First split up multiple concatenated fields with placeholder
	placeholder := " _PLACEHOLDER___ "
	preparedFormat := format
	concatenatedRe := regexp.MustCompile(`[A-Za-z0-9_]\$[A-Za-z0-9_]`)
	for concatenatedRe.MatchString(preparedFormat) {
		preparedFormat = regexp.MustCompile(`([A-Za-z0-9_])\$([A-Za-z0-9_]+)(\\?([^$\\A-Za-z0-9_]))`).ReplaceAllString(
			preparedFormat, fmt.Sprintf("${1}${3}%s$$${2}${3}", placeholder),
		)
	}

	formatRegex := regexp.MustCompile(`([^$ ]*)\$([a-z_]+)([^$ ]*)([ ]?)`)
	specialNginxRegexes := getSpecialNginxRegexes()
	fields := formatRegex.FindAllStringSubmatch(preparedFormat+" ", -1)
	re := formatRegex.ReplaceAllString(preparedFormat+" ", "$2$4")
	for _, field := range fields {
		terminateChar := field[3]
		if len([]rune(terminateChar)) == 0 {
			terminateChar = field[4]
		}
		if specialRegex, found := specialNginxRegexes[field[2]]; found {
			re = strings.Replace(re, field[2]+field[4], regexp.QuoteMeta(field[1])+"(?P<"+field[2]+">"+specialRegex+")"+regexp.QuoteMeta(field[3])+field[4], 1)
		} else {
			re = strings.Replace(re, field[2]+field[4], regexp.QuoteMeta(field[1])+"(?P<"+field[2]+">[^"+terminateChar+"]*)"+regexp.QuoteMeta(field[3]+field[4]), 1)
		}
	}

	// Finally remove placeholder
	re = regexp.MustCompile(fmt.Sprintf(".%s", placeholder)).ReplaceAllString(re, "")
	return &Parser{format, regexp.MustCompile(fmt.Sprintf("^%v", strings.Trim(re, " ")))}
}

// ParseString parses a log file line using internal format regexp. If a line
// does not match the given format an error will be returned.
func (parser *Parser) ParseString(line string) (entry *Entry, err error) {
	re := parser.regexp
	fields := re.FindStringSubmatch(line)
	if fields == nil {
		err = fmt.Errorf("access log line '%v' does not match given format '%v'", line, re)
		return
	}

	// Iterate over subexp foung and fill the map record
	entry = NewEmptyEntry()
	for i, name := range re.SubexpNames() {
		if i == 0 {
			continue
		}
		entry.SetField(name, fields[i])
	}
	return
}

// NewNginxParser parses the nginx conf file to find log_format with the given
// name and returns a parser for this format. It returns an error if cannot find
// the given log format.
func NewNginxParser(conf io.Reader, name string) (parser *Parser, err error) {
	scanner := bufio.NewScanner(conf)
	re := regexp.MustCompile(fmt.Sprintf(`^\s*log_format\s+%v\s+(.+)\s*$`, name))
	found := false
	var format string
	for scanner.Scan() {
		var line string
		if !found {
			// Find a log_format definition
			line = scanner.Text()
			formatDef := re.FindStringSubmatch(line)
			if formatDef == nil {
				continue
			}
			found = true
			line = formatDef[1]
		} else {
			line = scanner.Text()
		}
		// Look for a definition end
		re = regexp.MustCompile(`^\s*(.*?)\s*(;|$)`)
		lineSplit := re.FindStringSubmatch(line)
		if l := len(lineSplit[1]); l > 2 {
			format += lineSplit[1][1 : l-1]
		}
		if lineSplit[2] == ";" {
			break
		}
	}
	if !found {
		err = fmt.Errorf("`log_format %v` not found in given config", name)
	} else {
		err = scanner.Err()
	}
	parser = NewParser(format)
	return
}
