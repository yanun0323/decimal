package decimal

import "errors"

var (
	errInvalidScanType = errors.New("invalid scan type")
	errInvalidBSONType = errors.New("invalid bson type")
	errNilReceiver     = errors.New("nil receiver")
)
