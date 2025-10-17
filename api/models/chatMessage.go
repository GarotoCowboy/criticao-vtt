package models

import (
	"time"

	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ChatMessage struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time

	TableID uint `gorm:"not null;index"`

	DeletedAt gorm.DeletedAt `gorm:"index"`
	// public message
	TableUser   TableUser `gorm:"constraint:OnDelete:CASCADE"`
	TableUserID uint      `gorm:"not null;index"`

	//private message
	ToTableUserId *uint      `gorm:"index"`
	ToTableUser   *TableUser `gorm:"constraint:OnDelete:SET NULL"`

	Message          string               `json:"message" gorm:"not null"`
	MediaURL         *string              `json:"media_url"`
	Attachments      datatypes.JSON       `json:"attachments" gorm:"type:jsonb"`
	ReplyToMessageId *string              `gorm:"index"`
	MessageType      consts.MessageType   `json:"messageType" gorm:"not null"`
	MessageStatus    consts.MessageStatus `json:"messageStatus" gorm:"not null"`
}
