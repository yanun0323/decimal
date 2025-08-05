package decimal

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Decimal) Scan(value any) error {
	// first try to see if the data is stored in database as a Numeric datatype
	switch v := value.(type) {
	case float32:
		var err error
		*d, err = New(strconv.FormatFloat(float64(v), 'f', -1, 64))
		return err
	case float64:
		var err error
		*d, err = New(strconv.FormatFloat(v, 'f', -1, 64))
		return err
	case int64:
		// at least in sqlite3 when the value is 0 in db, the data is sent to us as an int64 instead of a float64 ...
		*d = Decimal(strconv.FormatInt(v, 10))
		return nil
	case string:
		var err error
		*d, err = New(v)
		return err
	case []byte:
		var err error
		*d, err = New(string(v))
		return err
	default:
		return fmt.Errorf("could not convert value '%+v' to decimal of type '%T'",
			value, value)
	}
}

// Value implements the driver.Valuer interface for database writes
func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}
