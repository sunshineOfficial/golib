package db

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NamedGet(tx *sqlx.Tx, dest any, query string, arg any) error {
	rows, err := tx.NamedQuery(query, arg)
	if err != nil {
		return fmt.Errorf("tx.NamedQuery: %w", err)
	}
	defer func() {
		err = errors.Join(err, fmt.Errorf("close rows: %w", rows.Close()))
	}()

	if !rows.Next() {
		err = errors.New("no rows returned")
		return err
	}

	if err = rows.StructScan(dest); err != nil {
		err = fmt.Errorf("rows.StructScan: %w", err)
		return err
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows.Err: %w", err)
		return err
	}

	return nil
}
