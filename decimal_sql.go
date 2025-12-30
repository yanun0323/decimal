package decimal

import "database/sql/driver"

// Scan implements the sql.Scanner interface.
func (d *Decimal128) Scan(src any) error {
	if d == nil {
		return errNilReceiver
	}
	if src == nil {
		*d = Decimal128{}
		return nil
	}
	switch v := src.(type) {
	case int64:
		*d = New128FromInt(v)
		return nil
	case float64:
		dec, err := New128FromFloat(v)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	case []byte:
		u, err := parseDecimalBytes128(v)
		if err != nil {
			return err
		}
		*d = Decimal128(u)
		return nil
	case string:
		dec, err := New128FromString(v)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	default:
		return errInvalidScanType
	}
}

// Value implements the driver.Valuer interface.
func (d Decimal128) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements the sql.Scanner interface.
func (d *Decimal256) Scan(src any) error {
	if d == nil {
		return errNilReceiver
	}
	if src == nil {
		*d = Decimal256{}
		return nil
	}
	switch v := src.(type) {
	case int64:
		*d = New256FromInt(v)
		return nil
	case float64:
		dec, err := New256FromFloat(v)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	case []byte:
		u, err := parseDecimalBytes(v)
		if err != nil {
			return err
		}
		*d = Decimal256(u)
		return nil
	case string:
		dec, err := New256FromString(v)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	default:
		return errInvalidScanType
	}
}

// Value implements the driver.Valuer interface.
func (d Decimal256) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements the sql.Scanner interface.
func (d *Decimal512) Scan(src any) error {
	if d == nil {
		return errNilReceiver
	}
	if src == nil {
		*d = Decimal512{}
		return nil
	}
	switch v := src.(type) {
	case int64:
		*d = New512FromInt(v)
		return nil
	case float64:
		dec, err := New512FromFloat(v)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	case []byte:
		u, err := parseDecimalBytes512(v)
		if err != nil {
			return err
		}
		*d = Decimal512(u)
		return nil
	case string:
		dec, err := New512FromString(v)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	default:
		return errInvalidScanType
	}
}

// Value implements the driver.Valuer interface.
func (d Decimal512) Value() (driver.Value, error) {
	return d.String(), nil
}
