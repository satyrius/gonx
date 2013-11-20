package gonx

import (
	"fmt"
)

// Parsed log record. Use Get method to retrieve a value by name instead of
// threating this as a map, because inner representation is in design.
type Entry map[string]string

// Return entry field value by name or empty string and error if it
// does not exist.
func (entry *Entry) Get(name string) (value string, err error) {
	value, ok := (*entry)[name]
	if !ok {
		err = fmt.Errorf("Field '%v' does not found in record %+v", name, *entry)
	}
	return
}
