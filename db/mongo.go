package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/sunshineOfficial/golib/goerr"

	opt "go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

// NewMongo creates mongo.Client instance, specified ctx is used to preparatory checks: connect and ping
func NewMongo(ctx context.Context, url string, options ...Option) (*mongo.Client, error) {
	if url == "" {
		return nil, errors.New("[mongo] empty url")
	}

	clientOptions := opt.Client().
		ApplyURI(url)

	holder := applyOptions(options...)
	if holder.traces {
		clientOptions = clientOptions.SetMonitor(otelmongo.NewMonitor())
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("[mongo] can't create client: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("[mongo] can't ping: %v", err)
	}

	return client, nil
}

// WrapMongoError оборачивает ошибки от драйвера mongodb в ошибки goerr, если это применимо.
// nil -> nil
// mongo.ErrNoDocuments -> goerr.ErrNotFound
// остальные ошибки -> без изменений
func WrapMongoError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return goerr.ErrNotFound
	}

	return err
}
