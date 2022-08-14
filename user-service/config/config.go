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
    PostServiceHost   string
    PostServicePort   int
    KafkaHost   string
    KafkaPort   int
}

// Load loads environment vars and inflates Config
func Load() Config {
    c := Config{}

    c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

    c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "dbpost"))
    c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5433))
    c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "users"))
    c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "hatsker"))
    c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "1"))
    c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "post_service"))
    c.PostServicePort = cast.ToInt(getOrReturnDefault("POST_SERVICE_HOST", 7007))
    c.KafkaHost = cast.ToString(getOrReturnDefault("KAFKA_HOST", "kafka"))
    c.KafkaPort = cast.ToInt(getOrReturnDefault("KAFKA_PORT", 9092))

    c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))

    c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":9000"))
 
    return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
    _, exists := os.LookupEnv(key)
    if exists {
        return os.Getenv(key)
    }

    return defaultValue
}
