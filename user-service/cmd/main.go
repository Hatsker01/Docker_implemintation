package main

import (
	"context"
	"net"

	"github.com/Hatsker01/Docker_implemintation/user-service/config"
	pb "github.com/Hatsker01/Docker_implemintation/user-service/genproto"
	"github.com/Hatsker01/Docker_implemintation/user-service/pkg/db"
	"github.com/Hatsker01/Docker_implemintation/user-service/pkg/logger"
	trace "github.com/Hatsker01/Docker_implemintation/user-service/pkg/trace"
	"github.com/Hatsker01/Docker_implemintation/user-service/service"
	grpcClient "github.com/Hatsker01/Docker_implemintation/user-service/service/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "template-service")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	grpcC, err := grpcClient.New(cfg)
	if err != nil {
		log.Error("error establishing grpc connection", logger.Error(err))
		return
	}

	userService := service.NewUserService(connDB, log, grpcC)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	ctx := context.Background()
	prv, err := trace.NewProvider(trace.ProviderConfig{
		JaegerEndpoint: "http://localhost:14268/api/traces",
		ServiceName:    "server",
		ServiceVersion: "2.0.0",
		Environment:    "dev",
		Disabled:       false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer prv.Close(ctx)

	// Bootstrap listener.
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	// Bootstrap gRPC server.
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(trace.NewGRPUnaryServerInterceptor()),
	)

	// Bootstrap bank gRPC service server and respond to requests.
	bankAccService := bank.AccountService{}
	bankpb.RegisterAccountServiceServer(grpcServer, bankAccService)

	log.Fatalln(grpcServer.Serve(listener))
}
