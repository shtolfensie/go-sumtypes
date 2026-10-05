package result

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		in   Result[any]
		want string
	}{
		{name: "zero value", in: Ok[any](0), want: `{"value":0}`},
		{name: "slice", in: Ok[any]([]int{1, 2}), want: `{"value":[1,2]}`},
		{name: "nil value", in: Ok[any](nil), want: `{"value":null}`},
		{name: "error", in: Err[any](errors.New(`bad "value"`)), want: `{"error":"bad \"value\""}`},
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
	var ok Result[[]int]
	if err := json.Unmarshal([]byte(`{"value":[1,2]}`), &ok); err != nil {
		t.Fatalf("Unmarshal(value) error = %v", err)
	}
	if ok.IsErr() {
		t.Fatal("Unmarshal(value) produced Err, want Ok")
	}
	got := ok.Expect("checked above")
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("Unmarshal(value) = %v, want [1 2]", got)
	}

	var failure Result[int]
	if err := json.Unmarshal([]byte(`{"error":"failed"}`), &failure); err != nil {
		t.Fatalf("Unmarshal(error) error = %v", err)
	}
	if failure.IsOk() {
		t.Fatal("Unmarshal(error) produced Ok, want Err")
	}
	encoded, err := json.Marshal(failure)
	if err != nil {
		t.Fatalf("Marshal(decoded error) error = %v", err)
	}
	if string(encoded) != `{"error":"failed"}` {
		t.Fatalf("Marshal(decoded error) = %s, want %s", encoded, `{"error":"failed"}`)
	}
}

func TestUnmarshalJSONRejectsInvalidVariants(t *testing.T) {
	tests := []string{
		`{}`,
		`{"value":1,"error":"failed"}`,
		`{"error":null}`,
		`[]`,
	}

	for _, data := range tests {
		t.Run(data, func(t *testing.T) {
			r := Ok(7)
			if err := json.Unmarshal([]byte(data), &r); err == nil {
				t.Fatal("Unmarshal() error = nil, want an error")
			}
			if !r.Eq(7) {
				t.Fatalf("result changed after failed unmarshal: %+v", r)
			}
		})
	}
}

func TestResultJSONInStruct(t *testing.T) {
	type payload struct {
		Result Result[string] `json:"result"`
	}

	data, err := json.Marshal(payload{Result: Ok("value")})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if string(data) != `{"result":{"value":"value"}}` {
		t.Fatalf("Marshal() = %s, want %s", data, `{"result":{"value":"value"}}`)
	}
}

func TestUnmarshalResultJSONInStruct(t *testing.T) {
	type payload struct {
		Result Result[string] `json:"result"`
	}

	var got payload
	if err := json.Unmarshal([]byte(`{"result":{"value":"success"}}`), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.Result.IsErr() {
		t.Fatal("Unmarshal() produced Err, want Ok")
	}
	if value := got.Result.Expect("checked above"); value != "success" {
		t.Fatalf("Unmarshal() value = %q, want %q", value, "success")
	}
}

func TestMarshalJSONPropagatesValueError(t *testing.T) {
	if _, err := json.Marshal(Ok(make(chan int))); err == nil {
		t.Fatal("Marshal() error = nil, want an error")
	}
}
