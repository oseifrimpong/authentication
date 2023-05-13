package service

import (
	"auth-service/auth/common"
	"auth-service/auth/dto"
	"auth-service/auth/repository"

	"go.uber.org/zap"
)

type Service interface {
	Register(req *dto.APIRequest) (res *dto.APIResponse, err error)
	Login(req *dto.APIRequest) (res *dto.APIResponse, err error)
	Logout(req *dto.AccessDetails) (res *dto.APIResponse, err error)
	UpdatePassword(data *dto.UpdatePasswordRequest) (res *dto.APIResponse, err error)
	RefreshToken(string) (res *dto.APIResponse, err error)
	VerifyToken(string) (res *dto.APIResponse, err error)
	RequestPasswordChange(emailStr string) (res *dto.APIResponse, err error)
}

type userService struct {
	repo   repository.Repository
	logger *zap.Logger
}

func NewUserService(repo repository.Repository, log *zap.Logger) Service {
	return &userService{
		repo:   repo,
		logger: log,
	}
}

func (svc userService) Register(req *dto.APIRequest) (res *dto.APIResponse, err error) {

	var response dto.APIResponse

	switch req.Type {
	case "email":
		reqData := &dto.AuthRequestData{
			Email:     req.Payload.Email,
			UserName:  req.Payload.UserName,
			FirstName: req.Payload.FirstName,
			LastName:  req.Payload.LastName,
			Password:  req.Payload.Password,
		}
		resp, err := svc.repo.Register(reqData, req.UserType)
		if err != nil {
			response.StatusCode = common.REGISTER_FAILED_CODE
			response.Message = err.Error()
			return &response, err
		}

		response.StatusCode = common.SUCCESS_CODE
		response.Message = common.REGISTER_SUCCESS
		response.Data = resp
		return &response, nil
	default:
		res := &dto.APIResponse{
			StatusCode: common.REGISTER_FAILED_CODE,
			Message:    common.WRONG_PAYLOAD_TYPE,
		}
		return res, nil
	}
}

func (svc *userService) Login(req *dto.APIRequest) (res *dto.APIResponse, err error) {
	var response dto.APIResponse

	switch req.Type {
	case "email":
		e := &dto.AuthRequestData{
			Email:    req.Payload.Email,
			Password: req.Payload.Password,
		}
		resp, err := svc.repo.Login(e)
		if err != nil {
			response.StatusCode = common.LOGIN_FAILED_CODE
			response.Message = err.Error()
			return &response, err
		}
		response.StatusCode = common.SUCCESS_CODE
		response.Message = common.LOGIN_SUCCESS
		response.Data = resp
		return &response, nil
	default:
		resp := &dto.APIResponse{
			StatusCode: common.WRONG_LOGIN_REQUEST_CODE,
			Message:    common.WRONG_PAYLOAD_TYPE,
		}
		return resp, nil
	}
}

func (svc *userService) Logout(t *dto.AccessDetails) (res *dto.APIResponse, err error) {
	var response *dto.APIResponse
	err = svc.repo.Logout(t)
	if err != nil {
		response.StatusCode = common.LOGOUT_FAILED_CODE
		response.Message = common.LOGOUT_FAILED + ":||" + err.Error()
	}

	response.StatusCode = common.SUCCESS_CODE
	response.Message = common.LOGOUT_SUCCESS

	return response, nil
}

func (svc *userService) RefreshToken(refreshToken string) (res *dto.APIResponse, err error) {
	var response dto.APIResponse
	resp, err := svc.repo.RefreshToken(refreshToken)
	if err != nil {
		response.StatusCode = common.REGISTER_FAILED_CODE
		response.Message = common.REFRESH_TOKEN_INCORRECT + ":||" + err.Error()
		return &response, err
	}

	response.StatusCode = common.SUCCESS_CODE
	response.Message = common.REFRESH_SUCCESS
	response.Data = resp

	return &response, nil
}

func (svc *userService) UpdatePassword(data *dto.UpdatePasswordRequest) (res *dto.APIResponse, err error) {
	var response dto.APIResponse
	err = svc.repo.UpdatePassword(data)
	if err != nil {
		response.StatusCode = common.UPDATE_PASSWORD_FAILED_CODE
		response.Message = err.Error()
		return &response, err
	}

	response.StatusCode = common.SUCCESS_CODE
	response.Message = "Password Updated Successfully"

	return &response, nil
}

func (svc *userService) RequestPasswordChange(emailStr string) (res *dto.APIResponse, err error) {
	var response dto.APIResponse
	err = svc.repo.RequestPasswordChange(emailStr)
	if err != nil {
		response.StatusCode = common.UPDATE_PASSWORD_FAILED_CODE
		response.Message = err.Error()
		return &response, err
	}

	response.StatusCode = common.SUCCESS_CODE
	response.Message = "Request for password reset is successful"

	return &response, nil
}

func (svc *userService) VerifyToken(key string) (res *dto.APIResponse, err error) {
	r := &dto.APIResponse{}

	verifyResult, err := svc.repo.FetchAuth(key)
	if err != nil {
		r.StatusCode = common.VERIFICATION_FAILED_CODE
		r.Message = common.VERIFICATION_FAILED
		return r, err
	}
	_ = verifyResult
	return r, nil
}
