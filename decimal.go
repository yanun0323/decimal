package decimal

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
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

// NewDecimal create a Decimal.
//
// acceptable symbol (+-.,_0123456789)
//
// Example:
//
//	d, err := decimal.NewDecimal("123,456,789.000")
func NewDecimal(s string) (Decimal, error) {
	buf := []byte(s)
	dot := false
	isSymbolDroppable := _isFirstSymbolDroppable
	for i, b := 0, byte(0); i < len(buf); isSymbolDroppable = _isSymbolDroppable {
		b = buf[i]
		droppable, valid := isSymbolDroppable[b]
		if !valid {
			return "", errors.New(fmt.Sprintf("invalid symbol (%s) in %s", string(b), s))
		}

		if b == '.' {
			if dot {
				return "", errors.New("duplicate dot")
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
		return "", errors.New("can't convert to Decimal empty string")
	}

	inserted, _ := findOrInsertDecimalPoint(buf)
	return Decimal(cleanZero(inserted)), nil
}

// RequireDecimal returns a new Decimal from a string representation or panics if NewDecimal would have returned an error.
//
// Example:
//
//	d := decimal.RequireDecimal("123,456")
//	d.String() // "123456"
//
//	d2 := decimal.RequireDecimal("") // Panic!!!
func RequireDecimal(s string) Decimal {
	d, err := NewDecimal(s)
	if err != nil {
		panic(err)
	}
	return d
}

type Decimal string

// String return string from Decimal
func (d Decimal) String() string {
	return string(d)
}

// Truncate truncates off digits from the number, without rounding.
//
// NOTE: precision is the last digit that will not be truncated (must be >= 0).
//
// Example:
//
//	d, _ := decimal.NewDecimal("123.456")
//	d.Truncate(2).String() // "123.45"
func (d Decimal) Truncate(i int) Decimal {
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

// SatoshiToDecimal convert satoshi to decimal.
//
// Example:
//
//	d, _ := decimal.NewDecimal("3")
//	d.SatoshiToDecimal().String() // "300000000"
func (d Decimal) SatoshiToDecimal() Decimal {
	s := string(d)
	switch s {
	case "0", "0.", "0.0":
		return "0"
	}

	ss := strings.Split(s, ".")

	var isMinus bool
	if len(ss[0]) != 0 && ss[0][0] == '-' {
		isMinus = true
		ss[0] = ss[0][1:]
	}

	switch len(ss) {
	case 1:
		buf := make([]byte, 0, len(ss[0])+8)
		if len(ss[0]) <= 8 {
			buf = append(buf, '0', '.')
			buf = append(buf, []byte(strings.Repeat("0", 8-len(ss[0])))...)
			buf = append(buf, ss[0]...)
		} else {
			buf = append(buf, ss[0][:len(ss[0])-8]...)
			buf = append(buf, '.')
			buf = append(buf, ss[0][len(ss[0])-8:]...)
		}
		return Decimal(getMinusString(isMinus) + cleanZero(buf))
	case 2:
		prefix, suffix := append([]byte("00000000"), ss[0]...), []byte(ss[1])
		buf := make([]byte, 0, len(prefix)+len(suffix)+1)
		buf = append(buf, prefix[:len(prefix)-8]...)
		buf = append(buf, '.')
		buf = append(buf, prefix[len(prefix)-8:]...)
		buf = append(buf, suffix...)
		return Decimal(getMinusString(isMinus) + cleanZero(buf))
	default:
		return Decimal("0")
	}
}

// DecimalToSatoshi convert decimal to satoshi.
//
// Example:
//
//	d, _ := decimal.NewDecimal("300_000_000")
//	d.DecimalToSatoshi().String() // "3"
func (d Decimal) DecimalToSatoshi() Decimal {
	s := string(d)
	switch s {
	case "0", "0.", "0.0":
		return "0"
	}

	ss := strings.Split(s, ".")

	var isMinus bool
	if len(ss[0]) != 0 && ss[0][0] == '-' {
		isMinus = true
		ss[0] = ss[0][1:]
	}

	switch len(ss) {
	case 1:
		return Decimal(getMinusString(isMinus) + ss[0] + "00000000")
	case 2:
		prefix, suffix := []byte(ss[0]), append([]byte(ss[1]), "00000000"...)
		buf := make([]byte, 0, len(prefix)+len(suffix)+1)
		buf = append(buf, prefix...)
		buf = append(buf, suffix[:8]...)
		buf = append(buf, '.')
		buf = append(buf, suffix[8:]...)
		return Decimal(getMinusString(isMinus) + cleanZero(buf))
	default:
		return "0"
	}
}

// Add add Decimal
//
// Example:
//
//	d1, _ := decimal.NewDecimal("100")
//	d2, _ := decimal.NewDecimal("90.99")
//	d1.Add(d2).String() // "190.01"
func (base Decimal) Add(addition Decimal) Decimal {
	b, a := []byte(base), []byte(addition)
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

// Sub subtract Decimal
//
// Example:
//
//	d1, _ := decimal.NewDecimal("100")
//	d2, _ := decimal.NewDecimal("90.99")
//	d1.Sub(d2).String() // "9.01"
func (base Decimal) Sub(subtraction Decimal) Decimal {
	b, a := []byte(base), []byte(subtraction)
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

// add two unsigned string with shifting
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

	return cleanZero(result)
}

// subtract two unsigned string with shifting
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

	return cleanZero(result)
}

// find the index of decimal point. (if no decimal point, it will insert it into the end of the number)
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

// clean the zero and dot of prefix and suffix
func cleanZero(num []byte) string {
	if len(num) == 0 {
		return "0"
	}

	for num[0] == '0' {
		num = num[1:]
		if len(num) == 0 {
			return "0"
		}
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

func getMinusString(isMinus bool) string {
	if isMinus {
		return "-"
	}
	return ""
}
