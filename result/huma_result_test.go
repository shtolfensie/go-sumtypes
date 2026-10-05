//go:build go_sumtypes_huma

package result

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

var _ huma.SchemaProvider = Result[int]{}

func TestHumaSchema(t *testing.T) {
	registry := huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer)
	schema := Ok(42).Schema(registry)

	if len(schema.OneOf) != 2 {
		t.Fatalf("len(Schema().OneOf) = %d, want 2", len(schema.OneOf))
	}

	valueSchema := schema.OneOf[0]
	if valueSchema.Type != huma.TypeObject {
		t.Fatalf("value schema type = %q, want %q", valueSchema.Type, huma.TypeObject)
	}
	if valueSchema.Properties["value"].Type != huma.TypeInteger {
		t.Fatalf("value property type = %q, want %q", valueSchema.Properties["value"].Type, huma.TypeInteger)
	}
	if len(valueSchema.Required) != 1 || valueSchema.Required[0] != "value" {
		t.Fatalf("value required properties = %v, want [value]", valueSchema.Required)
	}

	errorSchema := schema.OneOf[1]
	if errorSchema.Type != huma.TypeObject {
		t.Fatalf("error schema type = %q, want %q", errorSchema.Type, huma.TypeObject)
	}
	if errorSchema.Properties["error"].Type != huma.TypeString {
		t.Fatalf("error property type = %q, want %q", errorSchema.Properties["error"].Type, huma.TypeString)
	}
	if len(errorSchema.Required) != 1 || errorSchema.Required[0] != "error" {
		t.Fatalf("error required properties = %v, want [error]", errorSchema.Required)
	}
}
