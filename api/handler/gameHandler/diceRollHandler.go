package gameHandler

import (
	"net/http"
	"strconv"

	"github.com/GarotoCowboy/vttProject/api/dto/gameDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/service/gameService/dice"
	"github.com/GarotoCowboy/vttProject/api/service/user"
	"github.com/gin-gonic/gin"
)

func RollDiceHandler(ctx *gin.Context) {

	userIDValue, exists := ctx.Get("user_id")

	if !exists {
		handler.SendError(ctx, http.StatusBadRequest, "user_id not found in context")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		handler.SendError(ctx, http.StatusBadRequest, "invalid user_id type in context")
		return
	}

	tableIDStr := ctx.Param("tableID")

	if tableIDStr == "" {
		handler.SendError(ctx, http.StatusBadRequest, gameDTO.ErrParamIsRequired("tableID", "uint").Error())
		return
	}

	tableID, err := strconv.ParseUint(tableIDStr, 10, 64)
	if err != nil || tableID <= 0 {
		handler.SendError(ctx, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	request := gameDTO.RollResultRequest{}
	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("error rolling dice")
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	roll, err := dice.Roll(request.NumDices, request.Sides, request.Bonuses, uint(tableID), userID, handler.GetHandlerDB())
	if err != nil {
		handler.SendError(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	userData, _ := user.GetUser(handler.GetHandlerDB(), userID)

	resp := gameDTO.RollResultResponse{
		Bonuses:    roll.Bonuses,
		Rolls:      roll.Rolls,
		SumOfBonus: roll.SumOfBonus,
		SumOfRolls: roll.SumOfRolls,
		Total:      roll.Total,
		UserName:   userData.Username,
	}

	handler.SendSucess(ctx, "roll dice", resp)

}
