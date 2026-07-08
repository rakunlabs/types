package types

import (
	"encoding/json"
	"testing"
)

func TestJSON(t *testing.T) {
	type Value struct {
		Name string `json:"name"`
	}

	v := JSON[Value]{}

	t.Run("Scan", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			if err := v.Scan(nil); err != nil {
				t.Error(err)
			}
			if v.Valid {
				t.Error("expected false")
			}
		})
		t.Run("[]byte", func(t *testing.T) {
			if err := v.Scan([]byte(`{"name":"test"}`)); err != nil {
				t.Error(err)
			}
			if !v.Valid {
				t.Error("expected true")
			}

			if v.V.Name != "test" {
				t.Error("expected test")
			}
		})
		t.Run("string", func(t *testing.T) {
			if err := v.Scan(`{"name":"test"}`); err != nil {
				t.Error(err)
			}
			if !v.Valid {
				t.Error("expected true")
			}
		})
		t.Run("unsupported", func(t *testing.T) {
			if err := v.Scan(1); err == nil {
				t.Error("expected error")
			}
		})
		t.Run("json marshal", func(t *testing.T) {
			v := JSON[Value]{V: Value{Name: "test"}, Valid: true}

			vByte, err := json.Marshal(v)
			if err != nil {
				t.Error(err)
			}

			if string(vByte) != `{"name":"test"}` {
				t.Error("expected {\"name\":\"test\"}")
			}
		})
		t.Run("json unmarshal", func(t *testing.T) {
			v := JSON[Value]{}

			if err := json.Unmarshal([]byte(`{"name":"test"}`), &v); err != nil {
				t.Error(err)
			}

			if v.V.Name != "test" {
				t.Error("expected test")
			}

			if !v.Valid {
				t.Error("expected true")
			}
		})
		t.Run("json unmarshal null", func(t *testing.T) {
			v := JSON[Value]{}

			if err := json.Unmarshal([]byte("null"), &v); err != nil {
				t.Error(err)
			}

			if v.Valid {
				t.Error("expected false")
			}
		})

		t.Run("json marshal number", func(t *testing.T) {
			v := map[string]any{
				"name": json.Number("1"),
			}

			res, err := json.Marshal(v)
			if err != nil {
				t.Error(err)
			}

			if string(res) != `{"name":1}` {
				t.Error("expected false")
			}
		})
	})
}
