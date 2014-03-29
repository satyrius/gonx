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

func TestGetEntryKey(t *testing.T) {
	entry1 := NewEntry(Fields{"foo": "1", "bar": "Hello world #1", "name": "alpha"})
	entry2 := NewEntry(Fields{"foo": "2", "bar": "Hello world #2", "name": "alpha"})
	entry3 := NewEntry(Fields{"foo": "2", "bar": "Hello world #3", "name": "alpha"})
	entry4 := NewEntry(Fields{"foo": "3", "bar": "Hello world #4", "name": "beta"})

	reducer := NewGroupBy([]string{"name"})
	assert.Equal(t, reducer.GetEntryKey(entry1), reducer.GetEntryKey(entry2))
	assert.Equal(t, reducer.GetEntryKey(entry1), reducer.GetEntryKey(entry3))
	assert.NotEqual(t, reducer.GetEntryKey(entry1), reducer.GetEntryKey(entry4))

	reducer = NewGroupBy([]string{"name", "foo"})
	assert.NotEqual(t, reducer.GetEntryKey(entry1), reducer.GetEntryKey(entry2))
	assert.Equal(t, reducer.GetEntryKey(entry2), reducer.GetEntryKey(entry3))
	assert.NotEqual(t, reducer.GetEntryKey(entry1), reducer.GetEntryKey(entry4))
	assert.NotEqual(t, reducer.GetEntryKey(entry2), reducer.GetEntryKey(entry4))
}

func TestGroupByReducer(t *testing.T) {
	reducer := NewGroupBy(
		// Fields to group by
		[]string{"host"},
		// Result reducers
		&Sum{[]string{"foo", "bar"}},
		new(Count),
	)
	assert.Implements(t, (*Reducer)(nil), reducer)

	// Prepare import channel
	input := make(chan *Entry, 10)
	input <- NewEntry(Fields{
		"uri":  "/asd/fgh",
		"host": "alpha.example.com",
		"foo":  "1",
		"bar":  "2",
		"baz":  "3",
	})
	input <- NewEntry(Fields{
		"uri":  "/zxc/vbn",
		"host": "beta.example.com",
		"foo":  "4",
		"bar":  "5",
		"baz":  "6",
	})
	input <- NewEntry(Fields{
		"uri":  "/ijk/lmn",
		"host": "beta.example.com",
		"foo":  "7",
		"bar":  "8",
		"baz":  "9",
	})
	close(input)
	output := make(chan *Entry, 1) // Make it buffered to avoid deadlock
	reducer.Reduce(input, output)

	// TODO read entries from output and assert grouped data
}
