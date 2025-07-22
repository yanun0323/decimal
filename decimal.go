package decimal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

/*
有没有办法优化的更快速、并且不消耗额外的记忆体
*/

var (
	DivisionPrecision = 16
	zero              = Decimal("0")
)

// Zero return the zero decimal
func Zero() Decimal {
	return zero
}

// New create a Decimal.
//
// acceptable symbol (+-.,_0123456789)
//
// Example:
//
//	d, err := decimal.New("123,456,789.000")
func New(s string) (Decimal, error) {
	return newDecimal(s)
}

func newDecimal(s string) (Decimal, error) {
	if len(s) == 0 {
		return Zero(), nil
	}

	switch s {
	case "0", "0.", "0.0", ".0":
		return zero, nil
	}

	// Use strings.Builder to avoid repeated allocations
	var buf strings.Builder
	buf.Grow(len(s)) // Pre-allocate capacity

	dot := false
	firstChar := true

	for i := 0; i < len(s); i++ {
		b := s[i]

		if firstChar {
			// Handle first character
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '-':
				buf.WriteByte(b)
				if b == '.' {
					dot = true
				}
			case '+':
				// Drop '+' sign
				continue
			default:
				return Zero(), fmt.Errorf("invalid symbol (%c) in %s", b, s)
			}
			firstChar = false
		} else {
			// Handle subsequent characters
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				buf.WriteByte(b)
			case '.':
				if dot {
					return Zero(), errors.New("duplicate dot")
				}
				dot = true
				buf.WriteByte(b)
			case '_', ',':
				// Drop these characters
				continue
			default:
				return Zero(), fmt.Errorf("invalid symbol (%c) in %s", b, s)
			}
		}
	}

	if buf.Len() == 0 {
		return Zero(), errors.New("can't convert to Decimal empty string")
	}

	result := buf.String()

	return Decimal(tidyString(result, !dot)), nil
}

// Require returns a new Decimal from a string representation or panics if New would have returned an error.
//
// Example:
//
//	d := decimal.Require("123,456")
//	d.String() // "123456"
//
//	d2 := decimal.Require("") // Panic!!!
func Require(s string) Decimal {
	d, err := newDecimal(s)
	if err != nil {
		panic(err)
	}
	return d
}

type Decimal string

// String return string from Decimal
func (d Decimal) String() string {
	return string(verify(d))
}

// Abs returns the absolute value of the decimal.
func (d Decimal) Abs() Decimal {
	d = verify(d)

	if d[0] == '-' {
		return d[1:]
	}

	return d
}

// Neg returns -d
//
// Example:
//
//	d, _ := decimal.New("123.456")
//	d.Neg().String() // "-123.45"
func (d Decimal) Neg() Decimal {
	d = verify(d)

	if d[0] == '-' {
		return d[1:]
	}

	return "-" + d
}

// Truncate truncates off digits from the number, without rounding.
//
// NOTE: precision is the last digit that will not be truncated (must be >= 0).
//
// Example:
//
//	d, _ := decimal.New("123.456")
//	d.Truncate(2).String() // "123.45"
func (d Decimal) Truncate(i int) Decimal {
	d = verify(d)
	if i < 0 {
		return d
	}

	index := -1
	for j := range d {
		if d[j] == '.' {
			index = j
			break
		}
	}

	if index == -1 {
		return d
	}

	if i == 0 {
		return Decimal(d[:index])
	}

	p := index + i + 1
	if p >= len(d) {
		return d
	}
	return Decimal(d[:p])
}

const (
	_zero        = "0"
	_zeroDot     = "0."
	_zeroDotZero = "0.0"
	_dotZero     = ".0"
)

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
func (d Decimal) Shift(shift int) Decimal {
	d = verify(d)
	s := string(d)
	switch s {
	case _zero, _zeroDot, _zeroDotZero, _dotZero:
		return zero
	}

	if shift == 0 {
		return d
	}

	if shift > 0 {
		return shiftPositive(s, shift)
	}

	return shiftNegative(s, -shift)
}

func combineStrings(ss ...string) string {
	builder := strings.Builder{}

	for _, s := range ss {
		builder.WriteString(s)
	}

	return builder.String()
}

// shiftPositive shifts decimal left. example: 3 to 300
func shiftPositive(s string, shift int) Decimal {
	ss := strings.Split(s, ".")

	var isMinus bool
	if len(ss[0]) != 0 && ss[0][0] == '-' {
		isMinus = true
		ss[0] = ss[0][1:]
	}

	switch len(ss) {
	case 1:
		return Decimal(combineStrings(prefix(isMinus), ss[0], strings.Repeat(_zero, shift)))
	case 2:
		builder := strings.Builder{}

		var (
			prefixes = ss[0]
			suffixes = ss[1]
		)
		builder.Reset()
		builder.Grow(len(prefixes) + len(suffixes) + shift + 2)
		builder.WriteString(prefix(isMinus))
		builder.WriteString(prefixes)
		if len(suffixes) > shift {
			builder.WriteString(suffixes[:shift])
			builder.WriteByte('.')
			builder.WriteString(suffixes[shift:])
			return Decimal(tidyString(builder.String()))
		} else {
			// When shift >= len(suffixes), result is an integer
			builder.WriteString(suffixes)
			builder.WriteString(strings.Repeat(_zero, shift-len(suffixes)))
			return Decimal(tidyString(builder.String(), true)) // noDecimalPoint = true
		}
	default:
		return zero
	}
}

// shiftNegative shifts decimal right (division by 10^shift). example: 300 to 3.00
func shiftNegative(s string, shift int) Decimal {
	ss := strings.Split(s, ".")

	var isMinus bool
	if len(ss[0]) != 0 && ss[0][0] == '-' {
		isMinus = true
		ss[0] = ss[0][1:]
	}

	builder := strings.Builder{}

	switch len(ss) {
	case 1:
		// Integer case: e.g., "12345" shift 3 -> "12.345"
		intPart := ss[0]
		builder.Reset()
		builder.Grow(len(intPart) + shift + 3)
		builder.WriteString(prefix(isMinus))

		if len(intPart) <= shift {
			// e.g., "123" shift 5 -> "0.00123"
			builder.WriteString("0.")
			builder.WriteString(strings.Repeat("0", shift-len(intPart)))
			builder.WriteString(intPart)
		} else {
			// e.g., "12345" shift 3 -> "12.345"
			builder.WriteString(intPart[:len(intPart)-shift])
			builder.WriteByte('.')
			builder.WriteString(intPart[len(intPart)-shift:])
		}
		return Decimal(tidyString(builder.String()))

	case 2:
		// Decimal case: e.g., "123.456" shift 2 -> "1.23456"
		intPart := ss[0]
		fracPart := ss[1]

		// Combine all digits
		allDigits := intPart + fracPart

		builder.Reset()
		builder.Grow(len(allDigits) + shift + 3)
		builder.WriteString(prefix(isMinus))

		if len(allDigits) <= shift {
			// e.g., "1.23" shift 5 -> "0.0000123"
			builder.WriteString("0.")
			builder.WriteString(strings.Repeat("0", shift-len(allDigits)))
			builder.WriteString(allDigits)
		} else {
			// Insert decimal point from the right
			// e.g., "10012345678.9" -> "100123456789", shift 8 -> "100.123456789"
			splitPos := len(allDigits) - shift - 1
			builder.WriteString(allDigits[:splitPos])
			builder.WriteByte('.')
			builder.WriteString(allDigits[splitPos:])
		}
		return Decimal(tidyString(builder.String()))

	default:
		return Zero()
	}
}

// Add return d + d2
//
// Example:
//
//	d1, _ := decimal.New("100")
//	d2, _ := decimal.New("90.99")
//	d1.Add(d2).String() // "190.01"
func (d Decimal) Add(d2 Decimal) Decimal {
	d, d2 = verify(d), verify(d2)

	b, a := []byte(d), []byte(d2.String())
	baseNegative := b[0] == '-'
	additionNegative := a[0] == '-'
	if baseNegative && additionNegative {
		return "-" + Decimal(unsignedAdd(b[1:], a[1:]))
	}

	if baseNegative {
		return Decimal(unsignedSub(a, b[1:]))
	}

	if additionNegative {
		return Decimal(unsignedSub(b, a[1:]))
	}

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
	d, d2 = verify(d), verify(d2)

	b, a := []byte(d), []byte(d2.String())
	baseNegative := b[0] == '-'
	additionNegative := a[0] == '-'
	if baseNegative && additionNegative {
		return Decimal(unsignedSub(a[1:], b[1:]))
	}

	if baseNegative {
		return "-" + Decimal(unsignedAdd(b[1:], a))
	}

	if additionNegative {
		return Decimal(unsignedAdd(b, a[1:]))
	}

	return Decimal(unsignedSub(b, a))
}

// unsignedAdd add two unsigned string with shifting
//
//	example: 1234.001 add 12.00001
//	1234.001**
//	**12.00001
//	         ^ // <- pointer go forward
//	example: 1234.00001 add 12.1
//	1234.00001
//	**12.1****
//		     ^ // <- pointer go forward
func unsignedAdd(base, addition []byte) string {
	b, bDecimalPoint := findOrInsertDecimalPoint(base)
	a, aDecimalPoint := findOrInsertDecimalPoint(addition)

	maxLenAfterDecimalPoint := max(len(b)-bDecimalPoint-1, len(a)-aDecimalPoint-1)
	maxP := max(bDecimalPoint, aDecimalPoint)

	resultLen := maxP + maxLenAfterDecimalPoint + 2 // +2 for carry and decimal point
	result := make([]byte, resultLen)

	p := maxP + maxLenAfterDecimalPoint
	bShifting := maxP - bDecimalPoint
	aShifting := maxP - aDecimalPoint

	var (
		delta        byte
		bChar, aChar byte
		bP, aP       int
		resultIdx    int = resultLen - 1
	)

	for ; p >= 0; p-- {
		bChar, aChar = '0', '0'
		if bP = p - bShifting; bP >= 0 && bP < len(b) {
			bChar = b[bP]
		}
		if aP = p - aShifting; aP >= 0 && aP < len(a) {
			aChar = a[aP]
		}

		if bChar == '.' {
			result[resultIdx] = '.'
			resultIdx--
			continue
		}

		delta += (bChar - '0') + (aChar - '0')
		if delta <= 9 {
			result[resultIdx] = delta + '0'
			delta = 0
		} else {
			result[resultIdx] = delta - 10 + '0'
			delta = 1
		}
		resultIdx--
	}

	if delta > 0 {
		result[resultIdx] = delta + '0'
		return tidy(result[resultIdx:])
	}

	return tidy(result[resultIdx+1:])
}

// unsignedSub subtract two unsigned string with shifting
//
//	example: 1234.001 sub 12.00001
//	1234.001**
//	**12.00001
//	         ^ // <- pointer go forward
//	example: 1234.00001 sub 12.1
//	1234.00001
//	**12.1****
//		     ^ // <- pointer go forward
func unsignedSub(base, subtraction []byte) string {
	b, bDecimalPoint := findOrInsertDecimalPoint(base)
	s, sDecimalPoint := findOrInsertDecimalPoint(subtraction)

	maxLenAfterDecimalPoint := max(len(b)-bDecimalPoint-1, len(s)-sDecimalPoint-1)
	maxP := max(bDecimalPoint, sDecimalPoint)

	resultLen := maxP + maxLenAfterDecimalPoint + 1 // +1 for decimal point
	result := make([]byte, resultLen)

	p := maxP + maxLenAfterDecimalPoint
	bShifting := maxP - bDecimalPoint
	sShifting := maxP - sDecimalPoint

	// Quick check: if shifting difference indicates subtraction is larger
	if sShifting < bShifting {
		return "-" + unsignedSub(s, b)
	}

	var (
		bChar, sChar byte
		bP, sP       int
		resultIdx    int = resultLen - 1
		borrow       int8
	)

	// If equal shifting, need to compare digit by digit to determine sign
	if sShifting == bShifting {
		count := max(len(b)+bShifting, len(s)+sShifting)
		for p := 0; p < count; p++ {
			bChar, sChar = '0', '0'
			if bP = p - bShifting; bP >= 0 && bP < len(b) {
				bChar = b[bP]
			}
			if sP = p - sShifting; sP >= 0 && sP < len(s) {
				sChar = s[sP]
			}

			if bChar == sChar || bChar == '.' {
				continue
			}

			if sChar > bChar {
				return "-" + unsignedSub(s, b)
			}
			break
		}
	}

	// Perform subtraction from right to left
	for ; p >= 0; p-- {
		bChar, sChar = '0', '0'
		if bP = p - bShifting; bP >= 0 && bP < len(b) {
			bChar = b[bP]
		}
		if sP = p - sShifting; sP >= 0 && sP < len(s) {
			sChar = s[sP]
		}

		if bChar == '.' {
			result[resultIdx] = '.'
			resultIdx--
			continue
		}

		diff := int8(bChar-'0') - int8(sChar-'0') - borrow
		if diff < 0 {
			result[resultIdx] = byte(10+diff) + '0'
			borrow = 1
		} else {
			result[resultIdx] = byte(diff) + '0'
			borrow = 0
		}
		resultIdx--
	}

	return tidy(result[resultIdx+1:])
}

// findOrInsertDecimalPoint find the index of decimal point. (if no decimal point, it will insert it into the end of the number)
//
// return inserted number and index of decimal point
func findOrInsertDecimalPoint(num []byte) ([]byte, int) {
	for i, b := range num {
		if b == '.' {
			return num, i
		}
	}
	// No decimal point found, append it
	num = append(num, '.')
	return num, len(num) - 1
}

// clean the zero and dot of prefixes and suffixes
func tidy(num []byte, noDecimalPoint ...bool) string {
	return tidyString(string(num), noDecimalPoint...)
}

// clean the zero and dot of prefixes and suffixes
func tidyString(num string, noDecimalPoint ...bool) string {
	if len(num) == 0 {
		return _zero
	}

	// Handle sign prefix
	var sign string
	start := 0
	switch num[0] {
	case '+':
		start = 1
	case '-':
		sign = "-"
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

	// Determine if we should remove decimal point
	noDP := len(noDecimalPoint) > 0 && noDecimalPoint[0]

	// Find end position
	end := len(num)
	if !noDP {
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
		if sign != "" {
			return sign + _zero + result
		}
		return _zero + result
	}

	// Handle single decimal point
	if result == "." {
		return _zero
	}

	return sign + result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// prefix return '-' when isMinus is true, return "" when isMinus is false
func prefix(isMinus bool) string {
	if isMinus {
		return "-"
	}
	return ""
}

// verify makes sure Decimal is valid for calculation
func verify(d Decimal) Decimal {
	dd, err := newDecimal(string(d))
	if err != nil {
		panic(err)
	}

	return dd
}

// IsZero return d == 0
func (d Decimal) IsZero() bool {
	if len(d) == 0 {
		return true
	}

	d = verify(d)

	for _, c := range d {
		switch c {
		case '0', '.':
			continue
		default:
			return false
		}
	}
	return true
}

func (d Decimal) isZero() bool {
	if len(d) == 0 {
		return true
	}

	for _, c := range d {
		switch c {
		case '0', '.':
			continue
		default:
			return false
		}
	}
	return true
}

// IsPositive return d > 0
func (d Decimal) IsPositive() bool {
	return verify(d).isPositive()
}

func (d Decimal) isPositive() bool {
	return !d.isZero() && !d.isNegative()
}

// IsNegative return d < 0
func (d Decimal) IsNegative() bool {
	return verify(d).isNegative()
}

func (d Decimal) isNegative() bool {
	return d[0] == '-'
}

// Equal return d == d2
func (d Decimal) Equal(d2 Decimal) bool {
	return verify(d) == verify(d2)
}

// Greater return d > d2
func (d Decimal) Greater(d2 Decimal) bool {
	return greater(verify(d), verify(d2))
}

// greater return true if the d1 > d2
//
//	example: 1234.001 vs 12.00001
//	1234.001**
//	**12.00001
//	^ // <- pointer go backward
//	example: 1234.00001 vs 12.1
//	1234.00001
//	**12.1****
//	^ // <- pointer go backward
func greater(d1, d2 Decimal) bool {
	if d1.isPositive() && d2.isNegative() {
		return true
	}

	if d1.isNegative() && d2.isPositive() {
		return false
	}

	fb := []byte(d1)
	sb := []byte(d2)
	if fb[0] == '-' {
		fb = fb[1:]
		sb = sb[1:]
	}
	f, fDecimalPoint := findOrInsertDecimalPoint(fb)
	s, sDecimalPoint := findOrInsertDecimalPoint(sb)

	maxLenAfterDecimalPoint := max(len(f)-fDecimalPoint-1, len(s)-sDecimalPoint-1)

	if fDecimalPoint != sDecimalPoint {
		return fDecimalPoint > sDecimalPoint
	}

	count := fDecimalPoint + maxLenAfterDecimalPoint + 1
	for i := 0; i < count; i++ {
		if len(f) == i {
			return false
		}

		if len(s) == i {
			return true
		}

		if f[i] == '.' {
			continue
		}

		if f[i] != s[i] {
			return f[i] > s[i]
		}
	}

	return false
}

// Less return d < d2
func (d Decimal) Less(d2 Decimal) bool {
	return less(verify(d), verify(d2))
}

// less return true if the d1 < d2
//
//	example: 1234.001 vs 12.00001
//	1234.001**
//	**12.00001
//	^ // <- pointer go backward
//	example: 1234.00001 vs 12.1
//	1234.00001
//	**12.1****
//	^ // <- pointer go backward
func less(d1, d2 Decimal) bool {
	if d1.isNegative() && d2.isPositive() {
		return true
	}

	if d1.isPositive() && d2.isNegative() {
		return false
	}

	fb := []byte(d1)
	sb := []byte(d2)
	if fb[0] == '-' {
		fb = fb[1:]
		sb = sb[1:]
	}
	f, fDecimalPoint := findOrInsertDecimalPoint(fb)
	s, sDecimalPoint := findOrInsertDecimalPoint(sb)

	maxLenAfterDecimalPoint := max(len(f)-fDecimalPoint-1, len(s)-sDecimalPoint-1)

	if fDecimalPoint != sDecimalPoint {
		return fDecimalPoint < sDecimalPoint
	}

	count := fDecimalPoint + maxLenAfterDecimalPoint + 1
	for i := 0; i < count; i++ {
		if len(f) == i {
			return true
		}

		if len(s) == i {
			return false
		}

		if f[i] == '.' {
			continue
		}

		if f[i] != s[i] {
			return f[i] < s[i]
		}
	}

	return false
}

// GreaterOrEqual return d >= d2
func (d Decimal) GreaterOrEqual(d2 Decimal) bool {
	return !d.Less(d2)
}

// LessOrEqual return d <= d2
func (d Decimal) LessOrEqual(d2 Decimal) bool {
	return !d.Greater(d2)
}

// Mul return d * d2
func (d Decimal) Mul(d2 Decimal) Decimal {
	if d.isZero() || d2.isZero() {
		return Zero()
	}

	// d = verify(d)
	// d2 = verify(d2)

	a, right1 := removeDecimalPoint(string(d))
	b, right2 := removeDecimalPoint(string(d2))

	minus := false
	if a[0] == '-' {
		a = a[1:]
		minus = !minus
	}

	if b[0] == '-' {
		b = b[1:]
		minus = !minus
	}

	multiplied := multiplyPureNumber(a, b)
	idx := right1 + right2
	if idx == 0 {
		return Decimal(prefix(minus) + tidy(multiplied, true))
	}

	idx = len(multiplied) - idx
	buf := make([]byte, 0, len(multiplied)+1)
	buf = append(buf, multiplied[:idx]...)
	buf = append(buf, '.')
	buf = append(buf, multiplied[idx:]...)

	return Decimal(prefix(minus) + tidy(buf))
}

// removeDecimalPoint removes decimal point and return the count of the digit right the decimal
func removeDecimalPoint(s string) (result []byte, countOfRightSide int) {
	idx := -1
	for i := range s {
		if s[i] == '.' {
			idx = i
			break
		}
	}
	if idx == -1 {
		return []byte(s), 0
	}
	return append([]byte(s[:idx]), s[idx+1:]...), len(s) - idx - 1
}

// multiplyPureNumber return d1 * d2, d1 & d2 must contain only number 0~9
func multiplyPureNumber(d1 []byte, d2 []byte) []byte {
	if len(d1) < len(d2) {
		d1, d2 = d2, d1
	}

	len1, len2 := len(d1), len(d2)
	result := make([]byte, len1+len2)

	// Optimized multiplication using single loop with carry propagation
	for i := len2 - 1; i >= 0; i-- {
		val2 := d2[i] - '0'
		if val2 == 0 {
			continue
		}

		carry := byte(0)
		resultIdx := i + len1

		// Inner multiplication loop
		for j := len1 - 1; j >= 0; j-- {
			val1 := d1[j] - '0'
			prod := val1*val2 + result[resultIdx] + carry
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

	return result
}

// Div returns d / d2. If it doesn't divide exactly, the result will have DivisionPrecision digits after the decimal point.
func (d Decimal) Div(d2 Decimal) Decimal {
	decimal.DivisionPrecision = DivisionPrecision
	return Decimal(decimal.RequireFromString(string(d)).Div(decimal.RequireFromString(string(d2))).String())
}
