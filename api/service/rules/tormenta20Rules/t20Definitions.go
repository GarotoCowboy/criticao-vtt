package tormenta20Rules

type ClassDefinition struct {
	ID                string
	Name              string
	Description       string
	HitPointsPerLevel int
	ManaPerLevel      int
}

var AvaliableClasses = map[string]ClassDefinition{
	"Arcanist": {
		ID: "Arcanist",
		Name: "Arcanist",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Barbarian": {
		ID: "Barbarian",
		Name: "Barbarian",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Bard": {
		ID: "Bard",
		Name: "Bard",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Buccaneer": {
		ID: "Buccaneer",
		Name: "Buccaneer",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Hunter": {
		ID: "Hunter",
		Name: "Hunter",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Knight": {
		ID: "Knight",
		Name: "Knight",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Cleric": {
		ID: "Cleric",
		Name: "Cleric",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Druid": {
		ID: "Druid",
		Name: "Druid",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Warrior": {
		ID: "Warrior",
		Name: "Warrior",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Inventor": {
		ID: "Inventor",
		Name: "Inventor",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Rogue": {
		ID: "Rogue",
		Name: "Rogue",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Fighter": {
		ID: "Fighter",
		Name: "Fighter",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Noble": {
		ID: "Noble",
		Name: "Noble",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
	"Paladin": {
		ID: "Paladin",
		Name: "Paladin",
		Description: "",
		HitPointsPerLevel: 1,
		ManaPerLevel:      1,
	},
}
