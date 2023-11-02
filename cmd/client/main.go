package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/paulwizviz/datalake/internal/block"
	"github.com/paulwizviz/datalake/internal/dbmodel"
	"github.com/paulwizviz/datalake/internal/dbops"
	"google.golang.org/grpc"
)

func fetchByHash(ctx context.Context, conn *pgx.Conn, bcs block.BlockServiceClient, blkHash string) {
	req := &block.BlockHashRequest{
		BlockHash: blkHash,
	}

	resp, err := bcs.FetchBlockByHash(ctx, req)
	if err != nil {
		err := dbops.InsertIntoSyncEvent(ctx, conn, dbmodel.BlockNotFound, "")
		if err != nil {
			log.Println("Unable to insert sync event")
		}
		log.Fatal("Unable to fetch block", err)
	}

	err = dbops.InsertIntoBlockHeaders(ctx, conn, resp)
	if err != nil {
		err := dbops.InsertIntoSyncEvent(ctx, conn, dbmodel.BlockNotFound, "")
		if err != nil {
			log.Println("Unable to insert sync event")
		}
		log.Fatal("Unable to insert block headers", err)
	}

	err = dbops.InsertIntoSyncEvent(ctx, conn, dbmodel.BlockIndexVal, resp.Hash)
	if err != nil {
		log.Println("Unable to insert sync event")
	}
}

func fetchByNumber(ctx context.Context, conn *pgx.Conn, bcs block.BlockServiceClient, num string) {
	req := &block.BlockNumberRequest{
		BlockNumber: num,
	}

	resp, err := bcs.FetchBlockByNumber(ctx, req)
	if err != nil {
		err := dbops.InsertIntoSyncEvent(ctx, conn, dbmodel.BlockNotFound, "")
		if err != nil {
			log.Println("Unable to insert sync event")
		}
		log.Fatal("Unable to fetch block", err)
	}

	err = dbops.InsertIntoBlockHeaders(ctx, conn, resp)
	if err != nil {
		err := dbops.InsertIntoSyncEvent(ctx, conn, dbmodel.BlockNotFound, "")
		if err != nil {
			log.Println("Unable to insert sync event")
		}
		log.Fatal("Unable to insert block headers", err)
	}

	err = dbops.InsertIntoSyncEvent(ctx, conn, dbmodel.BlockIndexVal, resp.Hash)
	if err != nil {
		log.Println("Unable to insert sync event")
	}
}

func main() {

	url := "postgres://postgres:postgres@localhost:5432/postgres"
	dbconn, err := dbops.Connection(url)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close(context.Background())

	ctx := context.Background()
	err = dbops.CreateTable(ctx, dbconn)
	if err != nil {
		log.Println(err)
	}

	port := 9000
	conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to setup connection on Port: %v", port)
	}
	defer conn.Close()

	c := block.NewBlockServiceClient(conn)

	if len(os.Args) != 3 {
		fmt.Println(os.Args)
		log.Fatal("Insufficient Aregument")
	}
	if os.Args[1] == "hash" {
		fetchByHash(ctx, dbconn, c, os.Args[2])
	}
	if os.Args[1] == "number" {
		fetchByNumber(ctx, dbconn, c, os.Args[2])
	}

}
