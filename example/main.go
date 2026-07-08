package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/greatcloak/decimal"
	"github.com/rs/zerolog/log"
	"github.com/worldline-go/initializer"
	"github.com/worldline-go/types"
	"github.com/worldline-go/types/example/database"
	"github.com/worldline-go/types/example/handler"
)

func main() {
	initializer.Init(
		run,
		initializer.WithMsgf("types example"),
	)
}

func run(ctx context.Context) error {
	db, err := database.Connect(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := database.MigrateDB(ctx, db); err != nil {
		return err
	}

	dbHandler := handler.New(db)

	// ////////////////////////////////////////
	// Create a train
	price, err := decimal.NewFromString("1919.23")
	if err != nil {
		return err
	}

	id, err := dbHandler.CreateTrain(ctx, &handler.Train{
		Details: map[string]any{
			"from":  "Istanbul",
			"to":    "Amsterdam",
			"value": 123.65,
			"peron": 5,
		},
		Additionals:  types.RawJSON(`{"key": "value"}`),
		Price:        types.Null[json.Number]{Null: sql.Null[json.Number]{Valid: true, V: json.Number(price.String())}},
		LastPrice:    decimal.NullDecimal{Valid: true, Decimal: price},
		Rate:         types.Null[string]{Null: sql.Null[string]{Valid: true, V: "5.87"}},
		CustomNumber: types.Null[string]{Null: sql.Null[string]{Valid: true, V: "123456"}},
		Slice:        types.Slice[string]{"a", "b", "c"},
		Data: types.NewJSON(&handler.Data{
			X: 123,
		}),
		CreatedAt: types.NewNull(types.Time{Time: time.Now()}),
	})
	if err != nil {
		return err
	}

	log.Info().Int64("id", id).Msg("Train ID")

	// ////////////////////////////////////////
	// Get a train
	train, err := dbHandler.GetTrain(ctx, id)
	if err != nil {
		return err
	}

	log.Info().Interface("details", train.Details).Str("price", train.Price.V.String()).
		Stringer("value", train.Details["value"].(json.Number)).
		Str("peron", train.Details["peron"].(json.Number).String()).
		Str("rate", train.Rate.V).
		Str("custom_number", train.CustomNumber.V).
		Str("created_at", train.CreatedAt.V.String()).
		Msg("Train Get")

	jsonLog(train)

	additionals, err := train.Additionals.ToMap()
	if err != nil {
		return err
	}

	log.Info().Interface("additionals", additionals).Msg("Train Get Additional with ToMap")

	// ////////////////////////////////////////
	// Update train to set back as null in database
	train.Details = nil
	train.Additionals = nil
	// train.Data.V = nil
	train.Slice = nil
	train.CreatedAt.Valid = false

	trainRaw, err := json.Marshal(train)
	if err != nil {
		return err
	}

	log.Info().RawJSON("train", trainRaw).Msg("Train Data For Update")

	if err := dbHandler.UpdateTrain(ctx, id, train); err != nil {
		return err
	}

	// ////////////////////////////////////////
	// Get a train
	train, err = dbHandler.GetTrain(ctx, id)
	if err != nil {
		return err
	}

	// Details now is nil
	log.Info().Interface("details", train.Details).Interface("additionals", train.Additionals).
		Str("price", train.Price.V.String()).
		Str("last_price", train.LastPrice.Decimal.String()).
		Str("rate", train.Rate.V).
		Msg("Train Updated")

	jsonLog(train)

	return nil
}

func jsonLog(v any) error {
	trainJSON, err := json.Marshal(v)
	if err != nil {
		return err
	}

	log.Info().RawJSON("train", trainJSON).Msg("Train get with RawJSON")

	return nil
}
