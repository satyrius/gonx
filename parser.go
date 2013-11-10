package gonx

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

func NginxLogFormat(nginxConf io.Reader, formatName string) (format string, err error) {
	scanner := bufio.NewScanner(nginxConf)
	re := regexp.MustCompile(fmt.Sprintf(`^.*log_format\s+%v\s+(.+)\s*$`, formatName))
	found := false
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
	return
}

func FormatRegexp(format string) *regexp.Regexp {
	format = regexp.MustCompile(`\\\$([a-z_]+)(\\?(.))`).ReplaceAllString(
		regexp.QuoteMeta(format), "(?P<$1>[^$3]+)$2")
	return regexp.MustCompile(fmt.Sprintf("^%v$", format))
}

func GetRecord(line string, re *regexp.Regexp) (record map[string]string, err error) {
	// Parse line to fill map record. Return error if a line does not match given format
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
