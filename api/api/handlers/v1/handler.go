package v1

import (
	"errors"
	"net/http"

	"github.com/Hatsker01/Docker_implemintation/api/api/auth"
	"github.com/Hatsker01/Docker_implemintation/api/api/model"
	"github.com/Hatsker01/Docker_implemintation/api/config"
	"github.com/Hatsker01/Docker_implemintation/api/pkg/logger"
	"github.com/Hatsker01/Docker_implemintation/api/services"
	"github.com/Hatsker01/Docker_implemintation/api/storage/repo"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redisStorage   repo.RepositoryStorage
	JwtHandler     auth.JwtHandler
}

type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.RepositoryStorage
	jwtHandler     auth.JwtHandler
}

func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redisStorage:   c.Redis,
		JwtHandler:     c.jwtHandler,
	}
}

func CheckClaims(h *handlerV1, c *gin.Context) jwt.MapClaims {
	var (
		ErrUnauthorized = errors.New("unauthorized")
		authorization   model.JwtRequestModel
		claims          jwt.MapClaims
		err             error
	)

	authorization.Token = c.GetHeader("Authorization")
	if c.Request.Header.Get("Authorization") == "" {
		c.JSON(http.StatusUnauthorized, ErrUnauthorized)
		h.log.Error("Unauthorized request:", logger.Error(err))

	}
	h.JwtHandler.Token = authorization.Token
	claims, err = h.JwtHandler.ExtractClaims()
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrUnauthorized)
		h.log.Error("token is invalid:", logger.Error(err))
		return nil
	}
	return claims
}

// func New1(c *HandlerV1Config)*handlerV1{
// 	return &handlerV1{
// 		log: c.Logger,
// 		serviceManager: c.ServiceManager,
// 		cfg: c.Cfg,
// 	}
// }
