package dbops

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/paulwizviz/datalake/internal/dbmodel"
)

func Connection(url string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func CreateTable(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, dbmodel.CreateBlockHeadersStmt)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, dbmodel.CreateIndexBlockHeaderStmt)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, dbmodel.CreateSyncEventsStmt)
	if err != nil {
		return err
	}
	return nil
}
