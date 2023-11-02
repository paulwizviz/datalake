package main

import (
	"context"
	"fmt"
	"log"

	"github.com/paulwizviz/datalake/internal/block"
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
	req := block.BlockNumberRequest{
		BlockNumber: "123",
	}

	resp, err := c.FetchBlockByNumber(context.TODO(), &req)
	if err != nil {
		log.Fatalf("Unable to fetch block by number: %v", err)
	}
	fmt.Println(resp)
}
