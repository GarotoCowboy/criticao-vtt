package tablehandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListTablesHandler(ctx *gin.Context) {

	tables := []models.Table{}

	if err := handler.GetHandlerDB().
		Preload("Owner").
		Preload("Members").
		Preload("Members.User").
		Find(&tables).Error; err != nil {
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	handler.SendSucess(ctx, "list-tables", tables)
}
