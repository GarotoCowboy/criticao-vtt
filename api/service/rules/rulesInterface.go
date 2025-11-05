package rules

import "github.com/GarotoCowboy/vttProject/api/models"

type RulesServiceInterface interface {
	CalculateInitialAttributes()
	CalculateSkills()
	CharacterLevelUp()
	GenerateInitialSheetData() (models.T20Sheet, error)
}
