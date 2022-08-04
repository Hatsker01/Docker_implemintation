package main

import (
	"fmt"

	"github.com/Hatsker01/Docker_implemintation/api/api"
	"github.com/Hatsker01/Docker_implemintation/api/config"
	"github.com/Hatsker01/Docker_implemintation/api/pkg/logger"
	"github.com/Hatsker01/Docker_implemintation/api/services"
	"github.com/gomodule/redigo/redis"

	rds "github.com/Hatsker01/Docker_implemintation/api/storage/redis"
	"github.com/casbin/casbin/util"
	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	gormadapter "github.com/casbin/gorm-adapter/v2"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api")

	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	_, err := gormadapter.NewAdapter("postgres", psqlString, true)
	if err != nil {
		log.Error("new adapter error", logger.Error(err))

	}

	Enforcer, err := casbin.NewEnforcer(cfg.CasbinConfigPath, "./config/policy_defenition.csv")

	if err != nil {
		log.Error("new enforcer error", logger.Error(err))
		return
	}
	err = Enforcer.LoadPolicy()
	if err != nil {
		log.Error("new load policy error", logger.Error(err))
		return
	}

	fmt.Println(Enforcer)

	pool := redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	redisRepo := rds.NewRedisRepo(&pool)

	Enforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("KeyMatch", util.KeyMatch)
	Enforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("KeyMatch3", util.KeyMatch3)
	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		Casbin:         Enforcer,
		ServiceManager: serviceManager,
		RedisRepo:      redisRepo,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}

}
