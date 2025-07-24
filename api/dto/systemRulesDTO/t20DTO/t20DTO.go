package t20DTO

type AttributesDTO struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

// DTO to modify skills map
type SkillsDTO struct {
	DefaultBaseAttribute string
	CurrentBaseAttribute string
	Trained              bool
	Bonus                int
	OnlyTrained          bool
	ArmorPenalty         bool
	CraftName            string
}

// DTO to modify mana map
type ManaDTO struct {
	Actual   int
	MaxMana  int
	TempMana int
}

// DTO to modify hp map
type HpDTO struct {
	Actual int
	MaxHp  int
	TempHp int
}

type AbilitiesDTO struct {
	Name        string
	Description string
}

type ArmorDTO struct {
	Defense        int
	DexterityBonus bool
	ArmorBonus     int
	ShieldBonus    int
	OtherBonus     int
}

type EquipmentDTO struct {
	Name   string
	Amount int
	Weight float64
}

type ClassAndLevelDTO struct {
	Class string
	Level int
}

type AttacksDTO struct {
	Name        string
	AttackTest  string
	Damage      string
	Critical    string
	AttackRange string
}

type CharacterInfoDTO struct {
	Deity  string
	Notes  string
	Origin string
	Race   string
}
