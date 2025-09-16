package option

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)


type Option[T any] struct {
	value T
	present   bool
}

func (r Option[T]) IsSome() bool {
	return r.present
}

func (r Option[T]) IsNone() bool {
	return !r.present
}

func (r Option[T]) Expect(msg string) T {
	if !r.present {
		panic(msg)
	}
	return r.value
}

func (r Option[T]) ExpectE(msg error) T {
	if !r.present{
		panic(msg)
	}
	return r.value
}

func (r Option[T]) UnwrapOr(val T) T {
	if !r.present {
		return val
	}
	return r.value
}

func (r Option[T]) Eq(val T) bool {
	if !r.present {
		return false
	}
	v := reflect.ValueOf(r.value)
	w := reflect.ValueOf(val)
	if v.Comparable() && w.Comparable() {
		return v.Equal(w)
	}
	return false
}


func Map[T any, U any](r Option[T], f func(val T) U) Option[U] {
	if r.IsSome() {
		return Option[U]{value: f(r.value), present: true}
	}
	return Option[U]{}
}

func MapOr[T any, U any](r Option[T], fallback U, f func(val T) U) U {
	if r.IsSome() {
		return f(r.value)
	}
	return fallback
}

func None[T any]() Option[T] {
	return Option[T]{present: false}
}

func Some[T any](val T) Option[T] {
	return Option[T]{value: val, present: true}
}

// type jsonPayload[T any] struct {
// 	Value *T `json:"value,omitempty"`
// 	Err error `json:"error,omitempty"`
// }
//
// func (r Result[T]) MarshalJSON() ([]byte, error) {
// 	var p *jsonPayload[T]
// 	if r.IsOk() {
// 		p = &jsonPayload[T]{
// 			Value: &r.value,
// 		}
// 	} else {
// 		p = &jsonPayload[T]{
// 			Err: r.err,
// 		}
// 	}
//
// 	return json.Marshal(p)
// }


// Value implements driver.Valuer interface for database writes
func (o Option[T]) Value() (driver.Value, error) {
	if o.IsNone() {
		return nil, nil
	}
	return o.Expect("option was checked before"), nil
}

// Scan implements sql.Scanner interface for database reads
func (o *Option[T]) Scan(value interface{}) error {
	if value == nil {
		*o = None[T]()
		return nil
	}

	switch v := value.(type) {
	case T:
		*o = Some(v)
	default:
		var e T
		return fmt.Errorf("cannot scan %T into Option[%T]", value, e)
	}

	return nil
}
