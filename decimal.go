package decimal

import (
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

const (
	// DivisionPrecision is the number of decimal places for division operations.
	DivisionPrecision int = 16

	// Zero constant, to make computations faster.
	// Zero should never be compared with == or != directly, please use decimal.Equal or decimal.Cmp instead.
	Zero Decimal = Decimal("0")

	_zero = "0"
)

// New create a Decimal. If value is empty, return zero.
//
// Acceptable symbol (+-.,_0123456789)
func New(value ...string) (Decimal, error) {
	if len(value) == 0 {
		return Zero, nil
	}

	buf, err := newDecimal([]byte(value[0]))
	if err != nil {
		return Zero, err
	}

	return Decimal(buf), nil
}

// Require returns a new Decimal from a string representation or panics if New would have returned an error.
//
// Acceptable symbol (+-.,_0123456789)
//
// Example:
//
//	d := decimal.Require("123,456")
//	d2 := decimal.Require("")        // "0"
//	d3 := decimal.Require("&$")      // Panic!!!
func Require(value string) Decimal {
	return Decimal(normalize([]byte(value)))
}

// NewFromInt converts a int64 to Decimal.
func NewFromInt(value int64) Decimal {
	return Decimal(strconv.FormatInt(value, 10))
}

// NewFromInt32 converts a int32 to Decimal.
func NewFromInt32(value int32) Decimal {
	return Decimal(strconv.FormatInt(int64(value), 10))
}

// NewFromFloat create a Decimal from a float64.
//
// NOTE: this will create zero value on NaN, +/-inf
func NewFromFloat(value float64) Decimal {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return Zero
	}

	return Decimal(strconv.FormatFloat(value, 'f', -1, 64))
}

// NewFromFloat32 create a Decimal from a float32.
//
// NOTE: this will create zero value on NaN, +/-inf
func NewFromFloat32(value float32) Decimal {
	vf := float64(value)
	if math.IsNaN(vf) || math.IsInf(vf, 0) {
		return Zero
	}

	return Decimal(strconv.FormatFloat(vf, 'f', -1, 64))
}

// NewFromBigInt returns a new Decimal from a big.Int, value * 10 ^ exp
func NewFromBigInt(value *big.Int, exp int) Decimal {
	return Decimal(shift([]byte(value.String()), exp))
}

// NewFromString returns a new Decimal from a string representation.
// Trailing zeroes are not trimmed.
//
// Acceptable symbol (+-.,_0123456789)
//
// NOTE: This function is for compatibility with the shopspring/decimal package.
// Please use New instead.
//
// Deprecated: Use New(value) instead.
func NewFromString(value string) (Decimal, error) {
	return New(value)
}

// RequireFromString returns a new Decimal from a string representation
// or panics if NewFromString would have returned an error.
//
// Acceptable symbol (+-.,_0123456789)
//
// NOTE: This function is for compatibility with the shopspring/decimal package.
// Please use Require instead.
//
// Deprecated: Use Require(value) instead.
func RequireFromString(value string) Decimal {
	return Require(value)
}

// NewFromFormattedString returns a new Decimal from a formatted string representation.
// The second argument - replRegexp, is a regular expression that is used to find characters that should be
// removed from given decimal string representation. All matched characters will be replaced with an empty string.
//
// Example:
//
//	r := regexp.MustCompile("[$,]")
//	d1, err := NewFromFormattedString("$5,125.99", r)
//
//	r2 := regexp.MustCompile("[_]")
//	d2, err := NewFromFormattedString("1_000_000", r2)
//
//	r3 := regexp.MustCompile("[USD\\s]")
//	d3, err := NewFromFormattedString("5000 USD", r3)
func NewFromFormattedString(value string, replRegexp *regexp.Regexp) (Decimal, error) {
	parsedValue := replRegexp.ReplaceAllString(value, "")
	d, err := NewFromString(parsedValue)
	if err != nil {
		return Zero, err
	}
	return d, nil
}

// Min returns the smallest Decimal that was passed in the arguments.
//
// To call this function with an array, you must do:
//
//	Min(arr[0], arr[1:]...)
//
// This makes it harder to accidentally call Min with 0 arguments.
func Min(first Decimal, rest ...Decimal) Decimal {
	curr := first
	for i := 1; i < len(rest); i++ {
		if rest[i].LessThan(curr) {
			curr = rest[i]
		}
	}

	return curr
}

// Max returns the largest Decimal that was passed in the arguments.
//
// To call this function with an array, you must do:
//
//	Max(arr[0], arr[1:]...)
//
// This makes it harder to accidentally call Max with 0 arguments.
func Max(first Decimal, rest ...Decimal) Decimal {
	curr := first
	for i := 0; i < len(rest); i++ {
		if rest[i].GreaterThan(curr) {
			curr = rest[i]
		}
	}

	return curr
}

// Sum returns the combined total of the provided first and rest Decimals
func Sum(first Decimal, rest ...Decimal) Decimal {
	summed := first
	for _, d := range rest {
		summed = summed.Add(d)
	}
	return summed
}

// Avg returns the average value of the provided first and rest Decimals
func Avg(first Decimal, rest ...Decimal) Decimal {
	count := NewFromInt(int64(len(rest) + 1))
	sum := Sum(first, rest...)
	return sum.Div(count)
}

type Decimal string

// String returns the string representation of the decimal with the fixed point.
//
// Example:
//
//	d := New("-12.345")
//	println(d.String())
//
// Output:
//
//	-12.345
func (d Decimal) String() string {
	return string(normalize([]byte(d)))
}

// Copy returns a copy of decimal with the same value and exponent, but a different pointer to value.
func (d Decimal) Copy() Decimal {
	return d
}

// StringFixed returns a rounded fixed-point string with places digits after the decimal point.
//
// Example:
//
//	NewFromFloat(0).StringFixed(2)    // "0.00"
//	NewFromFloat(0).StringFixed(0)    // "0"
//	NewFromFloat(5.45).StringFixed(0) // "5"
//	NewFromFloat(5.45).StringFixed(1) // "5.5"
//	NewFromFloat(5.45).StringFixed(2) // "5.45"
//	NewFromFloat(5.45).StringFixed(3) // "5.450"
//	NewFromFloat(545).StringFixed(-1) // "550"
func (d Decimal) StringFixed(places int) string {
	buf := normalize(truncate([]byte(d), places))
	if places <= 0 {
		return string(buf)
	}

	dotIdx := findDotIndex(buf)
	if dotIdx == -1 {
		dotIdx = len(buf) - 1
	}

	rightDotIdx := len(buf) - 1 - dotIdx
	if places > rightDotIdx {
		buf = pushBackRepeat(buf, '0', places-rightDotIdx)
	}

	return string(buf)
}

// Abs returns the absolute value of the decimal.
func (d Decimal) Abs() Decimal {
	buf := normalize([]byte(d))

	if buf[0] == '-' {
		return Decimal(trimFront(buf, 1))
	}

	return Decimal(buf)
}

// Neg returns -d
//
// Example:
//
//	decimal.New("123.456").Neg().String() // "-123.45"
func (d Decimal) Neg() Decimal {
	buf := normalize([]byte(d))

	if buf[0] == '-' {
		return Decimal(trimFront(buf, 1))
	}

	return Decimal(pushFront(buf, '-'))
}

// Truncate truncates off digits from the number, without rounding.
//
// NOTE: precision is the last digit that will not be truncated (must be >= 0).
//
// Example:
//
//	decimal.New("123.456").Truncate(2).String() // "123.45"
func (d Decimal) Truncate(precision int) Decimal {
	if precision > len(d) {
		return d
	}

	if -precision > len(d) {
		return Zero
	}

	return Decimal(truncate(normalize([]byte(d)), precision))
}

// Shift shifts the decimal in base 10.
// It shifts left when shift is positive and right if shift is negative.
// In simpler terms, the given value for shift is added to the exponent
// of the decimal.
//
// Example:
//
//	d, _ := decimal.New("3")
//	d.Shift(3).String()  // "3000"
//	d.Shift(-3).String() // "0.003"
func (d Decimal) Shift(sf int) Decimal {
	return Decimal(shift(normalize([]byte(d)), sf))
}

func combineToDecimal(ss ...string) Decimal {
	builder := strings.Builder{}
	l := 0
	for i := range ss {
		l += len(ss[i])
	}
	builder.Grow(l)

	for _, s := range ss {
		builder.WriteString(s)
	}

	return Decimal(builder.String())
}

// Add return d + d2
//
// Example:
//
//	d1, _ := decimal.New("100")
//	d2, _ := decimal.New("90.99")
//	d1.Add(d2).String() // "190.01"
func (d Decimal) Add(d2 Decimal) Decimal {
	b, a := normalize([]byte(d)), normalize([]byte(d2))
	baseNegative := b[0] == '-'
	additionNegative := a[0] == '-'
	if baseNegative && additionNegative {
		b = trimFront(b, 1)
		a = trimFront(a, 1)
		// -b - -a = - (b+a)
		return Decimal(pushFront(unsignedAdd(b, a), '-'))
	}

	if baseNegative {
		b = trimFront(b, 1)
		// -b + a = a - b
		return Decimal(unsignedSub(a, b))
	}

	if additionNegative {
		a = trimFront(a, 1)
		// b + -a = b - a
		return Decimal(unsignedSub(b, a))
	}

	// b + a = b + a
	return Decimal(unsignedAdd(b, a))
}

// Sub return d - d2
//
// Example:
//
//	d1, _ := decimal.New("100")
//	d2, _ := decimal.New("90.99")
//	d1.Sub(d2).String() // "9.01"
func (d Decimal) Sub(d2 Decimal) Decimal {
	return Decimal(sub(normalize([]byte(d)), normalize([]byte(d2))))
}

func sub(b, a []byte) []byte {
	baseNegative := b[0] == '-'
	additionNegative := a[0] == '-'
	if baseNegative && additionNegative {
		b = trimFront(b, 1)
		a = trimFront(a, 1)
		// -b - -a = -b + a = a - b
		return unsignedSub(a, b)
	}

	if baseNegative {
		b = trimFront(b, 1)
		// -b - a = - (b + a)
		return pushFront(unsignedAdd(a, b), '-')
	}

	if additionNegative {
		a = trimFront(a, 1)
		// b - -a = b + a
		return unsignedAdd(b, a)
	}

	// b - a = b - a
	return unsignedSub(b, a)
}

// clean the zero and dot of prefixes and suffixes
func tidy(num []byte) Decimal {
	return tidyString(string(num))
}

// clean the zero and dot of prefixes and suffixes
func tidyString(num string) Decimal {
	if len(num) == 0 {
		return _zero
	}

	hasDot := false
	for _, c := range num {
		if c == '.' {
			hasDot = true
			break
		}
	}

	// Handle sign prefix
	var sign bool
	start := 0
	switch num[0] {
	case '+':
		start = 1
	case '-':
		sign = true
		start = 1
	}

	// Find start position (skip leading zeros)
	for start < len(num) && num[start] == '0' {
		start++
	}

	// If all zeros or empty after sign
	if start >= len(num) {
		return _zero
	}

	// Find end position
	end := len(num)
	if hasDot {
		// Remove trailing zeros (but not if noDecimalPoint is true)
		for end > start && num[end-1] == '0' {
			end--
		}

		// Handle edge cases
		if end == start {
			return _zero
		}

		// Remove trailing decimal point
		if end > start && num[end-1] == '.' {
			end--
		}
	}

	// Extract the significant part
	result := num[start:end]

	if len(result) == 0 {
		return _zero
	}

	if result[len(result)-1] == '.' {
		result = result[:len(result)-1]
	}

	// Handle leading decimal point
	if len(result) > 0 && result[0] == '.' {
		if sign {
			return combineToDecimal("-", _zero, result)
		}
		return combineToDecimal(_zero, result)
	}

	// Handle single decimal point
	if result == "." {
		return _zero
	}

	if sign {
		return combineToDecimal("-", result)
	}

	return Decimal(result)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Pow returns d to the power d2
func (d Decimal) Pow(d2 Decimal) Decimal {
	cache := make([]byte, 0, len(d)*len(d2))
	return Decimal(pow(normalize([]byte(d)), normalize([]byte(d2)), &cache))
}

func pow(a, b []byte, reused ...*[]byte) []byte {
	bIntPart := intPartInt64(b)
	if bIntPart == 0 {
		return oneBytes
	}

	if bIntPart == 1 {
		return a
	}

	temp := pow(a, div(b, twoBytes), reused...)
	if bIntPart%2 == 0 {
		return sqr(temp, reused...)
	}

	if bIntPart > 0 {
		return mul(sqr(temp, reused...), a, reused...)
	}

	return div(sqr(temp, reused...), a)
}

// Rat returns a rational number representation of the decimal.
func (d Decimal) Rat() *big.Rat {
	r := new(big.Rat)
	r.SetString(d.String())
	return r
}

// BigFloat returns decimal as BigFloat.
// Be aware that casting decimal to BigFloat might cause a loss of precision.
func (d Decimal) BigFloat() *big.Float {
	f := new(big.Float)
	f.SetString(d.String())
	return f
}

// BigInt returns integer component of the decimal as a BigInt.
func (d Decimal) BigInt() *big.Int {
	i := new(big.Int)
	i.SetString(string(intPart(normalize([]byte(d)))), 10)
	return i
}

// Float64 returns the nearest float64 value for d and a bool indicating
// whether f represents d exactly.
// For more details, see the documentation for big.Rat.Float64
func (d Decimal) Float64() (f float64, exact bool) {
	return d.Rat().Float64()
}

// IntPart returns the integer component of the decimal.
func (d Decimal) IntPart() int64 {
	return intPartInt64(normalize([]byte(d)))
}

// IsZero return d == 0
func (d Decimal) IsZero() bool {
	if len(d) == 0 {
		return true
	}

	return isZero(normalize([]byte(d)))
}

// IsInteger returns true when decimal can be represented as an integer value, otherwise, it returns false.
func (d Decimal) IsInteger() bool {
	if len(d) == 0 {
		return true
	}

	buf := normalize([]byte(d))
	dotIdx := findDotIndex(buf)
	if dotIdx == -1 {
		return true
	}

	for ; dotIdx < len(buf); dotIdx++ {
		if buf[dotIdx] != '0' {
			return false
		}
	}

	return true
}

// IsPositive return d > 0
func (d Decimal) IsPositive() bool {
	buf := normalize([]byte(d))

	return !isZero(buf) && !isNegative(buf)
}

// IsNegative return d < 0
func (d Decimal) IsNegative() bool {
	return isNegative([]byte(d))
}

// Cmp compares the numbers represented by d and d2 and returns:
//
//	-1 if d <  d2
//	 0 if d == d2
//	+1 if d >  d2
func (d Decimal) Cmp(d2 Decimal) int {
	b := normalize([]byte(d))
	b2 := normalize([]byte(d2))

	if string(b) == string(b2) {
		return 0
	}

	if great(b, b2) {
		return 1
	}

	return -1
}

// Equal return d == d2
func (d Decimal) Equal(d2 Decimal) bool {
	return string(normalize([]byte(d))) == string(normalize([]byte(d2)))
}

// GreaterThan return d > d2
func (d Decimal) GreaterThan(d2 Decimal) bool {
	return great(normalize([]byte(d)), normalize([]byte(d2)))
}

// LessThan return d < d2
func (d Decimal) LessThan(d2 Decimal) bool {
	return less(normalize([]byte(d)), normalize([]byte(d2)))
}

// GreaterThanOrEqual return d >= d2
func (d Decimal) GreaterThanOrEqual(d2 Decimal) bool {
	return !less(normalize([]byte(d)), normalize([]byte(d2)))
}

// LessThanOrEqual return d <= d2
func (d Decimal) LessThanOrEqual(d2 Decimal) bool {
	return !great(normalize([]byte(d)), normalize([]byte(d2)))
}

// Mul return d * d2
func (d Decimal) Mul(d2 Decimal) Decimal {
	return Decimal(mul(normalize([]byte(d)), normalize([]byte(d2))))
}

func sqr(a []byte, reused ...*[]byte) []byte {
	if isZero(a) {
		return zeroBytes
	}

	a, right := removeDecimalPoint(a)

	multiplied := multiplyPureNumber(a, a, reused...)
	rightSumDigit := right + right
	if rightSumDigit != 0 {
		if len(multiplied) <= rightSumDigit {
			multiplied = pushFrontRepeat(multiplied, '0', rightSumDigit-len(multiplied)+1)
		}
		digitIdx := len(multiplied) - rightSumDigit
		multiplied = insert(multiplied, digitIdx, '.')
	}

	return multiplied
}

func mul(a, b []byte, reused ...*[]byte) []byte {
	if isZero(a) || isZero(b) {
		return zeroBytes
	}

	a, right1 := removeDecimalPoint(a)
	b, right2 := removeDecimalPoint(b)

	minus := false
	if a[0] == '-' {
		a = trimFront(a, 1)
		minus = !minus
	}

	if b[0] == '-' {
		b = trimFront(b, 1)
		minus = !minus
	}

	multiplied := multiplyPureNumber(a, b, reused...)
	rightSumDigit := right1 + right2
	if rightSumDigit != 0 {
		if len(multiplied) <= rightSumDigit {
			multiplied = pushFrontRepeat(multiplied, '0', rightSumDigit-len(multiplied)+1)
		}
		digitIdx := len(multiplied) - rightSumDigit
		multiplied = insert(multiplied, digitIdx, '.')
	}

	if minus {
		multiplied = pushFront(multiplied, '-')
	}

	return multiplied
}

// removeDecimalPoint removes decimal point and return the count of the digit right the decimal
//
// OPTIMIZED: NO COPY
func removeDecimalPoint(s []byte) (result []byte, countOfRightSide int) {
	for i := range s {
		if s[i] == '.' {
			return remove(s, i), len(s) - i - 1
		}
	}
	return s, 0
}

// multiplyPureNumber return d1 * d2, d1 & d2 must contain only number 0~9
func multiplyPureNumber(d1 []byte, d2 []byte, reused ...*[]byte) []byte {
	if len(d1) < len(d2) {
		d1, d2 = d2, d1
	}

	var (
		extraCap   = 3 // for outside this func to append '-0.'
		len1, len2 = len(d1), len(d2)
		resultCap  = len1 + len2 + extraCap
		isReused   = len(reused) != 0
		resultIdx  int
		carry      byte
		val1, val2 byte
		prod       byte
		result     []byte // make([]byte, len1+len2+extraCap)
	)

	if isReused {
		result = extend(*reused[0], resultCap)
		for i := range result {
			result[i] = 0
		}

		for len(result) < resultCap {
			result = append(result, 0)
		}

		result = result[:resultCap]
	} else {
		result = make([]byte, resultCap)
	}

	// Optimized multiplication using single loop with carry propagation
	for i := len2 - 1; i >= 0; i-- {
		val2 = d2[i] - '0'
		if val2 == 0 {
			continue
		}

		carry = 0
		resultIdx = i + len1

		// Inner multiplication loop
		for j := len1 - 1; j >= 0; j-- {
			val1 = d1[j] - '0'
			prod = val1*val2 + result[resultIdx] + carry
			result[resultIdx] = prod % 10
			carry = prod / 10
			resultIdx--
		}

		// Handle remaining carry
		if carry > 0 {
			result[i] += carry
		}
	}

	// Convert to ASCII in-place - no additional memory allocation
	for i := range result {
		result[i] += '0'
	}

	return tidyBytes(trimBack(result, extraCap))
}

// Sign return the sign of the decimal
//
// return 1 if d > 0, 0 if d == 0, -1 if d < 0
func (d Decimal) Sign() int {
	return int(sign(normalize([]byte(d))))
}

// Round rounds the decimal to places decimal places.
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Example:
//
//	NewFromFloat(5.45).Round(1).String() // "5.5"
//	NewFromFloat(545).Round(-1).String() // "550"
func (d Decimal) Round(places int) Decimal {
	return Decimal(roundInCondition(normalize([]byte(d)), places, func(isDecimalNeg bool, bankChar, roundChar byte) bool {
		return roundChar >= '5'
	}))
}

// RoundBank rounds the decimal to places decimal places.
// If the final digit to round is equidistant from the nearest two integers the
// rounded value is taken as the even number
//
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Examples:
//
//	NewFromFloat(5.45).RoundBank(1).String() // "5.4"
//	NewFromFloat(545).RoundBank(-1).String() // "540"
//	NewFromFloat(5.46).RoundBank(1).String() // "5.5"
//	NewFromFloat(546).RoundBank(-1).String() // "550"
//	NewFromFloat(5.55).RoundBank(1).String() // "5.6"
//	NewFromFloat(555).RoundBank(-1).String() // "560"
func (d Decimal) RoundBank(places int) Decimal {
	return Decimal(roundInCondition(normalize([]byte(d)), places, func(isDecimalNeg bool, bankChar, roundChar byte) bool {
		return roundChar > '5' || (roundChar == '5' && (bankChar-'0')%2 == 1)
	}))
}

// RoundAwayFromZero rounds the decimal away from zero.
//
// Example:
//
//	NewFromFloat(545).RoundAwayFromZero(-2).String()   // "600"
//	NewFromFloat(-545).RoundAwayFromZero(-2).String()   // "-600"
//	NewFromFloat(1.1001).RoundAwayFromZero(2).String() // "1.11"
//	NewFromFloat(-1.454).RoundAwayFromZero(1).String() // "-1.5"
func (d Decimal) RoundAwayFromZero(places int) Decimal {
	return Decimal(roundInCondition(normalize([]byte(d)), places, func(isDecimalNeg bool, bankChar, roundChar byte) bool {
		return roundChar != '0'
	}))
}

// RoundTowardToZero rounds the decimal towards zero.
//
// Example:
//
//	NewFromFloat(545).RoundTowardToZero(-2).String()   // "500"
//	NewFromFloat(-545).RoundTowardToZero(-2).String()  // "-500"
//	NewFromFloat(1.1001).RoundTowardToZero(2).String() // "1.1"
//	NewFromFloat(-1.454).RoundTowardToZero(1).String() // "-1.4"
func (d Decimal) RoundTowardToZero(places int) Decimal {
	return Decimal(truncate(normalize([]byte(d)), places))
}

// RoundUp rounds the decimal away from zero.
//
// NOTE: This function is for compatibility with the shopspring/decimal package.
// Please use RoundAwayFromZero instead.
//
// Deprecated: Use RoundAwayFromZero(places) instead.
func (d Decimal) RoundUp(places int) Decimal {
	return d.RoundAwayFromZero(places)
}

// RoundDown rounds the decimal towards zero.
//
// NOTE: This function is for compatibility with the shopspring/decimal package.
// Please use RoundTowardToZero instead.
//
// Deprecated: Use RoundTowardToZero(places) instead.
func (d Decimal) RoundDown(places int) Decimal {
	return d.RoundTowardToZero(places)
}

// RoundFloor rounds the decimal towards zero.
//
// NOTE: This function is for compatibility with the shopspring/decimal package.
// Please use Floor instead.
//
// Deprecated: Use Floor(places) instead.
func (d Decimal) RoundFloor(places int) Decimal {
	return d.Floor(places)
}

// RoundCeil rounds the decimal towards +infinity.
//
// NOTE: This function is for compatibility with the shopspring/decimal package.
// Please use Ceil instead.
//
// Deprecated: Use Ceil(places) instead.
func (d Decimal) RoundCeil(places int) Decimal {
	return d.Ceil(places)
}

// Ceil rounds the decimal towards +infinity.
//
// Example:
//
//	NewFromFloat(545).Ceil(-2).String()   // "600"
//	NewFromFloat(-545).Ceil(-2).String()   // "-500"
//	NewFromFloat(1.1001).Ceil(2).String() // "1.11"
//	NewFromFloat(-1.454).Ceil(1).String() // "-1.4"
func (d Decimal) Ceil(places int) Decimal {
	return Decimal(roundInCondition(normalize([]byte(d)), places, func(isDecimalNeg bool, bankChar, roundChar byte) bool {
		if isDecimalNeg {
			return roundChar == '0'
		} else {
			return roundChar != '0'
		}
	}))
}

// Floor rounds the decimal towards -infinity.
//
// Example:
//
//	NewFromFloat(545).Floor(-2).String()   //  "500"
//	NewFromFloat(-545).Floor(-2).String()  //  "-600"
//	NewFromFloat(1.1001).Floor(2).String() //  "1.1"
//	NewFromFloat(-1.454).Floor(1).String() //  "-1.5"
func (d Decimal) Floor(places int) Decimal {
	return Decimal(roundInCondition(normalize([]byte(d)), places, func(isDecimalNeg bool, bankChar, roundChar byte) bool {
		if isDecimalNeg {
			return roundChar != '0'
		} else {
			return roundChar == '0'
		}
	}))
}

// roundInCondition rounds the decimal in the condition
func roundInCondition(buf []byte, place int, carryCond func(isDecimalNeg bool, bankChar, roundChar byte) bool) []byte {
	isDecimalNeg := len(buf) != 0 && buf[0] == '-'
	var result []byte
	if isDecimalNeg {
		result = roundInConditionWithoutSign(trimFront(buf, 1), place, isDecimalNeg, carryCond)
	} else {
		result = roundInConditionWithoutSign(buf, place, isDecimalNeg, carryCond)
	}

	if isZero(result) {
		return zeroBytes
	}

	if isDecimalNeg {
		return pushFront(result, '-')
	}

	return result
}

func roundInConditionWithoutSign(buf []byte, places int, isDecimalNeg bool, carryCond func(isDecimalNeg bool, bankChar, roundChar byte) bool) []byte {
	dotIdx := findDotIndex(buf)
	// Handle negative precision (round to left of decimal point)
	if places < 0 {
		var roundPos int
		negPlace := -places

		if dotIdx == -1 {
			roundPos = len(buf) - negPlace
		} else {
			roundPos = dotIdx - negPlace
		}

		if roundPos <= -1 {
			return zeroBytes
		}

		var (
			needCarry bool
			result    []byte
		)

		if roundPos == 0 {
			needCarry = roundPos < len(buf) && carryCond(isDecimalNeg, '1', buf[roundPos])
			result = buf[:roundPos]
		} else {
			needCarry = roundPos < len(buf) && carryCond(isDecimalNeg, buf[roundPos-1], buf[roundPos])
			result = buf[:roundPos]
		}

		if needCarry {
			result = addCarryToPosition(result, roundPos-1)
		}

		return pushBackRepeat(result, '0', negPlace)
	}

	// Handle zero or positive precision (round at or right of decimal point)
	if dotIdx == -1 {
		// No decimal point, no rounding needed for positive precision
		return buf
	}

	if places == 0 {
		var (
			needCarry bool
			result    []byte
		)
		// Round to integer
		if dotIdx == 0 {
			needCarry = dotIdx+1 < len(buf) && carryCond(isDecimalNeg, '1', buf[dotIdx+1])
			result = buf[:dotIdx]
		} else {
			needCarry = dotIdx+1 < len(buf) && carryCond(isDecimalNeg, buf[dotIdx-1], buf[dotIdx+1])
			result = buf[:dotIdx]
		}

		if needCarry {
			result = addCarryToPosition(result, len(result)-1)
		}

		return result
	}

	// Positive precision
	roundPos := dotIdx + places + 1
	if roundPos >= len(buf) {
		// Already has fewer digits than requested precision
		return buf
	}

	var (
		needCarry bool
		result    []byte
	)
	if roundPos == 0 {
		needCarry = carryCond(isDecimalNeg, '1', buf[roundPos])
		result = buf[:roundPos]
	} else {
		needCarry = carryCond(isDecimalNeg, buf[roundPos-1], buf[roundPos])
		result = buf[:roundPos]
	}

	if needCarry {
		result = addCarryToPosition(result, roundPos-1)
	}

	return result
}

// addCarryToPosition adds 1 to the digit at the specified position and handles carry propagation
func addCarryToPosition(buf []byte, pos int) []byte {
	if pos < 0 || len(buf) == 0 {
		return pushFront(buf, '1')
	}

	carry := byte(1)
	for i := pos; i >= 0 && carry > 0; i-- {
		if buf[i] == '.' {
			continue
		}

		digit := buf[i] - '0' + carry
		if digit <= 9 {
			buf[i] = digit + '0'
			carry = 0
		} else {
			buf[i] = '0'
			carry = 1
		}
	}

	// If there's still a carry, we need to add a digit at the front
	if carry > 0 {
		buf = pushFront(buf, '1')
	}

	return buf
}

func (d Decimal) Mod(d2 Decimal) Decimal {
	b1 := normalize([]byte(d))
	b2 := normalize([]byte(d2))
	divided := div(normalize([]byte(d)), b2)
	divided = truncate(divided, 0)
	return Decimal(sub(b1, mul(divided, b2)))
}
