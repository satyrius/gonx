package gonx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadAllReducer(t *testing.T) {
	reducer := new(ReadAll)
	assert.Implements(t, (*Reducer)(nil), reducer)

	// Prepare import channel
	input := make(chan *Entry, 1)
	entry := NewEmptyEntry()
	input <- entry
	close(input)

	output := make(chan *Entry, 1) // Make it buffered to avoid deadlock
	reducer.Reduce(input, output)

	// ReadAll reducer writes input channel to the output
	result, opened := <-output
	assert.True(t, opened)
	assert.Equal(t, result, entry)
}

func TestCountReducer(t *testing.T) {
	reducer := new(Count)
	assert.Implements(t, (*Reducer)(nil), reducer)

	// Prepare import channel
	input := make(chan *Entry, 2)
	input <- NewEmptyEntry()
	input <- NewEmptyEntry()
	close(input)

	output := make(chan *Entry, 1) // Make it buffered to avoid deadlock
	reducer.Reduce(input, output)

	result, opened := <-output
	assert.True(t, opened)
	count, err := result.Field("count")
	assert.NoError(t, err)
	assert.Equal(t, count, "2")
}

func TestSumReducer(t *testing.T) {
	reducer := &Sum{[]string{"foo", "bar"}}
	assert.Implements(t, (*Reducer)(nil), reducer)

	// Prepare import channel
	input := make(chan *Entry, 2)
	input <- NewEntry(Fields{
		"uri": "/asd/fgh",
		"foo": "123",
		"bar": "234",
		"baz": "345",
	})
	input <- NewEntry(Fields{
		"uri": "/zxc/vbn",
		"foo": "456",
		"bar": "567",
		"baz": "678",
	})
	close(input)
	output := make(chan *Entry, 1) // Make it buffered to avoid deadlock
	reducer.Reduce(input, output)

	result, opened := <-output
	assert.True(t, opened)
	value, err := result.FloatField("foo")
	assert.NoError(t, err)
	assert.Equal(t, value, 123.0+456)
	value, err = result.FloatField("bar")
	assert.NoError(t, err)
	assert.Equal(t, value, 234.0+567.0)
	_, err = result.Field("buz")
	assert.Error(t, err)
}

func TestAvgReducer(t *testing.T) {
	reducer := &Avg{[]string{"foo", "bar"}}
	assert.Implements(t, (*Reducer)(nil), reducer)

	// Prepare import channel
	input := make(chan *Entry, 2)
	input <- NewEntry(Fields{
		"uri": "/asd/fgh",
		"foo": "123",
		"bar": "234",
		"baz": "345",
	})
	input <- NewEntry(Fields{
		"uri": "/zxc/vbn",
		"foo": "456",
		"bar": "567",
		"baz": "678",
	})
	close(input)
	output := make(chan *Entry, 1) // Make it buffered to avoid deadlock
	reducer.Reduce(input, output)

	result, opened := <-output
	assert.True(t, opened)
	value, err := result.FloatField("foo")
	assert.NoError(t, err)
	assert.Equal(t, value, (123.0+456)/2.0)
	value, err = result.FloatField("bar")
	assert.NoError(t, err)
	assert.Equal(t, value, (234.0+567.0)/2.0)
	_, err = result.Field("buz")
	assert.Error(t, err)
}
