package gonx

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDatetimeFilter(t *testing.T) {
	filter := &Datetime{
		Field:  "timestamp",
		Format: time.RFC3339,
		Start:  time.Date(2015, time.February, 2, 2, 2, 2, 0, time.UTC),
		End:    time.Date(2015, time.May, 5, 5, 5, 5, 0, time.UTC),
	}
	assert.Implements(t, (*Filter)(nil), filter)

	// Prepare import channel
	input := make(chan *Entry, 5)
	input <- NewEntry(Fields{
		"timestamp": "2015-01-01T01:01:01Z",
		"foo":       "12",
	})
	input <- NewEntry(Fields{
		"timestamp": "2015-02-02T02:02:02Z",
		"foo":       "34",
	})
	input <- NewEntry(Fields{
		"timestamp": "2015-03-03T03:03:03Z",
		"foo":       "56",
	})
	input <- NewEntry(Fields{
		"timestamp": "2015-04-04T04:04:04Z",
		"foo":       "78",
	})
	input <- NewEntry(Fields{
		"timestamp": "2015-05-05T05:05:05Z",
		"foo":       "90",
	})
	close(input)

	output := make(chan *Entry, 5) // Make it buffered to avoid deadlock
	filter.Filter(input, output)

	expected := []string{
		"'timestamp'=2015-02-02T02:02:02Z;'foo'=34",
		"'timestamp'=2015-03-03T03:03:03Z;'foo'=56",
		"'timestamp'=2015-04-04T04:04:04Z;'foo'=78",
	}
	results := []string{}

	for result := range output {
		results = append(
			results,
			result.FieldsHash([]string{"timestamp", "foo"}),
		)
	}
	assert.Equal(t, results, expected)
}
