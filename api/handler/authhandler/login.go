package authhandler

import (
	"net/http"

	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/service/auth"
	"github.com/gin-gonic/gin"
)

func LoginHandler(ctx *gin.Context) {
	var request auth.LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		handler.SendError(ctx, http.StatusBadRequest, "invalid requisition body, fields 'username' and 'password' is necessary")
		return
	}

	token, expiresAt, err := auth.LoginService(handler.GetHandlerDB(), request)
	if err != nil {
		handler.SendError(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "login successful",
		"token":      token,
		"expires_at": expiresAt,
	})
}
