package decimal

import (
	"math/big"
)

var (
	shiftUint = 32
)

// Div returns d / d2. Ultra-optimized for trading performance.
func (d Decimal) Div(d2 Decimal) Decimal {
	d = verify(d)
	d2 = verify(d2)

	if d2.isZero() {
		panic("division by zero")
	}

	if d.isZero() {
		return Zero()
	}

	ib, iShift := removeDecimalPoint([]byte(d))
	ib = shift(ib, shiftUint)
	iShift += shiftUint

	i, ok := new(big.Int).SetString(string(ib), 10)
	if !ok {
		panic("convert decimal to big int")
	}

	ib2, i2Shift := removeDecimalPoint([]byte(d2))
	i2, ok := new(big.Int).SetString(string(ib2), 10)
	if !ok {
		panic("convert decimal to big int")
	}

	i = i.Div(i, i2)

	return Decimal(i.String()).shift(i2Shift - iShift).truncate(DivisionPrecision)
}
