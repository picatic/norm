package norm

import (
	"github.com/gocraft/dbr"
)

// Copy errors from dbr
var (
	ErrNotFound           = dbr.ErrNotFound
	ErrNotSupported       = dbr.ErrNotSupported
	ErrTableNotSpecified  = dbr.ErrTableNotSpecified
	ErrColumnNotSpecified = dbr.ErrColumnNotSpecified
	ErrInvalidPointer     = dbr.ErrInvalidPointer
	ErrPlaceholderCount   = dbr.ErrPlaceholderCount
	ErrInvalidSliceLength = dbr.ErrInvalidSliceLength
	ErrCantConvertToTime  = dbr.ErrCantConvertToTime
	ErrInvalidTimestring  = dbr.ErrInvalidTimestring
)
