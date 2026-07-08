package validator

import (
	"testing"

	"github.com/greatcloak/decimal"
)

func TestValidator(t *testing.T) {
	t.Run("init", func(t *testing.T) {
		if V == nil {
			t.Error("validator is nil")
		}
	})

	t.Run("validate", func(t *testing.T) {
		type Price struct {
			Amount   decimal.Decimal     `validate:"required,gte=0,lte=1000000"`
			Discount decimal.NullDecimal `validate:"omitempty,gte=0,lt=100"`
		}

		price := Price{
			Amount:   decimal.NewFromFloat(500.00),
			Discount: decimal.NewNullDecimal(decimal.NewFromFloat(10.00)),
		}

		if err := V.Struct(price); err != nil {
			t.Error(err)
		}
	})

	t.Run("validate error", func(t *testing.T) {
		type Price struct {
			Amount   decimal.Decimal     `validate:"required,gte=0,lte=1000000"`
			Discount decimal.NullDecimal `validate:"omitempty,gte=100,lt=1000"`
		}

		price := Price{
			Amount:   decimal.NewFromFloat(10.00),
			Discount: decimal.NewNullDecimal(decimal.NewFromFloat(999.00)),
		}

		if err := V.Struct(price); err != nil {
			t.Error(err)
		}
	})
}
