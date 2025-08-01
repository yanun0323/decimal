package decimal

import (
	"errors"
	"fmt"
	"slices"
)

var (
	zeroBytes = []byte{'0'}
)

func quickCheckZero(buf []byte) bool {
	switch string(buf) {
	case "",
		"0", ".",
		"0.", ".0", "+0", "-0", "00",
		"000", "0.0", ".00", "00.", "+00", "+0.", "-00", "-0.":
		return true
	default:
		return false
	}
}

func newDecimal(buf []byte) ([]byte, error) {
	if len(buf) == 0 {
		return zeroBytes, nil
	}

	if quickCheckZero(buf) {
		return zeroBytes, nil
	}

	dot := false
	firstChar := true
	i := 0
	var b byte

	for i < len(buf) {
		b = buf[i]

		if firstChar {
			// Handle first character
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '-':
				if b == '.' {
					dot = true
				}
			case '+':
				buf = trimFront(buf, 1)
				continue
			default:
				return zeroBytes, fmt.Errorf("invalid symbol (%c) in %s", b, string(buf))
			}
			firstChar = false
		} else {
			// Handle subsequent characters
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			case '.':
				if dot {
					return zeroBytes, errors.New("duplicate dot")
				}
				dot = true
			case '_', ',':
				buf = remove(buf, i)
				continue
			default:
				return zeroBytes, fmt.Errorf("invalid symbol (%c) in %s", b, string(buf))
			}
		}

		i++
	}

	if len(buf) == 0 {
		return zeroBytes, errors.New("can't convert to Decimal empty string")
	}

	return tidyBytes(buf), nil

}

// clean the zero and dot of prefixes and suffixes
func tidyBytes(num []byte) []byte {
	if len(num) == 0 {
		return zeroBytes
	}

	hasDot := slices.Contains(num, '.')

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
		return zeroBytes
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
			return zeroBytes
		}

		// Remove trailing decimal point
		if end > start && num[end-1] == '.' {
			end--
		}
	}

	// Extract the significant part
	result := trim(num, start, end)

	if len(result) == 0 {
		return zeroBytes
	}

	if result[len(result)-1] == '.' {
		result = result[:len(result)-1]
	}

	// Handle leading decimal point
	if len(result) > 0 && result[0] == '.' {
		if sign {
			return pushFront(result, '-', '0')
		}
		return pushFront(result, '0')
	}

	// Handle single decimal point
	if len(result) == 1 && result[0] == '.' {
		return zeroBytes
	}

	if sign {
		return pushFront(result, '-')
	}

	return result
}

func findDotIndex(buf []byte) int {
	for j := range buf {
		if buf[j] == '.' {
			return j
		}
	}

	return -1
}

func normalize(buf []byte) []byte {
	normalized, err := newDecimal(buf)
	if err != nil {
		panic(err)
	}

	return normalized
}

func truncate(buf []byte, i int) []byte {
	dotIdx := findDotIndex(buf)

	if i < 0 { // positive
		i = -i
		if dotIdx != -1 {
			buf = buf[:dotIdx]
		}

		if i == 0 {
			return buf
		}

		if i >= len(buf) {
			return zeroBytes
		}

		{ // fill zeros to the right
			for j := 1; j <= i; j++ {
				buf[len(buf)-j] = '0'
			}
			return buf
		}
	}

	// negative
	if dotIdx == -1 {
		return buf
	}

	if i == 0 {
		return buf[:dotIdx]
	}

	p := dotIdx + i + 1
	if p >= len(buf) {
		return buf
	}
	return buf[:p]
}

// 100.123456789 sf=8
// 100123456789
// idx=3 tidx=11
func shift(buf []byte, sf int) []byte {
	if sf == 0 {
		return buf
	}

	sign := false
	if buf[0] == '-' {
		sign = true
		buf = trimFront(buf, 1)
	}

	// 123.456
	// shift 5  -> 123456 -> 12345600
	// shift -5 -> 123456 -> 0.00123456
	// shift 2  -> 123456 -> 12345.6
	idx := findDotIndex(buf)
	if idx == -1 { // 123456
		idx = len(buf)
	}

	// idx: 3
	buf = remove(buf, idx)

	// targetIdx: 3+5 = 8
	// targetIdx: 3-5 = -2
	// targetIdx: 3+2 = 5
	targetIdx := idx + sf

	if targetIdx >= len(buf) { // 12345600
		buf = pushBackRepeat(buf, '0', targetIdx-len(buf))

		if sign {
			buf = pushFront(buf, '-')
		}

		return buf
	}

	if targetIdx < 0 { // 0.00123456
		if sign {
			buf = pushFrontRepeat(buf, '0', -targetIdx+3)
			buf[0] = '-'
			buf[2] = '.'
		} else {
			buf = pushFrontRepeat(buf, '0', -targetIdx+2)
			buf[1] = '.'
		}

		return buf
	}

	// 12345.6
	buf = insert(buf, targetIdx, '.')

	if sign {
		buf = pushFront(buf, '-')
	}

	return buf
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
func unsignedAdd(b, a []byte) []byte {
	bDecimalPoint := findDotIndex(b)
	if bDecimalPoint == -1 {
		b = pushBack(b, '.')
		bDecimalPoint = len(b) - 1
	}

	aDecimalPoint := findDotIndex(a)
	if aDecimalPoint == -1 {
		a = pushBack(a, '.')
		aDecimalPoint = len(a) - 1
	}

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
		return tidyBytes(trimFront(result, resultIdx))
	}

	return tidyBytes(trimFront(result, resultIdx+1))
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
func unsignedSub(b, s []byte) []byte {
	bDecimalPoint := findDotIndex(b)
	if bDecimalPoint == -1 {
		b = pushBack(b, '.')
		bDecimalPoint = len(b) - 1
	}

	sDecimalPoint := findDotIndex(s)
	if sDecimalPoint == -1 {
		s = pushBack(s, '.')
		sDecimalPoint = len(s) - 1
	}

	maxLenAfterDecimalPoint := max(len(b)-bDecimalPoint-1, len(s)-sDecimalPoint-1)
	maxP := max(bDecimalPoint, sDecimalPoint)

	resultLen := maxP + maxLenAfterDecimalPoint + 1 // +1 for decimal point
	result := make([]byte, resultLen)

	p := maxP + maxLenAfterDecimalPoint
	bShifting := maxP - bDecimalPoint
	sShifting := maxP - sDecimalPoint

	// Quick check: if shifting difference indicates subtraction is larger
	if sShifting < bShifting {
		return pushFront(unsignedSub(s, b), '-')
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
				return pushFront(unsignedSub(s, b), '-')
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

	return tidyBytes(trimFront(result, resultIdx+1))
}

func isZero(buf []byte) bool {
	for _, c := range buf {
		if c != '0' && c != '.' {
			return false
		}
	}
	return true
}

func isNegative(buf []byte) bool {
	return buf[0] == '-'
}
