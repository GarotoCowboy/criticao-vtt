package uploadHandler

import (
	"fmt"
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/service/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadFileMP3(ctx *gin.Context) {

	username := ctx.PostForm("username")
	messageUUID := ctx.PostForm("message_id")
	tableID := ctx.PostForm("table_id")

	file, err := ctx.FormFile("mp3")
	if err != nil {
		handler.GetHandlerLogger().ErrorF("failed to get file: %v", err)
		handler.SendError(ctx, http.StatusBadRequest, "error to upload MP3")
		return
	}

	filePath, fileName, err := upload.UploadMessageMP3File(file, tableID, username, messageUUID)
	if err != nil {
		handler.GetHandlerLogger().ErrorF("error validating upload: %v", err)
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
	}

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		handler.GetHandlerLogger().ErrorF("failed to save file: %v", err)
		handler.SendError(ctx, http.StatusInternalServerError, "error saving file")
		return
	}
	handler.SendSucess(ctx, "upload MP3 audio", gin.H{
		"fileName": fileName,
		"url":      fmt.Sprintf("/files/table_%s/chat/%s", tableID, fileName),
	})
}
