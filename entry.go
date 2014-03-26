package gonx

import (
	"fmt"
	"strconv"
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

// Return entry field value as float64. Rutuen nil if field does not exists
// and convertion error if cannot cast a type.
func (entry *Entry) GetFloat(name string) (value float64, err error) {
	tmp, err := entry.Get(name)
	if err == nil {
		value, err = strconv.ParseFloat(tmp, 64)
	}
	return
}
