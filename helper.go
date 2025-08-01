package decimal

func pushFront[T any](slice []T, values ...T) []T {
	if cap(slice) == len(slice)+len(values) {
		// Efficiently move existing slice elements backward by len(values) positions
		copy(slice[len(values):cap(slice)], slice)
		// Copy new values to the front
		copy(slice[:len(values)], values)
		// Return slice with updated length
		return slice[:len(slice)+len(values)]
	}

	buf := make([]T, len(slice)+len(values))
	copy(buf[len(values):], slice)
	copy(buf[:len(values)], values)
	return buf
}
