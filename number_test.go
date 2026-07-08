package types

import (
	"encoding/json"
	"testing"

	"github.com/greatcloak/decimal"
)

func TestNumber(t *testing.T) {
	t.Run("decimal", func(t *testing.T) {
		decimal.MarshalJSONWithoutQuotes = true
		v := Map[any]{
			"amount": decimal.RequireFromString("123.65"),
		}

		data, err := json.Marshal(v)
		if err != nil {
			t.Error(err)
		}

		if string(data) != `{"amount":123.65}` {
			t.Error("invalid json")
		}
	})

	t.Run("to decimal", func(t *testing.T) {
		data := []byte(`{"amount":123.6547382472397438472347328947473824723984723}`)
		v := struct {
			Amount decimal.Decimal `json:"amount"`
		}{}

		err := json.Unmarshal(data, &v)
		if err != nil {
			t.Error(err)
		}

		if v.Amount.String() != "123.6547382472397438472347328947473824723984723" {
			t.Error("invalid decimal")
		}
	})

	t.Run("Scan json", func(t *testing.T) {
		v := Map[any]{
			"amount": json.Number("123.65"),
		}

		data, err := json.Marshal(v)
		if err != nil {
			t.Error(err)
		}

		if string(data) != `{"amount":123.65}` {
			t.Error("invalid json")
		}
	})

	t.Run("Value json", func(t *testing.T) {
		v := struct {
			Amount json.Number `json:"amount"`
		}{}

		data := []byte(`{"amount":"123.65"}`)
		err := json.Unmarshal(data, &v)
		if err != nil {
			t.Error(err)
		}

		if v.Amount != "123.65" {
			t.Error("invalid amount")
		}
	})
}
