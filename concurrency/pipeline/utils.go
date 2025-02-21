package pipeline

// Index retrieves a response of type T at the specified index from the ResponseImplementor.
// It takes a generic type parameter T, a ResponseImplementor interface, and an index i (since response is in slice).
// Returns the value at index i if it exists and can be type asserted to T, along with a boolean
// indicating if the operation was successful. Returns zero value and false if index is out of bounds
// or type assertion fails.
func Index[T any](ri ResponsesImplementor, i int) (t T, valid bool) {
	// Get all responses from the ResponseImplementor
	rg := ri.Get()

	// Check if index is out of bounds
	if i > len(rg)-1 {
		return t, false
	}

	// Assert the response at index i to type T
	if res, ok := rg[i].(T); ok {
		// If successful, return the value and true
		return res, true
	}

	// If type assertion fails, return zero value and false
	return t, false
}

// Find searches and returns the first response of type T from the ResponseImplementor.
// It takes a generic type parameter T and a ResponseImplementor interface.
// Returns the found value of type T and a boolean indicating if a value was found.
func Find[T any](ri ResponsesImplementor) (t T, found bool) {
	// Get all responses from the ResponseImplementor
	rr := ri.Get()

	// Iterate through each response
	for _, r := range rr {
		// Try to type assert the response to type T
		if _t, ok := r.(T); ok {
			// If successful, return the type and true
			return _t, true
		}
	}

	// If no value of type T is found, return zero value and false
	return t, false
}

// Get retrieves the last response from the ResponseImplementor and attempts to convert it to type T.
// It takes a generic type parameter T and a ResponseImplementor interface.
// Returns the last response if it can be type asserted to T, otherwise returns zero value of T.
// This is useful for getting the most recent response from a pipeline sequence.
func Get[T any](ri ResponsesImplementor) (t T) {
	// Get all responses from the ResponseImplementor
	rr := ri.Get()

	// Try to type assert the last response to type T
	res, ok := rr[len(rr)-1].(T)

	// If successful, return the value
	if ok {
		return res
	}

	// If type assertion fails, return zero value
	return t
}
