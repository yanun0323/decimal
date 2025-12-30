// All overflow is truncated to 256-bit two's-complement arithmetic.
package decimal

import (
	"encoding/binary"
	"errors"
	"math"
	"math/bits"
	"strconv"
)

const (
	scaleDigits      = 32
	maxDecimalDigits = 77 // 10^77 < 2^256 < 10^78
)

var (
	errInvalidDecimal     = errors.New("invalid decimal256")
	errInvalidBinaryLen   = errors.New("invalid binary length")
	errInvalidFloat       = errors.New("invalid float64")
	errInvalidJSONDecimal = errors.New("invalid json decimal256")
	pow10U64              = [...]uint64{
		1,
		10,
		100,
		1000,
		10000,
		100000,
		1000000,
		10000000,
		100000000,
		1000000000,
		10000000000,
		100000000000,
		1000000000000,
		10000000000000,
		100000000000000,
		1000000000000000,
		10000000000000000,
		100000000000000000,
		1000000000000000000,
		10000000000000000000,
	}
)

// Decimal256 is a fixed-scale decimal stored as a 256-bit two's-complement integer.
//
// The scale is 10^32, so the numeric value is: raw / 10^32.
// The zero value represents 0 and is ready to use.
//
// Memory layout is fixed and little-endian across 4 uint64 words.
// This is required for binary compatibility.
type Decimal256 [4]uint64

// u256 is an internal unsigned 256-bit integer in little-endian word order.
type u256 [4]uint64

// u512 is an internal unsigned 512-bit integer in little-endian word order.
type u512 [8]uint64

var (
	scale = u256{0x85acef8100000000, 0x000004ee2d6d415b, 0x0, 0x0} // 10^32
	// 10^64, used by Inv().
	scaleSquared = u256{0x0, 0x6e38ed64bf6a1f01, 0xe93ff9f4daa797ed, 0x0000000000184f03}
	pow10        = buildPow10()
	// ln(2) scaled by 1e32.
	constLn2 = Decimal256(u256{0x797c31134d266499, 0x36adfeef0c4, 0x0, 0x0})
	// ln(10) scaled by 1e32.
	constLn10 = Decimal256(u256{0x686e1a5c5b723214, 0xb5a455ec490, 0x0, 0x0})
)

// New256 constructs a Decimal256 from integer and fractional parts.
//
// intPart keeps only the lowest 32 decimal digits (higher digits are dropped).
// decimalPart keeps only the highest 32 fractional digits (lower digits are dropped).
// decimalPart is interpreted as fractional digits with an implicit scale based on
// its decimal digit length (e.g. 987654321 -> 0.987654321). It is then scaled to
// 10^32 before combining with intPart. If decimalPart has more than 32 digits, it
// is truncated toward zero. The result is: intPart*10^32 + scaled(decimalPart),
// with two's-complement wrap on overflow.
func New256(intPart, decimalPart int64) Decimal256 {
	ip := mul256(u256FromInt64(intPart), scale)
	raw := lower256(ip)
	if decimalPart != 0 {
		abs, neg := absInt64(decimalPart)
		digits := decimalDigitsU64(abs)
		shift := scaleDigits - digits
		frac := u256{abs, 0, 0, 0}
		var scaled u256
		if shift == 0 {
			scaled = frac
		} else if shift < 0 {
			factor := pow10Value(int64(-shift))
			if isZero(factor) {
				scaled = u256{}
			} else {
				scaled = divByU256Trunc(frac, factor)
			}
		} else if shift < len(pow10U64) {
			scaled = mul256ByUint64(frac, pow10U64[shift])
		} else {
			p := mul256(frac, pow10Value(int64(shift)))
			scaled = lower256(p)
		}
		if neg {
			scaled = neg256(scaled)
		}
		raw = add256(raw, scaled)
	}
	return Decimal256(applyPrecision256(raw))
}

// New256FromString parses a decimal string with optional sign, dot, and exponent.
//
// It accepts leading/trailing ASCII whitespace and optional '_' separators.
// Exponent shifting is applied first, then integer digits beyond 32 are dropped and
// fractional digits beyond 32 are dropped. Excess fractional digits are truncated
// (toward zero) to the fixed 32-digit scale.
func New256FromString(s string) (Decimal256, error) {
	u, err := parseDecimalString(s)
	if err != nil {
		return Decimal256{}, err
	}
	return Decimal256(u), nil
}

// New256FromInt constructs a Decimal256 from an int64 integer value.
func New256FromInt(v int64) Decimal256 {
	p := mul256(u256FromInt64(v), scale)
	return Decimal256(applyPrecision256(lower256(p)))
}

// New256FromFloat converts a float64 to Decimal256 by truncating toward zero.
//
// NaN or Inf returns an error. Overflow wraps according to two's-complement truncation.
func New256FromFloat(v float64) (Decimal256, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return Decimal256{}, errInvalidFloat
	}
	if v == 0 {
		return Decimal256{}, nil
	}
	neg := v < 0
	if neg {
		v = -v
	}
	scaled := v * 1e32
	if math.IsInf(scaled, 0) {
		return Decimal256{}, errInvalidFloat
	}
	u := u256FromFloatTrunc(scaled)
	u = applyPrecision256(u)
	if neg {
		u = neg256(u)
	}
	return Decimal256(u), nil
}

// New256FromBinary decodes a 32-byte little-endian binary representation.
//
// Precision rules are applied after decoding.
func New256FromBinary(b []byte) (Decimal256, error) {
	if len(b) != 32 {
		return Decimal256{}, errInvalidBinaryLen
	}
	u := u256{
		binary.LittleEndian.Uint64(b[0:8]),
		binary.LittleEndian.Uint64(b[8:16]),
		binary.LittleEndian.Uint64(b[16:24]),
		binary.LittleEndian.Uint64(b[24:32]),
	}
	return Decimal256(applyPrecision256(u)), nil
}

// New256FromJSON decodes a JSON string or number into a Decimal256.
func New256FromJSON(b []byte) (Decimal256, error) {
	start, end := trimSpaceBytes(b)
	if start >= end {
		return Decimal256{}, errInvalidJSONDecimal
	}
	if b[start] == '"' {
		if end-start < 2 || b[end-1] != '"' {
			return Decimal256{}, errInvalidJSONDecimal
		}
		for i := start + 1; i < end-1; i++ {
			if b[i] == '\\' || b[i] < 0x20 {
				return Decimal256{}, errInvalidJSONDecimal
			}
		}
		u, err := parseDecimalBytes(b[start+1 : end-1])
		if err != nil {
			return Decimal256{}, err
		}
		return Decimal256(u), nil
	}
	u, err := parseDecimalBytes(b[start:end])
	if err != nil {
		return Decimal256{}, err
	}
	return Decimal256(u), nil
}

// Int64 returns the integer and fractional parts as int64 values.
//
// Both parts are truncated to int64 with two's-complement wrap if out of range.
// The fractional part is returned in 10^32 base units.
func (d Decimal256) Int64() (intPart, decimalPart int64) {
	u := u256(d)
	if isZero(u) {
		return 0, 0
	}
	neg := isNeg(u)
	abs := u
	if neg {
		abs = neg256(abs)
	}
	q, r := divMod256By256(abs, scale)
	qi := int64(q[0])
	ri := int64(r[0])
	if neg {
		qi = -qi
		ri = -ri
	}
	return qi, ri
}

// Float64 converts Decimal256 to float64.
//
// Precision is limited by float64; large values may overflow to Inf.
func (d Decimal256) Float64() float64 {
	u := u256(d)
	if isZero(u) {
		return 0
	}
	neg := isNeg(u)
	if neg {
		u = neg256(u)
	}
	f := u256ToFloat(u)
	f = f / 1e32
	if neg {
		f = -f
	}
	return f
}

// String returns the shortest decimal representation without trailing zeros.
func (d Decimal256) String() string {
	u := u256(d)
	if isZero(u) {
		return "0"
	}
	neg := isNeg(u)
	if neg {
		u = neg256(u)
	}
	q, r := divMod256By256(u, scale)
	intStr := u256ToDecimal(q)
	if isZero(r) {
		if neg {
			return "-" + intStr
		}
		return intStr
	}
	frac := u256ToDecimalFixed(r, scaleDigits)
	frac = trimRightZeros(frac)
	if neg {
		return "-" + intStr + "." + frac
	}
	return intStr + "." + frac
}

// StringFixed returns a decimal string with exactly n fractional digits.
//
// If n > 32 it is truncated to 32. If n <= 0, no fractional part is shown.
func (d Decimal256) StringFixed(n int) string {
	if n > scaleDigits {
		n = scaleDigits
	}
	if n <= 0 {
		u := u256(d)
		if isZero(u) {
			return "0"
		}
		neg := isNeg(u)
		if neg {
			u = neg256(u)
		}
		q, _ := divMod256By256(u, scale)
		intStr := u256ToDecimal(q)
		if neg {
			return "-" + intStr
		}
		return intStr
	}
	return d.stringFixedN(n)
}

// AppendString appends the shortest decimal representation without trailing zeros to dst.
func (d Decimal256) AppendString(dst []byte) []byte {
	u := u256(d)
	if isZero(u) {
		return append(dst, '0')
	}
	neg := isNeg(u)
	if neg {
		u = neg256(u)
	}
	q, r := divMod256By256(u, scale)
	if neg {
		dst = append(dst, '-')
	}
	dst = appendU256Decimal(dst, q)
	if isZero(r) {
		return dst
	}
	dst = append(dst, '.')
	fracStart := len(dst)
	dst = appendU256DecimalFixed(dst, r, scaleDigits)
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
// If n > 32 it is truncated to 32. If n <= 0, no fractional part is appended.
func (d Decimal256) AppendStringFixed(dst []byte, n int) []byte {
	if n > scaleDigits {
		n = scaleDigits
	}
	u := u256(d)
	if isZero(u) {
		if n <= 0 {
			return append(dst, '0')
		}
		dst = append(dst, '0', '.')
		for i := 0; i < n; i++ {
			dst = append(dst, '0')
		}
		return dst
	}
	neg := isNeg(u)
	if neg {
		u = neg256(u)
	}
	q, r := divMod256By256(u, scale)
	if neg {
		dst = append(dst, '-')
	}
	dst = appendU256Decimal(dst, q)
	if n > 0 {
		dst = append(dst, '.')
		var fracBuf [64]byte
		frac := appendU256DecimalFixed(fracBuf[:0], r, scaleDigits)
		if n < scaleDigits {
			frac = frac[:n]
		}
		dst = append(dst, frac...)
	}
	return dst
}

// IsZero reports whether the value is exactly zero.
func (d Decimal256) IsZero() bool {
	return isZero(u256(d))
}

// IsPositive reports whether the value is greater than zero.
func (d Decimal256) IsPositive() bool {
	u := u256(d)
	return !isZero(u) && !isNeg(u)
}

// IsNegative reports whether the value is less than zero.
func (d Decimal256) IsNegative() bool {
	u := u256(d)
	return !isZero(u) && isNeg(u)
}

// Sign returns 0 if zero, 1 if positive, and 2 if negative.
func (d Decimal256) Sign() int {
	u := u256(d)
	if isZero(u) {
		return 0
	}
	if isNeg(u) {
		return 2
	}
	return 1
}

// Neg returns the arithmetic negation of d.
func (d Decimal256) Neg() Decimal256 {
	return Decimal256(neg256(u256(d)))
}

// Inv returns the multiplicative inverse (1/d).
//
// For zero, it returns zero.
func (d Decimal256) Inv() Decimal256 {
	u := u256(d)
	if isZero(u) {
		return d
	}
	neg := isNeg(u)
	if neg {
		u = neg256(u)
	}
	q, _ := divMod256By256(scaleSquared, u)
	if neg {
		q = neg256(q)
	}
	return Decimal256(q)
}

// Abs returns the absolute value of d.
func (d Decimal256) Abs() Decimal256 {
	u := u256(d)
	if isNeg(u) {
		u = neg256(u)
	}
	return Decimal256(u)
}

// Truncate truncates to n fractional digits (banker-friendly truncation toward zero).
//
// If n > 32, it returns d unchanged. If n <= -32, it returns zero.
func (d Decimal256) Truncate(n int) Decimal256 {
	return d.truncateWithMode(n, roundModeTowardZero)
}

// Shift moves the decimal point by n digits.
//
// Positive n shifts left (multiply by 10^n), negative n shifts right (divide by 10^-n).
// If n > 32, it returns d unchanged. If n <= -32, it returns zero.
func (d Decimal256) Shift(n int) Decimal256 {
	if n > scaleDigits {
		return d
	}
	if n <= -scaleDigits {
		return Decimal256{}
	}
	if n == 0 {
		return d
	}
	if n > 0 {
		factor := pow10Mod(int64(n))
		p := mul256(u256(d), factor)
		return Decimal256(lower256(p))
	}
	factor := pow10Value(int64(-n))
	return Decimal256(divByU256Trunc(u256(d), factor))
}

// Round rounds to n fractional digits using banker's rounding.
//
// If n > 32, it returns d unchanged. If n <= -32, it returns zero.
func (d Decimal256) Round(n int) Decimal256 {
	return d.truncateWithMode(n, roundModeBanker)
}

// RoundAwayFromZero rounds to n fractional digits, away from zero.
//
// If n > 32, it returns d unchanged. If n <= -32, it returns zero.
func (d Decimal256) RoundAwayFromZero(n int) Decimal256 {
	return d.truncateWithMode(n, roundModeAwayFromZero)
}

// RoundTowardToZero truncates to n fractional digits toward zero.
//
// If n > 32, it returns d unchanged. If n <= -32, it returns zero.
func (d Decimal256) RoundTowardToZero(n int) Decimal256 {
	return d.truncateWithMode(n, roundModeTowardZero)
}

// Ceil rounds toward positive infinity with n fractional digits.
//
// If n > 32, it returns d unchanged. If n <= -32, it returns zero.
func (d Decimal256) Ceil(n int) Decimal256 {
	return d.truncateWithMode(n, roundModeCeil)
}

// Floor rounds toward negative infinity with n fractional digits.
//
// If n > 32, it returns d unchanged. If n <= -32, it returns zero.
func (d Decimal256) Floor(n int) Decimal256 {
	return d.truncateWithMode(n, roundModeFloor)
}

// Equal reports whether d == other.
func (d Decimal256) Equal(other Decimal256) bool {
	return d == other
}

// GreaterThan reports whether d > other.
func (d Decimal256) GreaterThan(other Decimal256) bool {
	return cmp256Signed(u256(d), u256(other)) > 0
}

// LessThan reports whether d < other.
func (d Decimal256) LessThan(other Decimal256) bool {
	return cmp256Signed(u256(d), u256(other)) < 0
}

// GreaterOrEqual reports whether d >= other.
func (d Decimal256) GreaterOrEqual(other Decimal256) bool {
	return cmp256Signed(u256(d), u256(other)) >= 0
}

// LessOrEqual reports whether d <= other.
func (d Decimal256) LessOrEqual(other Decimal256) bool {
	return cmp256Signed(u256(d), u256(other)) <= 0
}

// Add returns d + other with 256-bit truncation on overflow.
func (d Decimal256) Add(other Decimal256) Decimal256 {
	return Decimal256(add256(u256(d), u256(other)))
}

// Sub returns d - other with 256-bit truncation on overflow.
func (d Decimal256) Sub(other Decimal256) Decimal256 {
	return Decimal256(sub256(u256(d), u256(other)))
}

// Mul returns d * other with fixed 32-digit scale.
func (d Decimal256) Mul(other Decimal256) Decimal256 {
	u := u256(d)
	v := u256(other)
	if isZero(u) || isZero(v) {
		return Decimal256{}
	}
	neg := isNeg(u) != isNeg(v)
	if isNeg(u) {
		u = neg256(u)
	}
	if isNeg(v) {
		v = neg256(v)
	}
	p := mul256(u, v)
	q, _ := divMod512By256(p, scale)
	if neg {
		q = neg256(q)
	}
	return Decimal256(q)
}

// Div returns d / other with fixed 32-digit scale.
//
// If other is zero, it returns d unchanged.
func (d Decimal256) Div(other Decimal256) Decimal256 {
	u := u256(d)
	v := u256(other)
	if isZero(v) {
		return d
	}
	if isZero(u) {
		return Decimal256{}
	}
	neg := isNeg(u) != isNeg(v)
	if isNeg(u) {
		u = neg256(u)
	}
	if isNeg(v) {
		v = neg256(v)
	}
	p := mul256(u, scale)
	q, _ := divMod512By256(p, v)
	if neg {
		q = neg256(q)
	}
	return Decimal256(q)
}

// Mod returns d % other using truncation toward zero.
//
// If other is zero, it returns d unchanged.
func (d Decimal256) Mod(other Decimal256) Decimal256 {
	u := u256(d)
	v := u256(other)
	if isZero(v) {
		return d
	}
	if isZero(u) {
		return Decimal256{}
	}
	neg := isNeg(u)
	if neg {
		u = neg256(u)
	}
	if isNeg(v) {
		v = neg256(v)
	}
	_, r := divMod256By256(u, v)
	if neg {
		r = neg256(r)
	}
	return Decimal256(r)
}

// Pow returns d raised to an integer power specified by other.
//
// The exponent is truncated toward zero to an int64. Negative exponents use Inv().
func (d Decimal256) Pow(other Decimal256) Decimal256 {
	if other.IsZero() {
		return New256FromInt(1)
	}
	trunc := other.Truncate(0)
	exp, _ := trunc.Int64()
	if exp == 0 {
		return New256FromInt(1)
	}
	negExp := exp < 0
	if negExp {
		exp = -exp
	}
	result := New256FromInt(1)
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
func (d Decimal256) Sqrt() Decimal256 {
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
	gd, err := New256FromFloat(guess)
	if err != nil || gd.IsZero() {
		gd = New256FromInt(1)
	}
	prev := gd
	for i := 0; i < 32; i++ {
		inv := d.Div(gd)
		gd = gd.Add(inv)
		gd = divDecimalByUint64(gd, 2)
		if gd.Equal(prev) {
			break
		}
		prev = gd
	}
	return gd
}

// Exp returns e^d using range reduction and a Taylor series.
func (d Decimal256) Exp() Decimal256 {
	if d.IsZero() {
		return New256FromInt(1)
	}
	k := d.Div(constLn2).Round(0)
	kInt, _ := k.Int64()
	kLn2 := New256FromInt(kInt).Mul(constLn2)
	r := d.Sub(kLn2)

	term := New256FromInt(1)
	sum := term
	for i := uint64(1); i <= 96; i++ {
		term = term.Mul(r)
		term = divDecimalByUint64(term, i)
		if term.IsZero() {
			break
		}
		sum = sum.Add(term)
	}

	if kInt > 0 {
		sum = Decimal256(shl256(u256(sum), uint(kInt)))
	} else if kInt < 0 {
		sum = Decimal256(shr256(u256(sum), uint(-kInt)))
	}
	return sum
}

// Log returns the natural logarithm of d.
//
// For d <= 0, it returns d unchanged.
func (d Decimal256) Log() Decimal256 {
	if !d.IsPositive() {
		return d
	}
	u := u256(d)
	k := int64(bitLen256(u)) - int64(bitLen256(scale))
	var mScaled u256
	if k >= 0 {
		mScaled = shr256(u, uint(k))
	} else {
		mScaled = shl256(u, uint(-k))
	}
	m := Decimal256(mScaled)
	one := New256FromInt(1)
	mMinus := m.Sub(one)
	mPlus := m.Add(one)
	t := mMinus.Div(mPlus)
	t2 := t.Mul(t)
	term := t
	sum := t
	for i := uint64(3); i <= 199; i += 2 {
		term = term.Mul(t2)
		add := divDecimalByUint64(term, i)
		if add.IsZero() {
			break
		}
		sum = sum.Add(add)
	}
	lnm := sum.Mul(New256FromInt(2))
	kLn2 := New256FromInt(k).Mul(constLn2)
	return lnm.Add(kLn2)
}

// Log2 returns the base-2 logarithm of d.
//
// For d <= 0, it returns d unchanged.
func (d Decimal256) Log2() Decimal256 {
	ln := d.Log()
	if !d.IsPositive() {
		return ln
	}
	return ln.Div(constLn2)
}

// Log10 returns the base-10 logarithm of d.
//
// For d <= 0, it returns d unchanged.
func (d Decimal256) Log10() Decimal256 {
	ln := d.Log()
	if !d.IsPositive() {
		return ln
	}
	return ln.Div(constLn10)
}

// EncodeBinary encodes the raw 256-bit value into 32 bytes (little-endian).
func (d Decimal256) EncodeBinary() ([]byte, error) {
	var out [32]byte
	binary.LittleEndian.PutUint64(out[0:8], d[0])
	binary.LittleEndian.PutUint64(out[8:16], d[1])
	binary.LittleEndian.PutUint64(out[16:24], d[2])
	binary.LittleEndian.PutUint64(out[24:32], d[3])
	return out[:], nil
}

// AppendBinary appends the raw 256-bit value as 32 bytes (little-endian) to dst.
func (d Decimal256) AppendBinary(dst []byte) []byte {
	var out [32]byte
	binary.LittleEndian.PutUint64(out[0:8], d[0])
	binary.LittleEndian.PutUint64(out[8:16], d[1])
	binary.LittleEndian.PutUint64(out[16:24], d[2])
	binary.LittleEndian.PutUint64(out[24:32], d[3])
	return append(dst, out[:]...)
}

// EncodeJSON encodes the decimal as a JSON string.
func (d Decimal256) EncodeJSON() ([]byte, error) {
	s := d.String()
	buf := make([]byte, 0, len(s)+2)
	buf = strconv.AppendQuote(buf, s)
	return buf, nil
}

// AppendJSON appends the decimal as a JSON string to dst.
func (d Decimal256) AppendJSON(dst []byte) []byte {
	dst = append(dst, '"')
	dst = d.AppendString(dst)
	return append(dst, '"')
}

type roundMode int

const (
	roundModeTowardZero roundMode = iota
	roundModeAwayFromZero
	roundModeBanker
	roundModeCeil
	roundModeFloor
)

// truncateWithMode is an internal helper.
func (d Decimal256) truncateWithMode(n int, mode roundMode) Decimal256 {
	if n > scaleDigits {
		return d
	}
	if n <= -scaleDigits {
		return Decimal256{}
	}
	if n == scaleDigits {
		return d
	}
	u := u256(d)
	if isZero(u) {
		return d
	}
	neg := isNeg(u)
	if neg {
		u = neg256(u)
	}
	factor := pow10Value(int64(scaleDigits - n))
	q, r := divMod256By256(u, factor)
	if !isZero(r) {
		switch mode {
		case roundModeAwayFromZero:
			q = add256(q, u256{1, 0, 0, 0})
		case roundModeBanker:
			cmp := cmp256(add256(r, r), factor)
			if cmp > 0 || (cmp == 0 && (q[0]&1) == 1) {
				q = add256(q, u256{1, 0, 0, 0})
			}
		case roundModeCeil:
			if !neg {
				q = add256(q, u256{1, 0, 0, 0})
			}
		case roundModeFloor:
			if neg {
				q = add256(q, u256{1, 0, 0, 0})
			}
		}
	}
	res := mul256(q, factor)
	out := lower256(res)
	if neg {
		out = neg256(out)
	}
	return Decimal256(out)
}

// stringFixedN is an internal helper.
func (d Decimal256) stringFixedN(n int) string {
	u := u256(d)
	if isZero(u) {
		if n == 0 {
			return "0"
		}
		return "0." + repeatZero(n)
	}
	neg := isNeg(u)
	if neg {
		u = neg256(u)
	}
	q, r := divMod256By256(u, scale)
	intStr := u256ToDecimal(q)
	frac := u256ToDecimalFixed(r, scaleDigits)
	if n < scaleDigits {
		frac = frac[:n]
	}
	if neg {
		return "-" + intStr + "." + frac
	}
	return intStr + "." + frac
}

// buildPow10 is an internal helper.
func buildPow10() [65]u256 {
	var p [65]u256
	p[0] = u256{1, 0, 0, 0}
	for i := 1; i < len(p); i++ {
		p[i] = mul256ByUint64(p[i-1], 10)
	}
	return p
}

// pow10Mod is an internal helper.
func pow10Mod(n int64) u256 {
	if n <= 0 {
		return u256{1, 0, 0, 0}
	}
	if n <= 64 {
		return pow10[n]
	}
	result := u256{1, 0, 0, 0}
	base := u256{10, 0, 0, 0}
	for n > 0 {
		if (n & 1) == 1 {
			result = mul256Lo(result, base)
		}
		base = mul256Lo(base, base)
		n >>= 1
	}
	return result
}

// pow10Value is an internal helper.
func pow10Value(n int64) u256 {
	if n <= 0 {
		return u256{1, 0, 0, 0}
	}
	if n <= 64 {
		return pow10[n]
	}
	if n > maxDecimalDigits {
		return u256{}
	}
	return pow10Mod(n)
}

// parseDecimalString is an internal helper.
func parseDecimalString(s string) (u256, error) {
	start, end := trimSpaceString(s)
	if start >= end {
		return u256{}, errInvalidDecimal
	}
	idx := start
	sign := 1
	if s[idx] == '+' {
		idx++
	} else if s[idx] == '-' {
		sign = -1
		idx++
	}
	var val u256
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
				return u256{}, errInvalidDecimal
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
			return u256{}, errInvalidDecimal
		}
		sawDigit = true
		val = mul256ByUint64(val, 10)
		val = add256(val, u256{uint64(c - '0'), 0, 0, 0})
		if sawDot {
			fracDigits++
		}
		idx++
	}
	if !sawDigit {
		return u256{}, errInvalidDecimal
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
				return u256{}, errInvalidDecimal
			}
			if exp < (1 << 60) {
				exp = exp*10 + int64(c-'0')
			}
			idx++
		}
		if startExp == idx {
			return u256{}, errInvalidDecimal
		}
		exp *= expSign
	}
	shift := exp - fracDigits + scaleDigits
	if shift >= 0 {
		p := mul256(val, pow10Mod(shift))
		val = lower256(p)
	} else {
		factor := pow10Value(-shift)
		if isZero(factor) {
			val = u256{}
		} else {
			val = divByU256Trunc(val, factor)
		}
	}
	if sign < 0 {
		val = neg256(val)
	}
	return applyPrecision256(val), nil
}

// parseDecimalBytes is an internal helper.
func parseDecimalBytes(b []byte) (u256, error) {
	start, end := trimSpaceBytes(b)
	if start >= end {
		return u256{}, errInvalidDecimal
	}
	idx := start
	sign := 1
	if b[idx] == '+' {
		idx++
	} else if b[idx] == '-' {
		sign = -1
		idx++
	}
	var val u256
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
				return u256{}, errInvalidDecimal
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
			return u256{}, errInvalidDecimal
		}
		sawDigit = true
		val = mul256ByUint64(val, 10)
		val = add256(val, u256{uint64(c - '0'), 0, 0, 0})
		if sawDot {
			fracDigits++
		}
		idx++
	}
	if !sawDigit {
		return u256{}, errInvalidDecimal
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
				return u256{}, errInvalidDecimal
			}
			if exp < (1 << 60) {
				exp = exp*10 + int64(c-'0')
			}
			idx++
		}
		if startExp == idx {
			return u256{}, errInvalidDecimal
		}
		exp *= expSign
	}
	shift := exp - fracDigits + scaleDigits
	if shift >= 0 {
		p := mul256(val, pow10Mod(shift))
		val = lower256(p)
	} else {
		factor := pow10Value(-shift)
		if isZero(factor) {
			val = u256{}
		} else {
			val = divByU256Trunc(val, factor)
		}
	}
	if sign < 0 {
		val = neg256(val)
	}
	return applyPrecision256(val), nil
}

// applyPrecision256 drops integer digits beyond 32 and fractional digits beyond 32.
// It assumes the input is scaled by 10^32.
func applyPrecision256(u u256) u256 {
	if isZero(u) {
		return u
	}
	neg := isNeg(u)
	if neg {
		u = neg256(u)
	}
	intPart, fracPart := divMod256By256(u, scale)
	_, intRem := divMod256By256(intPart, scale)
	intPart = intRem
	// scaleDigits == 32, so fractional trimming is a no-op.
	raw := add256(lower256(mul256(intPart, scale)), fracPart)
	if neg {
		raw = neg256(raw)
	}
	return raw
}

// trimRightZeros is an internal helper.
func trimRightZeros(s string) string {
	end := len(s)
	for end > 0 && s[end-1] == '0' {
		end--
	}
	if end == 0 {
		return "0"
	}
	return s[:end]
}

// repeatZero is an internal helper.
func repeatZero(n int) string {
	if n <= 0 {
		return ""
	}
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = '0'
	}
	return string(buf)
}

// trimSpaceString is an internal helper.
func trimSpaceString(s string) (int, int) {
	start := 0
	end := len(s)
	for start < end && s[start] <= ' ' {
		start++
	}
	for end > start && s[end-1] <= ' ' {
		end--
	}
	return start, end
}

// trimSpaceBytes is an internal helper.
func trimSpaceBytes(b []byte) (int, int) {
	start := 0
	end := len(b)
	for start < end && b[start] <= ' ' {
		start++
	}
	for end > start && b[end-1] <= ' ' {
		end--
	}
	return start, end
}

// absInt64 is an internal helper.
func absInt64(v int64) (uint64, bool) {
	if v >= 0 {
		return uint64(v), false
	}
	return uint64(^v) + 1, true
}

// decimalDigitsU64 is an internal helper.
func decimalDigitsU64(v uint64) int {
	if v == 0 {
		return 1
	}
	n := 0
	for v != 0 {
		v /= 10
		n++
	}
	return n
}

// divDecimalByUint64 is an internal helper.
func divDecimalByUint64(d Decimal256, n uint64) Decimal256 {
	if n == 0 {
		return d
	}
	u := u256(d)
	if isZero(u) {
		return d
	}
	neg := isNeg(u)
	if neg {
		u = neg256(u)
	}
	q, _ := divMod256ByUint64(u, n)
	if neg {
		q = neg256(q)
	}
	return Decimal256(q)
}

// u256FromFloatTrunc is an internal helper.
func u256FromFloatTrunc(v float64) u256 {
	if v <= 0 {
		return u256{}
	}
	bits64 := math.Float64bits(v)
	exp := int((bits64>>52)&0x7ff) - 1023
	mant := bits64 & ((uint64(1) << 52) - 1)
	mant |= uint64(1) << 52
	if exp < 0 {
		return u256{}
	}
	shift := exp - 52
	if shift >= 256 {
		return u256{}
	}
	u := u256{mant, 0, 0, 0}
	if shift == 0 {
		return u
	}
	if shift > 0 {
		return shl256(u, uint(shift))
	}
	return shr256(u, uint(-shift))
}

// u256ToFloat is an internal helper.
func u256ToFloat(u u256) float64 {
	if isZero(u) {
		return 0
	}
	f := 0.0
	for i := 3; i >= 0; i-- {
		f = f*math.Exp2(64) + float64(u[i])
		if i == 0 {
			break
		}
	}
	return f
}

// appendU256Decimal is an internal helper.
func appendU256Decimal(dst []byte, v u256) []byte {
	if isZero(v) {
		return append(dst, '0')
	}
	const base = uint64(1_000_000_000_000_000_000) // 1e18
	var parts [5]uint64
	n := 0
	for !isZero(v) {
		var rem uint64
		v, rem = divMod256ByUint64(v, base)
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

// appendU256DecimalFixed is an internal helper.
func appendU256DecimalFixed(dst []byte, v u256, width int) []byte {
	if width <= 0 {
		return dst
	}
	var tmp [80]byte
	num := appendU256Decimal(tmp[:0], v)
	if len(num) >= width {
		return append(dst, num[len(num)-width:]...)
	}
	for i := 0; i < width-len(num); i++ {
		dst = append(dst, '0')
	}
	return append(dst, num...)
}

// u256ToDecimal is an internal helper.
func u256ToDecimal(v u256) string {
	if isZero(v) {
		return "0"
	}
	const base = uint64(1_000_000_000_000_000_000) // 1e18
	var parts [5]uint64
	n := 0
	for !isZero(v) {
		var rem uint64
		v, rem = divMod256ByUint64(v, base)
		parts[n] = rem
		n++
	}
	buf := make([]byte, 0, 80)
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

// u256ToDecimalFixed is an internal helper.
func u256ToDecimalFixed(v u256, width int) string {
	if width <= 0 {
		return ""
	}
	s := u256ToDecimal(v)
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

// divByU256Trunc is an internal helper.
func divByU256Trunc(v u256, d u256) u256 {
	if isZero(d) {
		return v
	}
	neg := isNeg(v)
	if neg {
		v = neg256(v)
	}
	q, _ := divMod256By256(v, d)
	if neg {
		q = neg256(q)
	}
	return q
}

// isZero is an internal helper.
func isZero(u u256) bool {
	return u[0]|u[1]|u[2]|u[3] == 0
}

// isNeg is an internal helper.
func isNeg(u u256) bool {
	return (u[3]>>63)&1 == 1
}

// add256 is an internal helper.
func add256(a, b u256) u256 {
	var out u256
	var c uint64
	out[0], c = bits.Add64(a[0], b[0], 0)
	out[1], c = bits.Add64(a[1], b[1], c)
	out[2], c = bits.Add64(a[2], b[2], c)
	out[3], _ = bits.Add64(a[3], b[3], c)
	return out
}

// sub256 is an internal helper.
func sub256(a, b u256) u256 {
	var out u256
	var c uint64
	out[0], c = bits.Sub64(a[0], b[0], 0)
	out[1], c = bits.Sub64(a[1], b[1], c)
	out[2], c = bits.Sub64(a[2], b[2], c)
	out[3], _ = bits.Sub64(a[3], b[3], c)
	return out
}

// neg256 is an internal helper.
func neg256(a u256) u256 {
	var out u256
	out[0] = ^a[0]
	out[1] = ^a[1]
	out[2] = ^a[2]
	out[3] = ^a[3]
	var c uint64
	out[0], c = bits.Add64(out[0], 1, 0)
	out[1], c = bits.Add64(out[1], 0, c)
	out[2], c = bits.Add64(out[2], 0, c)
	out[3], _ = bits.Add64(out[3], 0, c)
	return out
}

// cmp256 is an internal helper.
func cmp256(a, b u256) int {
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

// cmp256Signed is an internal helper.
func cmp256Signed(a, b u256) int {
	na := isNeg(a)
	nb := isNeg(b)
	if na && !nb {
		return -1
	}
	if !na && nb {
		return 1
	}
	return cmp256(a, b)
}

// mul256ByUint64 is an internal helper.
func mul256ByUint64(a u256, m uint64) u256 {
	var out u256
	var carry uint64
	for i := 0; i < 4; i++ {
		hi, lo := bits.Mul64(a[i], m)
		lo, c := bits.Add64(lo, carry, 0)
		out[i] = lo
		carry = hi + c
	}
	return out
}

// divMod256ByUint64 is an internal helper.
func divMod256ByUint64(a u256, d uint64) (u256, uint64) {
	var q u256
	var r uint64
	for i := 3; i >= 0; i-- {
		q[i], r = bits.Div64(r, a[i], d)
		if i == 0 {
			break
		}
	}
	return q, r
}

// mul256 is an internal helper.
func mul256(a, b u256) u512 {
	var p u512
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			hi, lo := bits.Mul64(a[i], b[j])
			k := i + j
			var c uint64
			p[k], c = bits.Add64(p[k], lo, 0)
			p[k+1], c = bits.Add64(p[k+1], hi, c)
			idx := k + 2
			for c != 0 && idx < 8 {
				p[idx], c = bits.Add64(p[idx], 0, c)
				idx++
			}
		}
	}
	return p
}

// mul256Lo is an internal helper.
func mul256Lo(a, b u256) u256 {
	return lower256(mul256(a, b))
}

// lower256 is an internal helper.
func lower256(p u512) u256 {
	return u256{p[0], p[1], p[2], p[3]}
}

// divMod256By256 is an internal helper.
func divMod256By256(n u256, d u256) (u256, u256) {
	if isZero(d) {
		return u256{}, u256{}
	}
	var n512 u512
	n512[0] = n[0]
	n512[1] = n[1]
	n512[2] = n[2]
	n512[3] = n[3]
	q, r := divMod512By256(n512, d)
	return q, r
}

// divMod512By256 is an internal helper.
func divMod512By256(n u512, d u256) (u256, u256) {
	if isZero(d) {
		return u256{}, u256{}
	}
	if isZeroU512(n) {
		return u256{}, u256{}
	}
	nBits := bitLen512(n)
	dBits := bitLen256(d)
	if nBits < dBits {
		return u256{}, lower256(n)
	}
	shift := nBits - dBits
	var q u256
	rem := n
	dShift := shl256To512(d, uint(shift))
	for i := shift; i >= 0; i-- {
		if cmp512(rem, dShift) >= 0 {
			rem = sub512(rem, dShift)
			if i < 256 {
				q[int(i/64)] |= 1 << uint(i%64)
			}
		}
		if i == 0 {
			break
		}
		dShift = shr1_512(dShift)
	}
	return q, lower256(rem)
}

// isZeroU512 is an internal helper.
func isZeroU512(u u512) bool {
	return u[0]|u[1]|u[2]|u[3]|u[4]|u[5]|u[6]|u[7] == 0
}

// bitLen256 is an internal helper.
func bitLen256(u u256) int {
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

// bitLen512 is an internal helper.
func bitLen512(u u512) int {
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

// cmp512 is an internal helper.
func cmp512(a, b u512) int {
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

// sub512 is an internal helper.
func sub512(a, b u512) u512 {
	var out u512
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

// shl256 is an internal helper.
func shl256(u u256, shift uint) u256 {
	if shift == 0 {
		return u
	}
	if shift >= 256 {
		return u256{}
	}
	wordShift := shift / 64
	bitShift := shift % 64
	var out u256
	for i := 0; i < 4; i++ {
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

// shr256 is an internal helper.
func shr256(u u256, shift uint) u256 {
	if shift == 0 {
		return u
	}
	if shift >= 256 {
		return u256{}
	}
	wordShift := shift / 64
	bitShift := shift % 64
	var out u256
	for i := int(wordShift); i < 4; i++ {
		src := i
		dst := i - int(wordShift)
		out[dst] |= u[src] >> bitShift
		if bitShift != 0 && src+1 < 4 {
			out[dst] |= u[src+1] << (64 - bitShift)
		}
	}
	return out
}

// shl256To512 is an internal helper.
func shl256To512(u u256, shift uint) u512 {
	if shift == 0 {
		return u512{u[0], u[1], u[2], u[3], 0, 0, 0, 0}
	}
	if shift >= 512 {
		return u512{}
	}
	wordShift := shift / 64
	bitShift := shift % 64
	var out u512
	for i := 0; i < 4; i++ {
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

// shr1_512 is an internal helper.
func shr1_512(u u512) u512 {
	var out u512
	var carry uint64
	for i := 7; i >= 0; i-- {
		out[i] = (u[i] >> 1) | (carry << 63)
		carry = u[i] & 1
		if i == 0 {
			break
		}
	}
	return out
}

// u256FromInt64 is an internal helper.
func u256FromInt64(v int64) u256 {
	if v >= 0 {
		return u256{uint64(v), 0, 0, 0}
	}
	uv := uint64(v)
	return u256{uv, ^uint64(0), ^uint64(0), ^uint64(0)}
}
