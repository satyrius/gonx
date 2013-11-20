package gonx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadAllReducer(t *testing.T) {
	reducer := new(ReadAll)
	assert.Implements(t, (*Reducer)(nil), reducer)

	// Prepare import chanel
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
