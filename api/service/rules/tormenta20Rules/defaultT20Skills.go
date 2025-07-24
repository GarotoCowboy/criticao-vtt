package tormenta20Rules

import "github.com/GarotoCowboy/vttProject/api/models"

var DefaultT20Skills = map[string]models.Skill{
	"Acrobatics":    {DefaultBaseAttribute: "dexterity",    CurrentBaseAttribute: "dexterity",    Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: true},
	"Dressage":      {DefaultBaseAttribute: "charisma",     CurrentBaseAttribute: "charisma",     Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: false},
	"Athletics":     {DefaultBaseAttribute: "strength",     CurrentBaseAttribute: "strength",     Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Performance":   {DefaultBaseAttribute: "charisma",     CurrentBaseAttribute: "charisma",     Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: false}, // CORRIGIDO: Atributo era dexterity
	"Riding":        {DefaultBaseAttribute: "dexterity",    CurrentBaseAttribute: "dexterity",    Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Knowledge":     {DefaultBaseAttribute: "intelligence", CurrentBaseAttribute: "intelligence", Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: false},
	"Healing":       {DefaultBaseAttribute: "wisdom",       CurrentBaseAttribute: "wisdom",       Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false}, // CORRIGIDO: Atributo era intelligence
	"Diplomacy":     {DefaultBaseAttribute: "charisma",     CurrentBaseAttribute: "charisma",     Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Deception":     {DefaultBaseAttribute: "charisma",     CurrentBaseAttribute: "charisma",     Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Fortitude":     {DefaultBaseAttribute: "constitution", CurrentBaseAttribute: "constitution", Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Stealth":       {DefaultBaseAttribute: "dexterity",    CurrentBaseAttribute: "dexterity",    Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: true}, // CORRIGIDO: Typo na chave
	"Warfare":       {DefaultBaseAttribute: "intelligence", CurrentBaseAttribute: "intelligence", Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: false},
	"Initiative":    {DefaultBaseAttribute: "dexterity",    CurrentBaseAttribute: "dexterity",    Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Intimidation":  {DefaultBaseAttribute: "charisma",     CurrentBaseAttribute: "charisma",     Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Intuition":     {DefaultBaseAttribute: "wisdom",       CurrentBaseAttribute: "wisdom",       Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Investigation": {DefaultBaseAttribute: "intelligence", CurrentBaseAttribute: "intelligence", Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Gambling":      {DefaultBaseAttribute: "charisma",     CurrentBaseAttribute: "charisma",     Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: false},
	"Thieving":      {DefaultBaseAttribute: "dexterity",    CurrentBaseAttribute: "dexterity",    Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: true},
	"Fighting":      {DefaultBaseAttribute: "strength",     CurrentBaseAttribute: "strength",     Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Mysticism":     {DefaultBaseAttribute: "intelligence", CurrentBaseAttribute: "intelligence", Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: false},
	"Nobility":      {DefaultBaseAttribute: "intelligence", CurrentBaseAttribute: "intelligence", Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: false},
	"Craft1":        {DefaultBaseAttribute: "intelligence", CurrentBaseAttribute: "intelligence", Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: false},
	"Craft2":        {DefaultBaseAttribute: "intelligence", CurrentBaseAttribute: "intelligence", Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: false},
	"Perception":    {DefaultBaseAttribute: "wisdom",       CurrentBaseAttribute: "wisdom",       Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Piloting":      {DefaultBaseAttribute: "dexterity",    CurrentBaseAttribute: "dexterity",    Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: false},
	"Aiming":        {DefaultBaseAttribute: "dexterity",    CurrentBaseAttribute: "dexterity",    Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Reflexes":      {DefaultBaseAttribute: "dexterity",    CurrentBaseAttribute: "dexterity",    Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Religion":      {DefaultBaseAttribute: "wisdom",       CurrentBaseAttribute: "wisdom",       Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: true,  ArmorPenalty: false},
	"Survival":      {DefaultBaseAttribute: "wisdom",       CurrentBaseAttribute: "wisdom",       Trained: false, Bonus: 0, OtherBonus:0, OnlyTrained: false, ArmorPenalty: false},
	"Willpower":     {DefaultBaseAttribute: "wisdom",        CurrentBaseAttribute: "wisdom",       Trained: false, Bonus: 0,OtherBonus:0 ,OnlyTrained: false, ArmorPenalty: false},
}
