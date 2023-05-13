package delivery

import (
	"auth-service/auth/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type keysController struct {
	svc    service.KeyService
	logger *zap.Logger
}

func NewKeysController(svc service.KeyService, logger *zap.Logger) *keysController {
	return &keysController{svc, logger}
}

func (u *keysController) PublicKeys(ctx *gin.Context) {
	res, err := u.svc.PublicKey()
	if err != nil {
		u.logger.Error("Error fetching public keys." + " || " + err.Error())
		ctx.SecureJSON(http.StatusNoContent, "Error fetching public keys.")
		return
	}
	u.logger.Info("Fetched Public Keys.")
	ctx.SecureJSON(http.StatusOK, res)
}
