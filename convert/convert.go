package convert

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/rakunlabs/types"
)

func NullToPtr[T any](v types.Null[T]) *T {
	if v.Valid {
		return &v.V
	}

	return nil
}

func SQLNullToPtr[T any](v sql.Null[T]) *T {
	if v.Valid {
		return &v.V
	}

	return nil
}

// ////////////////////////////////////////
// Time conversion functions

func TimeToStringPtr(v *time.Time, opts ...OptionTime) *string {
	if v == nil || v.IsZero() {
		return nil
	}

	str := v.Format(apply(opts).TimeFormat)

	return &str
}

func StringToTime(v string, opts ...OptionTime) (time.Time, error) {
	t, err := time.Parse(apply(opts).TimeFormat, v)

	return t, err
}

func StringPtrToTime(v *string, opts ...OptionTime) (time.Time, error) {
	if v == nil {
		return time.Time{}, nil
	}

	return StringToTime(*v, opts...)
}

func StringPtrToTimePtr(v *string, opts ...OptionTime) (*time.Time, error) {
	if v == nil {
		return nil, nil
	}

	t, err := StringToTime(*v, opts...)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// ////////////////////////////////////////

func PtrToZero[T any](v *T) T {
	if v == nil {
		return *new(T)
	}

	return *v
}

func Ptr[T any](v T) *T {
	return &v
}

func BytesToMap(v []byte) (types.Map[any], error) {
	if len(v) == 0 {
		return nil, nil
	}

	m := make(types.Map[any])

	if err := m.Scan(v); err != nil {
		return nil, err
	}

	return m, nil
}

// ////////////////////////////////////////

func IgnoreErr[T any](v T, _ error) T {
	return v
}

// ////////////////////////////////////////

// RawTo converts a byte slice to a type T using JSON decoding.
func RawTo[T any](v []byte) (*T, error) {
	t := new(T)

	if len(v) == 0 {
		return t, nil
	}

	// Parse the JSON data
	decoder := json.NewDecoder(bytes.NewReader(v))
	decoder.UseNumber()

	if err := decoder.Decode(t); err != nil {
		return nil, err
	}

	return t, nil
}

func RawToNull[T any](v []byte) (types.Null[T], error) {
	result := types.Null[T]{}
	if len(v) == 0 {
		return result, nil
	}

	decoder := json.NewDecoder(bytes.NewReader(v))
	decoder.UseNumber()

	if err := decoder.Decode(&result); err != nil {
		return types.Null[T]{}, err
	}

	return result, nil
}
