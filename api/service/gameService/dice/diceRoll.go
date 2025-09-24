package dice

import (
	"fmt"

	"github.com/GarotoCowboy/vttProject/api/utils"
)

type RollResult struct {
	Rolls      []int `json:"rolls"`
	Bonuses    []int `json:"bonuses"`
	SumOfRolls int   `json:"sum_of_rolls"`
	SumOfBonus int   `json:"sum_of_bonus"`
	Total      int   `json:"total"`
}

func Roll(numDice, sides int, bonuses []int) (*RollResult, error) {
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
