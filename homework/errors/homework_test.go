package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%v errors occured:\n", len(e.errs)))
	for _, er := range e.errs {
		buf.WriteString(fmt.Sprintf("\t* %v", er))
	}
	buf.WriteString("\n")
	return buf.String()
}

func (e *MultiError) Is(target error) bool {
	for _, er := range e.errs {
		if errors.Is(er, target) {
			return true
		}
	}
	return false
}

func (e *MultiError) Unwrap() error {
	if len(e.errs) == 1 {
		return e.errs[0]
	} else if len(e.errs) > 1 {
		multiErrors := make([]error, len(e.errs)-1)
		copy(multiErrors, e.errs)
		return &MultiError{
			errs: multiErrors,
		}
	}
	return nil
}

func (e *MultiError) As(target any) bool {
	for _, er := range e.errs {
		if errors.As(er, target) {
			return true
		}
	}
	return false
}

func Append(err error, errs ...error) *MultiError {
	if err == nil {
		return &MultiError{
			errs: errs,
		}
	}
	if er, ok := err.(*MultiError); ok {
		multiErrors := append(er.errs, errs...)
		return &MultiError{
			errs: multiErrors,
		}
	}
	multiErrs := make([]error, 0, len(errs)+1)
	multiErrs = append(multiErrs, errs...)
	multiErrs = append(multiErrs, err)
	return &MultiError{
		errs: multiErrs,
	}
}

type ErrNo3 struct{}

func (e ErrNo3) Error() string {
	return "error 3"
}

func TestMultiError(t *testing.T) {
	var err error
	err2 := errors.New("error 2")
	err = Append(err, errors.New("error 1"))
	err = Append(err, err2)
	err = Append(err, ErrNo3{})

	expectedMessage := "3 errors occured:\n\t* error 1\t* error 2\t* error 3\n"
	assert.EqualError(t, err, expectedMessage)

	assert.True(t, errors.Is(err, err2))
	assert.True(t, errors.As(err, &ErrNo3{}))

	unwrapped := errors.Unwrap(err)
	expectedUnwrappedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, unwrapped, expectedUnwrappedMessage)
}
