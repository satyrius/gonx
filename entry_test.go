package gonx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEntry(t *testing.T) {
	entry := NewEntry(Fields{"foo": "1"})

	// Get existings field
	val, err := entry.Field("foo")
	assert.NoError(t, err)
	assert.Equal(t, val, "1")

	// Get field that does not exist
	val, err = entry.Field("bar")
	assert.Error(t, err)
	assert.Equal(t, val, "")
}

func TestEntryFloatField(t *testing.T) {
	entry := NewEntry(Fields{"foo": "1", "bar": "not a number"})

	// Get existings field
	val, err := entry.FloatField("foo")
	assert.NoError(t, err)
	assert.Equal(t, val, 1.0)

	// Type casting eror
	val, err = entry.FloatField("bar")
	assert.Error(t, err)
	assert.Equal(t, val, 0.0)

	// Get field that does not exist
	val, err = entry.FloatField("baz")
	assert.Error(t, err)
	assert.Equal(t, val, 0.0)
}

func TestSetEntryField(t *testing.T) {
	entry := NewEmptyEntry()
	assert.Equal(t, len(entry.fields), 0)

	entry.SetField("foo", "123")
	value, err := entry.Field("foo")
	assert.NoError(t, err)
	assert.Equal(t, value, "123")

	entry.SetField("foo", "234")
	value, err = entry.Field("foo")
	assert.NoError(t, err)
	assert.Equal(t, value, "234")
}

func TestSetEntryFloatField(t *testing.T) {
	entry := NewEmptyEntry()
	entry.SetFloatField("foo", 123.4567)
	value, err := entry.Field("foo")
	assert.NoError(t, err)
	assert.Equal(t, value, "123.46")
}

func TestSetEntryUintField(t *testing.T) {
	entry := NewEmptyEntry()
	entry.SetUintField("foo", 123)
	value, err := entry.Field("foo")
	assert.NoError(t, err)
	assert.Equal(t, value, "123")
}

func TestMergeEntries(t *testing.T) {
	entry1 := NewEntry(Fields{"foo": "1", "bar": "hello"})
	entry2 := NewEntry(Fields{"foo": "2", "bar": "hello", "name": "alpha"})
	entry1.Merge(entry2)

	val, err := entry1.Field("foo")
	assert.NoError(t, err)
	assert.Equal(t, val, "2")

	val, err = entry1.Field("bar")
	assert.NoError(t, err)
	assert.Equal(t, val, "hello")

	val, err = entry1.Field("name")
	assert.NoError(t, err)
	assert.Equal(t, val, "alpha")
}

func TestGetEntryGroupHash(t *testing.T) {
	entry1 := NewEntry(Fields{"foo": "1", "bar": "Hello world #1", "name": "alpha"})
	entry2 := NewEntry(Fields{"foo": "2", "bar": "Hello world #2", "name": "alpha"})
	entry3 := NewEntry(Fields{"foo": "2", "bar": "Hello world #3", "name": "alpha"})
	entry4 := NewEntry(Fields{"foo": "3", "bar": "Hello world #4", "name": "beta"})

	fields := []string{"name"}
	assert.Equal(t, entry1.FieldsHash(fields), entry2.FieldsHash(fields))
	assert.Equal(t, entry1.FieldsHash(fields), entry3.FieldsHash(fields))
	assert.NotEqual(t, entry1.FieldsHash(fields), entry4.FieldsHash(fields))

	fields = []string{"name", "foo"}
	assert.NotEqual(t, entry1.FieldsHash(fields), entry2.FieldsHash(fields))
	assert.Equal(t, entry2.FieldsHash(fields), entry3.FieldsHash(fields))
	assert.NotEqual(t, entry1.FieldsHash(fields), entry4.FieldsHash(fields))
	assert.NotEqual(t, entry2.FieldsHash(fields), entry4.FieldsHash(fields))
}

func TestPartialEntry(t *testing.T) {
	entry := NewEntry(Fields{"foo": "1", "bar": "Hello world #1", "name": "alpha"})
	partial := entry.Partial([]string{"name", "foo"})

	assert.Equal(t, len(partial.fields), 2)
	val, _ := partial.Field("name")
	assert.Equal(t, val, "alpha")
	val, _ = partial.Field("foo")
	assert.Equal(t, val, "1")
}

func TestEntrySumFields(t *testing.T) {
	entry := NewEntry(Fields{"foo": "123", "bar": "234", "baz": "not a number"})
	val := entry.SumFields([]string{"foo", "bar"})
	assert.Equal(t, val, 123.0+234)
	val = entry.SumFields([]string{"bar", "baz"})
	assert.Equal(t, val, 234.0)
}
