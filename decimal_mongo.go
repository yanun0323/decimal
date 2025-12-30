package decimal

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (d Decimal128) MarshalBSONValue() (bson.Type, []byte, error) {
	dec, err := bson.ParseDecimal128(d.String())
	if err != nil {
		return 0, nil, err
	}
	return bson.MarshalValue(dec)
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (d *Decimal128) UnmarshalBSONValue(typ bson.Type, data []byte) error {
	if d == nil {
		return errNilReceiver
	}
	switch typ {
	case bson.TypeNull, bson.TypeUndefined:
		*d = Decimal128{}
		return nil
	case bson.TypeDecimal128:
		var v bson.Decimal128
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		dec, err := New128FromString(v.String())
		if err != nil {
			return err
		}
		*d = dec
		return nil
	case bson.TypeString:
		var s string
		if err := bson.UnmarshalValue(typ, data, &s); err != nil {
			return err
		}
		dec, err := New128FromString(s)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	case bson.TypeInt32:
		var v int32
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		*d = New128FromInt(int64(v))
		return nil
	case bson.TypeInt64:
		var v int64
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		*d = New128FromInt(v)
		return nil
	case bson.TypeDouble:
		var v float64
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		dec, err := New128FromFloat(v)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	default:
		return errInvalidBSONType
	}
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (d Decimal256) MarshalBSONValue() (bson.Type, []byte, error) {
	return bson.MarshalValue(d.String())
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (d *Decimal256) UnmarshalBSONValue(typ bson.Type, data []byte) error {
	if d == nil {
		return errNilReceiver
	}
	switch typ {
	case bson.TypeNull, bson.TypeUndefined:
		*d = Decimal256{}
		return nil
	case bson.TypeDecimal128:
		var v bson.Decimal128
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		dec, err := New256FromString(v.String())
		if err != nil {
			return err
		}
		*d = dec
		return nil
	case bson.TypeString:
		var s string
		if err := bson.UnmarshalValue(typ, data, &s); err != nil {
			return err
		}
		dec, err := New256FromString(s)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	case bson.TypeInt32:
		var v int32
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		*d = New256FromInt(int64(v))
		return nil
	case bson.TypeInt64:
		var v int64
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		*d = New256FromInt(v)
		return nil
	case bson.TypeDouble:
		var v float64
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		dec, err := New256FromFloat(v)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	default:
		return errInvalidBSONType
	}
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (d Decimal512) MarshalBSONValue() (bson.Type, []byte, error) {
	return bson.MarshalValue(d.String())
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (d *Decimal512) UnmarshalBSONValue(typ bson.Type, data []byte) error {
	if d == nil {
		return errNilReceiver
	}
	switch typ {
	case bson.TypeNull, bson.TypeUndefined:
		*d = Decimal512{}
		return nil
	case bson.TypeDecimal128:
		var v bson.Decimal128
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		dec, err := New512FromString(v.String())
		if err != nil {
			return err
		}
		*d = dec
		return nil
	case bson.TypeString:
		var s string
		if err := bson.UnmarshalValue(typ, data, &s); err != nil {
			return err
		}
		dec, err := New512FromString(s)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	case bson.TypeInt32:
		var v int32
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		*d = New512FromInt(int64(v))
		return nil
	case bson.TypeInt64:
		var v int64
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		*d = New512FromInt(v)
		return nil
	case bson.TypeDouble:
		var v float64
		if err := bson.UnmarshalValue(typ, data, &v); err != nil {
			return err
		}
		dec, err := New512FromFloat(v)
		if err != nil {
			return err
		}
		*d = dec
		return nil
	default:
		return errInvalidBSONType
	}
}
