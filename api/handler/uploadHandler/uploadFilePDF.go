package uploadHandler

import (
	"fmt"
	"net/http"

	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/service/upload"
	"github.com/gin-gonic/gin"
)

func UploadFilePDF(ctx *gin.Context) {
	username := ctx.PostForm("username")
	messageUUID := ctx.PostForm("message_id")
	tableID := ctx.PostForm("table_id")

	file, err := ctx.FormFile("pdf")
	if err != nil {
		handler.SendError(ctx, http.StatusBadRequest, "PDF file is required")
	}

	filepath, filename, err := upload.UploadMessagePDFFile(file, tableID, username, messageUUID)
	if err != nil {
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if err := ctx.SaveUploadedFile(file, filepath); err != nil {
		handler.SendError(ctx, http.StatusInternalServerError, "failed to save file")
	}
	handler.SendSucess(ctx, "upload PDF sucess", gin.H{
		"filename": filename,
		"url":      fmt.Sprintf("/files/table_%s/%s", tableID, filename),
	})
}
