// All overflow is truncated to 512-bit two's-complement arithmetic.
package decimal

import (
	"encoding/binary"
	"math"
	"math/bits"
	"strconv"
)

const (
	scaleDigits512            = 64
	maxDecimalDigits512       = 154 // 10^154 < 2^512 < 10^155
	decimal512PrecisionDigits = 64
)

// Decimal512 is a fixed-scale decimal stored as a 512-bit two's-complement integer.
//
// The scale is 10^64, so the numeric value is: raw / 10^64.
// The zero value represents 0 and is ready to use.
//
// Memory layout is fixed and little-endian across 8 uint64 words.
// This is required for binary compatibility.
type Decimal512 [8]uint64

// u512x is an internal unsigned 512-bit integer in little-endian word order.
type u512x [8]uint64

// u1024 is an internal unsigned 1024-bit integer in little-endian word order.
type u1024 [16]uint64

var (
	scale512  = u512x{0x0, 0x6e38ed64bf6a1f01, 0xe93ff9f4daa797ed, 0x0000000000184f03, 0x0, 0x0, 0x0, 0x0} // 10^64
	pow10_512 = buildPow10_512()
	// ln(2) scaled by 1e64.
	constLn2_512 = Decimal512(u512x{0x0ab61e0079bfe6de, 0x15f4c45211bf4ade, 0x2aa49d4009a60be2, 0x000000000010d977, 0x0, 0x0, 0x0, 0x0})
	// ln(10) scaled by 1e64.
	constLn10_512 = Decimal512(u512x{0xea71dc1617d12f21, 0xdbb8f1138566b86c, 0xafb42c6b2c8625b2, 0x000000000037f905, 0x0, 0x0, 0x0, 0x0})
)

// New512 constructs a Decimal512 from integer and fractional parts.
//
// intPart keeps only the lowest 64 decimal digits (higher digits are dropped).
// decimalPart keeps only the highest 64 fractional digits (lower digits are dropped).
// decimalPart is interpreted as fractional digits with an implicit scale based on
// its decimal digit length (e.g. 987654321 -> 0.987654321). It is then scaled to
// 10^64 before combining with intPart. If decimalPart has more than 64 digits, it
// is truncated toward zero. The result is: intPart*10^64 + scaled(decimalPart),
// with two's-complement wrap on overflow.
func New512(intPart, decimalPart int64) Decimal512 {
	ip := mul512(u512FromInt64(intPart), scale512)
	raw := lower512(ip)
	if decimalPart != 0 {
		abs, neg := absInt64(decimalPart)
		digits := decimalDigitsU64(abs)
		shift := scaleDigits512 - digits
		frac := u512x{abs, 0, 0, 0, 0, 0, 0, 0}
		var scaled u512x
		if shift == 0 {
			scaled = frac
		} else if shift < 0 {
			factor := pow10Value512(int64(-shift))
			if isZero512(factor) {
				scaled = u512x{}
			} else {
				scaled = divByU512Trunc(frac, factor)
			}
		} else if shift < len(pow10U64) {
			scaled = mul512ByUint64(frac, pow10U64[shift])
		} else {
			p := mul512(frac, pow10Value512(int64(shift)))
			scaled = lower512(p)
		}
		if neg {
			scaled = neg512(scaled)
		}
		raw = add512(raw, scaled)
	}
	return Decimal512(applyPrecision512(raw))
}

// New512FromString parses a decimal string with optional sign, dot, and exponent.
//
// It accepts leading/trailing ASCII whitespace and optional '_' separators.
// Exponent shifting is applied first, then integer digits beyond 64 are dropped and
// fractional digits beyond 64 are dropped. Excess fractional digits are truncated
// (toward zero) to the fixed 64-digit scale.
func New512FromString(s string) (Decimal512, error) {
	u, err := parseDecimalString512(s)
	if err != nil {
		return Decimal512{}, err
	}
	return Decimal512(u), nil
}

// New512FromInt constructs a Decimal512 from an int64 integer value.
func New512FromInt(v int64) Decimal512 {
	p := mul512(u512FromInt64(v), scale512)
	return Decimal512(applyPrecision512(lower512(p)))
}

// New512FromFloat converts a float64 to Decimal512 by truncating toward zero.
//
// NaN or Inf returns an error. Overflow wraps according to two's-complement truncation.
func New512FromFloat(v float64) (Decimal512, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return Decimal512{}, errInvalidFloat
	}
	if v == 0 {
		return Decimal512{}, nil
	}
	neg := v < 0
	if neg {
		v = -v
	}
	scaled := v * 1e64
	if math.IsInf(scaled, 0) {
		return Decimal512{}, errInvalidFloat
	}
	u := u512FromFloatTrunc(scaled)
	u = applyPrecision512(u)
	if neg {
		u = neg512(u)
	}
	return Decimal512(u), nil
}

// Int64 returns the integer and fractional parts as int64 values.
//
// Both parts are truncated to int64 with two's-complement wrap if out of range.
// The fractional part is returned in 10^64 base units.
func (d Decimal512) Int64() (intPart, decimalPart int64) {
	u := u512x(d)
	if isZero512(u) {
		return 0, 0
	}
	neg := isNeg512(u)
	abs := u
	if neg {
		abs = neg512(abs)
	}
	q, r := divMod512By512(abs, scale512)
	qi := int64(q[0])
	ri := int64(r[0])
	if neg {
		qi = -qi
		ri = -ri
	}
	return qi, ri
}

// Float64 converts Decimal512 to float64.
//
// Precision is limited by float64; large values may overflow to Inf.
func (d Decimal512) Float64() float64 {
	u := u512x(d)
	if isZero512(u) {
		return 0
	}
	neg := isNeg512(u)
	if neg {
		u = neg512(u)
	}
	f := u512ToFloat(u)
	f = f / 1e64
	if neg {
		f = -f
	}
	return f
}

// String returns the shortest decimal representation without trailing zeros.
func (d Decimal512) String() string {
	u := u512x(d)
	if isZero512(u) {
		return "0"
	}
	neg := isNeg512(u)
	if neg {
		u = neg512(u)
	}
	q, r := divMod512By512(u, scale512)
	intStr := u512ToDecimal(q)
	if isZero512(r) {
		if neg {
			return "-" + intStr
		}
		return intStr
	}
	frac := u512ToDecimalFixed(r, scaleDigits512)
	frac = trimRightZeros(frac)
	if neg {
		return "-" + intStr + "." + frac
	}
	return intStr + "." + frac
}

// StringFixed returns a decimal string with exactly n fractional digits.
//
// If n > 64 it is truncated to 64. If n <= 0, no fractional part is shown.
func (d Decimal512) StringFixed(n int) string {
	if n > scaleDigits512 {
		n = scaleDigits512
	}
	if n <= 0 {
		u := u512x(d)
		if isZero512(u) {
			return "0"
		}
		neg := isNeg512(u)
		if neg {
			u = neg512(u)
		}
		q, _ := divMod512By512(u, scale512)
		intStr := u512ToDecimal(q)
		if neg {
			return "-" + intStr
		}
		return intStr
	}
	return d.stringFixed512(n)
}

// AppendString appends the shortest decimal representation without trailing zeros to dst.
func (d Decimal512) AppendString(dst []byte) []byte {
	u := u512x(d)
	if isZero512(u) {
		return append(dst, '0')
	}
	neg := isNeg512(u)
	if neg {
		u = neg512(u)
	}
	q, r := divMod512By512(u, scale512)
	if neg {
		dst = append(dst, '-')
	}
	dst = appendU512Decimal(dst, q)
	if isZero512(r) {
		return dst
	}
	dst = append(dst, '.')
	fracStart := len(dst)
	dst = appendU512DecimalFixed(dst, r, scaleDigits512)
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
// If n > 64 it is truncated to 64. If n <= 0, no fractional part is appended.
func (d Decimal512) AppendStringFixed(dst []byte, n int) []byte {
	if n > scaleDigits512 {
		n = scaleDigits512
	}
	u := u512x(d)
	if isZero512(u) {
		if n <= 0 {
			return append(dst, '0')
		}
		dst = append(dst, '0', '.')
		for i := 0; i < n; i++ {
			dst = append(dst, '0')
		}
		return dst
	}
	neg := isNeg512(u)
	if neg {
		u = neg512(u)
	}
	q, r := divMod512By512(u, scale512)
	if neg {
		dst = append(dst, '-')
	}
	dst = appendU512Decimal(dst, q)
	if n > 0 {
		dst = append(dst, '.')
		var fracBuf [64]byte
		frac := appendU512DecimalFixed(fracBuf[:0], r, scaleDigits512)
		if n < scaleDigits512 {
			frac = frac[:n]
		}
		dst = append(dst, frac...)
	}
	return dst
}

// IsZero reports whether the value is exactly zero.
func (d Decimal512) IsZero() bool {
	return isZero512(u512x(d))
}

// IsPositive reports whether the value is greater than zero.
func (d Decimal512) IsPositive() bool {
	u := u512x(d)
	return !isZero512(u) && !isNeg512(u)
}

// IsNegative reports whether the value is less than zero.
func (d Decimal512) IsNegative() bool {
	u := u512x(d)
	return !isZero512(u) && isNeg512(u)
}

// Sign returns 0 if zero, 1 if positive, and 2 if negative.
func (d Decimal512) Sign() int {
	u := u512x(d)
	if isZero512(u) {
		return 0
	}
	if isNeg512(u) {
		return 2
	}
	return 1
}

// Neg returns the arithmetic negation of d.
func (d Decimal512) Neg() Decimal512 {
	return Decimal512(neg512(u512x(d)))
}

// Inv returns the multiplicative inverse (1/d).
//
// For zero, it returns zero.
func (d Decimal512) Inv() Decimal512 {
	u := u512x(d)
	if isZero512(u) {
		return d
	}
	neg := isNeg512(u)
	if neg {
		u = neg512(u)
	}
	q, _ := divMod512By512(scale512Squared(), u)
	if neg {
		q = neg512(q)
	}
	return Decimal512(q)
}

// Abs returns the absolute value of d.
func (d Decimal512) Abs() Decimal512 {
	u := u512x(d)
	if isNeg512(u) {
		u = neg512(u)
	}
	return Decimal512(u)
}

// Truncate truncates to n fractional digits (banker-friendly truncation toward zero).
//
// If n > 64, it returns d unchanged. If n <= -64, it returns zero.
func (d Decimal512) Truncate(n int) Decimal512 {
	return d.truncateWithMode512(n, roundModeTowardZero)
}

// Shift moves the decimal point by n digits.
//
// Positive n shifts left (multiply by 10^n), negative n shifts right (divide by 10^-n).
// If n > 64, it returns d unchanged. If n <= -64, it returns zero.
func (d Decimal512) Shift(n int) Decimal512 {
	if n > scaleDigits512 {
		return d
	}
	if n <= -scaleDigits512 {
		return Decimal512{}
	}
	if n == 0 {
		return d
	}
	if n > 0 {
		factor := pow10Mod512(int64(n))
		p := mul512(u512x(d), factor)
		return Decimal512(lower512(p))
	}
	factor := pow10Value512(int64(-n))
	return Decimal512(divByU512Trunc(u512x(d), factor))
}

// Round rounds to n fractional digits using banker's rounding.
//
// If n > 64, it returns d unchanged. If n <= -64, it returns zero.
func (d Decimal512) Round(n int) Decimal512 {
	return d.truncateWithMode512(n, roundModeBanker)
}

// RoundAwayFromZero rounds to n fractional digits, away from zero.
//
// If n > 64, it returns d unchanged. If n <= -64, it returns zero.
func (d Decimal512) RoundAwayFromZero(n int) Decimal512 {
	return d.truncateWithMode512(n, roundModeAwayFromZero)
}

// RoundTowardToZero truncates to n fractional digits toward zero.
//
// If n > 64, it returns d unchanged. If n <= -64, it returns zero.
func (d Decimal512) RoundTowardToZero(n int) Decimal512 {
	return d.truncateWithMode512(n, roundModeTowardZero)
}

// Ceil rounds toward positive infinity with n fractional digits.
//
// If n > 64, it returns d unchanged. If n <= -64, it returns zero.
func (d Decimal512) Ceil(n int) Decimal512 {
	return d.truncateWithMode512(n, roundModeCeil)
}

// Floor rounds toward negative infinity with n fractional digits.
//
// If n > 64, it returns d unchanged. If n <= -64, it returns zero.
func (d Decimal512) Floor(n int) Decimal512 {
	return d.truncateWithMode512(n, roundModeFloor)
}

// Equal reports whether d == other.
func (d Decimal512) Equal(other Decimal512) bool {
	return d == other
}

// GreaterThan reports whether d > other.
func (d Decimal512) GreaterThan(other Decimal512) bool {
	return cmp512Signed(u512x(d), u512x(other)) > 0
}

// LessThan reports whether d < other.
func (d Decimal512) LessThan(other Decimal512) bool {
	return cmp512Signed(u512x(d), u512x(other)) < 0
}

// GreaterOrEqual reports whether d >= other.
func (d Decimal512) GreaterOrEqual(other Decimal512) bool {
	return cmp512Signed(u512x(d), u512x(other)) >= 0
}

// LessOrEqual reports whether d <= other.
func (d Decimal512) LessOrEqual(other Decimal512) bool {
	return cmp512Signed(u512x(d), u512x(other)) <= 0
}

// Add returns d + other with 512-bit truncation on overflow.
func (d Decimal512) Add(other Decimal512) Decimal512 {
	return Decimal512(add512(u512x(d), u512x(other)))
}

// Sub returns d - other with 512-bit truncation on overflow.
func (d Decimal512) Sub(other Decimal512) Decimal512 {
	return Decimal512(sub512x(u512x(d), u512x(other)))
}

// Mul returns d * other with fixed 64-digit scale.
func (d Decimal512) Mul(other Decimal512) Decimal512 {
	u := u512x(d)
	v := u512x(other)
	if isZero512(u) || isZero512(v) {
		return Decimal512{}
	}
	neg := isNeg512(u) != isNeg512(v)
	if isNeg512(u) {
		u = neg512(u)
	}
	if isNeg512(v) {
		v = neg512(v)
	}
	p := mul512(u, v)
	q, _ := divMod1024By512(p, scale512)
	if neg {
		q = neg512(q)
	}
	return Decimal512(q)
}

// Div returns d / other with fixed 64-digit scale.
//
// If other is zero, it returns d unchanged.
func (d Decimal512) Div(other Decimal512) Decimal512 {
	u := u512x(d)
	v := u512x(other)
	if isZero512(v) {
		return d
	}
	if isZero512(u) {
		return Decimal512{}
	}
	neg := isNeg512(u) != isNeg512(v)
	if isNeg512(u) {
		u = neg512(u)
	}
	if isNeg512(v) {
		v = neg512(v)
	}
	p := mul512(u, scale512)
	q, _ := divMod1024By512(p, v)
	if neg {
		q = neg512(q)
	}
	return Decimal512(q)
}

// Mod returns d % other using truncation toward zero.
//
// If other is zero, it returns d unchanged.
func (d Decimal512) Mod(other Decimal512) Decimal512 {
	u := u512x(d)
	v := u512x(other)
	if isZero512(v) {
		return d
	}
	if isZero512(u) {
		return Decimal512{}
	}
	neg := isNeg512(u)
	if neg {
		u = neg512(u)
	}
	if isNeg512(v) {
		v = neg512(v)
	}
	_, r := divMod512By512(u, v)
	if neg {
		r = neg512(r)
	}
	return Decimal512(r)
}

// Pow returns d raised to an integer power specified by other.
//
// The exponent is truncated toward zero to an int64. Negative exponents use Inv().
func (d Decimal512) Pow(other Decimal512) Decimal512 {
	if other.IsZero() {
		return New512FromInt(1)
	}
	trunc := other.Truncate(0)
	exp, _ := trunc.Int64()
	if exp == 0 {
		return New512FromInt(1)
	}
	negExp := exp < 0
	if negExp {
		exp = -exp
	}
	result := New512FromInt(1)
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
func (d Decimal512) Sqrt() Decimal512 {
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
	gd, err := New512FromFloat(guess)
	if err != nil || gd.IsZero() {
		gd = New512FromInt(1)
	}
	prev := gd
	for i := 0; i < 32; i++ {
		inv := d.Div(gd)
		gd = gd.Add(inv)
		gd = divDecimal512ByUint64(gd, 2)
		if gd.Equal(prev) {
			break
		}
		prev = gd
	}
	return gd
}

// Exp returns e^d using range reduction and a Taylor series.
func (d Decimal512) Exp() Decimal512 {
	if d.IsZero() {
		return New512FromInt(1)
	}
	k := d.Div(constLn2_512).Round(0)
	kInt, _ := k.Int64()
	kLn2 := New512FromInt(kInt).Mul(constLn2_512)
	r := d.Sub(kLn2)

	term := New512FromInt(1)
	sum := term
	for i := uint64(1); i <= 96; i++ {
		term = term.Mul(r)
		term = divDecimal512ByUint64(term, i)
		if term.IsZero() {
			break
		}
		sum = sum.Add(term)
	}

	if kInt > 0 {
		sum = Decimal512(shl512(u512x(sum), uint(kInt)))
	} else if kInt < 0 {
		sum = Decimal512(shr512(u512x(sum), uint(-kInt)))
	}
	return sum
}

// Log returns the natural logarithm of d.
//
// For d <= 0, it returns d unchanged.
func (d Decimal512) Log() Decimal512 {
	if !d.IsPositive() {
		return d
	}
	u := u512x(d)
	k := int64(bitLen512u(u)) - int64(bitLen512u(scale512))
	var mScaled u512x
	if k >= 0 {
		mScaled = shr512(u, uint(k))
	} else {
		mScaled = shl512(u, uint(-k))
	}
	m := Decimal512(mScaled)
	one := New512FromInt(1)
	mMinus := m.Sub(one)
	mPlus := m.Add(one)
	t := mMinus.Div(mPlus)
	t2 := t.Mul(t)
	term := t
	sum := t
	for i := uint64(3); i <= 199; i += 2 {
		term = term.Mul(t2)
		add := divDecimal512ByUint64(term, i)
		if add.IsZero() {
			break
		}
		sum = sum.Add(add)
	}
	lnm := sum.Mul(New512FromInt(2))
	kLn2 := New512FromInt(k).Mul(constLn2_512)
	return lnm.Add(kLn2)
}

// Log2 returns the base-2 logarithm of d.
//
// For d <= 0, it returns d unchanged.
func (d Decimal512) Log2() Decimal512 {
	ln := d.Log()
	if !d.IsPositive() {
		return ln
	}
	return ln.Div(constLn2_512)
}

// Log10 returns the base-10 logarithm of d.
//
// For d <= 0, it returns d unchanged.
func (d Decimal512) Log10() Decimal512 {
	ln := d.Log()
	if !d.IsPositive() {
		return ln
	}
	return ln.Div(constLn10_512)
}

// EncodeBinary encodes the raw 512-bit value into 64 bytes (little-endian).
func (d Decimal512) EncodeBinary() ([]byte, error) {
	var out [64]byte
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint64(out[i*8:(i+1)*8], d[i])
	}
	return out[:], nil
}

// AppendBinary appends the raw 512-bit value as 64 bytes (little-endian) to dst.
func (d Decimal512) AppendBinary(dst []byte) []byte {
	var out [64]byte
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint64(out[i*8:(i+1)*8], d[i])
	}
	return append(dst, out[:]...)
}

// New512FromBinary decodes a 64-byte little-endian binary representation.
//
// Precision rules are applied after decoding.
func New512FromBinary(b []byte) (Decimal512, error) {
	if len(b) != 64 {
		return Decimal512{}, errInvalidBinaryLen
	}
	var out u512x
	for i := 0; i < 8; i++ {
		out[i] = binary.LittleEndian.Uint64(b[i*8 : (i+1)*8])
	}
	return Decimal512(applyPrecision512(out)), nil
}

// EncodeJSON encodes the decimal as a JSON string.
func (d Decimal512) EncodeJSON() ([]byte, error) {
	s := d.String()
	buf := make([]byte, 0, len(s)+2)
	buf = strconv.AppendQuote(buf, s)
	return buf, nil
}

// AppendJSON appends the decimal as a JSON string to dst.
func (d Decimal512) AppendJSON(dst []byte) []byte {
	dst = append(dst, '"')
	dst = d.AppendString(dst)
	return append(dst, '"')
}

// New512FromJSON decodes a JSON string or number into a Decimal512.
func New512FromJSON(b []byte) (Decimal512, error) {
	start, end := trimSpaceBytes(b)
	if start >= end {
		return Decimal512{}, errInvalidJSONDecimal
	}
	if b[start] == '"' {
		if end-start < 2 || b[end-1] != '"' {
			return Decimal512{}, errInvalidJSONDecimal
		}
		for i := start + 1; i < end-1; i++ {
			if b[i] == '\\' || b[i] < 0x20 {
				return Decimal512{}, errInvalidJSONDecimal
			}
		}
		u, err := parseDecimalBytes512(b[start+1 : end-1])
		if err != nil {
			return Decimal512{}, err
		}
		return Decimal512(u), nil
	}
	u, err := parseDecimalBytes512(b[start:end])
	if err != nil {
		return Decimal512{}, err
	}
	return Decimal512(u), nil
}

// truncateWithMode512 is an internal helper.
func (d Decimal512) truncateWithMode512(n int, mode roundMode) Decimal512 {
	if n > scaleDigits512 {
		return d
	}
	if n <= -scaleDigits512 {
		return Decimal512{}
	}
	if n == scaleDigits512 {
		return d
	}
	u := u512x(d)
	if isZero512(u) {
		return d
	}
	neg := isNeg512(u)
	if neg {
		u = neg512(u)
	}
	factor := pow10Value512(int64(scaleDigits512 - n))
	q, r := divMod512By512(u, factor)
	if !isZero512(r) {
		switch mode {
		case roundModeAwayFromZero:
			q = add512(q, u512x{1})
		case roundModeBanker:
			cmp := cmp512u(add512(r, r), factor)
			if cmp > 0 || (cmp == 0 && (q[0]&1) == 1) {
				q = add512(q, u512x{1})
			}
		case roundModeCeil:
			if !neg {
				q = add512(q, u512x{1})
			}
		case roundModeFloor:
			if neg {
				q = add512(q, u512x{1})
			}
		}
	}
	res := mul512(q, factor)
	out := lower512(res)
	if neg {
		out = neg512(out)
	}
	return Decimal512(out)
}

// stringFixed512 is an internal helper.
func (d Decimal512) stringFixed512(n int) string {
	u := u512x(d)
	if isZero512(u) {
		if n == 0 {
			return "0"
		}
		return "0." + repeatZero(n)
	}
	neg := isNeg512(u)
	if neg {
		u = neg512(u)
	}
	q, r := divMod512By512(u, scale512)
	intStr := u512ToDecimal(q)
	frac := u512ToDecimalFixed(r, scaleDigits512)
	if n < scaleDigits512 {
		frac = frac[:n]
	}
	if neg {
		return "-" + intStr + "." + frac
	}
	return intStr + "." + frac
}

// buildPow10_512 is an internal helper.
func buildPow10_512() [65]u512x {
	var p [65]u512x
	p[0] = u512x{1}
	for i := 1; i < len(p); i++ {
		p[i] = mul512ByUint64(p[i-1], 10)
	}
	return p
}

// pow10Mod512 is an internal helper.
func pow10Mod512(n int64) u512x {
	if n <= 0 {
		return u512x{1}
	}
	if n <= 64 {
		return pow10_512[n]
	}
	result := u512x{1}
	base := u512x{10}
	for n > 0 {
		if (n & 1) == 1 {
			result = mul512Lo(result, base)
		}
		base = mul512Lo(base, base)
		n >>= 1
	}
	return result
}

// pow10Value512 is an internal helper.
func pow10Value512(n int64) u512x {
	if n <= 0 {
		return u512x{1}
	}
	if n <= 64 {
		return pow10_512[n]
	}
	if n > maxDecimalDigits512 {
		return u512x{}
	}
	return pow10Mod512(n)
}

// parseDecimalString512 is an internal helper.
func parseDecimalString512(s string) (u512x, error) {
	start, end := trimSpaceString(s)
	if start >= end {
		return u512x{}, errInvalidDecimal
	}
	idx := start
	sign := 1
	if s[idx] == '+' {
		idx++
	} else if s[idx] == '-' {
		sign = -1
		idx++
	}
	var val u512x
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
				return u512x{}, errInvalidDecimal
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
			return u512x{}, errInvalidDecimal
		}
		sawDigit = true
		val = mul512ByUint64(val, 10)
		val = add512(val, u512x{uint64(c - '0')})
		if sawDot {
			fracDigits++
		}
		idx++
	}
	if !sawDigit {
		return u512x{}, errInvalidDecimal
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
				return u512x{}, errInvalidDecimal
			}
			if exp < (1 << 60) {
				exp = exp*10 + int64(c-'0')
			}
			idx++
		}
		if startExp == idx {
			return u512x{}, errInvalidDecimal
		}
		exp *= expSign
	}
	shift := exp - fracDigits + scaleDigits512
	if shift >= 0 {
		p := mul512(val, pow10Mod512(shift))
		val = lower512(p)
	} else {
		factor := pow10Value512(-shift)
		if isZero512(factor) {
			val = u512x{}
		} else {
			val = divByU512Trunc(val, factor)
		}
	}
	if sign < 0 {
		val = neg512(val)
	}
	return applyPrecision512(val), nil
}

// parseDecimalBytes512 is an internal helper.
func parseDecimalBytes512(b []byte) (u512x, error) {
	start, end := trimSpaceBytes(b)
	if start >= end {
		return u512x{}, errInvalidDecimal
	}
	idx := start
	sign := 1
	if b[idx] == '+' {
		idx++
	} else if b[idx] == '-' {
		sign = -1
		idx++
	}
	var val u512x
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
				return u512x{}, errInvalidDecimal
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
			return u512x{}, errInvalidDecimal
		}
		sawDigit = true
		val = mul512ByUint64(val, 10)
		val = add512(val, u512x{uint64(c - '0')})
		if sawDot {
			fracDigits++
		}
		idx++
	}
	if !sawDigit {
		return u512x{}, errInvalidDecimal
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
				return u512x{}, errInvalidDecimal
			}
			if exp < (1 << 60) {
				exp = exp*10 + int64(c-'0')
			}
			idx++
		}
		if startExp == idx {
			return u512x{}, errInvalidDecimal
		}
		exp *= expSign
	}
	shift := exp - fracDigits + scaleDigits512
	if shift >= 0 {
		p := mul512(val, pow10Mod512(shift))
		val = lower512(p)
	} else {
		factor := pow10Value512(-shift)
		if isZero512(factor) {
			val = u512x{}
		} else {
			val = divByU512Trunc(val, factor)
		}
	}
	if sign < 0 {
		val = neg512(val)
	}
	return applyPrecision512(val), nil
}

// applyPrecision512 drops integer digits beyond 64 and fractional digits beyond 64.
// It assumes the input is scaled by 10^64.
func applyPrecision512(u u512x) u512x {
	if isZero512(u) {
		return u
	}
	neg := isNeg512(u)
	if neg {
		u = neg512(u)
	}
	intPart, fracPart := divMod512By512(u, scale512)
	_, intRem := divMod512By512(intPart, pow10Value512(scaleDigits512))
	intPart = intRem
	// scaleDigits512 == decimal512PrecisionDigits, so fractional trimming is a no-op.
	raw := add512(lower512(mul512(intPart, scale512)), fracPart)
	if neg {
		raw = neg512(raw)
	}
	return raw
}

// divDecimal512ByUint64 is an internal helper.
func divDecimal512ByUint64(d Decimal512, n uint64) Decimal512 {
	if n == 0 {
		return d
	}
	u := u512x(d)
	if isZero512(u) {
		return d
	}
	neg := isNeg512(u)
	if neg {
		u = neg512(u)
	}
	q, _ := divMod512ByUint64(u, n)
	if neg {
		q = neg512(q)
	}
	return Decimal512(q)
}

// u512FromFloatTrunc is an internal helper.
func u512FromFloatTrunc(v float64) u512x {
	if v <= 0 {
		return u512x{}
	}
	bits64 := math.Float64bits(v)
	exp := int((bits64>>52)&0x7ff) - 1023
	mant := bits64 & ((uint64(1) << 52) - 1)
	mant |= uint64(1) << 52
	if exp < 0 {
		return u512x{}
	}
	shift := exp - 52
	if shift >= 512 {
		return u512x{}
	}
	u := u512x{mant}
	if shift == 0 {
		return u
	}
	if shift > 0 {
		return shl512(u, uint(shift))
	}
	return shr512(u, uint(-shift))
}

// u512ToFloat is an internal helper.
func u512ToFloat(u u512x) float64 {
	if isZero512(u) {
		return 0
	}
	f := 0.0
	for i := 7; i >= 0; i-- {
		f = f*math.Exp2(64) + float64(u[i])
		if i == 0 {
			break
		}
	}
	return f
}

// appendU512Decimal is an internal helper.
func appendU512Decimal(dst []byte, v u512x) []byte {
	if isZero512(v) {
		return append(dst, '0')
	}
	const base = uint64(1_000_000_000_000_000_000) // 1e18
	var parts [10]uint64
	n := 0
	for !isZero512(v) {
		var rem uint64
		v, rem = divMod512ByUint64(v, base)
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

// appendU512DecimalFixed is an internal helper.
func appendU512DecimalFixed(dst []byte, v u512x, width int) []byte {
	if width <= 0 {
		return dst
	}
	var tmp [96]byte
	num := appendU512Decimal(tmp[:0], v)
	if len(num) >= width {
		return append(dst, num[len(num)-width:]...)
	}
	for i := 0; i < width-len(num); i++ {
		dst = append(dst, '0')
	}
	return append(dst, num...)
}

// u512ToDecimal is an internal helper.
func u512ToDecimal(v u512x) string {
	if isZero512(v) {
		return "0"
	}
	const base = uint64(1_000_000_000_000_000_000) // 1e18
	var parts [10]uint64
	n := 0
	for !isZero512(v) {
		var rem uint64
		v, rem = divMod512ByUint64(v, base)
		parts[n] = rem
		n++
	}
	buf := make([]byte, 0, 96)
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

// u512ToDecimalFixed is an internal helper.
func u512ToDecimalFixed(v u512x, width int) string {
	if width <= 0 {
		return ""
	}
	s := u512ToDecimal(v)
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

// divByU512Trunc is an internal helper.
func divByU512Trunc(v u512x, d u512x) u512x {
	if isZero512(d) {
		return v
	}
	neg := isNeg512(v)
	if neg {
		v = neg512(v)
	}
	q, _ := divMod512By512(v, d)
	if neg {
		q = neg512(q)
	}
	return q
}

// isZero512 is an internal helper.
func isZero512(u u512x) bool {
	return u[0]|u[1]|u[2]|u[3]|u[4]|u[5]|u[6]|u[7] == 0
}

// isNeg512 is an internal helper.
func isNeg512(u u512x) bool {
	return (u[7]>>63)&1 == 1
}

// add512 is an internal helper.
func add512(a, b u512x) u512x {
	var out u512x
	var c uint64
	out[0], c = bits.Add64(a[0], b[0], 0)
	out[1], c = bits.Add64(a[1], b[1], c)
	out[2], c = bits.Add64(a[2], b[2], c)
	out[3], c = bits.Add64(a[3], b[3], c)
	out[4], c = bits.Add64(a[4], b[4], c)
	out[5], c = bits.Add64(a[5], b[5], c)
	out[6], c = bits.Add64(a[6], b[6], c)
	out[7], _ = bits.Add64(a[7], b[7], c)
	return out
}

// sub512 is an internal helper.
func sub512x(a, b u512x) u512x {
	var out u512x
	var c uint64
	out[0], c = bits.Sub64(a[0], b[0], 0)
	out[1], c = bits.Sub64(a[1], b[1], c)
	out[2], c = bits.Sub64(a[2], b[2], c)
	out[3], c = bits.Sub64(a[3], b[3], c)
	out[4], c = bits.Sub64(a[4], b[4], c)
	out[5], c = bits.Sub64(a[5], b[5], c)
	out[6], c = bits.Sub64(a[6], b[6], c)
	out[7], _ = bits.Sub64(a[7], b[7], c)
	return out
}

// neg512 is an internal helper.
func neg512(a u512x) u512x {
	var out u512x
	for i := 0; i < 8; i++ {
		out[i] = ^a[i]
	}
	var c uint64
	out[0], c = bits.Add64(out[0], 1, 0)
	for i := 1; i < 8; i++ {
		out[i], c = bits.Add64(out[i], 0, c)
	}
	return out
}

// cmp512u is an internal helper.
func cmp512u(a, b u512x) int {
	for i := 7; i >= 0; i-- {
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

// cmp512Signed is an internal helper.
func cmp512Signed(a, b u512x) int {
	na := isNeg512(a)
	nb := isNeg512(b)
	if na && !nb {
		return -1
	}
	if !na && nb {
		return 1
	}
	return cmp512u(a, b)
}

// mul512ByUint64 is an internal helper.
func mul512ByUint64(a u512x, m uint64) u512x {
	var out u512x
	var carry uint64
	for i := 0; i < 8; i++ {
		hi, lo := bits.Mul64(a[i], m)
		lo, c := bits.Add64(lo, carry, 0)
		out[i] = lo
		carry = hi + c
	}
	return out
}

// divMod512ByUint64 is an internal helper.
func divMod512ByUint64(a u512x, d uint64) (u512x, uint64) {
	var q u512x
	var r uint64
	for i := 7; i >= 0; i-- {
		q[i], r = bits.Div64(r, a[i], d)
		if i == 0 {
			break
		}
	}
	return q, r
}

// mul512 is an internal helper.
func mul512(a, b u512x) u1024 {
	var p u1024
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			hi, lo := bits.Mul64(a[i], b[j])
			k := i + j
			var c uint64
			p[k], c = bits.Add64(p[k], lo, 0)
			p[k+1], c = bits.Add64(p[k+1], hi, c)
			idx := k + 2
			for c != 0 && idx < 16 {
				p[idx], c = bits.Add64(p[idx], 0, c)
				idx++
			}
		}
	}
	return p
}

// mul512Lo is an internal helper.
func mul512Lo(a, b u512x) u512x {
	return lower512(mul512(a, b))
}

// lower512 is an internal helper.
func lower512(p u1024) u512x {
	return u512x{p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7]}
}

// divMod512By512 is an internal helper.
func divMod512By512(n u512x, d u512x) (u512x, u512x) {
	if isZero512(d) {
		return u512x{}, u512x{}
	}
	var n1024 u1024
	for i := 0; i < 8; i++ {
		n1024[i] = n[i]
	}
	q, r := divMod1024By512(n1024, d)
	return q, r
}

// divMod1024By512 is an internal helper.
func divMod1024By512(n u1024, d u512x) (u512x, u512x) {
	if isZero512(d) {
		return u512x{}, u512x{}
	}
	if isZeroU1024(n) {
		return u512x{}, u512x{}
	}
	nBits := bitLen1024(n)
	dBits := bitLen512u(d)
	if nBits < dBits {
		return u512x{}, lower512(n)
	}
	shift := nBits - dBits
	var q u512x
	rem := n
	dShift := shl512To1024(d, uint(shift))
	for i := shift; i >= 0; i-- {
		if cmp1024(rem, dShift) >= 0 {
			rem = sub1024(rem, dShift)
			if i < 512 {
				q[int(i/64)] |= 1 << uint(i%64)
			}
		}
		if i == 0 {
			break
		}
		dShift = shr1_1024(dShift)
	}
	return q, lower512(rem)
}

// isZeroU1024 is an internal helper.
func isZeroU1024(u u1024) bool {
	var acc uint64
	for i := 0; i < 16; i++ {
		acc |= u[i]
	}
	return acc == 0
}

// bitLen512u is an internal helper.
func bitLen512u(u u512x) int {
	for i := 7; i >= 0; i-- {
		if u[i] != 0 {
			return i*64 + bits.Len64(u[i])
		}
		if i == 0 {
			break
		}
	}
	return 0
}

// bitLen1024 is an internal helper.
func bitLen1024(u u1024) int {
	for i := 15; i >= 0; i-- {
		if u[i] != 0 {
			return i*64 + bits.Len64(u[i])
		}
		if i == 0 {
			break
		}
	}
	return 0
}

// cmp1024 is an internal helper.
func cmp1024(a, b u1024) int {
	for i := 15; i >= 0; i-- {
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

// sub1024 is an internal helper.
func sub1024(a, b u1024) u1024 {
	var out u1024
	var c uint64
	out[0], c = bits.Sub64(a[0], b[0], 0)
	out[1], c = bits.Sub64(a[1], b[1], c)
	out[2], c = bits.Sub64(a[2], b[2], c)
	out[3], c = bits.Sub64(a[3], b[3], c)
	out[4], c = bits.Sub64(a[4], b[4], c)
	out[5], c = bits.Sub64(a[5], b[5], c)
	out[6], c = bits.Sub64(a[6], b[6], c)
	out[7], c = bits.Sub64(a[7], b[7], c)
	out[8], c = bits.Sub64(a[8], b[8], c)
	out[9], c = bits.Sub64(a[9], b[9], c)
	out[10], c = bits.Sub64(a[10], b[10], c)
	out[11], c = bits.Sub64(a[11], b[11], c)
	out[12], c = bits.Sub64(a[12], b[12], c)
	out[13], c = bits.Sub64(a[13], b[13], c)
	out[14], c = bits.Sub64(a[14], b[14], c)
	out[15], _ = bits.Sub64(a[15], b[15], c)
	return out
}

// shl512 is an internal helper.
func shl512(u u512x, shift uint) u512x {
	if shift == 0 {
		return u
	}
	if shift >= 512 {
		return u512x{}
	}
	wordShift := shift / 64
	bitShift := shift % 64
	var out u512x
	for i := 0; i < 8; i++ {
		dst := i + int(wordShift)
		if dst >= 8 {
			continue
		}
		out[dst] |= u[i] << bitShift
		if bitShift != 0 && dst+1 < 8 {
			out[dst+1] |= u[i] >> (64 - bitShift)
		}
	}
	return out
}

// shr512 is an internal helper.
func shr512(u u512x, shift uint) u512x {
	if shift == 0 {
		return u
	}
	if shift >= 512 {
		return u512x{}
	}
	wordShift := shift / 64
	bitShift := shift % 64
	var out u512x
	for i := int(wordShift); i < 8; i++ {
		src := i
		dst := i - int(wordShift)
		out[dst] |= u[src] >> bitShift
		if bitShift != 0 && src+1 < 8 {
			out[dst] |= u[src+1] << (64 - bitShift)
		}
	}
	return out
}

// shl512To1024 is an internal helper.
func shl512To1024(u u512x, shift uint) u1024 {
	if shift == 0 {
		return u1024{u[0], u[1], u[2], u[3], u[4], u[5], u[6], u[7]}
	}
	if shift >= 1024 {
		return u1024{}
	}
	wordShift := shift / 64
	bitShift := shift % 64
	var out u1024
	for i := 0; i < 8; i++ {
		dst := i + int(wordShift)
		if dst >= 16 {
			continue
		}
		out[dst] |= u[i] << bitShift
		if bitShift != 0 && dst+1 < 16 {
			out[dst+1] |= u[i] >> (64 - bitShift)
		}
	}
	return out
}

// shr1_1024 is an internal helper.
func shr1_1024(u u1024) u1024 {
	var out u1024
	var carry uint64
	for i := 15; i >= 0; i-- {
		out[i] = (u[i] >> 1) | (carry << 63)
		carry = u[i] & 1
		if i == 0 {
			break
		}
	}
	return out
}

// u512FromInt64 is an internal helper.
func u512FromInt64(v int64) u512x {
	if v >= 0 {
		return u512x{uint64(v)}
	}
	uv := uint64(v)
	return u512x{uv, ^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0)}
}

// scale512Squared is an internal helper.
func scale512Squared() u512x {
	return u512x{0x0, 0x0, 0x03df99092e953e01, 0x2374e42f0f1538fd, 0xc404dc08d3cff5ec, 0xa6337f19bccdb0da, 0x0000024ee91f2603, 0x0}
}
