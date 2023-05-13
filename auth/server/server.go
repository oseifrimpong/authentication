package server

import (
	"auth-service/auth/config"
	"auth-service/auth/model"
	"os"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"go.uber.org/zap"

	cors "github.com/rs/cors/wrapper/gin"
)

type dserver struct {
	router *gin.Engine
	cont   *dig.Container
}

func NewServer(e *gin.Engine, c *dig.Container) *dserver {
	return &dserver{
		router: e,
		cont:   c,
	}
}

func (ds *dserver) SetupDB() error {

	db, redis, err := config.InitializeDB()
	if err != nil {
		return err
	}

	redis.Ping().Result()

	db.AutoMigrate(&model.User{})
	return nil
}

func (ds *dserver) Start() error {
	ds.router.SetTrustedProxies(nil)
	ds.router.Use(cors.AllowAll())
	logger, _ := zap.NewProduction()

	ds.router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	ds.router.Use(ginzap.RecoveryWithZap(logger, true))
	return ds.router.Run(":" + os.Getenv("APP_PORT"))
}
