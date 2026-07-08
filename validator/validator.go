package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/greatcloak/decimal"
)

var V *validator.Validate

func init() {
	V = New()
}

func New() *validator.Validate {
	v := validator.New()

	// Register decimal.Decimal
	v.RegisterCustomTypeFunc(validateDecimal, decimal.Decimal{})
	v.RegisterCustomTypeFunc(validateNullDecimal, decimal.NullDecimal{})

	return v
}

func validateDecimal(field reflect.Value) any {
	if dec, ok := field.Interface().(decimal.Decimal); ok {
		v, _ := dec.Float64()
		return v
	}

	return nil
}

func validateNullDecimal(field reflect.Value) any {
	if dec, ok := field.Interface().(decimal.NullDecimal); ok {
		if dec.Valid {
			v, _ := dec.Decimal.Float64()
			return v
		}
	}

	return nil
}

// Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func Struct(s any) error {
	return V.Struct(s)
}
