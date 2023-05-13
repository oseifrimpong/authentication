package main

import (
	"auth-service/auth/config"
	"auth-service/auth/middleware"
	"auth-service/auth/server"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Auth-Service API
// @version 0.01
// @description Authentication and Authorization Service

// @host xx
// @BasePath /auth/api
// @schemes http https
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(-1)
	}
}

var logger *zap.Logger

// TODO: abstract sentry configuration
func run() error {

	switch gin.Mode() {
	case gin.ReleaseMode:
		logger = config.Logger()
	default:
		logger = config.Logger()

		err := godotenv.Load()
		if err != nil {
			fmt.Println("error loading .env file")
		}
	}

	// key := os.Getenv("SENTRY_KEY")
	// project := os.Getenv("SENTRY_PROJECT")

	// sentryDSN := fmt.Sprintf("https://%s@sentry.io/%s", key, project)

	err := sentry.Init(sentry.ClientOptions{
		Dsn:           "https://88706db0b9084deba496af9576e6ef53@o1115648.ingest.sentry.io/4504457984606208",
		EnableTracing: true,
		Debug:         true,
		// TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	g := gin.Default()
	g.Use(middleware.CORSMiddleware())
	// g.Use(sentrygin.New(sentrygin.Options{}))

	d := config.BuildProject()
	svr := server.NewServer(g, d)

	svr.MapRoutes()
	if err := svr.SetupDB(); err != nil {
		logger.Error("Databases failed to start" + err.Error())
		return err
	}
	defer sentry.Flush(2 * time.Second)

	defer logger.Sync()

	return svr.Start()
}
