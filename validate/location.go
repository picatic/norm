package validate

import (
	"fmt"
)

type Locationer interface {
	Location() string
}

type EmbeddedLocation []Locationer

func (el EmbeddedLocation) Location() string {
	embeddedLoc := el[0].Location()

	for i := 1; i < len(el); i++ {
		switch l := el[i].(type) {
		case Index:
			embeddedLoc += l.Location()
		default:
			embeddedLoc += "." + l.Location()
		}
	}

	return embeddedLoc
}

type Index int

func (i Index) Location() string {
	return fmt.Sprintf("[%d]", int(i))
}
