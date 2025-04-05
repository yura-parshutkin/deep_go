package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

// we use this structure to implement Unwrap method, we create new queue without first error
// for each unwrap method call
type QueueError struct {
	Errors []error
}

func (e *QueueError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	return e.Errors[0].Error()
}

func (e *QueueError) Unwrap() error {
	if len(e.Errors) > 1 {
		return &QueueError{
			Errors: e.Errors[1:],
		}
	}
	return nil
}

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.Errors)))
	if len(e.Errors) == 0 {
		return b.String()
	}
	for _, err := range e.Errors {
		b.WriteString(fmt.Sprintf("\t* %s", err.Error()))
	}
	b.WriteByte('\n')
	return b.String()
}

func (e *MultiError) Unwrap() error {
	if len(e.Errors) == 0 {
		return nil
	}
	return &QueueError{Errors: e.Errors}
}

func (e *MultiError) As(target any) bool {
	for i := range e.Errors {
		if errors.As(e.Errors[i], target) {
			return true
		}
	}
	return false
}

func (e *MultiError) Is(target error) bool {
	for _, err := range e.Errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e *MultiError) add(err error) {
	if err == nil {
		return
	}
	var multiErr *MultiError
	if errors.As(err, &multiErr) {
		e.Errors = append(e.Errors, err.(*MultiError).Errors...)
	} else {
		e.Errors = append(e.Errors, err)
	}
}

func Append(err error, errs ...error) *MultiError {
	ae := &MultiError{}
	ae.add(err)
	for i := range errs {
		ae.add(errs[i])
	}
	return ae
}

// it some random errors we use it to check As method
type MyError1 struct {
	Message string
}

func (e *MyError1) Error() string {
	return fmt.Sprintf("it's my own error 1 for testing: %s", e.Message)
}

type MyError2 struct {
	Message string
}

func (e *MyError2) Error() string {
	return fmt.Sprintf("it's my own error 2 for testing: %s", e.Message)
}

type MyError3 struct {
	Message string
}

func (e *MyError3) Error() string {
	return fmt.Sprintf("it's my own error 3 for testing: %s", e.Message)
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestMultiErrorHappyPaths(t *testing.T) {
	t.Run("without arguments", func(t *testing.T) {
		var err error
		err = Append(err)
		expectedMessage := "0 errors occured:\n"
		assert.EqualError(t, err, expectedMessage)
	})

	t.Run("with only nil arguments", func(t *testing.T) {
		var err error
		err = Append(err, nil, nil, nil)
		expectedMessage := "0 errors occured:\n"
		assert.EqualError(t, err, expectedMessage)
	})

	t.Run("with valid or invalid arguments", func(t *testing.T) {
		var err error
		err = Append(err, nil, fmt.Errorf("real error"), nil)
		expectedMessage := "1 errors occured:\n\t* real error\n"
		assert.EqualError(t, err, expectedMessage)
	})

	t.Run("unwrap", func(t *testing.T) {
		err := Append(errors.New("error 1"), errors.New("error 2"), errors.New("error 3"))

		err1 := errors.Unwrap(err)
		fmt.Printf("%p\n", err1)
		assert.EqualError(t, err1, "error 1")

		err2 := errors.Unwrap(err1)
		fmt.Printf("%p\n", err2)
		assert.EqualError(t, err2, "error 2")

		err3 := errors.Unwrap(err2)
		fmt.Printf("%p\n", err3)
		assert.EqualError(t, err3, "error 3")

		err4 := errors.Unwrap(err3)
		fmt.Println(&err4)
		assert.Nil(t, err4)
	})

	t.Run("unwrap recursive", func(t *testing.T) {
		err := Append(
			errors.New("error 1"),
			Append(errors.New("error 2"), errors.New("error 3")),
			errors.New("error 4"),
		)

		err1 := errors.Unwrap(err)
		assert.EqualError(t, err1, "error 1")

		err2 := errors.Unwrap(err1)
		assert.EqualError(t, err2, "error 2")

		err3 := errors.Unwrap(err2)
		assert.EqualError(t, err3, "error 3")

		err4 := errors.Unwrap(err3)
		assert.EqualError(t, err4, "error 4")

		err5 := errors.Unwrap(err4)
		assert.Nil(t, err5)
	})

	t.Run("error AS", func(t *testing.T) {
		err := Append(&MyError1{}, &MyError2{})

		myError1 := &MyError1{}
		myError2 := &MyError2{}
		myError3 := &MyError3{}

		assert.True(t, errors.As(err, &myError1), "error should be of type MyError1")
		assert.True(t, errors.As(err, &myError2), "error should be of type MyError2")
		assert.False(t, errors.As(err, &myError3), "error should not be of type MyError3")
	})

	t.Run("error IS", func(t *testing.T) {
		var (
			err1 = errors.New("error 1")
			err2 = errors.New("error 2")
			err5 = errors.New("error 3")
			err4 = errors.New("error 4")
			err3 = &MyError1{Message: "error 5"}
		)

		err := Append(err1, err2, err3)

		assert.True(t, errors.Is(err, err1))
		assert.True(t, errors.Is(err, err2))
		assert.True(t, errors.Is(err, err3))

		assert.False(t, errors.Is(err, err5))
		assert.False(t, errors.Is(err, err4))
	})
}
