package grpcClient

import (
	"fmt"

	"github.com/Hatsker01/Docker_implemintation/post-service/config"
	pb "github.com/Hatsker01/Docker_implemintation/post-service/genproto"
	"google.golang.org/grpc"
)

//GrpcClientI ...
type GrpcClientI interface {
	UserServise() pb.UserServiceClient
}

//GrpcClient ...
type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

//New ...
func New(cfg config.Config) (*GrpcClient, error) {
	connuser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithInsecure(),
	)

	if err != nil {
		return nil, fmt.Errorf("user service dial host: %s port: %d err: %s",
			cfg.UserServiceHost, cfg.UserServicePort, err.Error())
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"user_service": pb.NewUserServiceClient(connuser),
		},
	}, nil
}

func (g *GrpcClient) UserServise() pb.UserServiceClient {
	return g.connections["user_service"].(pb.UserServiceClient)
}
