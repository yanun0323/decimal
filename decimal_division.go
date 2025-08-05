package decimal

import (
	"math/big"
)

var (
	shiftUint = 24
)

// Div returns d / d2. Ultra-optimized for trading performance.
func (d Decimal) Div(d2 Decimal) Decimal {
	a := normalize([]byte(d))
	b := normalize([]byte(d2))

	if isZero(b) {
		panic("division by zero")
	}

	if isZero(a) {
		return Zero()
	}

	ib, iShift := removeDecimalPoint(a)
	ib = shift(ib, shiftUint)
	iShift += shiftUint

	i, ok := new(big.Int).SetString(string(ib), 10)
	if !ok {
		panic("convert decimal to big int")
	}

	ib2, i2Shift := removeDecimalPoint(b)
	i2, ok := new(big.Int).SetString(string(ib2), 10)
	if !ok {
		panic("convert decimal to big int")
	}

	i = i.Div(i, i2)

	return Decimal(truncate(shift([]byte(i.String()), i2Shift-iShift), DivisionPrecision))
}
