package delivery

import (
	"auth-service/auth/common"
	"auth-service/auth/dto"
	"auth-service/auth/middleware"
	"auth-service/auth/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type authController struct {
	svcAuth service.Service
	logger  *zap.Logger
}

func NewAuthController(svc service.Service, logger *zap.Logger) *authController {
	return &authController{svc, logger}
}

// @Summary Register User
// @Description Register user with the accepted type email, wechat or phone
// @Description For type email: only send email and password
// @Description For type phone: only send phone_code and number
// @Tags v1
// @Accept json
// @Produce json
// @Param body body dto.APIRequest true "Request Body"
// @Success 200 {object} dto.APIResponse
// @Router /v1/register [post]
func (u *authController) Create(ctx *gin.Context) {
	req := &dto.APIRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		u.logger.Error(err.Error())
		ctx.SecureJSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := u.svcAuth.Register(req)
	if err != nil {
		u.logger.Error(res.Message + " || " + err.Error())
		ctx.SecureJSON(http.StatusBadRequest, res)
		return
	}

	if res.StatusCode != 2000 {
		u.logger.Error(res.Message)
		ctx.SecureJSON(http.StatusBadRequest, res)
		return
	}
	u.logger.Info(string(res.Message))
	ctx.SecureJSON(http.StatusCreated, res)
}

// @Summary Login User
// @Description Login user with the accepted type email, wechat or phone
// @Description For type email: only send email and password
// @Description For type phone: only send phone_code and number
// @Description For type wechat: only send code, encrypted_data and iv
// @Tags v1
// @Accept json
// @Produce json
// @Param body body dto.APIRequest true "Request Body"
// @Success 200 {object} dto.APIResponse
// @Router /v1/login [post]
func (u *authController) Login(ctx *gin.Context) {
	req := dto.APIRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.logger.Error(err.Error())
		ctx.SecureJSON(http.StatusBadRequest, common.INVALID_JSON)
		return
	}

	res, err := u.svcAuth.Login(&req)
	if err != nil {
		u.logger.Error(res.Message + " || " + err.Error())
		ctx.SecureJSON(http.StatusBadRequest, res)
		return
	}
	if res.StatusCode != 2000 {
		u.logger.Error(res.Message)
		ctx.SecureJSON(http.StatusBadRequest, res)
		return
	}
	u.logger.Info(string(res.Message))
	ctx.SecureJSON(http.StatusOK, res)
}

// @Summary Logout
// @Description Logout user with access token
// @Description Add Bearer prefix before Authorization value.
// @Tags v1
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer + Token"
// @Success 200 {object} dto.APIResponse
// @Router /v1/logout [post]
func (u *authController) Logout(ctx *gin.Context) {
	token, err := middleware.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		u.logger.Error(err.Error())
		ctx.SecureJSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	res, err := u.svcAuth.Logout(token)
	if err != nil {
		u.logger.Error(res.Message + " || " + err.Error())
		ctx.SecureJSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	if res.StatusCode != 2000 {
		u.logger.Error(res.Message)
		ctx.SecureJSON(http.StatusBadRequest, res)
		return
	}
	u.logger.Info(string(res.Message))
	ctx.SecureJSON(http.StatusOK, res)
}

// @Summary Update Password
// @Description Update Password
// @Tags v1
// @Accept json
// @Produce json
// @Param body body dto.UpdatePasswordRequest true "Request Body"
// @Success 200 {object} dto.APIResponse
// @Router /v1/password [put]
func (u *authController) UpdatePassword(ctx *gin.Context) {
	req := dto.UpdatePasswordRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.logger.Error(err.Error())
		ctx.SecureJSON(http.StatusBadRequest, common.INVALID_JSON)
		return
	}

	resp, err := u.svcAuth.UpdatePassword(&req)
	if err != nil {
		u.logger.Error(resp.Message + " | " + err.Error())
		ctx.SecureJSON(http.StatusBadRequest, resp)
		return
	}

	u.logger.Info(string(resp.Message))
	ctx.SecureJSON(http.StatusOK, resp)
}

// @Summary Request for password reset
// @Description Request for password reset
// @Tags v1
// @Accept json
// @Produce json
// @Param body body dto.RequestPasswordChange true "Request Body"
// @Success 200 {object} dto.APIResponse
// @Router /v1/password [post]
func (a *authController) RequestPasswordChange(ctx *gin.Context) {

	req := dto.RequestPasswordChange{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		a.logger.Error(err.Error())
		ctx.SecureJSON(http.StatusBadRequest, common.INVALID_JSON)
		return
	}

	resp, err := a.svcAuth.RequestPasswordChange(req.Email)
	if err != nil {
		a.logger.Error(resp.Message + " | " + err.Error())
		ctx.SecureJSON(http.StatusBadRequest, resp)
		return
	}

	a.logger.Info(string(resp.Message))
	ctx.SecureJSON(http.StatusOK, resp)
}

// @Summary Refresh tokens
// @Description Refresh user's access and refresh tokens
// @Tags v1
// @Accept json
// @Produce json
// @Param X-Request-Token header string true "refresh Token"
// @Param Authorization header string true "Bearer + Token"
// @Success 200 {object} dto.APIResponse
// @Router /v1/refresh [post]
func (u *authController) RefreshToken(ctx *gin.Context) {
	refreshToken := ctx.GetHeader("X-Refresh-Token")

	res, err := u.svcAuth.RefreshToken(refreshToken)
	if err != nil {
		u.logger.Error(res.Message + " || " + err.Error())
		ctx.SecureJSON(http.StatusBadRequest, res)
		return
	}
	u.logger.Info(res.Message)
	ctx.SecureJSON(http.StatusOK, res)
}

func (u *authController) VerifyToken(ctx *gin.Context) {
	token, err := middleware.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		u.logger.Error(err.Error())
		ctx.SecureJSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	verifyRes, err := u.svcAuth.VerifyToken(token.AccessUUID)
	if err != nil {
		u.logger.Error(verifyRes.Message + " || " + err.Error())
		ctx.SecureJSON(http.StatusUnauthorized, verifyRes)
		return
	}

	ctx.Header("X-User-ID", token.UserID)
	ctx.Header("X-User-Roles", token.Role)
	ctx.Header("X-Tenant-Ids", token.TenantsIds)
	ctx.Header("X-Profile-ID", token.ProfileID)

	res := dto.APIResponse{
		StatusCode: common.SUCCESS_CODE,
		Message:    token.Role,
	}

	ctx.SecureJSON(http.StatusOK, res)
}
