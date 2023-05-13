package repository

import (
	"auth-service/auth/common"
	"auth-service/auth/dto"
	"auth-service/auth/model"
	"auth-service/auth/pkg"
	"auth-service/auth/pkg/token_management"
	"auth-service/auth/service/email"
	"auth-service/auth/service/identity"
	"errors"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	Login(data *dto.AuthRequestData) (res map[string]interface{}, err error)
	Logout(data *dto.AccessDetails) error
	Register(data *dto.AuthRequestData, userType string) (res map[string]interface{}, err error)
	RefreshToken(data string) (res map[string]interface{}, err error)
	UpdatePassword(data *dto.UpdatePasswordRequest) error
	FetchAuth(key string) (result string, err error)
	RequestPasswordChange(email string) error
}

type authRepo struct {
	db     *gorm.DB
	rd     *redis.Client
	logger *zap.Logger
}

func NewUserRepo(db *gorm.DB, rd *redis.Client, l *zap.Logger) Repository {
	return &authRepo{
		db:     db,
		rd:     rd,
		logger: l,
	}
}

func (a *authRepo) UpdatePassword(data *dto.UpdatePasswordRequest) error {

	redisResults, err := a.rd.Get(data.CTX).Result()
	if err != nil {
		a.logger.Error(err.Error())
		return errors.New("sorry, password reset link has expired")
	}

	if len(redisResults) == 0 {
		a.logger.Info("sorry, password reset link has expired")
		return errors.New("sorry, password reset link has expired")
	}

	encryptedPassword := pkg.SecurePass(data.Password)

	var user model.User
	err = a.db.Model(&model.User{}).Where("email = ?", data.Email).First(&user).Error
	if err != nil {
		a.logger.Error(err.Error())
		return errors.New(common.NO_RECORD_FOUND)
	}

	pwdMatch, err := pkg.ComparePasswords(user.Password, data.Password)
	if err != nil {
		a.logger.Error(err.Error())
		return errors.New("sorry, old and new password as the same")
	}

	if pwdMatch {
		a.logger.Info("sorry, old and new password as the same")
		return errors.New("sorry, old and new password as the same")
	}

	err = a.db.Model(&model.User{}).Where("email = ?", data.Email).Updates(map[string]interface{}{"password": encryptedPassword}).Error
	if err != nil {
		a.logger.Error(err.Error())
		return errors.New(common.PASSWORD_RESET_FAILED)
	}

	return nil
}

func (a *authRepo) RequestPasswordChange(emailStr string) error {
	var user model.User
	err := a.db.Where("email = ?", emailStr).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		a.logger.Error(err.Error())
		return errors.New(common.NO_RECORD_FOUND)
	}

	nonce, _ := pkg.GenerateNonce(a.rd, strconv.FormatInt(user.ID, 10))
	email.PasswordResetEmail(user.Email, nonce, a.logger)

	return nil
}

func (a *authRepo) Register(request *dto.AuthRequestData, userType string) (res map[string]interface{}, err error) {

	userPayload := &dto.UserPayload{
		Email:     request.Email,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		UserType:  userType,
	}

	result, err := identity.Create(userPayload)
	if err != nil {
		a.logger.Error(err.Error())
		return nil, errors.New(common.PROFILE_CREATION_FAILED)
	}

	encryptedPassword := pkg.SecurePass(request.Password)

	user := &model.User{
		Email:    request.Email,
		Password: encryptedPassword,
	}

	err = a.db.Create(user).Error
	if err != nil {
		a.logger.Error(err.Error())
		return nil, errors.New(common.REGISTRATION_FAILED)
	}

	roles := strings.Join(result.Data.Groups, ", ")
	if len(roles) == 0 {
		roles = userType
	}

	tenantIds := strings.Join(result.Data.TenantIds, ", ")

	profileID := result.Data.ProfileID

	if len(result.Data.ID) != 0 {
		profileID = result.Data.ID
	}

	token, err := token_management.CreateToken(user.ID, roles, tenantIds, profileID, user.IsVerified)
	if err != nil {
		a.logger.Error(err.Error())
		return nil, errors.New(common.SESSION_CREATION_FAILED)
	}

	err = token_management.SaveTokenDetails(a.rd, user.ID, token)
	if err != nil {
		a.logger.Error(err.Error())
		return nil, errors.New(common.SESSION_CREATION_FAILED)
	}

	tokenDetails := make(map[string]interface{})

	tokenDetails["access_token"] = token.AccessToken
	tokenDetails["refresh_token"] = token.RefreshToken
	tokenDetails["expires_at"] = token.AtExpires

	response := make(map[string]interface{})
	response["tokens"] = tokenDetails

	// nonce, _ := pkg.GenerateNonce(a.rd, strconv.FormatInt(user.ID, 10))
	// email.VerificationEmail(user.Email, user.Email, nonce, a.logger)

	return response, nil
}

func (u *authRepo) Login(data *dto.AuthRequestData) (res map[string]interface{}, err error) {
	user := model.User{}

	err = u.db.Where("email = ?", data.Email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		u.logger.Error(err.Error())
		return nil, errors.New(common.NO_RECORD_FOUND)
	}

	pwdMatch, err := pkg.ComparePasswords(user.Password, data.Password)
	if !pwdMatch {
		u.logger.Error(err.Error())
		return nil, errors.New(common.WRONG_PASSWORD)
	}

	result, err := identity.Retrieve(user.Email)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, errors.New(common.PROFILE_QUERY_FAILED)
	}

	roles := strings.Join(result.Data.Groups, ", ")
	tenantIds := strings.Join(result.Data.TenantIds, ", ")

	profileID := result.Data.ProfileID

	if len(result.Data.ID) != 0 {
		profileID = result.Data.ID
	}

	token, err := token_management.CreateToken(user.ID, roles, tenantIds, profileID, user.IsVerified)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, errors.New(common.SESSION_CREATION_FAILED)
	}

	reErr := token_management.SaveTokenDetails(u.rd, user.ID, token)
	if reErr != nil {
		u.logger.Error(err.Error())
		return nil, errors.New(common.SESSION_CREATION_FAILED)
	}

	t := make(map[string]interface{})

	t["access_token"] = token.AccessToken
	t["refresh_token"] = token.RefreshToken
	t["expires_at"] = token.AtExpires

	response := make(map[string]interface{})
	response["tokens"] = t

	return response, nil
}

func (u *authRepo) Logout(token *dto.AccessDetails) error {
	check, err := u.FetchAuth(token.AccessUUID)
	if err != nil {
		return err
	}

	_ = check

	result, err := token_management.DeleteKV(u.rd, token.AccessUUID)
	if err != nil || result == 0 {
		u.logger.Error(err.Error())
		return err
	}

	return nil
}

func (u *authRepo) RefreshToken(refreshToken string) (res map[string]interface{}, err error) {
	refreshedToken, err := token_management.RefreshToken(u.rd, refreshToken)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}
	return refreshedToken, err
}

func (u *authRepo) FetchAuth(authD string) (string, error) {
	userID, err := token_management.RetrieveKV(u.rd, authD)
	if err != nil {
		u.logger.Error(err.Error())
		return "", err
	}
	return userID, nil
}
