package tablehandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateTableHandler(ctx *gin.Context) {

	request := UpdateTableRequest{}

	ctx.BindJSON(&request)

	if err := request.Validate(); err != nil {
		handler.GetHandlerLogger().ErrorF("Validation error: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	id := ctx.Query("id")
	if id == "" {
		handler.SendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "string").Error())
		return
	}

	table := models.Table{}
	if err := handler.GetHandlerDB().First(&table, id).Error; err != nil {
		handler.SendError(ctx, http.StatusNotFound, "table not found")
		return
	}
	//Update table

	if request.Name != "" {
		table.Name = request.Name
	}
	if request.Password != "" {
		hashedPassword, err := utils.HashPassword(table.Password)
		if err != nil {
			handler.GetHandlerLogger().ErrorF("Error hashing password: %v", err.Error())
			handler.SendError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		table.Password = string(hashedPassword)
	}

	//if request.Members != nil {
	//	table.Members = request.Members
	//}

	if err := handler.GetHandlerDB().Save(&table).Error; err != nil {
		handler.GetHandlerLogger().ErrorF("error updating table: %v", err.Error())
		return
	}
	handler.SendSucess(ctx, "update-table", table)
}
