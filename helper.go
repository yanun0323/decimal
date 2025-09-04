package decimal

// pushBackRepeat pushes values to the back of the slice
//
// NOTE: COPY WHEN CAPACITY NOT ENOUGH
func pushBackRepeat(slice []byte, value byte, count int) []byte {
	resultLen := len(slice) + count
	if cap(slice) >= resultLen {
		for ; count > 0; count-- {
			slice = append(slice, value)
		}

		return slice
	}

	buf := make([]byte, resultLen)
	copy(buf, slice)
	for i := len(slice); i < resultLen; i++ {
		buf[i] = value
	}

	return buf
}

// pushFront pushes values to the front of the slice
//
// NOTE: COPY WHEN CAPACITY NOT ENOUGH
func pushFrontRepeat(slice []byte, value byte, count int) []byte {
	resultLen := len(slice) + count
	if cap(slice) >= resultLen {
		for i := resultLen - len(slice); i > 0; i-- {
			slice = append(slice, value)
		}

		copy(slice[count:], slice)
		for i := 0; i < count; i++ {
			slice[i] = value
		}

		return slice
	}

	buf := make([]byte, resultLen)
	copy(buf[count:], slice)
	for i := 0; i < count; i++ {
		buf[i] = value
	}

	return buf
}

// pushFront pushes values to the front of the slice
//
// NOTE: COPY WHEN CAPACITY NOT ENOUGH
func pushFront(slice []byte, values ...byte) []byte {
	count := len(values)
	resultLen := len(slice) + count
	if cap(slice) >= resultLen {
		var zero byte
		for len(slice) < resultLen {
			slice = append(slice, zero)
		}

		copy(slice[count:], slice)
		copy(slice, values)

		return slice
	}

	buf := make([]byte, resultLen)
	copy(buf[count:], slice)
	copy(buf, values)

	return buf
}

// pushBack pushes values to the back of the slice
//
// NOTE: COPY WHEN CAPACITY NOT ENOUGH
func pushBack(slice []byte, values ...byte) []byte {
	return append(extend(slice, len(slice)+len(values)), values...)
}

// remove removes the element at the given index and returns the slice without the element.
//
// It keeps the capacity of the slice.
//
// NOTE: NO COPY
func remove(slice []byte, idx int) []byte {
	if idx < 0 || idx >= len(slice) {
		return slice
	}

	return append(slice[:idx], slice[idx+1:]...)
}

// extend extends the slice to the target capacity
//
// NOTE: COPY WHEN CAPACITY NOT ENOUGH
func extend(slice []byte, targetCap int) []byte {
	if cap(slice) >= targetCap {
		return slice
	}

	return append(make([]byte, 0, targetCap), slice...)
}

// insert inserts the value at the given index and returns the slice with the value.
//
// NOTE: COPY WHEN CAPACITY NOT ENOUGH
func insert(slice []byte, idx int, value byte) []byte {
	if idx < 0 || idx > len(slice) {
		return slice
	}

	if idx == len(slice) {
		return append(slice, value)
	}

	if len(slice) < cap(slice) { // slice is not filled
		slice = append(slice, value)
		copy(slice[idx+1:], slice[idx:])
		slice[idx] = value
		return slice
	}

	buf := make([]byte, len(slice)+1)
	copy(buf, slice[:idx])
	copy(buf[idx+1:], slice[idx:])
	buf[idx] = value

	return buf
}

// trimFront equals to slice[count:]
//
// It keeps the capacity of the slice.
//
// NOTE: NO COPY
func trimFront(slice []byte, start int) []byte {
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
//
// NOTE: NO COPY
func trimBack(slice []byte, count int) []byte {
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
//
// NOTE: NO COPY
func trim(slice []byte, start, end int) []byte {
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
