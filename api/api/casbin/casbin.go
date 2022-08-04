package casbin

import (
	"net/http"
	"strings"

	"github.com/Hatsker01/Docker_implemintation/api/api/auth"
	"github.com/Hatsker01/Docker_implemintation/api/api/model"
	"github.com/Hatsker01/Docker_implemintation/api/config"
	casbin "github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JwtRoleStruct struct {
	enforce    *casbin.Enforcer
	conf       config.Config
	jwtHandler auth.JwtHandler
}

func NewJwtRoleStruct(e *casbin.Enforcer, c config.Config, jwtHandler auth.JwtHandler) gin.HandlerFunc {
	conf := &JwtRoleStruct{
		enforce:    e,
		conf:       c,
		jwtHandler: jwtHandler,
	}
	return func(c *gin.Context) {
		allow, err := conf.CheckPermission(c.Request)
		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			if v.Errors == jwt.ValidationErrorExpired {
				conf.RequireRefresh(c)
			} else {
				conf.RequirePermission(c)
			}
		} else if !allow {
			conf.RequirePermission(c)
		}
	}
}
func (a *JwtRoleStruct) RequireRefresh(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, model.ResponseError{
		Error: model.ServerError{
			Status:  "UNATHORIZED",
			Message: "Token is expired",
		},
	})
	c.AbortWithStatus(401)
}
func (a *JwtRoleStruct) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(401)
}

func (a *JwtRoleStruct) CheckPermission(r *http.Request) (bool, error) {
	role, err := a.GetRole(r)
	if err != nil {
		return false, err
	}
	method := r.Method
	path := r.URL.Path

	allowed, err := a.enforce.Enforce(role, path, method)
	if err != nil {
		panic(err)
	}
	return allowed, nil
}
func (a *JwtRoleStruct) GetRole(r *http.Request) (string, error) {
	var (
		role   string
		claims jwt.MapClaims
		err    error
	)
	jwtToken := r.Header.Get("Authorization")
	if jwtToken == "" {
		return "unauthorized", nil
	} else if strings.Contains(jwtToken, "Basic") {
		return "unauthorized", nil
	}
	a.jwtHandler.Token = jwtToken
	claims, err = a.jwtHandler.ExtractClaims()
	if err != nil {
		return "", err
	}
	if claims["role"].(string) == "authorized" {
		role = "authorized"
	} else {
		role = "unknown"
	}
	return role, nil

}
