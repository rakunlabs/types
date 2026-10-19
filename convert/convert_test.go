package convert

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/rakunlabs/types"
)

func TestTimeFormatPtr(t *testing.T) {
	timeNow := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	type args struct {
		v    *time.Time
		opts []OptionTime
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantNil bool
	}{
		{
			name: "nil time",
			args: args{
				v:    nil,
				opts: nil,
			},
			wantNil: true,
		},
		{
			name: "zero time",
			args: args{
				v:    &time.Time{},
				opts: nil,
			},
			wantNil: true,
		},
		{
			name: "dateonly format",
			args: args{
				v:    &timeNow,
				opts: []OptionTime{WithTimeFormat(time.DateOnly)},
			},
			want: "2024-01-01",
		},
		{
			name: "default format",
			args: args{
				v:    &timeNow,
				opts: nil,
			},
			want: "2024-01-01T00:00:00Z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TimeToStringPtr(tt.args.v, tt.args.opts...)

			if tt.wantNil {
				if got != nil {
					t.Errorf("TimeFormatPtr() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				if !tt.wantNil {
					t.Errorf("TimeFormatPtr() = nil, want %v", tt.want)
				}

				return
			}

			if *got != tt.want {
				t.Errorf("TimeFormatPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPtrToZero(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		var vPtr *string
		v := PtrToZero(vPtr)
		if v != "" {
			t.Error("expected empty string")
		}
	})

	t.Run("empty time", func(t *testing.T) {
		var vPtr *time.Time
		v := PtrToZero(vPtr)
		if !v.IsZero() {
			t.Error("expected zero time")
		}
	})

	t.Run("empty int", func(t *testing.T) {
		var vPtr *int
		v := PtrToZero(vPtr)
		if v != 0 {
			t.Error("expected zero int")
		}
	})

	t.Run("non-empty string", func(t *testing.T) {
		vPtr := Ptr("hello")
		v := PtrToZero(vPtr)
		if v != "hello" {
			t.Error("expected 'hello'")
		}
	})
}
func TestRawTo_CustomStruct(t *testing.T) {
	type Custom struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Valid bool   `json:"valid"`
	}

	jsonData := []byte(`{"id":42,"name":"test","valid":true}`)

	got, err := RawTo[Custom](jsonData)
	if err != nil {
		t.Fatalf("RawTo() error = %v", err)
	}

	want := &Custom{ID: 42, Name: "test", Valid: true}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RawTo() = %+v, want %+v", got, want)
	}
}

func TestRawTo_CustomStructNil(t *testing.T) {
	type Custom struct {
		ID int `json:"id"`
	}

	jsonData := []byte(``)

	got, err := RawTo[Custom](jsonData)
	if err != nil {
		t.Fatalf("RawTo() error = %v", err)
	}

	if got.ID != 0 {
		t.Errorf("RawTo() = %+v, want %+v", got, &Custom{ID: 0})
	}
}

func TestRawTo_InvalidJSON(t *testing.T) {
	type Custom struct {
		ID int `json:"id"`
	}
	invalidJSON := []byte(`{"id":`)

	_, err := RawTo[Custom](invalidJSON)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestPtr(t *testing.T) {
	val := 123
	ptr := Ptr(val)
	if ptr == nil || *ptr != val {
		t.Errorf("Ptr() = %v, want %v", ptr, val)
	}
}

func TestIgnoreErr(t *testing.T) {
	val := "hello"
	got := IgnoreErr(val, nil)
	if got != val {
		t.Errorf("IgnoreErr() = %v, want %v", got, val)
	}
}

func TestBytesToMap(t *testing.T) {
	var jsonData []byte
	got, err := BytesToMap(jsonData)
	if err != nil {
		t.Fatalf("BytesToMap() error = %v", err)
	}
	if got != nil {
		t.Errorf("BytesToMap() = %v, want nil", got)
	}
}

func TestRawToNull(t *testing.T) {
	type tt struct {
		Name  *string            `json:"name"`
		Value types.Null[string] `json:"value"`
	}

	t.Run("null test", func(t *testing.T) {
		jsonData := []byte(`{"name": "example", "value": null}`)

		got, err := RawToNull[tt](jsonData)
		if err != nil {
			t.Fatalf("RawToNull() error = %v", err)
		}

		want := tt{
			Name:  Ptr("example"),
			Value: types.Null[string]{ParsedNull: true, Null: sql.Null[string]{Valid: false}},
		}

		if *got.V.Name != *want.Name {
			t.Errorf("RawToNull() Name = %v, want %v", got.V.Name, *want.Name)
		}

		if !reflect.DeepEqual(got.V.Value, want.Value) {
			t.Errorf("RawToNull() Value = %v, want %v", got.V.Value, want.Value)
		}
	})

	t.Run("empty byte", func(t *testing.T) {
		var jsonData []byte

		got, err := RawToNull[tt](jsonData)
		if err != nil {
			t.Fatalf("RawToNull() error = %v", err)
		}

		want := tt{}

		if got.V.Name != nil {
			t.Errorf("RawToNull() Name = %v, want nil", got.V.Name)
		}

		if got.Valid {
			t.Errorf("RawToNull() Valid = %v, want false", got.Valid)
		}

		if !reflect.DeepEqual(got.V.Value, want.Value) {
			t.Errorf("RawToNull() Value = %v, want %v", got.V.Value, want.Value)
		}
	})

	t.Run("null byte", func(t *testing.T) {
		jsonData := []byte(`null`)

		got, err := RawToNull[tt](jsonData)
		if err != nil {
			t.Fatalf("RawToNull() error = %v", err)
		}

		want := tt{}

		if got.V.Name != nil {
			t.Errorf("RawToNull() Name = %v, want nil", got.V.Name)
		}

		if got.Valid {
			t.Errorf("RawToNull() Valid = %v, want false", got.Valid)
		}

		if !got.ParsedNull {
			t.Errorf("RawToNull() ParsedNull = %v, want true", got.ParsedNull)
		}

		if !reflect.DeepEqual(got.V.Value, want.Value) {
			t.Errorf("RawToNull() Value = %v, want %v", got.V.Value, want.Value)
		}
	})

	t.Run("dummy byte", func(t *testing.T) {
		jsonData := []byte(`{"test": 1234}`)

		got, err := RawToNull[tt](jsonData)
		if err != nil {
			t.Fatalf("RawToNull() error = %v", err)
		}

		want := tt{}

		if got.V.Name != nil {
			t.Errorf("RawToNull() Name = %v, want nil", got.V.Name)
		}

		if !got.Valid {
			t.Errorf("RawToNull() Valid = %v, want true", got.Valid)
		}

		if got.ParsedNull {
			t.Errorf("RawToNull() ParsedNull = %v, want false", got.ParsedNull)
		}

		if !reflect.DeepEqual(got.V.Value, want.Value) {
			t.Errorf("RawToNull() Value = %v, want %v", got.V.Value, want.Value)
		}
	})

	t.Run("fail byte", func(t *testing.T) {
		jsonData := []byte(`ddddd`)

		_, err := RawToNull[tt](jsonData)
		if err == nil {
			t.Fatalf("RawToNull() error = %v", err)
		}
	})
}
