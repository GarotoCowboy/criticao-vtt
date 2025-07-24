package tormenta20Rules

import "github.com/GarotoCowboy/vttProject/api/models"

// that function generate a initial tormenta 20 sheet
func (s *RulesService) GenerateInitialSheetData() (*models.T20Sheet, error) {

	sheet := &models.T20Sheet{

		//Initial Attributes
		Attributes: models.Attributes{
			Constitution: 0,
			Dexterity:    0,
			Wisdom:       0,
			Strength:     0,
			Intelligence: 0,
			Charisma:     0,
		},

		//Initial Mana Points
		ManaPoints: models.ManaPoints{
			MaxMana:  0,
			Actual:   0,
			TempMana: 0,
		},
		//Initial HP Points
		HpPoints: models.HpPoints{
			MaxHp:  0,
			Actual: 0,
			TempHp: 0,
		},


		//Initial Armor
		Armor: models.Armor{
			Defense:        10,
			DexterityBonus: false,
			ShieldBonus:    0,
			OtherBonus:     0,
		},

		//Initial sheet not have abilities
		Abilities: []models.Ability{},

		//Initial Character not have attacks
		Attacks: []models.Attack{},

		//Initial Character not have equipaments
		EquipmentItems: []models.EquipmentItem{},

		Skills: map[string]models.Skill{},

		ClassAndLevel: models.ClassAndLevel{Class: "", Level: 1},
	}

	 //For to search all skills and your values
		for key, value := range DefaultT20Skills {
			sheet.Skills[key] = value
	}

		return sheet, nil
	}
