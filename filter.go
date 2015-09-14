package gonx

import "time"

// Filter interface for Entries channel limiting.
//
// Filter method should accept input channel of Entries, deicide to keep
// or not to keep given entry. Write entry to the output channel if positive.
type Filter interface {
	Filter(input chan *Entry, output chan *Entry)
}

// Implements Filter interface to filter Entries with timestamp fields within
// the specified datetime interval.
type Datetime struct {
	Field  string
	Format string
	Start  time.Time
	End    time.Time
}

// Check field value to be in desired datetime range.
func (i *Datetime) Filter(input chan *Entry, output chan *Entry) {
	for entry := range input {
		val, err := entry.Field(i.Field)
		if err != nil {
			// TODO handle error
			continue
		}
		t, err := time.Parse(i.Format, val)
		if err != nil {
			// TODO handle error
			continue
		}
		if i.withinBounds(t) {
			output <- entry
		}
	}
	close(output)
}

func (i *Datetime) withinBounds(t time.Time) bool {
	if t.Equal(i.Start) {
		return true
	}
	if t.After(i.Start) && t.Before(i.End) {
		return true
	}
	return false
}
