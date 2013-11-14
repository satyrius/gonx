package gonx

import (
	"bufio"
	"fmt"
	"io"
	//"os"
	"regexp"
	"sync"
)

type Reader struct {
	format  string
	re      *regexp.Regexp
	files   chan io.Reader
	entries chan Entry
	wg      sync.WaitGroup
}

func NewMap(files chan io.Reader, format string) *Reader {
	r := &Reader{
		format:  format,
		files:   files,
		entries: make(chan Entry, 10),
	}
	// Get regexp to trigger its creation now to avoid data race
	// in future (this is not the method I want to lock Mutex
	// every time.
	// TODO Unbind this method from Reader struct
	r.GetFormatRegexp()

	for file := range files {
		r.wg.Add(1)
		go r.readFile(file)
	}

	go func() {
		r.wg.Wait()
		close(r.entries)
	}()

	return r
}

func (r *Reader) handleError(err error) {
	//fmt.Fprintln(os.Stderr, err)
}

func NewReader(logFile io.Reader, format string) *Reader {
	files := make(chan io.Reader, 1)
	files <- logFile
	close(files)
	return NewMap(files, format)
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

func (r *Reader) readFile(file io.Reader) {
	// Iterate over log file lines and spawn new mapper goroutine
	// to parse it into given format
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r.wg.Add(1)
		go func(line string) {
			entry, err := r.parseRecord(line)
			if err != nil {
				r.handleError(err)
			}
			r.entries <- entry
			r.wg.Done()
		}(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		r.handleError(err)
	}
	r.wg.Done()
}

// Read next line from entries channel, and return parsed record. If channel is closed
// then method returns io.EOF error
func (r *Reader) Read() (entry Entry, err error) {
	entry, ok := <-r.entries
	if !ok {
		err = io.EOF
	}
	return
}

func (r *Reader) parseRecord(line string) (record Entry, err error) {
	// Parse line to fill map record. Return error if a line does not match given format
	re := r.GetFormatRegexp()
	fields := re.FindStringSubmatch(line)
	if fields == nil {
		err = fmt.Errorf("Access log line '%v' does not match given format '%v'", line, re)
		return
	}

	// Iterate over subexp foung and fill the map record
	record = make(Entry)
	for i, name := range re.SubexpNames() {
		if i == 0 {
			continue
		}
		record[name] = fields[i]
	}
	return
}
