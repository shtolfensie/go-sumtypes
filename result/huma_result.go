//go:build go_sumtypes_huma

package result

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

// Schema describes Result as either a value or error object.
func (r Result[T]) Schema(registry huma.Registry) *huma.Schema {
	return &huma.Schema{
		OneOf: []*huma.Schema{
			{
				Type: huma.TypeObject,
				Properties: map[string]*huma.Schema{
					"value": huma.SchemaFromType(registry, reflect.TypeFor[T]()),
				},
				Required: []string{"value"},
			},
			{
				Type: huma.TypeObject,
				Properties: map[string]*huma.Schema{
					"error": {Type: huma.TypeString},
				},
				Required: []string{"error"},
			},
		},
	}
}
