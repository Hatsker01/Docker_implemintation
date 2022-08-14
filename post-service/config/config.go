package config

import (
	"os"

    "github.com/spf13/cast"
)

// Config ...
type Config struct {
    Environment       string // develop, staging, production
    PostgresHost      string
    PostgresPort      int
    PostgresDatabase  string
    PostgresUser      string
    PostgresPassword  string
    LogLevel          string
    RPCPort           string
    UserServiceHost   string
    UserServicePort   int
    KafkaHost   string
    KafkaPort   int
}

// Load loads environment vars and inflates Config
func Load() Config {
    c := Config{}

    c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

    c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "dbpost"))
    c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5433))
    c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "postdb"))
    c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "hatsker"))
    c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "1"))
    c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVISE_HOST", "user_service"))
    c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 9000))
    c.KafkaHost=cast.ToString(getOrReturnDefault("KAFKA_HOST","kafka"))
    c.KafkaPort=cast.ToInt(getOrReturnDefault("KAFKA_PORT", 9092))

    c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))

    c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":7007"))

    
    return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
    _, exists := os.LookupEnv(key)
    if exists {
        return os.Getenv(key)
    }

    return defaultValue
}
