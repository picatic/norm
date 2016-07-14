package field

import (
	"database/sql/driver"
)

func ScanValuer(value interface{}) (interface{}, error) {
	valuer, ok := value.(driver.Valuer)
	if !ok {
		return value, nil
	}

	v, err := valuer.Value()
	if err != nil {
		return nil, err
	}

	return v, nil
}
