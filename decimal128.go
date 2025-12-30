// All overflow is truncated to 128-bit two's-complement arithmetic.
package decimal

import (
	"encoding/binary"
	"math"
	"math/bits"
	"strconv"
)

const (
	scaleDigits128            = 16
	maxDecimalDigits128       = 38 // 10^38 < 2^128 < 10^39
	decimal128PrecisionDigits = 16
)

// Decimal128 is a fixed-scale decimal stored as a 128-bit two's-complement integer.
//
// The scale is 10^16, so the numeric value is: raw / 10^16.
// The zero value represents 0 and is ready to use.
//
// Memory layout is fixed and little-endian across 2 uint64 words.
// This is required for binary compatibility.
type Decimal128 [2]uint64

// u128 is an internal unsigned 128-bit integer in little-endian word order.
type u128 [2]uint64

// u256_128 is an internal unsigned 256-bit integer in little-endian word order.
type u256_128 [4]uint64

var (
	scale128 = u128{0x002386f26fc10000, 0x0000000000000000} // 10^16
	// 10^32, used by Inv().
	scaleSquared128 = u256_128{0x85acef8100000000, 0x000004ee2d6d415b, 0x0, 0x0}
	pow10_128       = buildPow10_128()
	// ln(2) scaled by 1e16.
	constLn2_128 = Decimal128(u128{0x0018a0230abe4edd, 0x0000000000000000})
	// ln(10) scaled by 1e16.
	constLn10_128 = Decimal128(u128{0x0051cde3b15487e8, 0x0000000000000000})
)

// NewDecimal128 constructs a Decimal128 from integer and fractional parts.
//
// intPart keeps only the lowest 16 decimal digits (higher digits are dropped).
// decimalPart keeps only the highest 16 fractional digits (lower digits are dropped).
// decimalPart is interpreted as fractional digits with an implicit scale based on
// its decimal digit length (e.g. 987654321 -> 0.987654321). It is then scaled to
// 10^16 before combining with intPart. If decimalPart has more than 16 digits, it
// is truncated toward zero. The result is: intPart*10^16 + scaled(decimalPart),
// with two's-complement wrap on overflow.
func NewDecimal128(intPart, decimalPart int64) Decimal128 {
	intPart = truncateIntPart128(intPart)
	decimalPart = truncateDecimalPart128(decimalPart)
	ip := mul128(u128FromInt64(intPart), scale128)
	raw := lower128(ip)
	if decimalPart != 0 {
		abs, neg := absInt64(decimalPart)
		digits := decimalDigitsU64(abs)
		shift := scaleDigits128 - digits
		frac := u128{abs, 0}
		var scaled u128
		if shift == 0 {
			scaled = frac
		} else if shift < 0 {
			factor := pow10Value128(int64(-shift))
			if isZero128(factor) {
				scaled = u128{}
			} else {
				scaled = divByU128Trunc(frac, factor)
			}
		} else if shift < len(pow10U64) {
			scaled = mul128ByUint64(frac, pow10U64[shift])
		} else {
			p := mul128(frac, pow10Value128(int64(shift)))
			scaled = lower128(p)
		}
		if neg {
			scaled = neg128(scaled)
		}
		raw = add128(raw, scaled)
	}
	return Decimal128(applyPrecision128(raw))
}

// NewDecimal128FromString parses a decimal string with optional sign, dot, and exponent.
//
// It accepts leading/trailing ASCII whitespace and optional '_' separators.
// Exponent shifting is applied first, then integer digits beyond 16 are dropped and
// fractional digits beyond 16 are dropped. Excess fractional digits are truncated
// (toward zero) to the fixed 16-digit scale.
func NewDecimal128FromString(s string) (Decimal128, error) {
	u, err := parseDecimalString128(s)
	if err != nil {
		return Decimal128{}, err
	}
	return Decimal128(u), nil
}

// NewDecimal128FromInt constructs a Decimal128 from an int64 integer value.
//
// Only the lowest 16 decimal digits are kept.
func NewDecimal128FromInt(v int64) Decimal128 {
	v = truncateIntPart128(v)
	p := mul128(u128FromInt64(v), scale128)
	return Decimal128(applyPrecision128(lower128(p)))
}

// NewDecimal128FromFloat converts a float64 to Decimal128 by truncating toward zero.
//
// The integer part keeps only the lowest 16 decimal digits. The fractional part
// keeps only the highest 16 digits (closest to the decimal point).
// NaN or Inf returns an error. Overflow wraps according to two's-complement truncation.
func NewDecimal128FromFloat(v float64) (Decimal128, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return Decimal128{}, errInvalidFloat
	}
	if v == 0 {
		return Decimal128{}, nil
	}
	neg := v < 0
	if neg {
		v = -v
	}
	intPart, frac := math.Modf(v)
	if intPart >= 1e16 {
		intPart = math.Mod(intPart, 1e16)
	}
	fracPart := math.Floor(frac * 1e16)
	if fracPart < 0 {
		fracPart = 0
	}
	if fracPart >= 1e16 {
		fracPart = 1e16 - 1
	}
	ip := mul128(u128{uint64(intPart), 0}, scale128)
	fp := u128{uint64(fracPart), 0}
	u := add128(lower128(ip), fp)
	u = applyPrecision128(u)
	if neg {
		u = neg128(u)
	}
	return Decimal128(u), nil
}

// Int64 returns the integer and fractional parts as int64 values.
//
// Both parts are truncated to int64 with two's-complement wrap if out of range.
// The fractional part is returned in 10^16 base units.
func (d Decimal128) Int64() (intPart, decimalPart int64) {
	u := u128(d)
	if isZero128(u) {
		return 0, 0
	}
	neg := isNeg128(u)
	abs := u
	if neg {
		abs = neg128(abs)
	}
	q, r := divMod128By128(abs, scale128)
	qi := int64(q[0])
	ri := int64(r[0])
	if neg {
		qi = -qi
		ri = -ri
	}
	return qi, ri
}

// Float64 converts Decimal128 to float64.
//
// Precision is limited by float64; large values may overflow to Inf.
func (d Decimal128) Float64() float64 {
	u := u128(d)
	if isZero128(u) {
		return 0
	}
	neg := isNeg128(u)
	if neg {
		u = neg128(u)
	}
	f := u128ToFloat(u)
	f = f / 1e16
	if neg {
		f = -f
	}
	return f
}

// String returns the shortest decimal representation without trailing zeros.
func (d Decimal128) String() string {
	u := u128(d)
	if isZero128(u) {
		return "0"
	}
	neg := isNeg128(u)
	if neg {
		u = neg128(u)
	}
	q, r := divMod128By128(u, scale128)
	intStr := u128ToDecimal(q)
	if isZero128(r) {
		if neg {
			return "-" + intStr
		}
		return intStr
	}
	frac := u128ToDecimalFixed(r, scaleDigits128)
	frac = trimRightZeros(frac)
	if neg {
		return "-" + intStr + "." + frac
	}
	return intStr + "." + frac
}

// StringFixed returns a decimal string with exactly n fractional digits.
//
// If n > 16 it is truncated to 16. If n <= 0, no fractional part is shown.
func (d Decimal128) StringFixed(n int) string {
	if n > scaleDigits128 {
		n = scaleDigits128
	}
	if n <= 0 {
		u := u128(d)
		if isZero128(u) {
			return "0"
		}
		neg := isNeg128(u)
		if neg {
			u = neg128(u)
		}
		q, _ := divMod128By128(u, scale128)
		intStr := u128ToDecimal(q)
		if neg {
			return "-" + intStr
		}
		return intStr
	}
	return d.stringFixed128(n)
}

// AppendString appends the shortest decimal representation without trailing zeros to dst.
func (d Decimal128) AppendString(dst []byte) []byte {
	u := u128(d)
	if isZero128(u) {
		return append(dst, '0')
	}
	neg := isNeg128(u)
	if neg {
		u = neg128(u)
	}
	q, r := divMod128By128(u, scale128)
	if neg {
		dst = append(dst, '-')
	}
	dst = appendU128Decimal(dst, q)
	if isZero128(r) {
		return dst
	}
	dst = append(dst, '.')
	fracStart := len(dst)
	dst = appendU128DecimalFixed(dst, r, scaleDigits128)
	end := len(dst)
	for end > fracStart && dst[end-1] == '0' {
		end--
	}
	if end == fracStart {
		return dst[:fracStart-1]
	}
	return dst[:end]
}

// AppendStringFixed appends a decimal string with exactly n fractional digits to dst.
//
// If n > 16 it is truncated to 16. If n <= 0, no fractional part is appended.
func (d Decimal128) AppendStringFixed(dst []byte, n int) []byte {
	if n > scaleDigits128 {
		n = scaleDigits128
	}
	u := u128(d)
	if isZero128(u) {
		if n <= 0 {
			return append(dst, '0')
		}
		dst = append(dst, '0', '.')
		for i := 0; i < n; i++ {
			dst = append(dst, '0')
		}
		return dst
	}
	neg := isNeg128(u)
	if neg {
		u = neg128(u)
	}
	q, r := divMod128By128(u, scale128)
	if neg {
		dst = append(dst, '-')
	}
	dst = appendU128Decimal(dst, q)
	if n > 0 {
		dst = append(dst, '.')
		var fracBuf [64]byte
		frac := appendU128DecimalFixed(fracBuf[:0], r, scaleDigits128)
		if n < scaleDigits128 {
			frac = frac[:n]
		}
		dst = append(dst, frac...)
	}
	return dst
}

// IsZero reports whether the value is exactly zero.
func (d Decimal128) IsZero() bool {
	return isZero128(u128(d))
}

// IsPositive reports whether the value is greater than zero.
func (d Decimal128) IsPositive() bool {
	u := u128(d)
	return !isZero128(u) && !isNeg128(u)
}

// IsNegative reports whether the value is less than zero.
func (d Decimal128) IsNegative() bool {
	u := u128(d)
	return !isZero128(u) && isNeg128(u)
}

// Sign returns 0 if zero, 1 if positive, and 2 if negative.
func (d Decimal128) Sign() int {
	u := u128(d)
	if isZero128(u) {
		return 0
	}
	if isNeg128(u) {
		return 2
	}
	return 1
}

// Neg returns the arithmetic negation of d.
func (d Decimal128) Neg() Decimal128 {
	return Decimal128(neg128(u128(d)))
}

// Inv returns the multiplicative inverse (1/d).
//
// For zero, it returns zero.
func (d Decimal128) Inv() Decimal128 {
	u := u128(d)
	if isZero128(u) {
		return d
	}
	neg := isNeg128(u)
	if neg {
		u = neg128(u)
	}
	q, _ := divMod256_128By128(scaleSquared128, u)
	if neg {
		q = neg128(q)
	}
	return Decimal128(q)
}

// Abs returns the absolute value of d.
func (d Decimal128) Abs() Decimal128 {
	u := u128(d)
	if isNeg128(u) {
		u = neg128(u)
	}
	return Decimal128(u)
}

// Truncate truncates to n fractional digits (banker-friendly truncation toward zero).
//
// If n > 16, it returns d unchanged. If n <= -16, it returns zero.
func (d Decimal128) Truncate(n int) Decimal128 {
	return d.truncateWithMode128(n, roundModeTowardZero)
}

// Shift moves the decimal point by n digits.
//
// Positive n shifts left (multiply by 10^n), negative n shifts right (divide by 10^-n).
// If n > 16, it returns d unchanged. If n <= -16, it returns zero.
func (d Decimal128) Shift(n int) Decimal128 {
	if n > scaleDigits128 {
		return d
	}
	if n <= -scaleDigits128 {
		return Decimal128{}
	}
	if n == 0 {
		return d
	}
	if n > 0 {
		factor := pow10Mod128(int64(n))
		p := mul128(u128(d), factor)
		return Decimal128(lower128(p))
	}
	factor := pow10Value128(int64(-n))
	return Decimal128(divByU128Trunc(u128(d), factor))
}

// Round rounds to n fractional digits using banker's rounding.
//
// If n > 16, it returns d unchanged. If n <= -16, it returns zero.
func (d Decimal128) Round(n int) Decimal128 {
	return d.truncateWithMode128(n, roundModeBanker)
}

// RoundAwayFromZero rounds to n fractional digits, away from zero.
//
// If n > 16, it returns d unchanged. If n <= -16, it returns zero.
func (d Decimal128) RoundAwayFromZero(n int) Decimal128 {
	return d.truncateWithMode128(n, roundModeAwayFromZero)
}

// RoundTowardToZero truncates to n fractional digits toward zero.
//
// If n > 16, it returns d unchanged. If n <= -16, it returns zero.
func (d Decimal128) RoundTowardToZero(n int) Decimal128 {
	return d.truncateWithMode128(n, roundModeTowardZero)
}

// Ceil rounds toward positive infinity with n fractional digits.
//
// If n > 16, it returns d unchanged. If n <= -16, it returns zero.
func (d Decimal128) Ceil(n int) Decimal128 {
	return d.truncateWithMode128(n, roundModeCeil)
}

// Floor rounds toward negative infinity with n fractional digits.
//
// If n > 16, it returns d unchanged. If n <= -16, it returns zero.
func (d Decimal128) Floor(n int) Decimal128 {
	return d.truncateWithMode128(n, roundModeFloor)
}

// Equal reports whether d == other.
func (d Decimal128) Equal(other Decimal128) bool {
	return d == other
}

// GreaterThan reports whether d > other.
func (d Decimal128) GreaterThan(other Decimal128) bool {
	return cmp128Signed(u128(d), u128(other)) > 0
}

// LessThan reports whether d < other.
func (d Decimal128) LessThan(other Decimal128) bool {
	return cmp128Signed(u128(d), u128(other)) < 0
}

// GreaterOrEqual reports whether d >= other.
func (d Decimal128) GreaterOrEqual(other Decimal128) bool {
	return cmp128Signed(u128(d), u128(other)) >= 0
}

// LessOrEqual reports whether d <= other.
func (d Decimal128) LessOrEqual(other Decimal128) bool {
	return cmp128Signed(u128(d), u128(other)) <= 0
}

// Add returns d + other with 128-bit truncation on overflow.
func (d Decimal128) Add(other Decimal128) Decimal128 {
	return Decimal128(add128(u128(d), u128(other)))
}

// Sub returns d - other with 128-bit truncation on overflow.
func (d Decimal128) Sub(other Decimal128) Decimal128 {
	return Decimal128(sub128(u128(d), u128(other)))
}

// Mul returns d * other with fixed 16-digit scale.
func (d Decimal128) Mul(other Decimal128) Decimal128 {
	u := u128(d)
	v := u128(other)
	if isZero128(u) || isZero128(v) {
		return Decimal128{}
	}
	neg := isNeg128(u) != isNeg128(v)
	if isNeg128(u) {
		u = neg128(u)
	}
	if isNeg128(v) {
		v = neg128(v)
	}
	p := mul128(u, v)
	q, _ := divMod256_128By128(p, scale128)
	if neg {
		q = neg128(q)
	}
	return Decimal128(q)
}

// Div returns d / other with fixed 16-digit scale.
//
// If other is zero, it returns d unchanged.
func (d Decimal128) Div(other Decimal128) Decimal128 {
	u := u128(d)
	v := u128(other)
	if isZero128(v) {
		return d
	}
	if isZero128(u) {
		return Decimal128{}
	}
	neg := isNeg128(u) != isNeg128(v)
	if isNeg128(u) {
		u = neg128(u)
	}
	if isNeg128(v) {
		v = neg128(v)
	}
	p := mul128(u, scale128)
	q, _ := divMod256_128By128(p, v)
	if neg {
		q = neg128(q)
	}
	return Decimal128(q)
}

// Mod returns d % other using truncation toward zero.
//
// If other is zero, it returns d unchanged.
func (d Decimal128) Mod(other Decimal128) Decimal128 {
	u := u128(d)
	v := u128(other)
	if isZero128(v) {
		return d
	}
	if isZero128(u) {
		return Decimal128{}
	}
	neg := isNeg128(u)
	if neg {
		u = neg128(u)
	}
	if isNeg128(v) {
		v = neg128(v)
	}
	_, r := divMod128By128(u, v)
	if neg {
		r = neg128(r)
	}
	return Decimal128(r)
}

// Pow returns d raised to an integer power specified by other.
//
// The exponent is truncated toward zero to an int64. Negative exponents use Inv().
func (d Decimal128) Pow(other Decimal128) Decimal128 {
	if other.IsZero() {
		return NewDecimal128FromInt(1)
	}
	trunc := other.Truncate(0)
	exp, _ := trunc.Int64()
	if exp == 0 {
		return NewDecimal128FromInt(1)
	}
	negExp := exp < 0
	if negExp {
		exp = -exp
	}
	result := NewDecimal128FromInt(1)
	base := d
	for exp > 0 {
		if exp&1 == 1 {
			result = result.Mul(base)
		}
		base = base.Mul(base)
		exp >>= 1
	}
	if negExp {
		return result.Inv()
	}
	return result
}

// Sqrt returns the square root of d using Newton's method.
//
// For negative values, it returns d unchanged.
func (d Decimal128) Sqrt() Decimal128 {
	if d.IsNegative() {
		return d
	}
	if d.IsZero() {
		return d
	}
	guess := d.Float64()
	if guess <= 0 || math.IsInf(guess, 0) || math.IsNaN(guess) {
		return d
	}
	guess = math.Sqrt(guess)
	gd, err := NewDecimal128FromFloat(guess)
	if err != nil || gd.IsZero() {
		gd = NewDecimal128FromInt(1)
	}
	prev := gd
	for i := 0; i < 32; i++ {
		inv := d.Div(gd)
		gd = gd.Add(inv)
		gd = divDecimal128ByUint64(gd, 2)
		if gd.Equal(prev) {
			break
		}
		prev = gd
	}
	return gd
}

// Exp returns e^d using range reduction and a Taylor series.
func (d Decimal128) Exp() Decimal128 {
	if d.IsZero() {
		return NewDecimal128FromInt(1)
	}
	k := d.Div(constLn2_128).Round(0)
	kInt, _ := k.Int64()
	kLn2 := NewDecimal128FromInt(kInt).Mul(constLn2_128)
	r := d.Sub(kLn2)

	term := NewDecimal128FromInt(1)
	sum := term
	for i := uint64(1); i <= 96; i++ {
		term = term.Mul(r)
		term = divDecimal128ByUint64(term, i)
		if term.IsZero() {
			break
		}
		sum = sum.Add(term)
	}

	if kInt > 0 {
		sum = Decimal128(shl128(u128(sum), uint(kInt)))
	} else if kInt < 0 {
		sum = Decimal128(shr128(u128(sum), uint(-kInt)))
	}
	return sum
}

// Log returns the natural logarithm of d.
//
// For d <= 0, it returns d unchanged.
func (d Decimal128) Log() Decimal128 {
	if !d.IsPositive() {
		return d
	}
	u := u128(d)
	k := int64(bitLen128(u)) - int64(bitLen128(scale128))
	var mScaled u128
	if k >= 0 {
		mScaled = shr128(u, uint(k))
	} else {
		mScaled = shl128(u, uint(-k))
	}
	m := Decimal128(mScaled)
	one := NewDecimal128FromInt(1)
	mMinus := m.Sub(one)
	mPlus := m.Add(one)
	t := mMinus.Div(mPlus)
	t2 := t.Mul(t)
	term := t
	sum := t
	for i := uint64(3); i <= 199; i += 2 {
		term = term.Mul(t2)
		add := divDecimal128ByUint64(term, i)
		if add.IsZero() {
			break
		}
		sum = sum.Add(add)
	}
	lnm := sum.Mul(NewDecimal128FromInt(2))
	kLn2 := NewDecimal128FromInt(k).Mul(constLn2_128)
	return lnm.Add(kLn2)
}

// Log2 returns the base-2 logarithm of d.
//
// For d <= 0, it returns d unchanged.
func (d Decimal128) Log2() Decimal128 {
	ln := d.Log()
	if !d.IsPositive() {
		return ln
	}
	return ln.Div(constLn2_128)
}

// Log10 returns the base-10 logarithm of d.
//
// For d <= 0, it returns d unchanged.
func (d Decimal128) Log10() Decimal128 {
	ln := d.Log()
	if !d.IsPositive() {
		return ln
	}
	return ln.Div(constLn10_128)
}

// EncodeBinary encodes the raw 128-bit value into 16 bytes (little-endian).
func (d Decimal128) EncodeBinary() ([]byte, error) {
	var out [16]byte
	binary.LittleEndian.PutUint64(out[0:8], d[0])
	binary.LittleEndian.PutUint64(out[8:16], d[1])
	return out[:], nil
}

// AppendBinary appends the raw 128-bit value as 16 bytes (little-endian) to dst.
func (d Decimal128) AppendBinary(dst []byte) []byte {
	var out [16]byte
	binary.LittleEndian.PutUint64(out[0:8], d[0])
	binary.LittleEndian.PutUint64(out[8:16], d[1])
	return append(dst, out[:]...)
}

// NewDecimal128FromBinary decodes a 16-byte little-endian binary representation.
//
// Precision rules are applied after decoding.
func NewDecimal128FromBinary(b []byte) (Decimal128, error) {
	if len(b) != 16 {
		return Decimal128{}, errInvalidBinaryLen
	}
	u := u128{
		binary.LittleEndian.Uint64(b[0:8]),
		binary.LittleEndian.Uint64(b[8:16]),
	}
	return Decimal128(applyPrecision128(u)), nil
}

// EncodeJSON encodes the decimal as a JSON string.
func (d Decimal128) EncodeJSON() ([]byte, error) {
	s := d.String()
	buf := make([]byte, 0, len(s)+2)
	buf = strconv.AppendQuote(buf, s)
	return buf, nil
}

// AppendJSON appends the decimal as a JSON string to dst.
func (d Decimal128) AppendJSON(dst []byte) []byte {
	dst = append(dst, '"')
	dst = d.AppendString(dst)
	return append(dst, '"')
}

// NewDecimal128FromJSON decodes a JSON string or number into a Decimal128.
func NewDecimal128FromJSON(b []byte) (Decimal128, error) {
	start, end := trimSpaceBytes(b)
	if start >= end {
		return Decimal128{}, errInvalidJSONDecimal
	}
	if b[start] == '"' {
		if end-start < 2 || b[end-1] != '"' {
			return Decimal128{}, errInvalidJSONDecimal
		}
		for i := start + 1; i < end-1; i++ {
			if b[i] == '\\' || b[i] < 0x20 {
				return Decimal128{}, errInvalidJSONDecimal
			}
		}
		u, err := parseDecimalBytes128(b[start+1 : end-1])
		if err != nil {
			return Decimal128{}, err
		}
		return Decimal128(u), nil
	}
	u, err := parseDecimalBytes128(b[start:end])
	if err != nil {
		return Decimal128{}, err
	}
	return Decimal128(u), nil
}

// truncateWithMode128 is an internal helper.
func (d Decimal128) truncateWithMode128(n int, mode roundMode) Decimal128 {
	if n > scaleDigits128 {
		return d
	}
	if n <= -scaleDigits128 {
		return Decimal128{}
	}
	if n == scaleDigits128 {
		return d
	}
	u := u128(d)
	if isZero128(u) {
		return d
	}
	neg := isNeg128(u)
	if neg {
		u = neg128(u)
	}
	factor := pow10Value128(int64(scaleDigits128 - n))
	q, r := divMod128By128(u, factor)
	if !isZero128(r) {
		switch mode {
		case roundModeAwayFromZero:
			q = add128(q, u128{1, 0})
		case roundModeBanker:
			cmp := cmp128(add128(r, r), factor)
			if cmp > 0 || (cmp == 0 && (q[0]&1) == 1) {
				q = add128(q, u128{1, 0})
			}
		case roundModeCeil:
			if !neg {
				q = add128(q, u128{1, 0})
			}
		case roundModeFloor:
			if neg {
				q = add128(q, u128{1, 0})
			}
		}
	}
	res := mul128(q, factor)
	out := lower128(res)
	if neg {
		out = neg128(out)
	}
	return Decimal128(out)
}

// stringFixed128 is an internal helper.
func (d Decimal128) stringFixed128(n int) string {
	u := u128(d)
	if isZero128(u) {
		if n == 0 {
			return "0"
		}
		return "0." + repeatZero(n)
	}
	neg := isNeg128(u)
	if neg {
		u = neg128(u)
	}
	q, r := divMod128By128(u, scale128)
	intStr := u128ToDecimal(q)
	frac := u128ToDecimalFixed(r, scaleDigits128)
	if n < scaleDigits128 {
		frac = frac[:n]
	}
	if neg {
		return "-" + intStr + "." + frac
	}
	return intStr + "." + frac
}

// buildPow10_128 is an internal helper.
func buildPow10_128() [65]u128 {
	var p [65]u128
	p[0] = u128{1, 0}
	for i := 1; i < len(p); i++ {
		p[i] = mul128ByUint64(p[i-1], 10)
	}
	return p
}

// pow10Mod128 is an internal helper.
func pow10Mod128(n int64) u128 {
	if n <= 0 {
		return u128{1, 0}
	}
	if n <= 64 {
		return pow10_128[n]
	}
	result := u128{1, 0}
	base := u128{10, 0}
	for n > 0 {
		if (n & 1) == 1 {
			result = mul128Lo(result, base)
		}
		base = mul128Lo(base, base)
		n >>= 1
	}
	return result
}

// pow10Value128 is an internal helper.
func pow10Value128(n int64) u128 {
	if n <= 0 {
		return u128{1, 0}
	}
	if n <= 64 {
		return pow10_128[n]
	}
	if n > maxDecimalDigits128 {
		return u128{}
	}
	return pow10Mod128(n)
}

// parseDecimalString128 is an internal helper.
func parseDecimalString128(s string) (u128, error) {
	start, end := trimSpaceString(s)
	if start >= end {
		return u128{}, errInvalidDecimal
	}
	idx := start
	sign := 1
	if s[idx] == '+' {
		idx++
	} else if s[idx] == '-' {
		sign = -1
		idx++
	}
	var val u128
	fracDigits := int64(0)
	sawDigit := false
	sawDot := false
	for idx < end {
		c := s[idx]
		if c == '_' {
			idx++
			continue
		}
		if c == '.' {
			if sawDot {
				return u128{}, errInvalidDecimal
			}
			sawDot = true
			idx++
			continue
		}
		if c == 'e' || c == 'E' {
			idx++
			break
		}
		if c < '0' || c > '9' {
			return u128{}, errInvalidDecimal
		}
		sawDigit = true
		val = mul128ByUint64(val, 10)
		val = add128(val, u128{uint64(c - '0'), 0})
		if sawDot {
			fracDigits++
		}
		idx++
	}
	if !sawDigit {
		return u128{}, errInvalidDecimal
	}
	var exp int64
	if idx < end {
		expSign := int64(1)
		if s[idx] == '+' {
			idx++
		} else if s[idx] == '-' {
			expSign = -1
			idx++
		}
		startExp := idx
		for idx < end {
			c := s[idx]
			if c == '_' {
				idx++
				continue
			}
			if c < '0' || c > '9' {
				return u128{}, errInvalidDecimal
			}
			if exp < (1 << 60) {
				exp = exp*10 + int64(c-'0')
			}
			idx++
		}
		if startExp == idx {
			return u128{}, errInvalidDecimal
		}
		exp *= expSign
	}
	shift := exp - fracDigits + scaleDigits128
	if shift >= 0 {
		p := mul128(val, pow10Mod128(shift))
		val = lower128(p)
	} else {
		factor := pow10Value128(-shift)
		if isZero128(factor) {
			val = u128{}
		} else {
			val = divByU128Trunc(val, factor)
		}
	}
	if sign < 0 {
		val = neg128(val)
	}
	return applyPrecision128(val), nil
}

// parseDecimalBytes128 is an internal helper.
func parseDecimalBytes128(b []byte) (u128, error) {
	start, end := trimSpaceBytes(b)
	if start >= end {
		return u128{}, errInvalidDecimal
	}
	idx := start
	sign := 1
	if b[idx] == '+' {
		idx++
	} else if b[idx] == '-' {
		sign = -1
		idx++
	}
	var val u128
	fracDigits := int64(0)
	sawDigit := false
	sawDot := false
	for idx < end {
		c := b[idx]
		if c == '_' {
			idx++
			continue
		}
		if c == '.' {
			if sawDot {
				return u128{}, errInvalidDecimal
			}
			sawDot = true
			idx++
			continue
		}
		if c == 'e' || c == 'E' {
			idx++
			break
		}
		if c < '0' || c > '9' {
			return u128{}, errInvalidDecimal
		}
		sawDigit = true
		val = mul128ByUint64(val, 10)
		val = add128(val, u128{uint64(c - '0'), 0})
		if sawDot {
			fracDigits++
		}
		idx++
	}
	if !sawDigit {
		return u128{}, errInvalidDecimal
	}
	var exp int64
	if idx < end {
		expSign := int64(1)
		if b[idx] == '+' {
			idx++
		} else if b[idx] == '-' {
			expSign = -1
			idx++
		}
		startExp := idx
		for idx < end {
			c := b[idx]
			if c == '_' {
				idx++
				continue
			}
			if c < '0' || c > '9' {
				return u128{}, errInvalidDecimal
			}
			if exp < (1 << 60) {
				exp = exp*10 + int64(c-'0')
			}
			idx++
		}
		if startExp == idx {
			return u128{}, errInvalidDecimal
		}
		exp *= expSign
	}
	shift := exp - fracDigits + scaleDigits128
	if shift >= 0 {
		p := mul128(val, pow10Mod128(shift))
		val = lower128(p)
	} else {
		factor := pow10Value128(-shift)
		if isZero128(factor) {
			val = u128{}
		} else {
			val = divByU128Trunc(val, factor)
		}
	}
	if sign < 0 {
		val = neg128(val)
	}
	return applyPrecision128(val), nil
}

// truncateIntPart128 keeps only the lowest 16 decimal digits.
func truncateIntPart128(v int64) int64 {
	if v == 0 {
		return 0
	}
	abs, neg := absInt64(v)
	if decimalDigitsU64(abs) > decimal128PrecisionDigits {
		abs = abs % pow10U64[decimal128PrecisionDigits]
	}
	if neg {
		return -int64(abs)
	}
	return int64(abs)
}

// truncateDecimalPart128 keeps only the highest 16 fractional digits.
func truncateDecimalPart128(v int64) int64 {
	if v == 0 {
		return 0
	}
	abs, neg := absInt64(v)
	digits := decimalDigitsU64(abs)
	if digits > decimal128PrecisionDigits {
		abs = abs / pow10U64[digits-decimal128PrecisionDigits]
	}
	if neg {
		return -int64(abs)
	}
	return int64(abs)
}

// applyPrecision128 drops integer digits beyond 16 and fractional digits beyond 16.
// It assumes the input is scaled by 10^16.
func applyPrecision128(u u128) u128 {
	if isZero128(u) {
		return u
	}
	neg := isNeg128(u)
	if neg {
		u = neg128(u)
	}
	intPart, fracPart := divMod128By128(u, scale128)
	_, intRem := divMod128ByUint64(intPart, pow10U64[decimal128PrecisionDigits])
	intPart = u128{intRem, 0}
	// scaleDigits128 == decimal128PrecisionDigits, so fractional trimming is a no-op.
	raw := add128(lower128(mul128(intPart, scale128)), fracPart)
	if neg {
		raw = neg128(raw)
	}
	return raw
}

// divDecimal128ByUint64 is an internal helper.
func divDecimal128ByUint64(d Decimal128, n uint64) Decimal128 {
	if n == 0 {
		return d
	}
	u := u128(d)
	if isZero128(u) {
		return d
	}
	neg := isNeg128(u)
	if neg {
		u = neg128(u)
	}
	q, _ := divMod128ByUint64(u, n)
	if neg {
		q = neg128(q)
	}
	return Decimal128(q)
}

// u128FromFloatTrunc is an internal helper.
func u128FromFloatTrunc(v float64) u128 {
	if v <= 0 {
		return u128{}
	}
	bits64 := math.Float64bits(v)
	exp := int((bits64>>52)&0x7ff) - 1023
	mant := bits64 & ((uint64(1) << 52) - 1)
	mant |= uint64(1) << 52
	if exp < 0 {
		return u128{}
	}
	shift := exp - 52
	if shift >= 128 {
		return u128{}
	}
	u := u128{mant, 0}
	if shift == 0 {
		return u
	}
	if shift > 0 {
		return shl128(u, uint(shift))
	}
	return shr128(u, uint(-shift))
}

// u128ToFloat is an internal helper.
func u128ToFloat(u u128) float64 {
	if isZero128(u) {
		return 0
	}
	f := 0.0
	for i := 1; i >= 0; i-- {
		f = f*math.Exp2(64) + float64(u[i])
		if i == 0 {
			break
		}
	}
	return f
}

// appendU128Decimal is an internal helper.
func appendU128Decimal(dst []byte, v u128) []byte {
	if isZero128(v) {
		return append(dst, '0')
	}
	const base = uint64(1_000_000_000_000_000_000) // 1e18
	var parts [3]uint64
	n := 0
	for !isZero128(v) {
		var rem uint64
		v, rem = divMod128ByUint64(v, base)
		parts[n] = rem
		n++
	}
	for i := n - 1; i >= 0; i-- {
		if i == n-1 {
			dst = strconv.AppendUint(dst, parts[i], 10)
		} else {
			var tmp [20]byte
			chunk := strconv.AppendUint(tmp[:0], parts[i], 10)
			for j := len(chunk); j < 18; j++ {
				dst = append(dst, '0')
			}
			dst = append(dst, chunk...)
		}
		if i == 0 {
			break
		}
	}
	return dst
}

// appendU128DecimalFixed is an internal helper.
func appendU128DecimalFixed(dst []byte, v u128, width int) []byte {
	if width <= 0 {
		return dst
	}
	var tmp [48]byte
	num := appendU128Decimal(tmp[:0], v)
	if len(num) >= width {
		return append(dst, num[len(num)-width:]...)
	}
	for i := 0; i < width-len(num); i++ {
		dst = append(dst, '0')
	}
	return append(dst, num...)
}

// u128ToDecimal is an internal helper.
func u128ToDecimal(v u128) string {
	if isZero128(v) {
		return "0"
	}
	const base = uint64(1_000_000_000_000_000_000) // 1e18
	var parts [3]uint64
	n := 0
	for !isZero128(v) {
		var rem uint64
		v, rem = divMod128ByUint64(v, base)
		parts[n] = rem
		n++
	}
	buf := make([]byte, 0, 48)
	for i := n - 1; i >= 0; i-- {
		if i == n-1 {
			buf = strconv.AppendUint(buf, parts[i], 10)
		} else {
			var tmp [20]byte
			chunk := strconv.AppendUint(tmp[:0], parts[i], 10)
			for j := len(chunk); j < 18; j++ {
				buf = append(buf, '0')
			}
			buf = append(buf, chunk...)
		}
		if i == 0 {
			break
		}
	}
	return string(buf)
}

// u128ToDecimalFixed is an internal helper.
func u128ToDecimalFixed(v u128, width int) string {
	if width <= 0 {
		return ""
	}
	s := u128ToDecimal(v)
	if len(s) >= width {
		return s[len(s)-width:]
	}
	pad := width - len(s)
	buf := make([]byte, 0, width)
	for i := 0; i < pad; i++ {
		buf = append(buf, '0')
	}
	buf = append(buf, s...)
	return string(buf)
}

// divByU128Trunc is an internal helper.
func divByU128Trunc(v u128, d u128) u128 {
	if isZero128(d) {
		return v
	}
	neg := isNeg128(v)
	if neg {
		v = neg128(v)
	}
	q, _ := divMod128By128(v, d)
	if neg {
		q = neg128(q)
	}
	return q
}

// isZero128 is an internal helper.
func isZero128(u u128) bool {
	return u[0]|u[1] == 0
}

// isNeg128 is an internal helper.
func isNeg128(u u128) bool {
	return (u[1]>>63)&1 == 1
}

// add128 is an internal helper.
func add128(a, b u128) u128 {
	var out u128
	var c uint64
	out[0], c = bits.Add64(a[0], b[0], 0)
	out[1], _ = bits.Add64(a[1], b[1], c)
	return out
}

// sub128 is an internal helper.
func sub128(a, b u128) u128 {
	var out u128
	var c uint64
	out[0], c = bits.Sub64(a[0], b[0], 0)
	out[1], _ = bits.Sub64(a[1], b[1], c)
	return out
}

// neg128 is an internal helper.
func neg128(a u128) u128 {
	var out u128
	out[0] = ^a[0]
	out[1] = ^a[1]
	var c uint64
	out[0], c = bits.Add64(out[0], 1, 0)
	out[1], _ = bits.Add64(out[1], 0, c)
	return out
}

// cmp128 is an internal helper.
func cmp128(a, b u128) int {
	if a[1] > b[1] {
		return 1
	}
	if a[1] < b[1] {
		return -1
	}
	if a[0] > b[0] {
		return 1
	}
	if a[0] < b[0] {
		return -1
	}
	return 0
}

// cmp128Signed is an internal helper.
func cmp128Signed(a, b u128) int {
	na := isNeg128(a)
	nb := isNeg128(b)
	if na && !nb {
		return -1
	}
	if !na && nb {
		return 1
	}
	return cmp128(a, b)
}

// mul128ByUint64 is an internal helper.
func mul128ByUint64(a u128, m uint64) u128 {
	var out u128
	var carry uint64
	for i := 0; i < 2; i++ {
		hi, lo := bits.Mul64(a[i], m)
		lo, c := bits.Add64(lo, carry, 0)
		out[i] = lo
		carry = hi + c
	}
	return out
}

// divMod128ByUint64 is an internal helper.
func divMod128ByUint64(a u128, d uint64) (u128, uint64) {
	var q u128
	var r uint64
	for i := 1; i >= 0; i-- {
		q[i], r = bits.Div64(r, a[i], d)
		if i == 0 {
			break
		}
	}
	return q, r
}

// mul128 is an internal helper.
func mul128(a, b u128) u256_128 {
	var p u256_128
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			hi, lo := bits.Mul64(a[i], b[j])
			k := i + j
			var c uint64
			p[k], c = bits.Add64(p[k], lo, 0)
			p[k+1], c = bits.Add64(p[k+1], hi, c)
			idx := k + 2
			for c != 0 && idx < 4 {
				p[idx], c = bits.Add64(p[idx], 0, c)
				idx++
			}
		}
	}
	return p
}

// mul128Lo is an internal helper.
func mul128Lo(a, b u128) u128 {
	return lower128(mul128(a, b))
}

// lower128 is an internal helper.
func lower128(p u256_128) u128 {
	return u128{p[0], p[1]}
}

// divMod128By128 is an internal helper.
func divMod128By128(n u128, d u128) (u128, u128) {
	if isZero128(d) {
		return u128{}, u128{}
	}
	var n256 u256_128
	n256[0] = n[0]
	n256[1] = n[1]
	q, r := divMod256_128By128(n256, d)
	return q, r
}

// divMod256_128By128 is an internal helper.
func divMod256_128By128(n u256_128, d u128) (u128, u128) {
	if isZero128(d) {
		return u128{}, u128{}
	}
	if isZeroU256_128(n) {
		return u128{}, u128{}
	}
	nBits := bitLen256_128(n)
	dBits := bitLen128(d)
	if nBits < dBits {
		return u128{}, lower128(n)
	}
	shift := nBits - dBits
	var q u128
	rem := n
	dShift := shl128To256(d, uint(shift))
	for i := shift; i >= 0; i-- {
		if cmp256_128(rem, dShift) >= 0 {
			rem = sub256_128(rem, dShift)
			if i < 128 {
				q[int(i/64)] |= 1 << uint(i%64)
			}
		}
		if i == 0 {
			break
		}
		dShift = shr1_256_128(dShift)
	}
	return q, lower128(rem)
}

// isZeroU256_128 is an internal helper.
func isZeroU256_128(u u256_128) bool {
	return u[0]|u[1]|u[2]|u[3] == 0
}

// bitLen128 is an internal helper.
func bitLen128(u u128) int {
	if u[1] != 0 {
		return 64 + bits.Len64(u[1])
	}
	if u[0] != 0 {
		return bits.Len64(u[0])
	}
	return 0
}

// bitLen256_128 is an internal helper.
func bitLen256_128(u u256_128) int {
	for i := 3; i >= 0; i-- {
		if u[i] != 0 {
			return i*64 + bits.Len64(u[i])
		}
		if i == 0 {
			break
		}
	}
	return 0
}

// cmp256_128 is an internal helper.
func cmp256_128(a, b u256_128) int {
	for i := 3; i >= 0; i-- {
		if a[i] > b[i] {
			return 1
		}
		if a[i] < b[i] {
			return -1
		}
		if i == 0 {
			break
		}
	}
	return 0
}

// sub256_128 is an internal helper.
func sub256_128(a, b u256_128) u256_128 {
	var out u256_128
	var c uint64
	out[0], c = bits.Sub64(a[0], b[0], 0)
	out[1], c = bits.Sub64(a[1], b[1], c)
	out[2], c = bits.Sub64(a[2], b[2], c)
	out[3], _ = bits.Sub64(a[3], b[3], c)
	return out
}

// shl128 is an internal helper.
func shl128(u u128, shift uint) u128 {
	if shift == 0 {
		return u
	}
	if shift >= 128 {
		return u128{}
	}
	wordShift := shift / 64
	bitShift := shift % 64
	var out u128
	for i := 0; i < 2; i++ {
		dst := i + int(wordShift)
		if dst >= 2 {
			continue
		}
		out[dst] |= u[i] << bitShift
		if bitShift != 0 && dst+1 < 2 {
			out[dst+1] |= u[i] >> (64 - bitShift)
		}
	}
	return out
}

// shr128 is an internal helper.
func shr128(u u128, shift uint) u128 {
	if shift == 0 {
		return u
	}
	if shift >= 128 {
		return u128{}
	}
	wordShift := shift / 64
	bitShift := shift % 64
	var out u128
	for i := int(wordShift); i < 2; i++ {
		src := i
		dst := i - int(wordShift)
		out[dst] |= u[src] >> bitShift
		if bitShift != 0 && src+1 < 2 {
			out[dst] |= u[src+1] << (64 - bitShift)
		}
	}
	return out
}

// shl128To256 is an internal helper.
func shl128To256(u u128, shift uint) u256_128 {
	if shift == 0 {
		return u256_128{u[0], u[1], 0, 0}
	}
	if shift >= 256 {
		return u256_128{}
	}
	wordShift := shift / 64
	bitShift := shift % 64
	var out u256_128
	for i := 0; i < 2; i++ {
		dst := i + int(wordShift)
		if dst >= 4 {
			continue
		}
		out[dst] |= u[i] << bitShift
		if bitShift != 0 && dst+1 < 4 {
			out[dst+1] |= u[i] >> (64 - bitShift)
		}
	}
	return out
}

// shr1_256_128 is an internal helper.
func shr1_256_128(u u256_128) u256_128 {
	var out u256_128
	var carry uint64
	for i := 3; i >= 0; i-- {
		out[i] = (u[i] >> 1) | (carry << 63)
		carry = u[i] & 1
		if i == 0 {
			break
		}
	}
	return out
}

// u128FromInt64 is an internal helper.
func u128FromInt64(v int64) u128 {
	if v >= 0 {
		return u128{uint64(v), 0}
	}
	uv := uint64(v)
	return u128{uv, ^uint64(0)}
}
