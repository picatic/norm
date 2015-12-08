package norm

import (
	"github.com/gocraft/dbr"
)

// Copy errors from dbr
var (
	ErrArgumentMismatch   = dbr.ErrArgumentMismatch
	ErrInvalidSliceLength = dbr.ErrInvalidSliceLength
	ErrInvalidSliceValue  = dbr.ErrInvalidSliceValue
	ErrInvalidValue       = dbr.ErrInvalidValue
	ErrNotFound           = dbr.ErrNotFound
	ErrNotUTF8            = dbr.ErrNotUTF8
)
