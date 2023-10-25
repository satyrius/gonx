package gonx

import "strings"

type SpliterPaser struct {
	pattern           string
	spliter           string
	splitedPattern    []string
	lenSplitedPattern int
}

func NewSpliteParser(pattern string, spliter string) *SpliterPaser {
	splitedPattern := strings.Split(pattern, spliter)
	return &SpliterPaser{pattern: pattern, splitedPattern: splitedPattern, lenSplitedPattern: len(splitedPattern), spliter: spliter}
}

func (sp *SpliterPaser) ParseString(line string) (entry *Entry, err error) {
	inputS := strings.SplitN(line, sp.spliter, sp.lenSplitedPattern)
	minL := min(len(inputS), sp.lenSplitedPattern)
	e := NewWithSize(minL)
	for i := 0; i < minL; i++ {
		if strings.Contains(sp.splitedPattern[i], " ") {
			e.Merge(parseWithSplit(inputS[i], sp.splitedPattern[i], " "))
		} else if strings.Contains(sp.splitedPattern[i], "-") {
			e.Merge(parseWithSplit(inputS[i], sp.splitedPattern[i], "-"))
		} else {
			e.SetField(strings.TrimPrefix(sp.splitedPattern[i], "$"), inputS[i])
		}
	}
	return e, nil
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func parseWithSplit(input, pattern, split string) *Entry {
	inputSubS := strings.Split(input, split)
	patternSubP := strings.Split(pattern, split)
	minL := min(len(inputSubS), len(patternSubP))
	e := NewWithSize(minL)
	for i := 0; i < minL; i++ {
		e.SetField(strings.TrimPrefix(patternSubP[i], "$"), inputSubS[i])
	}
	return e
}
