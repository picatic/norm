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
