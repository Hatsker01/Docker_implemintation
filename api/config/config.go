package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment string

	UserServiceHost string
	UserServicePort int

	PostgresHost string
	PostgresPort int
	PostgresUser string
	PostgresPassword string
	PostgresDatabase string

	PostSericeHost  string
	PostSericePort  int
	RedisHost       string
	RedisPort       int
	CasbinConfigPath string

	CtxTimeout int

	LogLevel  string
	HTTPPort  string
	SigninKey string
}

func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8090"))
	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "user_service"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 9000))
	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "redisdb"))
	c.RedisPort = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6378))
	c.SigninKey = cast.ToString(getOrReturnDefault("SIGNIN_KEY", "sijxoxyffnfxemfhuoehmniihgs"))
	c.PostgresHost=cast.ToString(getOrReturnDefault("POSTGRES_HOST","localhost"))
	c.PostgresPort=cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase=cast.ToString(getOrReturnDefault("POSTGRES_DB","home"))
	c.PostgresUser=cast.ToString(getOrReturnDefault("POSTGRES_USER","postgres"))
	c.PostgresPassword=cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD","1"	))

	c.CasbinConfigPath= cast.ToString(getOrReturnDefault("CASBIN_CONFIG_PATH","./config/rbac_model.conf"))
	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))
	return c

}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
