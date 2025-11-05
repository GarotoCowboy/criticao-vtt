package character

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/GarotoCowboy/vttProject/api/grpc/pb/character"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"github.com/GarotoCowboy/vttProject/api/service/rules/tormenta20Rules"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// this grpc function creates an character
func (c *CharacterService) CreateCharacter(ctx context.Context, req *character.CreateCharacterRequest) (*character.CreateCharacterResponse, error) {

	//Validate the fields of characterModel
	if err := Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid Request Body: %v", err.Error())
	}

	//I need change this, it's temporary, but i use to filter the sheet based on the system
	//I need implement the sheet logic too
	if req.SystemKey != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "SystemKey %d not supported yet", req.SystemKey)
	}

	//Temporary Code
	rules := tormenta20Rules.RulesService{}
	t20Sheet, err := rules.GenerateInitialSheetData()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot generate initial sheet: %v", err)
	}
	//Marshalling the sheet to byte to convert to jsonB on postgres
	sheetBytes, err := json.Marshal(t20Sheet)
	if err != nil {
		return &character.CreateCharacterResponse{}, fmt.Errorf("error marshalling sheet data: %w", err)
	}

	tableUser := models.TableUser{}

	if err := c.Db.Where("id = ?", req.TableUserId).First(&tableUser).Error; err != nil {
		return &character.CreateCharacterResponse{}, status.Error(codes.NotFound, "Table tableUser not found")
	}

	var characterModel = models.Character{
		Name:        req.CharacterName,
		PlayerName:  tableUser.User.Username,
		SystemKey:   consts.SystemKey(req.SystemKey),
		TableUserID: uint(req.TableUserId),
		SheetData:   sheetBytes,
	}
	//create characterModel
	if err := c.Db.Create(&characterModel).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "error creating characterModel: %w", err)
	}

	c.Logger.InfoF("Character created: %v", characterModel)

	return &character.CreateCharacterResponse{
		CharacterName: req.CharacterName,
		//Corrigir Depois o codigo abaixo
		//SheetData:   (sheetBytes),
		SystemKey:  string(req.SystemKey),
		PlayerName: req.PlayerName,
	}, nil

}

func (c *CharacterService) UpdateSheet(stream character.CharacterService_UpdateSheetServer) error {
	ctx := stream.Context()
	var subID string
	var charID uint

	opts := protojson.MarshalOptions{
		EmitUnpopulated: true,
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if subID != "" {
				c.unsubscribe(charID, subID)
			}
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "recv error: %v", err)
		}

		if subID == "" {
			charID = uint(req.GetCharacterId())
			subID = c.subscribe(charID, stream)
			defer c.unsubscribe(charID, subID)
		}

		//search for a character from the database
		charReq := &character.GetCharacterRequest{
			CharacterId: uint32(charID),
			TableId:     req.GetTableId(),
		}

		charResp, err := c.GetCharacter(ctx, charReq)
		if err != nil {
			return status.Errorf(codes.NotFound, "character not found: %v", err)
		}

		//sheet := charResp.GetSheet()

		//checks if the character sheet attributes are different from null for the update
		if req.Sheet.Attributes != nil {
			charResp.Sheet.Attributes = req.Sheet.Attributes
		}

		if req.Sheet.Skills != nil {
			charResp.Sheet.Skills = req.Sheet.Skills
		}

		if req.Sheet.ClassAndLevel != nil {
			charResp.Sheet.ClassAndLevel = req.Sheet.ClassAndLevel
		}

		if req.Sheet.Armor != nil {
			charResp.Sheet.Armor = req.Sheet.Armor
		}

		if req.Sheet.HpPoints != nil {
			charResp.Sheet.HpPoints = req.Sheet.HpPoints
		}

		if req.Sheet.EquipmentItems != nil {
			charResp.Sheet.EquipmentItems = req.Sheet.EquipmentItems
		}

		if req.Sheet.Attacks != nil {
			charResp.Sheet.Attacks = req.Sheet.Attacks
		}

		if req.Sheet.Abilities != nil {
			charResp.Sheet.Abilities = req.Sheet.Abilities
		}

		if req.Sheet.ManaPoints != nil {
			charResp.Sheet.ManaPoints = req.Sheet.ManaPoints
		}

		if req.Sheet.CharacterInfo != nil {
			charResp.Sheet.CharacterInfo = req.Sheet.CharacterInfo
		}

		if req.CharacterName != "" {
			charResp.Name = req.CharacterName
		}

		//applies business rules to calculate character sheet bonuses automatically

		bonusSheet, err := c.rules.CalculateSheetSkillsAutomatically(charResp.Sheet)
		if err != nil {
			return status.Errorf(codes.Internal, "could not calculate skill automatically: %v", err)
		}

		bonusSheet, err = c.rules.CalculateSheetDefenseAutomatically(charResp.Sheet)
		if err != nil {
			return status.Errorf(codes.Internal, "could not calculate armor bonus automatically: %v", err)
		}

		//marshall the form with all the information
		sheetBytes, err := opts.Marshal(bonusSheet)
		if err != nil {
			return status.Errorf(codes.Internal, "error marshalling sheet data: %v", err)
		}

		//Creates a character model with all the new information to update the character
		charModel := models.Character{
			Model:       gorm.Model{ID: charID},
			TableUserID: uint(req.GetTableId()),
			Name:        charResp.Name,
			SheetData:   sheetBytes,
		}

		if err := c.Db.WithContext(ctx).Save(&charModel).Error; err != nil {
			return status.Errorf(codes.Internal, "error updating character: %v", err)
		}

		resp := &character.CharacterUpdateResponse{
			CharacterName: charResp.Name,
			Sheet:         bonusSheet,
			LastModfield:  timestamppb.Now(),
		}

		c.mu.RLock()
		for id, sub := range c.subscribers[charID] {
			if err := sub.Send(resp); err != nil {
				c.unsubscribe(charID, id)
			}
		}
		c.mu.RUnlock()
		c.Logger.InfoF("broadcast update for character: %v", charID)

		//if err := stream.Send(resp); err != nil {
		//	return status.Errorf(codes.Internal, "Error to sending message: %v ", err.Error())
		//}
		//c.Logger.InfoF("character Updated, Name: %v",resp.CharacterName)
	}

}

// Search a sheet
func (c *CharacterService) GetCharacter(ctx context.Context, req *character.GetCharacterRequest) (*character.GetCharacterResponse, error) {
	if req.CharacterId <= 0 || req.TableId <= 0 {
		return &character.GetCharacterResponse{}, status.Errorf(codes.InvalidArgument, "character_Id or table_Id is invalid")
	}
	var characterId = req.CharacterId
	var characterModel models.Character
	var sheet character.Sheet

	if err := c.Db.Where("id=? AND table_user_id = ? AND deleted_at IS NULL", characterId, req.TableId).First(&characterModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &character.GetCharacterResponse{}, status.Errorf(codes.NotFound, "characterModel not found in this table", err)
		}
		return &character.GetCharacterResponse{}, err
	}
	if err := json.Unmarshal(characterModel.SheetData, &sheet); err != nil {
		return &character.GetCharacterResponse{}, status.Errorf(codes.Internal, "error unmarshalling sheet data: %v", err)
	}

	return &character.GetCharacterResponse{
		Sheet: &sheet,
		Name:  characterModel.Name,
	}, nil

}

// ListCharacter is a function that client make an requisiton and the stream server provides a list of Characters
// on a determinate table
func (c *CharacterService) ListCharacter(req *character.ListCharacterRequest, stream grpc.ServerStreamingServer[character.GetCharacterResponse]) error {
	if req.TableId <= 0 {
		return status.Errorf(codes.InvalidArgument, "table_Id is invalid")
	}

	var sheet character.Sheet
	var tableID = req.TableId
	var listCharacter []models.Character

	//DB request
	if err := c.Db.Where("table_user_id = ? AND deleted_at IS NULL", tableID).Find(&listCharacter).Error; err != nil {
		return status.Errorf(codes.Internal, "error fetching characters: %v", err)
	}

	//Verify if the list of characters is null
	if len(listCharacter) == 0 {
		return status.Errorf(codes.NotFound, "character not found in this table")
	}

	for _, characterFor := range listCharacter {
		if err := json.Unmarshal(characterFor.SheetData, &sheet); err != nil {
			return status.Errorf(codes.Internal, "error unmarshalling sheet data: %v", err)
		}
		characterResponse := &character.GetCharacterResponse{
			Name:  characterFor.Name,
			Sheet: &sheet,
		}

		if err := stream.Send(characterResponse); err != nil {
			return status.Errorf(codes.Internal, "error sending character: %v", err)
		}
		//todo Make to Test...
		//time.Sleep(time.Second * 1)
	}
	return nil
}

func (c *CharacterService) DeleteSheet(ctx context.Context, req *character.GetCharacterRequest) (*character.DeleteCharacterResponse, error) {

	//Verify if the Ids are valid
	if req.TableId <= 0 || req.CharacterId <= 0 {
		return &character.DeleteCharacterResponse{
			MessageStatus: "400",
			Message:       "character and table id is invalid",
		}, status.Errorf(codes.NotFound, "character and table id is invalid")
	}

	var characterModel = models.Character{}

	//Search first character with inputted id
	if err := c.Db.Where("id = ? AND table_user_id = ?", req.CharacterId, req.TableId).
		First(&characterModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &character.DeleteCharacterResponse{
				MessageStatus: "404",
				Message:       "character not found in the table",
			}, status.Errorf(codes.NotFound, "character not found in the table")
		}
		return nil, err
	}

	//Delete character with inputted id
	if err := c.Db.Delete(&characterModel).Error; err != nil {
		return &character.DeleteCharacterResponse{
			MessageStatus: "500",
			Message:       "cannot delete character",
		}, status.Errorf(codes.Internal, "cannot delete character: %v", err)
	}

	//sucess response(no content)
	return &character.DeleteCharacterResponse{
		MessageStatus: "204",
		Message:       "no content",
	}, nil
}
