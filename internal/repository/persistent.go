package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Ping(DBConn string) error {
	conn, err := pgx.Connect(context.Background(), DBConn)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		return err
	}

	return nil
}
