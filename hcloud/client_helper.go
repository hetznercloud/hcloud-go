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

// iterPages fetches each pages using the list function, and returns the result.
func iterPages[T any](listFn func(int) ([]*T, *Response, error)) ([]*T, error) {
	page := 1
	result := []*T{}

	for {
		pageResult, resp, err := listFn(page)
		if err != nil {
			return nil, err
		}

		result = append(result, pageResult...)

		if resp.Meta.Pagination == nil || resp.Meta.Pagination.NextPage == 0 {
			return result, nil
		}
		page = resp.Meta.Pagination.NextPage
	}
}

// firstBy fetches a list of items using the list function, and returns the first item
// of the list if present otherwise nil.
func firstBy[T any](listFn func() ([]*T, *Response, error)) (*T, *Response, error) {
	items, resp, err := listFn()
	if len(items) == 0 {
		return nil, resp, err
	}

	return items[0], resp, err
}

// firstByName is a wrapper around [firstBy], that checks if the provided name is not
// empty.
func firstByName[T any](name string, listFn func() ([]*T, *Response, error)) (*T, *Response, error) {
	if name == "" {
		return nil, nil, nil
	}

	return firstBy(listFn)
}
