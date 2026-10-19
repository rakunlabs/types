package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/greatcloak/decimal"
	"github.com/rakunlabs/into"
	"github.com/rakunlabs/logi"
	"github.com/rakunlabs/types"
	"github.com/rakunlabs/types/example/database"
	"github.com/rakunlabs/types/example/handler"
)

func main() {
	into.Init(
		run,
		into.WithLogger(logi.InitializeLog()),
		into.WithMsgf("types example"),
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

	slog.Info("Train ID", "id", id)

	// ////////////////////////////////////////
	// Get a train
	train, err := dbHandler.GetTrain(ctx, id)
	if err != nil {
		return err
	}

	slog.Info("Train Get",
		"details", train.Details,
		"price", train.Price.V.String(),
		"value", train.Details["value"].(json.Number).String(),
		"peron", train.Details["peron"].(json.Number).String(),
		"rate", train.Rate.V,
		"custom_number", train.CustomNumber.V,
		"created_at", train.CreatedAt.V.String(),
	)

	jsonLog(train)

	additionals, err := train.Additionals.ToMap()
	if err != nil {
		return err
	}

	slog.Info("Train Get Additional with ToMap", "additionals", additionals)

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

	slog.Info("Train Data For Update", "train", string(trainRaw))

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
	slog.Info("Train Updated",
		"details", train.Details,
		"additionals", train.Additionals,
		"price", train.Price.V.String(),
		"last_price", train.LastPrice.Decimal.String(),
		"rate", train.Rate.V,
	)

	jsonLog(train)

	return nil
}

func jsonLog(v any) error {
	trainJSON, err := json.Marshal(v)
	if err != nil {
		return err
	}

	slog.Info("Train get with RawJSON", "train", string(trainJSON))

	return nil
}
