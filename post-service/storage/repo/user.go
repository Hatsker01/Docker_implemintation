package repo

import (
	pb "github.com/Hatsker01/Docker_implemintation/post-service/genproto"
)

//PostStorageI ...
type PostStorageI interface {
	CreatePost(*pb.Post) (*pb.Post, error)
	GetPostById(id string) (*pb.Post, error)
	GetAllUserPosts(userID string) ([]*pb.Post, error)
	GetUserByPostId(postID string) (*pb.GetUserByPostIdResponse, error)
	DeletePost(userID string) (*pb.Emptya, error)
}
