package dice

import (
	"errors"
	"fmt"

	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"gorm.io/gorm"
)

type RollResult struct {
	Rolls      []int `json:"rolls"`
	Bonuses    []int `json:"bonuses"`
	SumOfRolls int   `json:"sum_of_rolls"`
	SumOfBonus int   `json:"sum_of_bonus"`
	Total      int   `json:"total"`
}

func Roll(numDice, sides int, bonuses []int, tableID, userID uint, db *gorm.DB) (*RollResult, error) {

	var membership models.TableUser
	if err := db.Where("table_id = ? AND user_id = ?", tableID, userID).First(&membership).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tableUser is not a member of this table")
		}
		return nil, err
	}

	if numDice <= 0 {
		return nil, fmt.Errorf("numDice must be positive")
	}

	if sides <= 1 {
		return nil, fmt.Errorf("sides must be greater than 1")
	}

	rolls := make([]int, numDice)
	sumOfRolls := 0
	for i := 0; i < numDice; i++ {
		roll := utils.SeededRand.Intn(sides) + 1
		rolls[i] = roll
		sumOfRolls += roll
	}

	sumOfBonus := 0
	for _, bonus := range bonuses {
		sumOfBonus += bonus
	}
	result := &RollResult{
		Rolls:      rolls,
		Bonuses:    bonuses,
		SumOfRolls: sumOfRolls,
		SumOfBonus: sumOfBonus,
		Total:      sumOfRolls + sumOfBonus,
	}
	return result, nil
}
