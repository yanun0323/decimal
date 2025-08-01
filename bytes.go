package decimal

import (
	"errors"
	"fmt"
)

var (
	zeroBytes = []byte{'0'}
)

func newDecimalBytes(s string) ([]byte, error) {
	if len(s) == 0 {
		return zeroBytes, nil
	}

	switch s {
	case "", "0", "0.", "0.0", ".0":
		return zeroBytes, nil
	}

	buf := []byte(s)
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
				buf = buf[1:]
				continue
			default:
				return zeroBytes, fmt.Errorf("invalid symbol (%c) in %s", b, s)
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
				buf = append(buf[:i], buf[i+1:]...)
				continue
			default:
				return zeroBytes, fmt.Errorf("invalid symbol (%c) in %s", b, s)
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
	result := num[start:end]

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
