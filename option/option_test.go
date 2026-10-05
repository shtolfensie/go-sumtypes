package option

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		in   Option[any]
		want string
	}{
		{name: "none", in: None[any](), want: "null"},
		{name: "string", in: Some[any]("value"), want: `"value"`},
		{name: "zero", in: Some[any](0), want: "0"},
		{name: "slice", in: Some[any]([]int{1, 2}), want: "[1,2]"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := json.Marshal(test.in)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			if string(got) != test.want {
				t.Fatalf("Marshal() = %s, want %s", got, test.want)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var some Option[[]int]
	if err := json.Unmarshal([]byte(`[1,2]`), &some); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if some.IsNone() {
		t.Fatal("Unmarshal() produced None, want Some")
	}
	got := some.Expect("checked above")
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("Unmarshal() value = %v, want [1 2]", got)
	}

	if err := json.Unmarshal([]byte(" null "), &some); err != nil {
		t.Fatalf("Unmarshal(null) error = %v", err)
	}
	if some.IsSome() {
		t.Fatal("Unmarshal(null) produced Some, want None")
	}
}

func TestUnmarshalJSONDoesNotModifyOptionOnError(t *testing.T) {
	o := Some(7)
	if err := json.Unmarshal([]byte(`"invalid"`), &o); err == nil {
		t.Fatal("Unmarshal() error = nil, want an error")
	}
	if !o.Eq(7) {
		t.Fatalf("option changed after failed unmarshal: %+v", o)
	}
}

func TestOptionJSONInStruct(t *testing.T) {
	type payload struct {
		Value Option[string] `json:"value"`
	}

	data, err := json.Marshal(payload{Value: None[string]()})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if string(data) != `{"value":null}` {
		t.Fatalf("Marshal() = %s, want %s", data, `{"value":null}`)
	}
}

func TestUnmarshalOptionJSONInStruct(t *testing.T) {
	type payload struct {
		Value Option[string] `json:"value"`
	}

	var got payload
	if err := json.Unmarshal([]byte(`{"value":"present"}`), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.Value.IsNone() {
		t.Fatal("Unmarshal() produced None, want Some")
	}
	if value := got.Value.Expect("checked above"); value != "present" {
		t.Fatalf("Unmarshal() value = %q, want %q", value, "present")
	}
}

func TestMarshalJSONPropagatesValueError(t *testing.T) {
	if _, err := json.Marshal(Some(make(chan int))); err == nil {
		t.Fatal("Marshal() error = nil, want an error")
	}
}
