package decimal

import (
	"strings"
)

/*
在不影响逻辑正确性的前提下, 优化此函数，让他运行的更快速、消耗更少记忆体
*/

var (
	DivisionPrecision = 16
	zero              = Decimal("0")
	one               = Decimal("1")
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
	buf, err := newDecimal([]byte(s))
	if err != nil {
		return zero, err
	}

	return Decimal(buf), nil
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
	return Decimal(normalize([]byte(s)))
}

type Decimal string

// String return string from Decimal
func (d Decimal) String() string {
	return string(normalize([]byte(d)))
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
//	d, _ := decimal.New("123.456")
//	d.Neg().String() // "-123.45"
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
//	d, _ := decimal.New("123.456")
//	d.Truncate(2).String() // "123.45"
func (d Decimal) Truncate(i int) Decimal {
	return Decimal(truncate(normalize([]byte(d)), i))
}

func (d Decimal) truncate(i int) Decimal {
	dotIdx := -1
	for j := range d {
		if d[j] == '.' {
			dotIdx = j
			break
		}
	}

	if i < 0 { // positive
		i = -i
		if dotIdx != -1 {
			d = d[:dotIdx]
		}

		if i == 0 {
			return d
		}

		if i >= len(d) {
			return zero
		}

		return Decimal(combineToDecimal(string(d[:len(d)-i]), strings.Repeat("0", i)))
	}

	// negative
	if dotIdx == -1 {
		return d
	}

	if i == 0 {
		return d[:dotIdx]
	}

	p := dotIdx + i + 1
	if p >= len(d) {
		return d
	}
	return d[:p]
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
func (d Decimal) Shift(sf int) Decimal {
	return Decimal(shift(normalize([]byte(d)), sf))
}

func (d Decimal) shift(shift int) Decimal {
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
		return combineToDecimal(prefix(isMinus), ss[0], strings.Repeat(_zero, shift))
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
			return tidyString(builder.String())
		} else {
			// When shift >= len(suffixes), result is an integer
			builder.WriteString(suffixes)
			builder.WriteString(strings.Repeat(_zero, shift-len(suffixes)))
			return tidyString(builder.String())
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
		return tidyString(builder.String())

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
		return tidyString(builder.String())

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
	b, a := normalize([]byte(d)), normalize([]byte(d2))
	baseNegative := b[0] == '-'
	additionNegative := a[0] == '-'
	if baseNegative && additionNegative {
		b = trimFront(b, 1)
		a = trimFront(a, 1)
		// -b - -a = -b + a = a - b
		return Decimal(unsignedSub(a, b))
	}

	if baseNegative {
		b = trimFront(b, 1)
		// -b - a = - (b + a)
		return Decimal(pushFront(unsignedAdd(a, b), '-'))
	}

	if additionNegative {
		a = trimFront(a, 1)
		// b - -a = b + a
		return Decimal(unsignedAdd(b, a))
	}

	// b - a = b - a
	return Decimal(unsignedSub(b, a))
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

// prefix return '-' when isMinus is true, return "" when isMinus is false
func prefix(isMinus bool) string {
	if isMinus {
		return "-"
	}
	return ""
}

// DecimalPattern is a regex pattern to validate decimal string format
//
// Rules:
// 1. First character can only contain "+-1234567890"
// 2. Last character can only contain "1234567890"
// 3. Middle characters can only contain "1234567890."
// 4. Cannot have consecutive "."
// var decimalPattern = regexp.MustCompile(`^[+\-0-9]([0-9]*\.?[0-9]*)*[0-9]$|^[+\-]?[0-9]$|^[+\-]?[0-9]*\.[0-9]+$`)

// verify makes sure Decimal is valid for calculation
func verify(d Decimal) Decimal {
	dd, err := newDecimal([]byte(d))
	if err != nil {
		panic(err)
	}

	return Decimal(dd)
}

// IsZero return d == 0
func (d Decimal) IsZero() bool {
	if len(d) == 0 {
		return true
	}

	return isZero(normalize([]byte(d)))
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
	buf := normalize([]byte(d))

	return !isZero(buf) && !isNegative(buf)
}

func (d Decimal) isPositive() bool {
	return isPositive(normalize([]byte(d)))
}

// IsNegative return d < 0
func (d Decimal) IsNegative() bool {
	return isNegative([]byte(d))
}

func (d Decimal) isNegative() bool {
	return d[0] == '-'
}

// Equal return d == d2
func (d Decimal) Equal(d2 Decimal) bool {
	return string(normalize([]byte(d))) == string(normalize([]byte(d2)))
}

// Greater return d > d2
func (d Decimal) Greater(d2 Decimal) bool {
	return greater(normalize([]byte(d)), normalize([]byte(d2)))
}

// Less return d < d2
func (d Decimal) Less(d2 Decimal) bool {
	return less(normalize([]byte(d)), normalize([]byte(d2)))
}

// GreaterOrEqual return d >= d2
func (d Decimal) GreaterOrEqual(d2 Decimal) bool {
	return !less(normalize([]byte(d)), normalize([]byte(d2)))
}

// LessOrEqual return d <= d2
func (d Decimal) LessOrEqual(d2 Decimal) bool {
	return !greater(normalize([]byte(d)), normalize([]byte(d2)))
}

// Mul return d * d2
func (d Decimal) Mul(d2 Decimal) Decimal {
	a := normalize([]byte(d))
	b := normalize([]byte(d2))

	if isZero(a) || isZero(b) {
		return zero
	}

	right1 := findDotIndex(a)
	if right1 == -1 {
		right1 = 0
	} else {
		a = remove(a, right1)
		right1 = len(a) - right1
	}

	right2 := findDotIndex(b)
	if right2 == -1 {
		right2 = 0
	} else {
		b = remove(b, right2)
		right2 = len(b) - right2
	}

	minus := false
	if a[0] == '-' {
		a = trimFront(a, 1)
		minus = !minus
	}

	if b[0] == '-' {
		b = trimFront(b, 1)
		minus = !minus
	}

	// 200ns
	multiplied := multiplyPureNumber(a, b)
	idx := right1 + right2
	if idx == 0 {
		if minus {
			return "-" + Decimal(multiplied)
		}

		return Decimal(multiplied)
	}

	idx = len(multiplied) - idx
	multiplied = insert(multiplied, idx, '.')

	if minus {
		return "-" + Decimal(multiplied)
	}

	return Decimal(multiplied)
}

// removeDecimalPoint removes decimal point and return the count of the digit right the decimal
func removeDecimalPoint(s []byte) (result []byte, countOfRightSide int) {
	for i := range s {
		if s[i] == '.' {
			return remove(s, i), len(s) - i - 1
		}
	}
	return s, 0
}

// multiplyPureNumber return d1 * d2, d1 & d2 must contain only number 0~9
func multiplyPureNumber(d1 []byte, d2 []byte) []byte {
	if len(d1) < len(d2) {
		d1, d2 = d2, d1
	}

	var (
		extraCap   = 2
		len1, len2 = len(d1), len(d2)
		result     = make([]byte, len1+len2+extraCap)
		resultIdx  int
		carry      byte
		val1, val2 byte
		prod       byte
	)

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

	return trimBack(result, extraCap)
}
