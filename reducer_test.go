package gonx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadAllReducer(t *testing.T) {
	reducer := new(ReadAll)
	assert.Implements(t, (*Reducer)(nil), reducer)

	// Prepare import channel
	input := make(chan Entry, 1)
	input <- Entry{}

	output := make(chan interface{}, 1) // Make it buffered to avoid deadlock
	reducer.Reduce(input, output)

	// ReadAll reducer writes input channel to the output
	result, opened := <-output
	assert.True(t, opened)
	_, ok := result.(chan Entry)
	assert.True(t, ok)
}

func TestSumReducer(t *testing.T) {
	reducer := &Sum{[]string{"foo", "bar"}}
	assert.Implements(t, (*Reducer)(nil), reducer)

	// Prepare import channel
	input := make(chan Entry, 2)
	input <- Entry{
		"uri": "/asd/fgh",
		"foo": "123",
		"bar": "234",
		"baz": "345",
	}
	input <- Entry{
		"uri": "/zxc/vbn",
		"foo": "456",
		"bar": "567",
		"baz": "678",
	}
	close(input)
	output := make(chan interface{}, 1) // Make it buffered to avoid deadlock
	reducer.Reduce(input, output)

	// ReadAll reducer writes a map with sums into output
	result, opened := <-output
	assert.True(t, opened)
	sum, ok := result.(map[string]float64)
	assert.True(t, ok)
	// The result should contain sums for "foo" and "bar" fields
	value, ok := sum["foo"]
	assert.True(t, ok)
	assert.Equal(t, value, 579.0)
	value, ok = sum["bar"]
	assert.True(t, ok)
	assert.Equal(t, value, 801.0)
	_, ok = sum["baz"]
	assert.False(t, ok)
}
