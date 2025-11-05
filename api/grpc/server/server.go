package server

import (
	"net"

	"github.com/GarotoCowboy/vttProject/api/grpc/routes"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync/broker"
	"github.com/GarotoCowboy/vttProject/api/middleware"
	"github.com/GarotoCowboy/vttProject/config"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

const (
	port = ":50051"
	host = "localhost"
)

// RunGRPCServer This function starts the grpc Server with tcp and the host = localhost:50051
func RunGRPCServer(db *gorm.DB, logger *config.Logger, broker *broker.Broker) {
	listen, err := net.Listen("tcp", host+port)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middleware.GrpcAuthInterceptor,
			middleware.GrpcTableMemberInterceptor(db),
		)),
	)
	reflection.Register(grpcServer)
	//Put the grpcRoutes
	routes.Routes(grpcServer, db, logger, broker)
	logger.InfoF("GRPC server listening on port 50051")
	if err := grpcServer.Serve(listen); err != nil {
		panic(err)
	}
}
