package uploadHandler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadAttachmentsHandler(ctx *gin.Context) {
	file, err := ctx.FormFile("attachment")
	if err != nil {
		handler.SendError(ctx, http.StatusBadRequest, "no file received")
		return
	}

	tableID := ctx.Param("table_id")
	if tableID == "" {
		handler.SendError(ctx, http.StatusBadRequest, "no table id received")
		return
	}

	fileType := ctx.Param("file_type")

	contentType := file.Header.Get("Content-Type")

	switch fileType {
	case "audio":
		if contentType != "audio/mpeg" && contentType != "audio/mp3" && contentType != "audio/ogg" {
			msg := fmt.Sprintf("Content-Type must be audio/mpeg or audio/mp3 : %s", contentType)
			handler.SendError(ctx, http.StatusBadRequest, msg)
			return
		}
		if file.Size > 20<<20 {
			msg := fmt.Sprintf("file too large %d bytes", file.Size)
			handler.SendError(ctx, http.StatusBadRequest, msg)
			return
		}
	case "document":
		if contentType != "application/pdf" {
			msg := fmt.Sprintf("Content-Type must be application/pdf : %s", contentType)
			handler.SendError(ctx, http.StatusBadRequest, msg)
			return
		}
		if file.Size > 50<<20 {
			msg := fmt.Sprintf("file too large %d bytes", file.Size)
			handler.SendError(ctx, http.StatusBadRequest, msg)
			return
		}
	default:
		if file.Size > 50<<20 {
			msg := fmt.Sprintf("file too large %d bytes", file.Size)
			handler.SendError(ctx, http.StatusBadRequest, msg)
			return
		}
	}

	extension := filepath.Ext(file.Filename)
	uniqueFileName := uuid.New().String() + extension

	subFolder := "general"
	if fileType == "audio" {
		subFolder = "audio"
	} else if fileType == "document" {
		subFolder = "documents"
	}

	dirPath := filepath.Join("./vttData", tableID, subFolder)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		handler.GetHandlerLogger().ErrorF("Failed to create directory %s", dirPath)
		handler.SendError(ctx, http.StatusInternalServerError, "server error preparing upload location")
		return
	}

	destinationPath := filepath.Join(dirPath, uniqueFileName)
	if err := ctx.SaveUploadedFile(file, destinationPath); err != nil {

		return
	}

	host := "http://localhost:8080"
	publicURL := fmt.Sprintf("%s/files/table_%s/uploads/%s", host, tableID, uniqueFileName)

	handler.SendSucess(ctx, "file-upload-successful", gin.H{
		"fileURL": publicURL,
	})
}
