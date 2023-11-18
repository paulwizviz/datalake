package dbops

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/paulwizviz/datalake/internal/block"
	"github.com/paulwizviz/datalake/internal/dbmodel"
)

// Connection returns postgres server connection
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

func InsertIntoBlockHeaders(ctx context.Context, conn *pgx.Conn, blk *block.Block) error {
	blktm := time.Unix(int64(blk.Timestamp), 0)
	_, err := conn.Exec(ctx, dbmodel.InsertIntoBlockHeaderStmt, blk.Hash, blk.ParentHash, blk.Number, blktm)
	if err != nil {
		return err
	}
	return nil
}

func InsertIntoSyncEvent(ctx context.Context, conn *pgx.Conn, topic string, blockHash string) error {
	uid := uuid.New()
	_, err := conn.Exec(ctx, dbmodel.InsertIntoSyncEventStmt, uid, topic, blockHash, time.Now())
	if err != nil {
		return err
	}
	return nil
}
