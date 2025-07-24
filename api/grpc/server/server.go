package server

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/routes"
	"github.com/GarotoCowboy/vttProject/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"net"
)

const (
	port = ":50051"
	host = "localhost"
)

// This function starts the grpc Server with tcp and the host = localhost:50051
func RunGRPCServer(db *gorm.DB, logger *config.Logger) {
	listen, err := net.Listen("tcp", host+port)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	//Put the grpcRoutes
	routes.Routes(grpcServer, db, logger)
	logger.InfoF("GRPC server listening on port 50051")
	if err := grpcServer.Serve(listen); err != nil {
		panic(err)
	}
}
