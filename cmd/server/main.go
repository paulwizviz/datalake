package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/paulwizviz/datalake/internal/block"
	"google.golang.org/grpc"
)

type Handler struct{}

func (h *Handler) FetchBlockByNumber(ctx context.Context, in *block.BlockNumberRequest) (*block.Block, error) {
	return &block.Block{
		Difficulty: "1",
	}, nil
}
func (h *Handler) FetchBlockByHash(ctx context.Context, in *block.BlockHashRequest) (*block.Block, error) {
	return &block.Block{
		Difficulty: "1",
	}, nil
}

func main() {
	port := 9000
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatalf("Fail to get listener on port: %v. Reason: %v", port, err)
	}
	log.Printf("Obtaining tcp listener on port: %v", port)

	h := Handler{}
	grpcServer := grpc.NewServer()
	block.RegisterBlockServiceServer(grpcServer, &h)

	log.Printf("Starting GRPC server on port: %v", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("GRPC server failed to start on port: %v. Error: %v", port, err)
	}

}
