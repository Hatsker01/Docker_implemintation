package v1

import (
	"context"
	"encoding/json"

	//"encoding/json"
	"fmt"
	"net/http"
	"time"

	pb "github.com/Hatsker01/Docker_implemintation/api/genproto"
	l "github.com/Hatsker01/Docker_implemintation/api/pkg/logger"
	"github.com/Hatsker01/Docker_implemintation/api/pkg/utils"

	//"github.com/Hatsker01/Docker_implemintation/api/pkg/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	//"github.com/Hatsker01/Docker_implemintation/api/api/model"
)

//CreateUser creates user
//@Summary Create user summary
//Description This api for creating user
//@Tags user
//@Accept json
//@Produce json
//@Param user body User true "user body"
//@Success 200 {string} Success!
//@Router /v1/users [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        pb.User
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to blind json", l.Error(err))
		return
	}
	//fmt.Println(&body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().CreateUser(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}
	bodyByte, err := json.Marshal(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert json", l.Error(err))
		return

	}
	err = h.redisStorage.Set(body.FirstName, string(bodyByte))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert json2", l.Error(err))
		return

	}
	c.JSON(http.StatusCreated, response)
}

//GetUser by id
//@Summary Get user summary
//Description This api for getting user by id
//@Tags user
//@Accept json
//@Produce json
//@Param id path string true "User_id"
//@Success 200 {string} User!
//@Router /v1/users/{id} [get]
func (h *handlerV1) GetUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetUserById(
		ctx, &pb.GetUserByIdRequest{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return

	}
	c.JSON(http.StatusOK, response)

}

//Delete User by id
//@Summary Delete user summary
//Description This api for delete user by id
//@Tags user
//@Accept json
//@Produce json
//@Param id path string true "User_id"
//@Success 200 Succesfully deleted!
//@Router /v1/users/delete/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().DeleteUser(
		ctx, &pb.GetUserByIdRequest{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return

	}

	_, err = h.redisStorage.Delete(guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert json2", l.Error(err))
		return

	}

	c.JSON(http.StatusOK, response)

}

//CreateUser creates user
//@Summary Create user summary
//Description This api for creating user
//@Tags user
//@Accept json
//@Produce json
//@Param user body User true "user body"
//@Success 200 {string} User!
//@Router /v1/users/update/:id [put]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		body        pb.User
		jspbMarshal protojson.MarshalOptions
	)
	//var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to blind json", l.Error(err))
		return
	}
	fmt.Println(&body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().UpdateUser(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}
	c.JSON(http.StatusCreated, response)

}

//CreateUser creates user
//@Summary Create user summary
//Description This api for creating user
//@Tags user
//@Accept json
//@Produce json
//@Success 200 {string} []User!
//@Router /v1/users/alluser [get]
func (h *handlerV1) GetAllUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetAllUser(
		ctx, &pb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return

	}
	c.JSON(http.StatusOK, response)

}

// //Search from Users
// //@Summary Get user summary
// //@Description Search User
// //@Tags user
// //@Accept json
// //@Produce json
// //@Param first_name query string true "First_Name"
// //@Success 200 {string} User!
// //@Router /v1/users/lala/:first_name [get]
// func (h *handlerV1) SearchUser(c *gin.Context) {
// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	first_name := c.Param("first_name")
// 	fmt.Println(first_name)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.UserService().GetAllUser(
// 		ctx, &pb.Empty{})
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err,
// 		})
// 		h.log.Fatal("failed get all users searching", l.Error(err))
// 		return
// 	}

// 	for _, use := range response.Users {

// 		bodyByte, err := json.Marshal(response)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"error": err.Error(),
// 			})
// 			h.log.Error("failed to convert json", l.Error(err))
// 			return

// 		}
// 		err = h.redisStorage.Set(use.FirstName, string(bodyByte))
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"error": err.Error(),
// 			})
// 			h.log.Error("failed to convert json2", l.Error(err))
// 			return

// 		}

// 	}
// 	a, err := h.redisStorage.Search(first_name)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed while get first-name", l.Error(err))
// 		return
// 	}
// 	if a != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed while get a", l.Error(err))
// 		return
// 	}
// 	getuse,err:=h.redisStorage.Get(first_name)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed while get first-name", l.Error(err))
// 		return
// 	}
// 	c.JSON(http.StatusOK,getuse)

// }

//Get UserList  by limit and page
//@Summary Get userlist summary
//@Security APiKeyAuth
//@Description This api for delete user by id
//@Tags user
//@Accept json
//@Produce json
//@Param limit query int true "Limit"
//@Param page query int true "Page"
//@Success 200 {string} []User!
//@Router /v1/users/users [get]
func (h *handlerV1) GetListUsers(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	er := CheckClaims(h, c)
	if er == nil {
		fmt.Println(er)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while comparing token",
		})
		h.log.Fatal("failed in token check")
		return

	}

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Fatal("failed json params" + errStr[0])
		return
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetListUsers(
		ctx, &pb.GetUserRequest{
			Limit: params.Limit,
			Page:  params.Page,
		},
	)
	//fmt.Println(response)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		h.log.Fatal("failed json params", l.Error(err))
		return

	}
	// c.Header().Set("Content-Type","application/json")
	c.JSON(http.StatusOK, response)

}

// func (h *handlerV1) GetAllUser(c *gin.Context){
// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames=true

// 	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response,err:=h.serviceManager.UserService().GetAllUser(
// 		ctx,&pb.
// 		)

// }

//Get UserList  by limit and page
//@Summary SignIn summary
//@Security BearerAuth
//Description This api for delete user by id
//@Tags user
//@Accept json
//@Produce json
//@Success 200 {string} []User!
//@Router /v1/user/user/user [get]
func (h *handlerV1) Useruse(c *gin.Context) {
	claims := CheckClaims(h, c)
	userID := claims["sub"].(string)
	fmt.Println(userID)

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetUserById(ctx, &pb.GetUserByIdRequest{
		Id: userID,
	})
	//fmt.Println(response)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		h.log.Fatal("failed json params", l.Error(err))
		return

	}
	// c.Header().Set("Content-Type","application/json")
	c.JSON(http.StatusOK, response)

}

type User struct {
	Id           string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	FirstName    string     `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name"`
	LastName     string     `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name"`
	Username     string     `json:"username"`
	Email        string     `protobuf:"bytes,4,opt,name=email,proto3" json:"email"`
	Bio          string     `protobuf:"bytes,5,opt,name=bio,proto3" json:"bio"`
	PhoneNumbers []string   `protobuf:"bytes,6,rep,name=phone_numbers,json=phoneNumbers,proto3" json:"phone_numbers"`
	Address      []*Address `protobuf:"bytes,7,rep,name=address,proto3" json:"address"`
	Status       string     `protobuf:"bytes,8,opt,name=status,proto3" json:"status"`
	CreatedAt    string     `protobuf:"bytes,9,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	UpdatedAt    string     `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
	DeletedAt    string     `protobuf:"bytes,11,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at"`
	Posts        []*Post    `protobuf:"bytes,12,rep,name=posts,proto3" json:"posts"`
	Code         string     `json:"code"`
	Password     string     `json:"password"`
}

type Address struct {
	City       string `protobuf:"bytes,1,opt,name=city,proto3" json:"city"`
	Country    string `protobuf:"bytes,2,opt,name=country,proto3" json:"country"`
	District   string `protobuf:"bytes,3,opt,name=district,proto3" json:"district"`
	PostalCode int64  `protobuf:"varint,4,opt,name=postal_code,json=postalCode,proto3" json:"postal_code"`
}

type Post struct {
	Id          string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Name        string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	Description string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description"`
	UserId      string   `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id"`
	Medias      []*Media `protobuf:"bytes,5,rep,name=medias,proto3" json:"medias"`
}
type Media struct {
	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type"`
	Link string `protobuf:"bytes,3,opt,name=link,proto3" json:"link"`
}
type QueryParams struct {
	Page  int
	Limit int
}
