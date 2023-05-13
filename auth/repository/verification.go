package repository

import (
	"auth-service/auth/model"
	"auth-service/auth/pkg"
	"auth-service/auth/service/email"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type VerificationRepo interface {
	VerifyEmail(contextStr string, email string) error
	ConfirmEmail(userID string) error
}

type verificationRepo struct {
	db     *gorm.DB
	rd     *redis.Client
	logger *zap.Logger
}

func NewVerificationRepo(db *gorm.DB, rd *redis.Client, l *zap.Logger) VerificationRepo {
	return &verificationRepo{db: db,
		rd:     rd,
		logger: l,
	}
}

func (v *verificationRepo) VerifyEmail(ctxStr string, email string) error {

	redisResults, err := v.rd.Get(ctxStr).Result()
	if err != nil {
		return errors.New("sorry, your confirmation link has expired")
	}

	if len(redisResults) == 0 {
		return errors.New("sorry, an error occurred while confirming your email")
	}

	if err := v.db.Model(&model.User{}).Where("email =?", email).Updates(map[string]interface{}{"is_verified": true}).Error; err != nil {
		return err
	}

	return nil
}

func (v *verificationRepo) ConfirmEmail(userID string) error {
	var user model.User
	id, _ := strconv.ParseInt(userID, 10, 36)
	if err := v.db.Model(&model.User{}).Where("id =?", id).Find(&user).Error; err != nil {
		return err
	}

	nonce, _ := pkg.GenerateNonce(v.rd, strconv.FormatInt(user.ID, 10))
	email.VerificationEmail(user.Email, user.Email, nonce, v.logger)
	return nil
}
