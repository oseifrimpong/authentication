package server

import (
	"auth-service/auth/config"
	"auth-service/auth/delivery"
	"auth-service/auth/middleware"
	"auth-service/auth/service"
	"os"

	"github.com/dimiro1/health"
	"github.com/dimiro1/health/db"
	"github.com/dimiro1/health/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "auth-service/docs"
)

func (ds *dserver) MapRoutes() {
	versionGroup := ds.router.Group("api/v1")
	unVer := ds.router.Group("api")

	ds.healthCheck(unVer)
	ds.authRoutes(versionGroup)
	ds.verificationRoutes(versionGroup)
}

func (ds *dserver) verificationRoutes(router *gin.RouterGroup) {

	route := router.Group("/")
	{
		var logger *zap.Logger
		ds.cont.Invoke(func(l *zap.Logger) {
			logger = l
		})
		var verificationSvc service.Verification
		ds.cont.Invoke(func(u service.Verification) {
			verificationSvc = u
		})

		controller := delivery.NewVerificationController(verificationSvc, logger)

		route.POST("/verify", controller.VerifyEmail)
		route.POST("/confirm", controller.ConfirmEmail)
	}
}

func (ds *dserver) authRoutes(api *gin.RouterGroup) {
	a := api.Group("/")
	{
		var logger *zap.Logger
		ds.cont.Invoke(func(l *zap.Logger) {
			logger = l
		})
		var userSvc service.Service
		var keySvc service.KeyService

		ds.cont.Invoke(func(u service.Service, k service.KeyService) {
			userSvc = u
			keySvc = k
		})
		auth := delivery.NewAuthController(userSvc, logger)
		keys := delivery.NewKeysController(keySvc, logger)

		a.POST("/register", auth.Create)
		a.POST("/login", auth.Login)
		a.PUT("/password", auth.UpdatePassword)
		a.POST("/password", auth.RequestPasswordChange)

		a.POST("/logout", middleware.TokenAuthMiddleware(), auth.Logout)
		a.POST("/refresh", auth.RefreshToken)
		a.GET("/verify", middleware.TokenAuthMiddleware(), auth.VerifyToken)

		a.GET("/.wellknown/jwk.json", keys.PublicKeys)
		a.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func (ds *dserver) healthCheck(api *gin.RouterGroup) {
	h := api.Group("/")
	{
		pgDb, _, _ := config.InitializeDB()

		psql, _ := pgDb.DB()
		postgres := db.NewPostgreSQLChecker(psql)

		handler := health.NewHandler()
		handler.AddChecker("database", postgres)
		handler.AddChecker("redis", redis.NewChecker("tcp", os.Getenv("REDIS_SERVICE_URL")))
		handler.AddInfo("api", "Service is alive")
		handler.AddInfo("version", os.Getenv("VERSION"))

		h.GET("healthcheck", gin.WrapH(handler))
	}
}
