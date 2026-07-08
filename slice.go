package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Slice[T any] []T

func (s *Slice[T]) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		// Parse the JSON data
		decoder := json.NewDecoder(bytes.NewReader(v))
		decoder.UseNumber()

		if err := decoder.Decode(s); err != nil {
			return err
		}

		return nil
	case string:
		// Parse the JSON string
		decoder := json.NewDecoder(bytes.NewReader([]byte(v)))
		decoder.UseNumber()

		if err := decoder.Decode(s); err != nil {
			return err
		}

		return nil
	case nil:
		*s = nil

		return nil
	default:
		return fmt.Errorf("%T, %w", value, ErrUnsupportedType)
	}
}

func (s Slice[T]) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}

	// Convert the slice to JSON
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}
