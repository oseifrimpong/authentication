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

type verificationController struct {
	svcVerify service.Verification
	logger    *zap.Logger
}

func NewVerificationController(svc service.Verification, logger *zap.Logger) *verificationController {
	return &verificationController{svc, logger}
}

// Verify email
// @Summary Verify User's Email
// @Description Logout user with access token
// @Description Add Bearer prefix before Authorization value.
// @Tags v1
// @Accept json
// @Produce json
// @Param body body model.VerifyEmailRequest true "Request Body"
// @Success 200 {object} model.APIResponse
// @Router /v1/verify [post]
func (v *verificationController) VerifyEmail(ctx *gin.Context) {
	req := dto.VerifyEmailRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v.logger.Error(err.Error())
		ctx.SecureJSON(http.StatusBadRequest, common.INVALID_JSON)
		return
	}

	res, err := v.svcVerify.VerifyEmail(req.CTX, req.Email)
	if err != nil {
		v.logger.Error(res.Message + " || " + err.Error())
		ctx.SecureJSON(http.StatusBadRequest, res)
		return
	}
	if res.StatusCode != 2000 {
		v.logger.Error(res.Message)
		ctx.SecureJSON(http.StatusBadRequest, res)
		return
	}
	v.logger.Info(string(res.Message))
	ctx.SecureJSON(http.StatusOK, res)
}

func (v *verificationController) ConfirmEmail(ctx *gin.Context) {
	tokenDetails, err := middleware.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		v.logger.Error(err.Error())
		ctx.SecureJSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	res, err := v.svcVerify.ConfirmEmail(tokenDetails.UserID)
	if err != nil {
		v.logger.Error(res.Message + " || " + err.Error())
		ctx.SecureJSON(http.StatusBadRequest, res)
		return
	}
	if res.StatusCode != 2000 {
		v.logger.Error(res.Message)
		ctx.SecureJSON(http.StatusBadRequest, res)
		return
	}
	v.logger.Info(string(res.Message))
	ctx.SecureJSON(http.StatusOK, res)
}
