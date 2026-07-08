package types

import "github.com/greatcloak/decimal"

// Set default decimal.MarshalJSONWithoutQuotes to true.
//   - This will remove the quotes from the json output.
func init() {
	decimal.MarshalJSONWithoutQuotes = true
}

type (
	Decimal     = decimal.Decimal
	NullDecimal = decimal.NullDecimal
)
