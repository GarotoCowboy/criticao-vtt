package tormenta20Rules

import (
	"fmt"

	"github.com/GarotoCowboy/vttProject/api/grpc/pb/character"

	"strings"
)

type RulesService struct {
}

type InitialCharacterInputsT20 struct {
}

type Attribute string

const (
	Intelligence Attribute = "Intelligence"
	Charisma     Attribute = "Charisma"
	Strength     Attribute = "Strength"
	Dexterity    Attribute = "Dexterity"
	Constitution Attribute = "Constitution"
	Wisdom       Attribute = "Wisdom"
)

func NewRulesService() *RulesService {
	return &RulesService{}
}

// create a function that search an attribute based on const Attribute
//func (s *RulesService) getT20Attribute(sheet *models.T20Sheet ,attribute Attribute) (int, error) {
//
//	switch attribute {
//	case Intelligence:
//		return sheet.Attributes.Intelligence, nil
//	case Charisma:
//		return sheet.Attributes.Charisma, nil
//	case Strength:return sheet.Attributes.Strength, nil
//	case Dexterity:return sheet.Attributes.Dexterity, nil
//	case Constitution:return sheet.Attributes.Constitution, nil
//	case Wisdom:return sheet.Attributes.Wisdom, nil
//	default:return 0, fmt.Errorf("Value not Found")
//	}
//
//}

// create a function that search an attribute based on const Attribute
func (s *RulesService) getT20Attribute(sheet *character.Sheet, attribute Attribute) (int, error) {

	switch attribute {
	case Intelligence:
		return int(sheet.Attributes.Intelligence), nil
	case Charisma:
		return int(sheet.Attributes.Charisma), nil
	case Strength:
		return int(sheet.Attributes.Strength), nil
	case Dexterity:
		return int(sheet.Attributes.Dexterity), nil
	case Constitution:
		return int(sheet.Attributes.Constitution), nil
	case Wisdom:
		return int(sheet.Attributes.Wisdom), nil
	default:
		return 0, fmt.Errorf("Value not Found")
	}

}

// function that pick the level and return the bonus based if he's trained on that skill
func (s *RulesService) calculateTrainedBonus(level int, isTrained bool) (int, error) {

	if !isTrained {
		return 0, nil
	}

	if level >= 1 && level <= 6 {

		return 2, nil
	} else if level > 6 && level <= 14 {

		return 4, nil
	} else if level > 14 {

		return 6, nil
	}
	return 0, fmt.Errorf("invalid level, level must be equals or higher than 1")
}

//func (s *RulesService) CalculateSheetSkillsAutomatically(sheet *models.T20Sheet) (*models.T20Sheet, error) {
//
//	if sheet == nil || sheet.Skills == nil {
//		return nil, fmt.Errorf("sheet or skills is nil")
//	}
//
//	levelBonus := sheet.ClassAndLevel.Level / 2
//
//	for skillName, skillData := range sheet.Skills {
//
//		baseAttribute, ok := skillAttributeMap[skillName]
//		if !ok {
//			rules.GetSystemRulesLogger().WarningF("The skill not have a attribute rule defined: %s", skillName)
//			continue
//		}
//
//		attribute, err := s.getT20Attribute(sheet,baseAttribute)
//		if err != nil {
//			return nil, err
//		}
//
//		trainedBonus, err := s.calculateTrainedBonus(sheet.ClassAndLevel.Level, skillData.Trained)
//		if err != nil {
//			return nil, err
//		}
//
//		skillData.Bonus = attribute + levelBonus + trainedBonus
//		sheet.Skills[skillName] = skillData
//	}
//
//	return sheet, nil
//}

func (s *RulesService) CalculateSheetSkillsAutomatically(sheet *character.Sheet) (*character.Sheet, error) {

	if sheet == nil || sheet.Skills == nil {
		return nil, fmt.Errorf("sheet or skills is nil")
	}

	levelBonus := int(sheet.ClassAndLevel.Level / 2)

	for skillName, skillData := range sheet.Skills {

		//baseAttribute, ok := skillAttributeMap[skillName]
		//if !ok {
		//	rules.GetSystemRulesLogger().WarningF("The skill not have a attribute rule defined: %s", skillName)
		//	continue
		//}

		normalizedAttr, err := normalizeAttributeName(skillData.CurrentBaseAttribute)
		if err != nil {
			return nil, fmt.Errorf("skill '%s' have an invalid attribute: %v", skillName, err)
		}

		skillData.CurrentBaseAttribute = normalizedAttr

		baseAttribute := Attribute(normalizedAttr)

		attribute, err := s.getT20Attribute(sheet, baseAttribute)
		if err != nil {
			return nil, fmt.Errorf("error to pick attribute %s: %v\n", baseAttribute, err)
		}

		trainedBonus, err := s.calculateTrainedBonus(int(sheet.ClassAndLevel.Level), skillData.Trained)
		if err != nil {
			return nil, fmt.Errorf("error to calculate trained Bonus %s: %v\n", baseAttribute, err)
		}

		skillData.Bonus = int32(attribute+levelBonus+trainedBonus) + skillData.OtherBonus
		sheet.Skills[skillName] = skillData
	}

	return sheet, nil
}

func (s *RulesService) CalculateSheetDefenseAutomatically(sheet *character.Sheet) (*character.Sheet, error) {

	attributeBonus, err := s.getT20Attribute(sheet, Dexterity)
	if err != nil {
		return nil, fmt.Errorf("error to pick attribute %s: %v\n", attributeBonus, err)
	}

	baseBonus := int32(10)
	sheet.Armor.Defense = baseBonus + sheet.Armor.ArmorBonus + sheet.Armor.OtherBonus + sheet.Armor.ShieldBonus
	if sheet.Armor.DexterityBonus == true {
		sheet.Armor.Defense += int32(attributeBonus)
	}
	return sheet, nil
}

func normalizeAttributeName(attr string) (string, error) {
	switch attrLower := strings.ToLower(attr); attrLower {
	case "strength":
		return "Strength", nil
	case "dexterity":
		return "Dexterity", nil
	case "constitution":
		return "Constitution", nil
	case "intelligence":
		return "Intelligence", nil
	case "wisdom":
		return "Wisdom", nil
	case "charisma":
		return "Charisma", nil
	default:
		return "", fmt.Errorf("invalid attribute: '%s'", attr)
	}
}
