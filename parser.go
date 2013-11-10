package gonx

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type Reader struct {
	format  string
	re      *regexp.Regexp
	scanner *bufio.Scanner
}

func NewReader(logFile io.Reader, format string) *Reader {
	return &Reader{
		format:  format,
		scanner: bufio.NewScanner(logFile),
	}
}

func NewNginxReader(logFile io.Reader, nginxConf io.Reader, formatName string) (reader *Reader, err error) {
	scanner := bufio.NewScanner(nginxConf)
	re := regexp.MustCompile(fmt.Sprintf(`^.*log_format\s+%v\s+(.+)\s*$`, formatName))
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
		err = fmt.Errorf("`log_format %v` not found in given config", formatName)
	} else {
		err = scanner.Err()
	}
	reader = NewReader(logFile, format)
	return
}

func (r *Reader) GetFormat() string {
	return r.format
}

func (r *Reader) GetFormatRegexp() *regexp.Regexp {
	if r.re != nil {
		return r.re
	}
	format := regexp.MustCompile(`\\\$([a-z_]+)(\\?(.))`).ReplaceAllString(
		regexp.QuoteMeta(r.format), "(?P<$1>[^$3]+)$2")
	r.re = regexp.MustCompile(fmt.Sprintf("^%v$", format))
	return r.re
}

// Read next line from log file, and return parsed record. If all lines read
// method return ni, io.EOF
func (r *Reader) Read() (record map[string]string, err error) {
	if r.scanner.Scan() {
		record, err = r.parseRecord(r.scanner.Text())
	} else {
		err = r.scanner.Err()
		if err == nil {
			err = io.EOF
		}
	}
	return
}

func (r *Reader) parseRecord(line string) (record map[string]string, err error) {
	// Parse line to fill map record. Return error if a line does not match given format
	re := r.GetFormatRegexp()
	fields := re.FindStringSubmatch(line)
	if fields == nil {
		err = fmt.Errorf("Access log line '%v' does not match given format '%v'", line, re)
		return
	}

	// Iterate over subexp foung and fill the map record
	record = make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i == 0 {
			continue
		}
		record[name] = fields[i]
	}
	return
}
