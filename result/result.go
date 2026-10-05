package result

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type Result[T any] struct {
	value T
	err   error
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) Expect(msg string) T {
	if r.err != nil {
		panic(msg)
	}
	return r.value
}

func (r Result[T]) ExpectE(msg error) T {
	if r.err != nil {
		panic(msg)
	}
	return r.value
}

func (r Result[T]) UnwrapOr(val T) T {
	if r.err != nil {
		return val
	}
	return r.value
}

func (r Result[T]) Eq(val T) bool {
	if r.err != nil {
		return false
	}

	v := reflect.ValueOf(r.value)
	w := reflect.ValueOf(val)
	if v.Comparable() && w.Comparable() {
		return v.Equal(w)
	}
	return false
}

func Map[T any, U any](r Result[T], f func(val T) U) Result[U] {
	if r.IsOk() {
		return Result[U]{value: f(r.value)}
	}
	return Result[U]{err: r.err}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func Ok[T any](val T) Result[T] {
	return Result[T]{value: val}
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.IsOk() {
		return json.Marshal(struct {
			Value T `json:"value"`
		}{Value: r.value})
	}

	return json.Marshal(struct {
		Error string `json:"error"`
	}{Error: r.err.Error()})
}

func (r *Result[T]) UnmarshalJSON(data []byte) error {
	if r == nil {
		return fmt.Errorf("result: UnmarshalJSON on nil pointer")
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	valueJSON, hasValue := payload["value"]
	errorJSON, hasError := payload["error"]
	if hasValue == hasError {
		return fmt.Errorf("result: JSON object must contain exactly one of value or error")
	}

	if hasValue {
		var value T
		if err := json.Unmarshal(valueJSON, &value); err != nil {
			return err
		}
		*r = Ok(value)
		return nil
	}

	var message string
	if bytes.Equal(bytes.TrimSpace(errorJSON), []byte("null")) {
		return fmt.Errorf("result: error must be a JSON string")
	}
	if err := json.Unmarshal(errorJSON, &message); err != nil {
		return err
	}
	*r = Err[T](errors.New(message))
	return nil
}
