package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendError(ctx *gin.Context, code int, msg string) {
	ctx.Header("Content-type", "application/json")
	ctx.JSON(code, gin.H{
		"message":   msg,
		"errorCode": code,
	})
}

func SendSucess(ctx *gin.Context, op string, data interface{}) {
	ctx.Header("Content-type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("operation from handler: %s sucessful", op),
		"data":    data,
	})
}

//func GetError(ctx *gin.Context,code int,msg string){
//	sendError(ctx,code,msg)
//}
