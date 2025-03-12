package hcloud

// allFromSchemaFunc transform each item in the list using the FromSchema function, and
// returns the result.
func allFromSchemaFunc[T, V any](all []T, fn func(T) V) []V {
	result := make([]V, len(all))
	for i, t := range all {
		result[i] = fn(t)
	}

	return result
}
