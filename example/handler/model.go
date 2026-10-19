package handler

import (
	"encoding/json"

	"github.com/greatcloak/decimal"
	"github.com/rakunlabs/types"
)

type Train struct {
	ID           int64                   `db:"id"            goqu:"skipinsert"`
	Details      types.Map[any]          `db:"details"       json:"details,omitempty"`
	Additionals  types.RawJSON           `db:"additionals"   json:"additionals,omitempty"`
	Price        types.Null[json.Number] `db:"price"`
	LastPrice    decimal.NullDecimal     `db:"last_price"` // sql.Null[decimal.Decimal] is also usable
	Rate         types.Null[string]      `db:"rate"`
	CustomNumber types.Null[string]      `db:"custom_number"`
	Data         types.JSON[*Data]       `db:"data"`
	Slice        types.Slice[string]     `db:"slice"`
	CreatedAt    types.Null[types.Time]  `db:"created_at"`
}

type Data struct {
	X int `json:"x"`
}
