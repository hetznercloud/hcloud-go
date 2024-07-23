package sliceutil

// Batches splits a slice into multiple batches of a desired size.
func Batches[T any](all []T, size int) (batches [][]T) {
	batches = make([][]T, 0, (len(all)/size)+1)
	for size < len(all) {
		// Set the capacity of each chunk so that appending to a chunk does not
		// modify the original slice.
		all, batches = all[size:], append(batches, all[:size:size])
	}
	return append(batches, all)
}
