package decimal

import (
	"math/big"
	"sync"
)

// Div2 returns d / d2 with better performance compared with Div.
//
// The algorithm:
//
//	Suppose d  -> A / 10^iShift,  d2 -> B / 10^i2Shift.
//	We want result scaled to DivisionPrecision decimal digits.
//	Compute shiftExp = DivisionPrecision + i2Shift - iShift.
//	If shiftExp >= 0:
//	    resultInt = (A * 10^shiftExp) / B
//	else:
//	    resultInt = A / (B * 10^{-shiftExp})
//	The scaled integer is then shifted back by DivisionPrecision using the
//	existing shift helper.
func (d Decimal) Div(d2 Decimal) Decimal {
	return Decimal(div(normalize([]byte(d)), normalize([]byte(d2))))
}

func getDivisionPrecision() int {
	return DivisionPrecision
}

func div(a, b []byte) []byte {
	if isZero(b) {
		panic("division by zero")
	}

	if isZero(a) {
		return zeroBytes
	}

	// Remove decimal point to get pure integer representations
	ib, iShift := removeDecimalPoint(a)
	ib2, i2Shift := removeDecimalPoint(b)

	// Convert to big.Int (base 10)
	bigA, ok := new(big.Int).SetString(string(ib), 10)
	if !ok {
		panic("convert decimal to big int")
	}
	bigB, ok := new(big.Int).SetString(string(ib2), 10)
	if !ok {
		panic("convert decimal to big int")
	}

	dp := getDivisionPrecision()
	// Calculate scaling factor to preserve DivisionPrecision digits
	shiftExp := dp + i2Shift - iShift

	// Scale numerator or denominator accordingly
	var scaled big.Int
	scaled.Set(bigA)

	if shiftExp >= 0 {
		scaled.Mul(&scaled, pow10(shiftExp))
		scaled.Div(&scaled, bigB)
	} else {
		var denom big.Int
		denom.Mul(bigB, pow10(-shiftExp))
		scaled.Div(&scaled, &denom)
	}

	return tidyBytes(shift([]byte(scaled.String()), -dp))
}

// ------------------------- helper -------------------------

var (
	pow10Once  sync.Once
	pow10Table []*big.Int
	pow10Mu    sync.RWMutex
)

// initPow10Table initializes the power-of-ten lookup with 0 -> 1.
func initPow10Table() {
	pow10Table = []*big.Int{big.NewInt(1)}
}

// pow10 returns 10^n using a shared cache to avoid repetitive Exp calculations.
func pow10(n int) *big.Int {
	if n < 0 {
		panic("pow10: negative exponent")
	}

	// Ensure the table is initialized exactly once.
	pow10Once.Do(initPow10Table)

	pow10Mu.RLock()
	if n < len(pow10Table) {
		v := pow10Table[n]
		pow10Mu.RUnlock()
		return v
	}
	pow10Mu.RUnlock()

	// Upgrade lock to write
	pow10Mu.Lock()
	defer pow10Mu.Unlock()

	// Re-check to avoid duplicate work after lock upgrade
	for len(pow10Table) <= n {
		next := new(big.Int).Mul(pow10Table[len(pow10Table)-1], big.NewInt(10))
		pow10Table = append(pow10Table, next)
	}
	return pow10Table[n]
}
