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

	entry := NewEntry(Fields{
		"timestamp": "2015-01-01T01:01:01Z",
	})
	assert.Nil(t, filter.Filter(entry))

	entry = NewEntry(Fields{
		"timestamp": "2015-02-02T02:02:02Z",
	})
	assert.Equal(t, filter.Filter(entry), entry)
}

func TestDatetimeFilterStart(t *testing.T) {
	// filter contains lower border only
	filter := &Datetime{
		Field:  "timestamp",
		Format: time.RFC3339,
		Start:  time.Date(2015, time.February, 2, 2, 2, 2, 0, time.UTC),
	}
	assert.Implements(t, (*Filter)(nil), filter)

	entry := NewEntry(Fields{
		"timestamp": "2015-01-01T01:01:01Z",
	})
	assert.Nil(t, filter.Filter(entry))

	entry = NewEntry(Fields{
		"timestamp": "2015-02-02T02:02:02Z",
	})
	assert.Equal(t, filter.Filter(entry), entry)
}

func TestDatetimeFilterEnd(t *testing.T) {
	// filter contains upper border only
	filter := &Datetime{
		Field:  "timestamp",
		Format: time.RFC3339,
		End:    time.Date(2015, time.May, 5, 5, 5, 5, 0, time.UTC),
	}
	assert.Implements(t, (*Filter)(nil), filter)

	entry := NewEntry(Fields{
		"timestamp": "2015-01-01T01:01:01Z",
	})
	assert.Equal(t, filter.Filter(entry), entry)

	entry = NewEntry(Fields{
		"timestamp": "2015-05-05T05:05:05Z",
	})
	assert.Nil(t, filter.Filter(entry))
}

func TestDatetimeReduce(t *testing.T) {
	filter := &Datetime{
		Field:  "timestamp",
		Format: time.RFC3339,
		Start:  time.Date(2015, time.February, 2, 2, 2, 2, 0, time.UTC),
		End:    time.Date(2015, time.May, 5, 5, 5, 5, 0, time.UTC),
	}
	assert.Implements(t, (*Reducer)(nil), filter)

	// Prepare input channel
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
	filter.Reduce(input, output)

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

func TestChainFilterWithRedicer(t *testing.T) {
	// Prepare input channel
	input := make(chan *Entry, 5)
	input <- NewEntry(Fields{
		"timestamp": "2015-01-01T01:01:01Z",
		"foo":       "12",
		"bar":       "34",
		"baz":       "56",
	})
	input <- NewEntry(Fields{
		"timestamp": "2015-02-02T02:02:02Z",
		"foo":       "34",
		"bar":       "56",
		"baz":       "78",
	})
	input <- NewEntry(Fields{
		"timestamp": "2015-04-04T04:04:04Z",
		"foo":       "78",
		"bar":       "90",
		"baz":       "12",
	})
	input <- NewEntry(Fields{
		"timestamp": "2015-05-05T05:05:05Z",
		"foo":       "90",
		"bar":       "34",
		"baz":       "56",
	})
	close(input)

	filter := &Datetime{
		Field:  "timestamp",
		Format: time.RFC3339,
		Start:  time.Date(2015, time.February, 2, 2, 2, 2, 0, time.UTC),
		End:    time.Date(2015, time.May, 5, 5, 5, 5, 0, time.UTC),
	}
	chain := NewChain(filter, &Avg{[]string{"foo", "bar"}}, &Count{})

	output := make(chan *Entry, 5) // Make it buffered to avoid deadlock
	chain.Reduce(input, output)

	result, ok := <-output
	assert.True(t, ok)

	value, err := result.FloatField("foo")
	assert.NoError(t, err)
	assert.Equal(t, value, (34.0+78)/2.0)

	value, err = result.FloatField("bar")
	assert.NoError(t, err)
	assert.Equal(t, value, (56.0+90)/2.0)

	count, err := result.Field("count")
	assert.NoError(t, err)
	assert.Equal(t, count, "2")

	_, err = result.Field("buz")
	assert.Error(t, err)
}
