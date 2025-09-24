package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

//TODO: FAZER VERIFICAÇÕES PARA GARANTIR A INTEGRIDADE DA ENTREGA DE IMAGEM!!!

func UploadUserImg(ctx *gin.Context) {

	//Get file
	file, err := ctx.FormFile("image")
	if err != nil {
		handler.GetHandlerLogger().ErrorF("error to upload image: %v ", err)
		handler.SendError(ctx, http.StatusInternalServerError, "error to upload image")
		return
	}

	//Validate type of libraryImg to PNG or JPEG
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/png" && contentType != "image/jpeg" {
		handler.GetHandlerLogger().ErrorF("invalid image type: %v ", err)
		handler.SendError(ctx, http.StatusInternalServerError, "Only Jpeg or Png file are allowed")
		return
	}

	//Validate with libraryImg with max of 5mb
	if file.Size > 5<<20 {
		handler.GetHandlerLogger().ErrorF("Image too large, exceded limit: %d bytes", file.Size)
		handler.SendError(ctx, http.StatusBadRequest, "exceed limit")
		return
	}

	filePath := "./libraryImg/" + file.Filename

	//Upload the file to specific folder
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		handler.GetHandlerLogger().ErrorF("error to save image to pathfile: %v ", err)
		handler.SendError(ctx, http.StatusBadRequest, "error to save image to pathfile")
		return
	}

	handler.SendSucess(ctx, "upload file", file)
}
