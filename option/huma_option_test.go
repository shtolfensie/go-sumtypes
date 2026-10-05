//go:build go_sumtypes_huma

package option

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

var (
	_ huma.SchemaProvider = Option[int]{}
	_ huma.ParamWrapper   = (*Option[int])(nil)
	_ huma.ParamReactor   = (*Option[int])(nil)
)

func TestHumaSchema(t *testing.T) {
	registry := huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer)
	schema := None[int]().Schema(registry)

	if schema.Type != huma.TypeInteger {
		t.Fatalf("Schema().Type = %q, want %q", schema.Type, huma.TypeInteger)
	}
	if !schema.Nullable {
		t.Fatal("Schema().Nullable = false, want true")
	}
}

func TestHumaParameterReceiver(t *testing.T) {
	var o Option[int]
	receiver := o.Receiver()
	if !receiver.CanSet() {
		t.Fatal("Receiver().CanSet() = false, want true")
	}

	receiver.SetInt(42)
	o.OnParamSet(true, int64(42))
	if !o.Eq(42) {
		t.Fatalf("parsed option = %+v, want Some(42)", o)
	}
}
