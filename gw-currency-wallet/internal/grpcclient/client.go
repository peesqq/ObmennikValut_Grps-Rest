package grpcclient

import (
	"log"

	"github.com/peesqq/proto-exchange/proto"

	"google.golang.org/grpc"
)

type GRPCClient struct {
	ExchangeService proto.ExchangeServiceClient
}

func NewGRPCClient(serverAddr string) *GRPCClient {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	return &GRPCClient{
		ExchangeService: proto.NewExchangeServiceClient(conn),
	}
}
