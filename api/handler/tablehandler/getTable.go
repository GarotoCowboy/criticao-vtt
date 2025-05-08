package tablehandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTableHandler(ctx *gin.Context) {
	id := ctx.Query("id")

	if id == "" {
		handler.SendError(ctx, http.StatusBadRequest, errParamIsRequired("id",
			"queryParameter").Error())
		return
	}

	table := models.Table{}

	if err := handler.GetHandlerDB().
		Preload("Owner").
		Preload("Members").
		Preload("Members.User").
		First(&table).Error; err != nil {
		handler.SendError(ctx, http.StatusNotFound, err.Error())
		return
	}
	handler.SendSucess(ctx, "get-table", table)
}
