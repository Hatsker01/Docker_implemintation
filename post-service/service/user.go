package service

import (
	"context"
	"fmt"

	pb "github.com/Hatsker01/Docker_implemintation/post-service/genproto"
	l "github.com/Hatsker01/Docker_implemintation/post-service/pkg/logger"
	cl "github.com/Hatsker01/Docker_implemintation/post-service/service/grpc_client"
	"github.com/Hatsker01/Docker_implemintation/post-service/storage"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	kafka "github.com/segmentio/kafka-go"
)

//PostService ...
type PostService struct {
	storage storage.IStorage
	logger  l.Logger
	client  cl.GrpcClientI
}

//NewPostService ...
func NewPostService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}
func (s *PostService) Consume(ctx context.Context, a *pb.Emptya) (*pb.Emptya, error) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	const (
		topic          = "user.user"
		broker1Address = "localhost:9092"
	)
	t := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address},
		Topic:   topic,
	})

	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := t.ReadMessage(ctx)
		if err != nil {
			return nil, err
		}
		var user pb.User

		err = user.Unmarshal(msg.Value)
		if err != nil {
			return nil, err
		}

		for _, post := range user.Posts {
			id, err := uuid.NewV4()
			if err != nil {
				s.logger.Error("failed while generating uuid for new post", l.Error(err))
				return nil, status.Error(codes.Internal, "failed while generating uuid")
			}
			post.Id = id.String()
			_, err = s.storage.Post().CreatePost(post)
			if err != nil {
				return nil, err
			}
		}
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))

	}

	//Output:

}
func (s *PostService) CreatePost(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid for new post", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while generating uuid")
	}
	req.Id = id.String()
	user, err := s.storage.Post().CreatePost(req)
	if err != nil {
		s.logger.Error("failed while inserting post", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while inserting post")
	}

	return user, nil
}

func (s *PostService) GetPostById(ctx context.Context, req *pb.GetPostByIdRequest) (*pb.Post, error) {
	user, err := s.storage.Post().GetPostById(req.UserId)
	if err != nil {
		s.logger.Error("failed get post", l.Error(err))
		return nil, status.Error(codes.Internal, "failed get user")
	}

	return user, err
}

func (s *PostService) GetAllUserPosts(ctx context.Context, req *pb.GetUserPostsrequest) (*pb.GetUserPosts, error) {
	posts, err := s.storage.Post().GetAllUserPosts(req.UserId)
	if err != nil {
		s.logger.Error("failed get all user posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed get all user posts")
	}

	// user, err := s.client.UserServise().GetUserById(ctx, &pb.GetUserByIdRequest{
	// 	Id: req.UserId,
	// })

	// if err != nil {
	// 	s.logger.Error("failed get a user by user_id in posts", l.Error(err))
	// 	return nil, status.Error(codes.Internal, "failed get a user by user_id in posts")
	// }

	// user.Posts = posts

	return &pb.GetUserPosts{
		Posts: posts,
	}, err
}

func (s *PostService) GetUserByPostId(ctx context.Context, req *pb.GetUserByPostIdRequest) (*pb.GetUserByPostIdResponse, error) {
	post, err := s.storage.Post().GetUserByPostId(req.Post_Id)
	if err != nil {
		s.logger.Error("failed get a post", l.Error(err))
		return nil, status.Error(codes.Internal, "failed get a post")
	}

	user, err := s.client.UserServise().GetUserById(ctx, &pb.GetUserByIdRequest{
		Id: post.UserId,
	})

	if err != nil {
		s.logger.Error("failed get a user by user_id in posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed get a user by user_id in posts")
	}

	post.UserFirstname = user.FirstName
	post.UserLastname = user.LastName

	return post, err
}

func (s *PostService) DeletePost(ctx context.Context, req *pb.GetUserByPostIdRequest) (*pb.Emptya, error) {
	_, err := s.storage.Post().DeletePost(req.Post_Id)
	if err != nil {
		s.logger.Error("failed delete posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed delete posts")
	}
	return &pb.Emptya{}, nil
}
