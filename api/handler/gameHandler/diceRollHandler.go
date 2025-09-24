package gameHandler

import (
	"net/http"
	"strconv"

	"github.com/GarotoCowboy/vttProject/api/dto/gameDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/service/gameService/dice"
	"github.com/GarotoCowboy/vttProject/api/service/tableUser"
	"github.com/gin-gonic/gin"
)

func RollDiceHandler(ctx *gin.Context) {

	tableUserIdStr := ctx.Param("tableUser")

	if tableUserIdStr == "" {
		handler.SendError(ctx, http.StatusBadRequest, gameDTO.ErrParamIsRequired("tableUserId", "uint").Error())
	}

	tableUserId, err := strconv.Atoi(tableUserIdStr)
	if err != nil || tableUserId <= 0 {
		handler.SendError(ctx, http.StatusBadRequest, "id must be a positive integer")
	}

	tableUserData, err := tableUser.GetTableUser(handler.GetHandlerDB(), uint(tableUserId))
	if err != nil {
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	request := gameDTO.RollResultRequest{}
	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("error rolling dice")
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	roll, err := dice.Roll(request.NumDices, request.Sides, request.Bonuses)
	if err != nil {
		return
	}

	resp := gameDTO.RollResultResponse{
		Bonuses:    roll.Bonuses,
		Rolls:      roll.Rolls,
		SumOfBonus: roll.SumOfBonus,
		SumOfRolls: roll.SumOfRolls,
		Total:      roll.Total,
		UserName:   tableUserData.User.Username,
	}

	handler.SendSucess(ctx, "roll dice", resp)

}
