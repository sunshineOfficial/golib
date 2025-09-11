package db

import (
	"database/sql"
	"errors"

	"github.com/sunshineOfficial/golib/goerr"
)

// WrapSqlError оборачивает ошибки драйверов sql в ошибки goerr, если это применимо.
// nil -> nil
// sql.ErrNoRows -> goerr.ErrNotFound
// остальные ошибки -> без изменений
func WrapSqlError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return goerr.ErrNotFound
	}

	return err
}
