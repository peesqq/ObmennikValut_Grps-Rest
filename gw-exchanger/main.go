package main

import (
	"fmt"
	"log"
	"net"

	"github.com/peesqq/gw-exchanger/gw-exchanger/proto"
	"github.com/peesqq/gw-exchanger/internal/config"
	"github.com/peesqq/gw-exchanger/internal/db"
	"github.com/peesqq/gw-exchanger/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	conn, err := db.ConnectDB(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to database")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterExchangeServiceServer(grpcServer, &service.ExchangeService{})

	// Включаем gRPC Reflection API
	reflection.Register(grpcServer)

	fmt.Println("gRPC server running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
