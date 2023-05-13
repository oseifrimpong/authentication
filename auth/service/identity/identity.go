package identity

import (
	"auth-service/auth/dto"
	"errors"
	"os"

	"github.com/imroc/req"
)

func Create(user *dto.UserPayload) (res *dto.AuthResponse, err error) {
	user_service_url := os.Getenv("USER_SERVICE_URL")
	user_service_path := os.Getenv("USER_SERVICE_PATH")

	url := user_service_url + user_service_path

	userRequest := &dto.UserServiceAPIRequest{
		User: user,
	}

	r, err := req.Post(url, req.BodyJSON(&userRequest))
	if err != nil {
		return nil, err
	}

	result := &dto.AuthResponse{}
	r.ToJSON(&result)

	if !(r.Response().StatusCode == 200 || r.Response().StatusCode == 201) {
		return nil, errors.New(result.Message)
	}

	return result, nil
}

func Retrieve(email string) (res *dto.AuthResponse, err error) {
	user_service_url := os.Getenv("USER_SERVICE_URL")
	user_service_path := os.Getenv("USER_SERVICE_PATH")

	reqUrl := user_service_url + user_service_path

	param := req.Param{
		"email": email,
	}

	r, err := req.Get(reqUrl, param)
	if err != nil {
		return nil, err
	}

	result := dto.AuthResponse{}
	r.ToJSON(&result)

	if r.Response().StatusCode != 200 {
		return nil, errors.New(result.Message)
	}

	return &result, nil
}
