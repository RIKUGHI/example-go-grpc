package main

import (
	"go-grpc/cmd/config"
	"go-grpc/cmd/services"
	"go-grpc/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	netListen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listened %v", err.Error())
	}

	db := config.ConnnectDatabase()

	grpcServer := grpc.NewServer()
	productService := services.ProductService{DB: db}
	pb.RegisterProductServiceServer(grpcServer, &productService)

	log.Printf("Server started at %v", netListen.Addr())
	if err := grpcServer.Serve(netListen); err != nil {
		log.Fatal("Failed to serve %v", err.Error())
	}
}
