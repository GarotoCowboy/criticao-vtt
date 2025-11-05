package table

import (
	"errors"

	"github.com/GarotoCowboy/vttProject/api/dto/tableDTO"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"gorm.io/gorm"
)

func CreateTable(db *gorm.DB, ownerID uint, req tableDTO.CreateTableRequest) (models.Table, error) {

	//Validate if request bod is valid
	if err := req.Validate(); err != nil {
		return models.Table{}, err
	}

	generatedInviteLink, err := createAndVerifyLinkCodeNotExists(db)
	if err != nil {
		return models.Table{}, err
	}

	var table = models.Table{
		Name:     req.Name,
		Password: req.Password,
		OwnerID:  ownerID,
	}
	table.InviteLink = generatedInviteLink

	//hash password to bcrypt
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.Table{}, err
	}

	table.Password = string(hashedPassword)

	//Save Table
	if err := db.Create(&table).Error; err != nil {
		return models.Table{}, err
	}

	ownerMember := models.TableUser{
		TableID: table.ID,
		UserID:  ownerID,
		Role:    consts.Role(2),
	}

	if err := db.Create(&ownerMember).Error; err != nil {
		return models.Table{}, err
	}

	var fullTable models.Table
	if err := db.Preload("Owner").First(&fullTable, table.ID).Error; err != nil {
		return models.Table{}, err
	}
	return fullTable, nil
}

func DeleteTable(db *gorm.DB, tableID, userID uint) (models.Table, error) {

	//Verify if tableID is valid
	if tableID == 0 {
		return models.Table{}, errors.New("invalid table_id")
	}

	var tableUserModel = models.TableUser{}

	if err := db.Where("table_id = ? AND user_id = ?", tableID, userID).First(&tableUserModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Table{}, errors.New("tableUser is not a member of this table")
		}
	}
	if tableUserModel.Role != 2 {
		return models.Table{}, errors.New("only masters can delete a table")
	}

	var table = models.Table{}

	//Preload all members registred at table
	if err := db.Preload("Members").First(&table, tableID).Error; err != nil {
		return models.Table{}, err
	}

	//Select all members
	if err := db.Select("Members").Delete(&table).Error; err != nil {
		return models.Table{}, err
	}
	//Delete table and members with inputted tableID

	if err := db.Delete(&table).Error; err != nil {
		return models.Table{}, err
	}
	return table, nil
}

func GetTable(db *gorm.DB, tableID, userID uint) (models.Table, error) {
	if tableID == 0 {
		return models.Table{}, errors.New("invalid table_id")
	}

	var membership models.TableUser
	if err := db.Where("table_id = ? AND user_id = ?", tableID, userID).First(&membership).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Table{}, errors.New("table not found")
		}
		return models.Table{}, err
	}

	var table = models.Table{}

	if err := db.Preload("Owner").Where("id = ?", tableID).First(&table).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Table{}, errors.New("table not found")
		}
		return models.Table{}, err
	}

	return table, nil
}

func ListTables(db *gorm.DB, userID uint) ([]models.Table, error) {

	var tables []models.Table

	if err := db.Joins("JOIN table_users ON table_users.table_id = tables.id").
		Where("table_users.user_id = ?", userID).
		Find(&tables).Error; err != nil {
		return nil, err
	}

	//if err := db.Preload("Owner").Find(&tables).Error; err != nil {
	//	return tables, err
	//}
	return tables, nil
}

func UpdateTable(db *gorm.DB, tableID, userID uint, req tableDTO.UpdateTableRequest) (models.Table, error) {

	if tableID == 0 {
		return models.Table{}, errors.New("invalid table_id")
	}

	var tableUser models.TableUser

	if err := db.Where("table_id = ? AND user_id = ?", tableID, userID).First(&tableUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Table{}, errors.New("tableUser is not a member of this table")
		}
		return models.Table{}, err
	}

	if tableUser.Role != 2 {
		return models.Table{}, errors.New("only masters can change a table")
	}

	tableData, err := GetTable(db, tableID, userID)
	if err != nil {
		return models.Table{}, err
	}

	if err := req.Validate(); err != nil {
		return models.Table{}, err
	}

	//Update userDTO

	if req.Name != "" {
		tableData.Name = req.Name
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return models.Table{}, err
		}

		tableData.Password = string(hashedPassword)
	}

	if err := db.Save(&tableData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Table{}, errors.New("table not found")
		}
		return models.Table{}, err
	}
	return tableData, nil
}

func createAndVerifyLinkCodeNotExists(db *gorm.DB) (string, error) {
	var count int64
	var link string
	for {
		link = utils.StringWithCharset(10)

		err := db.
			Model(&models.Table{}).
			Where("invite_link = ?", link).
			Count(&count).Error
		if err != nil {
			return "", err
		}

		if count == 0 {
			break
		}
	}
	return link, nil

}
