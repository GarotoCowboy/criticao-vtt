package tablehandler

import (
	"fmt"
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteTableHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		handler.SendError(ctx, http.StatusBadRequest, errParamIsRequired("id",
			"queryParameter").Error())
		return
	}
	table := models.Table{}

	if err := handler.GetHandlerDB().
		First(&table, id).Error; err != nil {
		handler.SendError(ctx, http.StatusNotFound, fmt.Sprintf("Table with id: %s not found", id))
		return
	}
	if err := handler.GetHandlerDB().
		Preload("Owner").
		Preload("Members").
		Preload("Members.Owner").
		Delete(&table).Error; err != nil {
		handler.SendError(ctx, http.StatusNotFound, fmt.Sprintf("error deleting table with id: %s", id))
		return
	}
	handler.SendSucess(ctx, "delete-user", table)

}
