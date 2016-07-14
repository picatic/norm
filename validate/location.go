package validate

import (
	"fmt"
)

type Locationer interface {
	Location() string
}

type Index int

func (i Index) Location() string {
	return fmt.Sprintf("[%d]", int(i))
}

type EmbeddedLocation []Locationer

func (el EmbeddedLocation) Location() string {
	result := el[0].Location()
	for i := 1; i < len(el); i++ {
		if l, ok := el[i].(Index); ok {
			result += l.Location()
			continue
		}

		result += "." + el[i].Location()
	}

	return result
}
