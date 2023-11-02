package main

import (
	"context"
	"fmt"
	"log"

	"github.com/paulwizviz/datalake/internal/block"
	"github.com/paulwizviz/datalake/internal/dbmodel"
	"github.com/paulwizviz/datalake/internal/dbops"
	"google.golang.org/grpc"
)

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
	req := &block.BlockHashRequest{
		BlockHash: "testdata/0x0e07da487d1c634a6ad96f62cef0f9cf52fc6a3f9df4d90e4f9bf58f844dc25c",
	}

	resp, err := c.FetchBlockByHash(context.TODO(), req)
	if err != nil {
		err := dbops.InsertIntoSyncEvent(context.TODO(), dbconn, dbmodel.BlockNotFound, "")
		if err != nil {
			log.Println("Unable to insert sync event")
		}
		log.Fatal("Unable to fetch block", err)
	}

	err = dbops.InsertIntoBlockHeaders(context.Background(), dbconn, resp)
	if err != nil {
		err := dbops.InsertIntoSyncEvent(context.TODO(), dbconn, dbmodel.BlockNotFound, "")
		if err != nil {
			log.Println("Unable to insert sync event")
		}
		log.Fatal("Unable to insert block headers", err)
	}

	err = dbops.InsertIntoSyncEvent(context.TODO(), dbconn, dbmodel.BlockIndexVal, resp.Hash)
	if err != nil {
		log.Println("Unable to insert sync event")
	}
}
