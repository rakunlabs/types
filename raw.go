package types

import (
	"database/sql/driver"
	"fmt"
)

// RawJSON same functionality with json.RawMessage and scan and value methods.
type RawJSON []byte

func (r *RawJSON) Scan(value any) error {
	var source []byte
	switch v := value.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	case nil:
		*r = nil

		return nil
	default:
		return fmt.Errorf("%T, %w", value, ErrUnsupportedType)
	}

	*r = append((*r)[0:0], source...)

	return nil
}

func (r RawJSON) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}

	return []byte(r), nil
}

func (r RawJSON) ToMap() (Map[any], error) {
	if r == nil {
		return nil, nil
	}

	m := make(Map[any])

	if err := m.Scan([]byte(r)); err != nil {
		return nil, err
	}

	return m, nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (r RawJSON) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}

	return r, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (r *RawJSON) UnmarshalJSON(data []byte) error {
	if r == nil {
		return fmt.Errorf("nil pointer, %w", ErrUnsupportedType)
	}

	*r = append((*r)[0:0], data...)

	return nil
}
