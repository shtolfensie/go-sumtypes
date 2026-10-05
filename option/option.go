package option

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

type Option[T any] struct {
	value   T
	present bool
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
	if !r.present {
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

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if o == nil {
		return fmt.Errorf("option: UnmarshalJSON on nil pointer")
	}
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*o = None[T]()
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*o = Some(value)
	return nil
}

// Value implements driver.Valuer interface for database writes
func (o Option[T]) Value() (driver.Value, error) {
	if o.IsNone() {
		return nil, nil
	}
	v := o.Expect("option was checked before")
	var val any
	switch tv := any(v).(type) {
	case int:
		val = int64(tv)
	case int8:
		val = int64(tv)
	case int16:
		val = int64(tv)
	case int32:
		val = int64(tv)
	default:
		val = tv
	}
	return val, nil
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
	case int64:
		var tval T
		switch any(tval).(type) {
		case int:
			intVal := any(int(v)).(T)
			*o = Some(intVal)
		case int8:
			intVal := any(int8(v)).(T)
			*o = Some(intVal)
		case int16:
			intVal := any(int16(v)).(T)
			*o = Some(intVal)
		case int32:
			intVal := any(int32(v)).(T)
			*o = Some(intVal)
		default:
			var e T
			return fmt.Errorf("cannot scan %T into Option[%T]", value, e)
		}
	default:
		var e T
		return fmt.Errorf("cannot scan %T into Option[%T]", value, e)
	}

	return nil
}
