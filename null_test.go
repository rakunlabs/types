package types

import (
	"encoding/json"
	"testing"
)

func TestNull(t *testing.T) {
	type tt struct {
		Name  *string      `json:"name"`
		Value Null[string] `json:"value"`
	}

	type test struct {
		json  []byte
		check func(t *testing.T, n tt)
	}

	tests := []test{
		{
			json: []byte(`{"name": "example", "value": "test"}`),
			check: func(t *testing.T, n tt) {
				if n.Name == nil || *n.Name != "example" {
					t.Errorf("expected Name to be 'example', got %v", n.Name)
				}

				if !n.Value.Valid || n.Value.V != "test" {
					t.Errorf("expected Value to be valid with 'test', got %v", n.Value)
				}
			},
		},
		{
			json: []byte(`{"name": null, "value": null}`),
			check: func(t *testing.T, n tt) {
				if n.Name != nil {
					t.Errorf("expected Name to be nil, got %v", n.Name)
				}

				if !n.Value.ParsedNull {
					t.Errorf("expected Value.ParsedNull to be true, got false")
				}
			},
		},
		{
			json: []byte(`{"name": "example"}`),
			check: func(t *testing.T, n tt) {
				if n.Name == nil || *n.Name != "example" {
					t.Errorf("expected Name to be 'example', got %v", n.Name)
				}

				if n.Value.Valid {
					t.Errorf("expected Value.Valid to be false, got true")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(string(test.json), func(t *testing.T) {
			n := tt{}

			if err := json.Unmarshal(test.json, &n); err != nil {
				t.Fatal(err)
			}

			test.check(t, n)
		})
	}
}

func TestNullMerge(t *testing.T) {
	old := NewNull("stored")

	tests := []struct {
		name   string
		update Null[string]
		want   Null[string]
	}{
		{
			name:   "value keeps update",
			update: NewNull("updated"),
			want:   NewNull("updated"),
		},
		{
			name:   "explicit null keeps update",
			update: Null[string]{ParsedNull: true},
			want:   Null[string]{ParsedNull: true},
		},
		{
			name:   "omitted uses old value",
			update: Null[string]{},
			want:   old,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.update.Merge(old)

			if test.update.V != test.want.V || test.update.Valid != test.want.Valid || test.update.ParsedNull != test.want.ParsedNull {
				t.Errorf("expected %+v, got %+v", test.want, test.update)
			}
		})
	}
}

func TestNullPresent(t *testing.T) {
	tests := []struct {
		name string
		null Null[string]
		want bool
	}{
		{name: "value", null: NewNull("updated"), want: true},
		{name: "explicit null", null: Null[string]{ParsedNull: true}, want: true},
		{name: "omitted", null: Null[string]{}, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.null.Present(); got != test.want {
				t.Errorf("expected %t, got %t", test.want, got)
			}
		})
	}
}
