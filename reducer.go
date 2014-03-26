package gonx

// Reducer interface for Entries channel redure.
//
// Each Reduce method should accept input channel of Entries, do it's job and
// the result should be written to the output channel.
//
// It does not return values because usually it runs in a separate
// goroutine and it is handy to use channel for reduced data retrieval.
type Reducer interface {
	Reduce(input chan *Entry, output chan interface{})
}

// Implements Reducer interface for simple input entries redirection to
// the output channel.
type ReadAll struct {
}

// Redirect input Entries channel directly to the output without any
// modifications. It is useful when you want jast to read file fast
// using asynchronous with mapper routines.
func (r *ReadAll) Reduce(input chan *Entry, output chan interface{}) {
	output <- input
}

// Implements Reducer interface to count entries
type Count struct {
}

// Simply count entrries and write a sum to the output channel
func (r *Count) Reduce(input chan *Entry, output chan interface{}) {
	count := 0
	for {
		_, ok := <-input
		if !ok {
			break
		}
		count++
	}
	output <- count
}

// Implements Reducer interface for summarize Entry values for the given fields
type Sum struct {
	Fields []string
}

// Summarize given Entry fields and return a map with result for each field.
func (r *Sum) Reduce(input chan *Entry, output chan interface{}) {
	sum := make(map[string]float64)
	for entry := range input {
		for _, name := range r.Fields {
			val, err := entry.FloatField(name)
			if err == nil {
				sum[name] += val
			}
		}
	}
	output <- sum
}

// Implements Reducer interface for average entries values calculation
type Avg struct {
	Fields []string
}

// Calculate average value for input channel Entries, using configured Fields
// of the struct. Write result to the output channel as map[string]float64
func (r *Avg) Reduce(input chan *Entry, output chan interface{}) {
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
	output <- avg
}
