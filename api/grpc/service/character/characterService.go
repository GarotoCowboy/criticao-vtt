package character

import (
	"sync"

	"github.com/GarotoCowboy/vttProject/api/grpc/pb/character"
	"github.com/GarotoCowboy/vttProject/api/service/rules/tormenta20Rules"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type CharacterService struct {
	character.UnimplementedCharacterServiceServer
	Db     *gorm.DB
	Logger *config.Logger
	mu     sync.RWMutex
	rules  *tormenta20Rules.RulesService
	subscribers map[uint]map[string]character.CharacterService_UpdateSheetServer
}

// function that initialize the CharacterService struct
func NewCharacterService(db *gorm.DB, logger *config.Logger) *CharacterService {
	return &CharacterService{
		Db:     db,
		Logger: logger,
		rules:  tormenta20Rules.NewRulesService(),
		subscribers: make(map[uint]map[string]character.CharacterService_UpdateSheetServer),
	}
}
