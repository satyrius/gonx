package gonx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntry(t *testing.T) {
	entry := Entry{"foo": "1"}

	// Get existings field
	val, err := entry.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, val, "1")

	// Get field that does not exist
	val, err = entry.Get("bar")
	assert.Error(t, err)
	assert.Equal(t, val, "")
}

func TestEntryGetFloat(t *testing.T) {
	entry := Entry{"foo": "1", "bar": "not a number"}

	// Get existings field
	val, err := entry.GetFloat("foo")
	assert.NoError(t, err)
	assert.Equal(t, val, 1.0)

	// Type casting eror
	val, err = entry.GetFloat("bar")
	assert.Error(t, err)
	assert.Equal(t, val, 0.0)

	// Get field that does not exist
	val, err = entry.GetFloat("baz")
	assert.Error(t, err)
	assert.Equal(t, val, 0.0)
}
