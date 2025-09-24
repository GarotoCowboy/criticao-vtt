package models

import "gorm.io/gorm"

type (
	Attributes struct {
		Strength     int `json:"strength"`
		Dexterity    int `json:"dexterity"`
		Constitution int `json:"constitution"`
		Intelligence int `json:"intelligence"`
		Wisdom       int `json:"wisdom"`
		Charisma     int `json:"charisma"`
	}
	Skill struct {
		DefaultBaseAttribute string `json:"defaultBaseAttribute"`
		CurrentBaseAttribute string `json:"currentBaseAttribute"`
		Trained              bool   `json:"trained"`
		Bonus                int    `json:"bonus"`
		OnlyTrained          bool   `json:"onlyTrained"`
		ArmorPenalty         bool   `json:"armorPenalty"`
		OtherBonus           int    `json:"otherBonus"`
	}
	Ability struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	Armor struct {
		Defense        int  `json:"defense"`
		DexterityBonus bool `json:"dexterityBonus"`
		ArmorBonus     int  `json:"armorBonus"`
		ShieldBonus    int  `json:"shieldBonus"`
		OtherBonus     int  `json:"otherBonus"`
	}
	HpPoints struct {
		Actual int `json:"actual"`
		MaxHp  int `json:"maxHp"`
		TempHp int `json:"tempHp"`
	}
	ManaPoints struct {
		Actual   int `json:"actual"`
		MaxMana  int `json:"maxMana"`
		TempMana int `json:"tempMana"`
	}
	EquipmentItem struct {
		Name   string  `json:"name"`
		Amount int     `json:"amount"`
		Weight float64 `json:"weight"`
	}
	ClassAndLevel struct {
		Class string `json:"class"`
		Level int    `json:"level"`
	}
	Attack struct {
		Name       string `json:"name"`
		AttackTest string `json:"test"`
		Damage     string `json:"damage"`
		Critical   string `json:"critical"`
		DamageType string `json:"damageType"`
		Range      string `json:"range"`
	}

	CharacterInfo struct {
		Deity  string `json:"deity"`
		Notes  string `json:"notes"`
		Origin string `json:"origin"`
		Race   string `json:"race"`
	}
)

type T20Sheet struct {
	gorm.Model
	//CharacterID uint `json:"characterID"`
	Attributes     Attributes       `gorm:"embedded;embeddedPrefix:attr_" json:"attributes"`
	HpPoints       HpPoints         `gorm:"embedded;embeddedPrefix:hp_" json:"hpPoints"`
	ManaPoints     ManaPoints       `gorm:"embedded;embeddedPrefix:mana_" json:"manaPoints"`
	Armor          Armor            `gorm:"embedded;embeddedPrefix:armor_" json:"armor"`
	CharacterInfo  CharacterInfo    `gorm:"embedded;embeddedPrefix:CInfo_" json:"characterInfo"`
	ClassAndLevel  ClassAndLevel    `gorm:"embedded;embeddedPrefix:CaL_" json:"classAndLevel"`
	Attacks        []Attack         `json:"attacks" gorm:"type:jsonb"`
	Abilities      []Ability        `json:"abilities" gorm:"type:jsonb"`
	Skills         map[string]Skill `json:"skills" gorm:"type:jsonb"`
	EquipmentItems []EquipmentItem  `json:"equipmentItems" gorm:"type:jsonb"`
}
