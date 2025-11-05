package consts

type PermissionLevel uint

const (
	PermissionNone PermissionLevel = iota
	PermissionMaster
	PermissionOwnerAndMaster
	PermissionAllPlayers
)

//type Permission struct {
//
//	OwnerId []uint `json:"ownerID"`
//
//	CanBeModifiedBy PermissionLevel `json:"can_be_modified_by" gorm:"default:1"`
//
//}
