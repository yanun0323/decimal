package decimal

// pushFront pushes values to the front of the slice
func pushBackRepeat[T any](slice []T, value T, count int) []T {
	resultLen := len(slice) + count
	if cap(slice) >= resultLen {
		for ; count > 0; count-- {
			slice = append(slice, value)
		}

		return slice
	}

	buf := make([]T, resultLen)
	copy(buf, slice)
	for i := len(slice); i < resultLen; i++ {
		buf[i] = value
	}

	return buf
}

// pushFront pushes values to the front of the slice
func pushFrontRepeat[T any](slice []T, value T, count int) []T {
	resultLen := len(slice) + count
	if cap(slice) >= resultLen {
		for len(slice) < resultLen {
			slice = append(slice, value)
		}

		copy(slice[count:], slice)
		for i := 0; i < count; i++ {
			slice[i] = value
		}

		return slice
	}

	buf := make([]T, resultLen)
	copy(buf[count:], slice)
	for i := 0; i < count; i++ {
		buf[i] = value
	}

	return buf
}

// pushFront pushes values to the front of the slice
func pushFront[T any](slice []T, values ...T) []T {
	count := len(values)
	resultLen := len(slice) + count
	if cap(slice) >= resultLen {
		var zero T
		for len(slice) < resultLen {
			slice = append(slice, zero)
		}

		copy(slice[count:], slice)
		copy(slice, values)

		return slice
	}

	buf := make([]T, resultLen)
	copy(buf[count:], slice)
	copy(buf, values)

	return buf
}

func pushBack[T any](slice []T, values ...T) []T {
	extend(slice, len(slice)+len(values))
	return append(slice, values...)
}

// remove removes the element at the given index and returns the slice without the element.
//
// It keeps the capacity of the slice.
func remove[T any](slice []T, idx int) []T {
	if idx < 0 || idx >= len(slice) {
		return slice
	}

	return append(slice[:idx], slice[idx+1:]...)
}

func extend[T any](slice []T, targetCap int) []T {
	if cap(slice) >= targetCap {
		return slice
	}

	return append(make([]T, 0, targetCap), slice...)
}

func insert[T any](slice []T, idx int, value T) []T {
	if idx < 0 || idx > len(slice) {
		return slice
	}

	if idx == len(slice) {
		return append(slice, value)
	}

	if len(slice) < cap(slice) { // slice is not filled
		slice = append(slice, value)
		copy(slice, slice[:idx])
		copy(slice[idx+1:], slice[idx:])
		slice[idx] = value
		return slice
	}

	buf := make([]T, len(slice)+1)
	copy(buf, slice[:idx])
	copy(buf[idx+1:], slice[idx:])
	buf[idx] = value

	return buf
}

// trimFront equals to slice[count:]
//
// It keeps the capacity of the slice.
func trimFront[T any](slice []T, start int) []T {
	if start <= 0 {
		return slice
	}

	if start >= len(slice) {
		return slice[:0]
	}

	return append(slice[:0], slice[start:]...)
}

// trimBack equals to slice[:len(slice)-end]
//
// It keeps the capacity of the slice.
func trimBack[T any](slice []T, count int) []T {
	if count <= 0 {
		return slice
	}

	if count >= len(slice) {
		return slice[:0]
	}

	return slice[:len(slice)-count]
}

// trim equals to slice[start:end]
//
// It keeps the capacity of the slice.
func trim[T any](slice []T, start, end int) []T {
	if start >= len(slice) || end <= 0 {
		return slice[:0]
	}

	if start < 0 {
		if end >= len(slice) {
			return slice
		}

		start = 0
	}

	if end > len(slice) {
		end = len(slice)
	}

	copy(slice, slice[start:end])

	return slice[:end-start]
}
