package repo

import (
	pb "github.com/Hatsker01/Docker_implemintation/user-service/genproto"
)

//UserStorageI ...
type UserStorageI interface {
	CreateUser(*pb.User) (*pb.User, error)
	UpdateUser(*pb.User) (string, error)
	GetUserById(id string) (*pb.User, error)
	GetAllUser() ([]*pb.User, error)
	GetUserFromPost(userID string) (*pb.GetUserFromPostResponse, error)
	DeleteUser(userId string) (pb.Empty, error)
	GetListUsers(limit, page int64) (*pb.GetUserResponse, error)
	CheckField(field, value string) (bool, error)
	Login(email string) (*pb.User, error)
}
type ErrInvalidField interface {
}
