package consts

import "fmt"

type Role uint8

const (
	Unspecified Role = iota
	Player
	Master
)

func SetRole(s Role) (Role, error) {
	switch s {
	case Master, Player:
		return s, nil

	default:
		return 0, fmt.Errorf("role does not exist")
	}
}
