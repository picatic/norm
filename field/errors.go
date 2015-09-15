package field

import (
	"errors"
	"fmt"
)

var (
	ErrorUnintializedShadow error = errors.New("Uninitialized ShadowValue")
	ErrorValueWasNotSet     error = errors.New("Value was not set")
)

func ErrorCouldNotScan(t string, v interface{}) error {
	return fmt.Errorf("Could not scan type: %s from value: %+v", t, v)
}
