package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func InitializeDB() (*gorm.DB, *redis.Client, error) {
	redis, err := newRedis()
	if err != nil {
		return nil, nil, err
	}

	postgresDB, err := newPostgres()
	if err != nil {
		return nil, nil, err
	}

	return postgresDB, redis, nil
}

func newPostgres() (*gorm.DB, error) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_DATABASE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
		PrepareStmt:                              true,
		Logger:                                   newLogger,
		FullSaveAssociations:                     true,
	})
	if err != nil {
		return nil, err
	}

	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTIONS"))

	db.Use(
		dbresolver.Register(dbresolver.Config{}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(maxIdleConns).
			SetMaxOpenConns(maxOpenConns),
	)
	return db, nil
}

var client *redis.Client

func newRedis() (*redis.Client, error) {
	dsn := os.Getenv("REDIS_SERVICE_URL")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client, err
}
