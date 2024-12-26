package gonx

// Reducer interface for Entries channel redure.
//
// Each Reduce method should accept input channel of Entries, do it's job and
// the result should be written to the output channel.
//
// It does not return values because usually it runs in a separate
// goroutine and it is handy to use channel for reduced data retrieval.
type Reducer interface {
	Reduce(input chan *Entry, output chan *Entry)
}

// ReadAll implements the Reducer interface for simple input entries redirected to
// the output channel.
type ReadAll struct {
}

// Reduce redirects input Entries channel directly to the output without any
// modifications. It is useful when you want jast to read file fast
// using asynchronous with mapper routines.
func (r *ReadAll) Reduce(input chan *Entry, output chan *Entry) {
	for entry := range input {
		output <- entry
	}
	close(output)
}

// Count implements the Reducer interface to count entries
type Count struct {
}

// Reduce simply counts entries and write a sum to the output channel
func (r *Count) Reduce(input chan *Entry, output chan *Entry) {
	var count uint64
	for {
		_, ok := <-input
		if !ok {
			break
		}
		count++
	}
	entry := NewEmptyEntry()
	entry.SetUintField("count", count)
	output <- entry
	close(output)
}

// Sum implements the Reducer interface for summarize Entry values for the given fields
type Sum struct {
	Fields []string
}

// Reduce summarizes given Entry fields and return a map with result for each field.
func (r *Sum) Reduce(input chan *Entry, output chan *Entry) {
	sum := make(map[string]float64)
	for entry := range input {
		for _, name := range r.Fields {
			val, err := entry.FloatField(name)
			if err == nil {
				sum[name] += val
			}
		}
	}
	entry := NewEmptyEntry()
	for name, val := range sum {
		entry.SetFloatField(name, val)
	}
	output <- entry
	close(output)
}

// Avg implements the Reducer interface for average entries values calculation
type Avg struct {
	Fields []string
}

// Reduce calculates the average value for input channel Entries, using configured Fields
// of the struct. Write result to the output channel as map[string]float64
func (r *Avg) Reduce(input chan *Entry, output chan *Entry) {
	avg := make(map[string]float64)
	count := 0.0
	for entry := range input {
		for _, name := range r.Fields {
			val, err := entry.FloatField(name)
			if err == nil {
				avg[name] = (avg[name]*count + val) / (count + 1)
			}
		}
		count++
	}
	entry := NewEmptyEntry()
	for name, val := range avg {
		entry.SetFloatField(name, val)
	}
	output <- entry
	close(output)
}

// Min implements the Reducer interface for min values calculation
type Min struct {
	Fields []string
}

// Reduce calculates the min values for input channel Entries, using configured Fields
// of the struct. Write result to the output channel as map[string]float64
func (r *Min) Reduce(input chan *Entry, output chan *Entry) {
	min := make(map[string]float64)
	for entry := range input {
		for _, name := range r.Fields {
			val, err := entry.FloatField(name)
			if err == nil {
				if val < min[name] || min[name] == 0 {
					min[name] = val
				}
			}
		}
	}
	entry := NewEmptyEntry()
	for name, val := range min {
		entry.SetFloatField(name, val)
	}
	output <- entry
	close(output)
}

// Max implements the Reducer interface for min values calculation
type Max struct {
	Fields []string
}

// Reduce calculates the min values for input channel Entries, using configured Fields
// of the struct. Write result to the output channel as map[string]float64
func (r *Max) Reduce(input chan *Entry, output chan *Entry) {
	max := make(map[string]float64)
	for entry := range input {
		for _, name := range r.Fields {
			val, err := entry.FloatField(name)
			if err == nil {
				if val > max[name] {
					max[name] = val
				}
			}
		}
	}
	entry := NewEmptyEntry()
	for name, val := range max {
		entry.SetFloatField(name, val)
	}
	output <- entry
	close(output)
}

// Chain implements the Reducer interface for chaining other reducers
type Chain struct {
	filters  []Filter
	reducers []Reducer
}

// NewChain creates a new chain of Reducers
func NewChain(reducers ...Reducer) *Chain {
	chain := new(Chain)
	for _, r := range reducers {
		if f, ok := interface{}(r).(Filter); ok {
			chain.filters = append(chain.filters, f)
		} else {
			chain.reducers = append(chain.reducers, r)
		}
	}
	return chain
}

// Reduce applies a chain of reducers to the input channel of entries and merge results
func (r *Chain) Reduce(input chan *Entry, output chan *Entry) {
	// Make input and output channel for each reducer
	subInput := make([]chan *Entry, len(r.reducers))
	subOutput := make([]chan *Entry, len(r.reducers))
	for i, reducer := range r.reducers {
		subInput[i] = make(chan *Entry, cap(input))
		subOutput[i] = make(chan *Entry, cap(output))
		go reducer.Reduce(subInput[i], subOutput[i])
	}

	// Read reducer master input channel
	for entry := range input {
		for _, f := range r.filters {
			entry = f.Filter(entry)
			if entry == nil {
				break
			}
		}
		// Publish input entry for each sub-reducers to process
		if entry != nil {
			for _, sub := range subInput {
				sub <- entry
			}
		}
	}
	for _, ch := range subInput {
		close(ch)
	}

	// Merge all results
	entry := NewEmptyEntry()
	for _, result := range subOutput {
		entry.Merge(<-result)
	}

	output <- entry
	close(output)
}

// GroupBy implements the Reducer interface to apply other reducers and get data grouped by
// given fields.
type GroupBy struct {
	Fields   []string
	reducers []Reducer
}

// NewGroupBy creates a new GroupBy Reducer
func NewGroupBy(fields []string, reducers ...Reducer) *GroupBy {
	return &GroupBy{
		Fields:   fields,
		reducers: reducers,
	}
}

// Reduce applies related reducers and group data by Fields.
func (r *GroupBy) Reduce(input chan *Entry, output chan *Entry) {
	subInput := make(map[string]chan *Entry)
	subOutput := make(map[string]chan *Entry)

	// Read reducer master input channel and create discinct input chanel
	// for each entry key we group by
	for entry := range input {
		key := entry.FieldsHash(r.Fields)
		if _, ok := subInput[key]; !ok {
			subInput[key] = make(chan *Entry, cap(input))
			subOutput[key] = make(chan *Entry, cap(output)+1)
			subOutput[key] <- entry.Partial(r.Fields)
			go NewChain(r.reducers...).Reduce(subInput[key], subOutput[key])
		}
		subInput[key] <- entry
	}
	for _, ch := range subInput {
		close(ch)
	}
	for _, ch := range subOutput {
		entry := <-ch
		entry.Merge(<-ch)
		output <- entry
	}
	close(output)
}
