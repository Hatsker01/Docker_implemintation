package service

import (
	"context"
	"fmt"

	//"reflect"

	//"strconv"

	// "os/user"

	pb "github.com/Hatsker01/Docker_implemintation/user-service/genproto"
	l "github.com/Hatsker01/Docker_implemintation/user-service/pkg/logger"
	cl "github.com/Hatsker01/Docker_implemintation/user-service/service/grpc_client"
	"github.com/Hatsker01/Docker_implemintation/user-service/storage"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	//"golang.org/x/tools/go/analysis/passes/nilfunc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//UserService ...
type UserService struct {
	storage storage.IStorage
	logger  l.Logger
	client  cl.GrpcClientI
}

//NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid for user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while generating uuid for user")
	}
	req.Id = id.String()
	fmt.Println(id)
	user, err := s.storage.User().CreateUser(req)
	if err != nil {
		s.logger.Error("failed while creating user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while creating user")
	}
	fmt.Println(req)
	if req.Posts != nil {
		for _, post := range req.Posts {
			post.UserId = req.Id
			createdPosts, err := s.client.PostService().CreatePost(context.Background(), post)
			if err != nil {
				s.logger.Error("failed while inserting user post", l.Error(err))
				return nil, status.Error(codes.Internal, "failed while inserting user post")
			}
			fmt.Println(createdPosts)
		}
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.User) (*pb.UpdateUserResponse, error) {
	id, err := s.storage.User().UpdateUser(req)
	if err != nil {
		s.logger.Error("failed while updating user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while updating user")
	}
	return &pb.UpdateUserResponse{
		Id: id,
	}, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.User, error) {
	user, err := s.storage.User().GetUserById(req.Id)
	if err != nil {
		s.logger.Error("failed while getting by Id user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting by Id user")
	}

	posts, err := s.client.PostService().GetAllUserPosts(ctx, &pb.GetUserPostsrequest{UserId: req.Id})

	if err != nil {
		s.logger.Error("failed while getting user posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting user posts")
	}

	user.Posts = posts.Posts
	return user, err
}

func (s *UserService) GetAllUser(ctx context.Context, req *pb.Empty) (*pb.GetAllResponse, error) {
	users, err := s.storage.User().GetAllUser()
	if err != nil {
		s.logger.Error("failed while getting All users", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting All users")
	}

	for _, user := range users {
		posts, err := s.client.PostService().GetAllUserPosts(
			ctx,
			&pb.GetUserPostsrequest{
				UserId: user.Id,
			},
		)
		if err != nil {
			s.logger.Error("failed while getting user posts", l.Error(err))
			return nil, status.Error(codes.Internal, "failed while getting user posts")
		}

		user.Posts = posts.Posts
	}

	return &pb.GetAllResponse{
		Users: users,
	}, err
}

func (s *UserService) GetUserFromPost(ctx context.Context, req *pb.GetUserFromPostRequest) (*pb.GetUserFromPostResponse, error) {
	user, err := s.storage.User().GetUserFromPost(req.UserId)
	if err != nil {
		s.logger.Error("failed while getting a user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting a user")
	}

	return user, nil
}

func (s *UserService) GetAllUserPosts(ctx context.Context, req *pb.GetUserPostsrequest) (*pb.GetUserPosts, error) {
	res, err := s.client.PostService().GetAllUserPosts(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, err
}
func (s *UserService) DeleteUser(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.Empty, error) {
	_, err := s.storage.User().DeleteUser(req.Id)

	if err != nil {
		s.logger.Error("failed while getting a user", l.Error(err))
		return &pb.Empty{}, status.Error(codes.Internal, "failed while getting a user")
	}
	//GetUserByPostIdRequest
	_, err = s.client.PostService().DeletePost(ctx, &pb.GetUserByPostIdRequest{
		Post_Id: req.Id,
	})
	if err != nil {
		s.logger.Error("failed while deletepost", l.Error(err))
		return &pb.Empty{}, status.Error(codes.Internal, "failed while post delete")
	}
	return &pb.Empty{}, nil

}
func (s *UserService) GetListUsers(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	users, err := s.storage.User().GetListUsers(req.Limit, req.Page)
	if err != nil {
		s.logger.Error("failed while getting user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting user")
	}
	// for _,user:=range users.Users{
	// 	post,err:=s.client.PostService().GetAllUserPosts(ctx,*pb.GetUserPostsrequest{
	// 		UserId: user.Id,
	// 	})
	// 	if err!=nil{

	// 	}
	// }
	return users, nil
}

func (s *UserService) CheckField(ctx context.Context, req *pb.CheckFieldRequest) (*pb.CheckFieldReponse, error) {
	check, err := s.storage.User().CheckField(req.Field, req.Value)
	if err != nil {
		s.logger.Error("failed while getting user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting user")

	}
	return &pb.CheckFieldReponse{
		Check: check,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.User, error) {
	user, err := s.storage.User().Login(req.Email)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err != nil {
		s.logger.Error("failed while login user", l.Error(err))
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.logger.Error("failed getting  user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting user")

	}
	return user, nil
}
