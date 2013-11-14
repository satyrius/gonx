package gonx

import (
	"fmt"
)

type Entry map[string]string

func (entry *Entry) Get(name string) (value string, err error) {
	value, ok := (*entry)[name]
	if !ok {
		err = fmt.Errorf("Field '%v' does not found in record %+v", name, *entry)
	}
	return
}
