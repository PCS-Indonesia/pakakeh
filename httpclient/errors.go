package httpclient

import (
	"errors"
	"strings"
	"sync"
)

// Errors implements error interface. This instance of MultiError has zero or more errors.
type Errors struct {
	mutex sync.Mutex // i am using mutex to handle race-cond if anything goes wrong
	errs  []error    // list of errors
}

// Push adds an error to MultiError.
func (e *Errors) Push(errString string) {
	// prevent race condition (if any)
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.errs = append(e.errs, errors.New(errString))
}

// HasError checks if Errors struct has any error.
func (e *Errors) HasError() error {
	// prevent race condition (if any)
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if len(e.errs) == 0 {
		return nil
	}

	return e
}

// Error implements error interface.
func (e *Errors) Error() string {
	// prevent race condition (if any)
	e.mutex.Lock()
	defer e.mutex.Unlock()

	formattedError := make([]string, len(e.errs))
	for i, e := range e.errs {
		formattedError[i] = e.Error()
	}

	// join errors separated by space after comma
	return strings.Join(formattedError, ", ")
}
