package main

import (
	"net"

	"github.com/Hatsker01/Docker_implemintation/post-service/config"
	pb "github.com/Hatsker01/Docker_implemintation/post-service/genproto"
	"github.com/Hatsker01/Docker_implemintation/post-service/pkg/db"
	"github.com/Hatsker01/Docker_implemintation/post-service/pkg/logger"
	"github.com/Hatsker01/Docker_implemintation/post-service/service"
	grpcClient "github.com/Hatsker01/Docker_implemintation/post-service/service/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"fmt"
	"os"

	gtrace "github.com/moxiaomomo/grpc-jaeger"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "post-service")
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

	postService := service.NewPostService(connDB, log, grpcC)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, postService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	tracer, _, err := gtrace.NewJaegerTracer("testCli", "127.0.0.1:6831")
	if err != nil {
		fmt.Printf("new tracer err: %+v\n", err)
		os.Exit(-1)
	}

	if tracer != nil {
		dialOpts = append(dialOpts, gtrace.DialOption(tracer))
	}
	// do rpc-call with dialOpts
	rpcCli(dialOpts)
}

func rpcCli(dialOpts []grpc.DialOption) {
	conn, err := grpc.Dial("127.0.0.1:8001", dialOpts...)
	if err != nil {
		fmt.Printf("grpc connect failed, err:%+v\n", err)
		return
	}
	defer conn.Close()

	// TODO: do some rpc-call
	// ...
}
