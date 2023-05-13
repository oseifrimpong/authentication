package service

import (
	"auth-service/auth/common"
	"auth-service/auth/dto"
	"auth-service/auth/repository"

	"go.uber.org/zap"
)

type Verification interface {
	VerifyEmail(contextStr string, email string) (res *dto.APIResponse, err error)
	ConfirmEmail(userID string) (res *dto.APIResponse, err error)
}

type verificationService struct {
	repo   repository.VerificationRepo
	logger *zap.Logger
}

// NewUserService with user repo
func NewVerificationService(repo repository.VerificationRepo, log *zap.Logger) Verification {
	return &verificationService{
		repo:   repo,
		logger: log,
	}
}

func (svc *verificationService) VerifyEmail(contextStr string, email string) (res *dto.APIResponse, err error) {
	var response dto.APIResponse

	err = svc.repo.VerifyEmail(contextStr, email)
	if err != nil {
		response.StatusCode = common.VERIFICATION_FAILED_CODE
		response.Message = "Email Verification Failed: " + err.Error()
		return &response, err
	}

	response.StatusCode = common.SUCCESS_CODE
	response.Message = "Successful Verification"

	return &response, nil
}

func (svc *verificationService) ConfirmEmail(userID string) (res *dto.APIResponse, err error) {
	var response dto.APIResponse

	err = svc.repo.ConfirmEmail(userID)
	if err != nil {
		response.StatusCode = common.VERIFICATION_FAILED_CODE
		response.Message = "Email Verification Failed: " + err.Error()
		return &response, err
	}

	response.StatusCode = common.SUCCESS_CODE
	response.Message = "Successful Verification"

	return &response, nil
}
