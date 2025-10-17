package userhandler

import (
	"net/http"

	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/gin-gonic/gin"
)

//TODO: FAZER VERIFICAÇÕES PARA GARANTIR A INTEGRIDADE DA ENTREGA DE IMAGEM!!!

//TODO: TENHO QUE REFATORAR TODO ESSE CODIGO PARA IMPLEMENTAR O JWT E SALVAR A IMAGEM USANDO O ID DO USUARIO QUE PEDIU.

func UploadUserImg(ctx *gin.Context) {

	//userIDValue, exists := ctx.Get("user_id")
	//if !exists {
	//	handler.SendError(ctx, http.StatusBadRequest, "user_id not found in context")
	//	return
	//}

	//userID,ok := userIDValue.(uint)
	//if !ok {
	//	handler.SendError(ctx, http.StatusBadRequest, "invalid user_id type in context")
	//	return
	//}

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
