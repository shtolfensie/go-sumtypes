//go:build go_sumtypes_huma

package option

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

// Schema describes Option as its nullable wrapped type.
func (o Option[T]) Schema(r huma.Registry) *huma.Schema {
	schema := huma.SchemaFromType(r, reflect.TypeFor[T]())
	if schema != nil && schema.Type != "" {
		schema.Nullable = true
	}
	return schema
}

// Receiver exposes the wrapped value to Huma's parameter parser.
func (o *Option[T]) Receiver() reflect.Value {
	return reflect.ValueOf(&o.value).Elem()
}

// OnParamSet records whether Huma received the parameter.
func (o *Option[T]) OnParamSet(isSet bool, _ any) {
	o.present = isSet
}
