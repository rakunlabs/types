package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSON[T any] struct {
	V     T
	Valid bool
}

func NewJSON[T any](v T) JSON[T] {
	return JSON[T]{V: v, Valid: true}
}

func (s *JSON[T]) Scan(value any) error {
	if value == nil {
		s.V, s.Valid = *new(T), false

		return nil
	}
	s.Valid = true

	switch v := value.(type) {
	case []byte:
		// Parse the JSON data
		decoder := json.NewDecoder(bytes.NewReader(v))
		decoder.UseNumber()

		if err := decoder.Decode(&s.V); err != nil {
			return err
		}

		return nil
	case string:
		// Parse the JSON string
		decoder := json.NewDecoder(bytes.NewReader([]byte(v)))
		decoder.UseNumber()

		if err := decoder.Decode(&s.V); err != nil {
			return err
		}

		return nil
	default:
		return fmt.Errorf("%T, %w", value, ErrUnsupportedType)
	}
}

func (s JSON[T]) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}

	// Convert the JSON to JSON
	b, err := json.Marshal(s.V)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(b, []byte("null")) {
		return nil, nil
	}

	return b, nil
}

func (s JSON[T]) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(s.V)
}

func (s *JSON[T]) UnmarshalJSON(data []byte) error {
	if data == nil || bytes.Equal(data, []byte("null")) {
		s.V, s.Valid = *new(T), false

		return nil
	}

	// Parse the JSON data
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(&s.V); err != nil {
		return err
	}

	s.Valid = true

	return nil
}
