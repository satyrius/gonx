package gonx

import (
	"fmt"
	"strconv"
	"strings"
)

// Fields is a shortcut for the map of strings
type Fields map[string]string

// Entry is a parsed log record. Use Get method to retrieve a value by name instead of
// threating this as a map, because inner representation is in design.
type Entry struct {
	fields Fields
}

// NewEmptyEntry creates an empty Entry to be filled later
func NewEmptyEntry() *Entry {
	return &Entry{make(Fields)}
}

// NewEntry creates an Entry with fiven fields
func NewEntry(fields Fields) *Entry {
	return &Entry{fields}
}

// Fields returns all fields of an entry
func (entry *Entry) Fields() Fields {
	return entry.fields
}

// Field returns an entry field value by name or empty string and error if it
// does not exist.
func (entry *Entry) Field(name string) (value string, err error) {
	value, ok := entry.fields[name]
	if !ok {
		err = fmt.Errorf("field '%v' does not found in record %+v", name, *entry)
	}
	return
}

// FloatField returns an entry field value as float64. Return nil if field does not exist
// and conversion error if cannot cast a type.
func (entry *Entry) FloatField(name string) (value float64, err error) {
	tmp, err := entry.Field(name)
	if err == nil {
		value, err = strconv.ParseFloat(tmp, 64)
	}
	return
}

// SetField sets the value of a field
func (entry *Entry) SetField(name string, value string) {
	entry.fields[name] = value
}

// SetFloatField is a Float field value setter. It accepts float64, but still store it as a
// string in the same fields map. The precision is 2, its enough for log
// parsing task
func (entry *Entry) SetFloatField(name string, value float64) {
	entry.SetField(name, strconv.FormatFloat(value, 'f', 2, 64))
}

// SetUintField is a Integer field value setter. It accepts float64, but still store it as a
// string in the same fields map.
func (entry *Entry) SetUintField(name string, value uint64) {
	entry.SetField(name, strconv.FormatUint(uint64(value), 10))
}

// Merge two entries by updating values for master entry with given.
func (entry *Entry) Merge(merge *Entry) {
	for name, value := range merge.fields {
		entry.SetField(name, value)
	}
}

// FieldsHash returns a hash of all fields
func (entry *Entry) FieldsHash(fields []string) string {
	var key []string
	for _, name := range fields {
		value, err := entry.Field(name)
		if err != nil {
			value = "NULL"
		}
		key = append(key, fmt.Sprintf("'%v'=%v", name, value))
	}
	return strings.Join(key, ";")
}

// Partial returns a partial field entry with the specified fields
func (entry *Entry) Partial(fields []string) *Entry {
	partial := NewEmptyEntry()
	for _, name := range fields {
		value, _ := entry.Field(name)
		partial.SetField(name, value)
	}
	return partial
}
