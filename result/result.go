package result

import (
	"encoding/json"
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

func Err[T comparable](err error) Result[T] {
	return Result[T]{err: err}
}

func Ok[T comparable](val T) Result[T] {
	return Result[T]{value: val}
}

type jsonPayload[T any] struct {
	Value *T `json:"value,omitempty"`
	Err error `json:"error,omitempty"`
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	var p *jsonPayload[T]
	if r.IsOk() {
		p = &jsonPayload[T]{
			Value: &r.value,
		}
	} else {
		p = &jsonPayload[T]{
			Err: r.err,
		}
	}

	return json.Marshal(p)
}
