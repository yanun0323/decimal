package decimal

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var (
	DivisionPrecision = 16
)

var (
	_isFirstSymbolDroppable = map[byte]bool{
		'0': false,
		'1': false,
		'2': false,
		'3': false,
		'4': false,
		'5': false,
		'6': false,
		'7': false,
		'8': false,
		'9': false,
		'.': false,
		'-': false,
		'+': true,
	}

	_isSymbolDroppable = map[byte]bool{
		'0': false,
		'1': false,
		'2': false,
		'3': false,
		'4': false,
		'5': false,
		'6': false,
		'7': false,
		'8': false,
		'9': false,
		'.': false,
		'_': true,
		',': true,
	}
)

// Zero return the zero decimal
func Zero() Decimal {
	return Decimal("0")
}

// New create a Decimal.
//
// acceptable symbol (+-.,_0123456789)
//
// Example:
//
//	d, err := decimal.New("123,456,789.000")
func New(s string) (Decimal, error) {
	if len(s) == 0 {
		return Zero(), nil
	}

	buf := []byte(s)
	dot := false
	isSymbolDroppable := _isFirstSymbolDroppable
	for i, b := 0, byte(0); i < len(buf); isSymbolDroppable = _isSymbolDroppable {
		b = buf[i]
		droppable, valid := isSymbolDroppable[b]
		if !valid {
			return Zero(), errors.New(fmt.Sprintf("invalid symbol (%s) in %s", string(b), s))
		}

		if b == '.' {
			if dot {
				return Zero(), errors.New("duplicate dot")
			}
			dot = true
		}
		if droppable {
			buf = append(buf[:i], buf[i+1:]...)
		} else {
			i++
		}
	}

	if len(buf) == 0 {
		return Zero(), errors.New("can't convert to Decimal empty string")
	}

	inserted, _ := findOrInsertDecimalPoint(buf)
	return Decimal(tidy(inserted)), nil
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
	d, err := New(s)
	if err != nil {
		panic(err)
	}
	return d
}

type Decimal string

// String return string from Decimal
func (d Decimal) String() string {
	d = blanker(d)
	return string(d)
}

// Abs returns the absolute value of the decimal.
func (d Decimal) Abs() Decimal {
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
	if d.IsZero() {
		return d
	}

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
	d = blanker(d)
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
	d = blanker(d)
	s := string(d)
	switch s {
	case "0", "0.", "0.0":
		return Zero()
	}

	if shift == 0 {
		return d
	}

	if shift > 0 {
		return shiftPositive(string(d), shift)
	}

	return shiftNegative(string(d), -shift)
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
		return Decimal(prefix(isMinus) + ss[0] + strings.Repeat("0", shift))
	case 2:
		prefixes, suffixes := []byte(ss[0]), append([]byte(ss[1]), bytes.Repeat([]byte{'0'}, shift)...)
		buf := make([]byte, 0, len(prefixes)+len(suffixes)+1)
		buf = append(buf, prefixes...)
		buf = append(buf, suffixes[:shift]...)
		buf = append(buf, '.')
		buf = append(buf, suffixes[shift:]...)
		return Decimal(prefix(isMinus) + tidy(buf))
	default:
		return "0"
	}
}

// shiftNegative shifts decimal left. example: 3 to 0.03
func shiftNegative(s string, shift int) Decimal {
	ss := strings.Split(s, ".")

	var isMinus bool
	if len(ss[0]) != 0 && ss[0][0] == '-' {
		isMinus = true
		ss[0] = ss[0][1:]
	}

	switch len(ss) {
	case 1:
		buf := make([]byte, 0, len(ss[0])+shift)
		if len(ss[0]) <= shift {
			buf = append(buf, '0', '.')
			buf = append(buf, []byte(strings.Repeat("0", shift-len(ss[0])))...)
			buf = append(buf, ss[0]...)
		} else {
			buf = append(buf, ss[0][:len(ss[0])-shift]...)
			buf = append(buf, '.')
			buf = append(buf, ss[0][len(ss[0])-shift:]...)
		}
		return Decimal(prefix(isMinus) + tidy(buf))
	case 2:
		zeros := bytes.Repeat([]byte{'0'}, shift)
		prefixes, suffixes := append(zeros, ss[0]...), []byte(ss[1])
		buf := make([]byte, 0, len(prefixes)+len(suffixes)+1)
		buf = append(buf, prefixes[:len(prefixes)-shift]...)
		buf = append(buf, '.')
		buf = append(buf, prefixes[len(prefixes)-shift:]...)
		buf = append(buf, suffixes...)
		return Decimal(prefix(isMinus) + tidy(buf))
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
	d = blanker(d)
	d2 = blanker(d2)

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
	d = blanker(d)
	d2 = blanker(d2)

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

	p := maxP + maxLenAfterDecimalPoint
	bShifting := (maxP - bDecimalPoint)
	aShifting := (maxP - aDecimalPoint)

	var (
		delta        byte
		buf          bytes.Buffer
		bChar, aChar byte
		bP, aP       int
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
			buf.WriteByte('.')
			continue
		}

		delta += (bChar - '0') + (aChar - '0')
		if delta <= 9 {
			buf.WriteByte(delta + '0')
			delta = 0
		} else {
			buf.WriteByte(delta - 10 + '0')
			delta = 1
		}

	}
	buf.WriteByte(delta + '0')

	reversed := buf.Bytes()
	result := make([]byte, 0, len(reversed))
	for i := len(reversed) - 1; i >= 0; i-- {
		result = append(result, reversed[i])
	}

	return tidy(result)
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

	bShifting := (maxP - bDecimalPoint)
	sShifting := (maxP - sDecimalPoint)
	if sShifting < bShifting {
		return "-" + unsignedSub(s, b)
	}

	var (
		bChar, sChar byte
		bP, sP       int
	)
	if sShifting == bShifting {
		count := max(len(b)+bShifting, len(s)+sShifting)
		isSubtractionLarger := false
		for p := 0; p < count; p++ {
			bChar, sChar = '0', '0'
			if bP = p - bShifting; bP >= 0 && bP < len(b) {
				bChar = b[bP]
			}
			if sP = p - sShifting; sP >= 0 && sP < len(s) {
				sChar = s[sP]
			}

			if bChar == sChar {
				continue
			}

			isSubtractionLarger = sChar > bChar
			break
		}

		if isSubtractionLarger {
			return "-" + unsignedSub(s, b)
		}
	}
	var (
		delta int8
		buf   bytes.Buffer
	)
	p := maxP + maxLenAfterDecimalPoint
	for ; p >= 0; p-- {
		bChar, sChar = '0', '0'
		if bP = p - bShifting; bP >= 0 && bP < len(b) {
			bChar = b[bP]
		}
		if sP = p - sShifting; sP >= 0 && sP < len(s) {
			sChar = s[sP]
		}

		if bChar == '.' {
			buf.WriteByte('.')
			continue
		}

		delta += int8(bChar-'0') - int8(sChar-'0')
		if delta < 0 {
			buf.WriteByte(byte(10+delta) + '0')
			delta = -1
		} else {
			buf.WriteByte(byte(delta) + '0')
			delta /= 10
		}
	}
	reversed := buf.Bytes()
	result := make([]byte, 0, len(reversed))
	for i := len(reversed) - 1; i >= 0; i-- {
		result = append(result, reversed[i])
	}

	return tidy(result)
}

// findOrInsertDecimalPoint find the index of decimal point. (if no decimal point, it will insert it into the end of the number)
//
// return inserted number and index of decimal point
func findOrInsertDecimalPoint(num []byte) ([]byte, int) {
	p := 0
	for range num {
		if num[p] == '.' {
			continue
		}
		p++
	}
	if p == len(num) {
		num = append(num, '.')
	}
	return num, p
}

// clean the zero and dot of prefixes and suffixes
func tidy(num []byte, noDecimalPoint ...bool) string {
	if len(num) == 0 {
		return "0"
	}

	for num[0] == '0' {
		num = num[1:]
		if len(num) == 0 {
			return "0"
		}
	}

	noDP := false
	if len(noDecimalPoint) != 0 {
		noDP = noDecimalPoint[0]
	}

	if noDP {
		return string(num)
	}

	for num[len(num)-1] == '0' {
		num = num[:len(num)-1]
		if len(num) == 0 {
			return "0"
		}
	}

	if len(num) == 1 && num[0] == '.' {
		return "0"
	}

	if num[len(num)-1] == '.' {
		num = num[:len(num)-1]
	}

	if num[0] == '.' {
		return "0" + string(num)
	}
	return string(num)
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

// blanker makes sure Decimal is not empty string
func blanker(d Decimal) Decimal {
	if len(d) == 0 {
		return Zero()
	}
	return d
}

// IsZero return d == 0
func (d Decimal) IsZero() bool {
	return len(d) == 0 || d.String() == "0"
}

// IsPositive return d > 0
func (d Decimal) IsPositive() bool {
	return d[0] != '-' && !d.IsZero()
}

// IsNegative return d < 0
func (d Decimal) IsNegative() bool {
	return d[0] == '-'
}

// Equal return d == d2
func (d Decimal) Equal(d2 Decimal) bool {
	return blanker(d) == blanker(d2)
}

// Greater return d > d2
func (d Decimal) Greater(d2 Decimal) bool {
	return greater(blanker(d), blanker(d2))
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
	if d1.IsPositive() && d2.IsNegative() {
		return true
	}

	if d1.IsNegative() && d2.IsPositive() {
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
	return less(blanker(d), blanker(d2))
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
	if d1.IsNegative() && d2.IsPositive() {
		return true
	}

	if d1.IsPositive() && d2.IsNegative() {
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
	return d.Equal(d2) || d.Greater(d2)
}

// LessOrEqual return d <= d2
func (d Decimal) LessOrEqual(d2 Decimal) bool {
	return d.Equal(d2) || d.Less(d2)
}

// Mul return d * d2
func (d Decimal) Mul(d2 Decimal) Decimal {
	if d.IsZero() || d2.IsZero() {
		return Zero()
	}

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
	idx := (right1 + right2)
	if idx == 0 {
		return Decimal(prefix(minus) + tidy(multiplied, true))
	}

	idx = len(multiplied) - idx
	buf := make([]byte, 0, len(multiplied)+1)
	buf = append(buf, multiplied[:idx]...)
	buf = append(buf, '.')
	buf = append(buf, multiplied[idx:]...)

	// buf := strings.Builder{}
	// buf.Write(multiplied[:idx])
	// buf.WriteByte('.')
	// buf.Write(multiplied[idx:])

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
		return multiplyPureNumber(d2, d1)
	}
	result := make([]byte, len(d1)+len(d2))

	var (
		symbol2, val2, cache, temp byte
		idx, j                     int
	)

	for i := len(d2) - 1; i >= 0; i-- {
		symbol2 = d2[i]
		if symbol2 == '0' {
			continue
		}
		val2 = (symbol2 - '0')
		idx = i + len(d1)
		j = idx - i - 1
		for ; idx >= 0; idx, j = idx-1, j-1 {
			if j >= 0 {
				cache += (d1[j]-'0')*val2 + result[idx]
			} else if cache == 0 {
				break
			}
			temp = cache / 10
			result[idx] = byte(cache - temp*10)
			cache = temp
		}
	}

	for i, v := range result {
		result[i] = v + '0'
	}

	return result
}

// Div returns d / d2. If it doesn't divide exactly, the result will have DivisionPrecision digits after the decimal point.
func (d Decimal) Div(d2 Decimal) Decimal {
	decimal.DivisionPrecision = DivisionPrecision
	return Decimal(decimal.RequireFromString(string(d)).Div(decimal.RequireFromString(string(d2))).String())
}

func (d Decimal) DivWorking(d2 Decimal) Decimal {
	decimal.DivisionPrecision = DivisionPrecision
	d1Pure, rightCount1 := removeDecimalPoint(string(d))
	d2Pure, rightCount2 := removeDecimalPoint(string(d2))

	offset := (len(d) - rightCount1 - 1) - (len(d2) - rightCount2 - 1)
	if offset > 0 {
		d2Pure = append(d2Pure, bytes.Repeat([]byte{'0'}, offset)...)
	} else {
		d1Pure = append(d1Pure, bytes.Repeat([]byte{'0'}, -offset)...)
	}

	var findTimes func(t int)
	findTimes = func(t int) {

	}
	_ = findTimes
	return d
}
