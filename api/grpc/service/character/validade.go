package character

import (
	"fmt"
	"github.com/GarotoCowboy/vttProject/api/grpc/proto/character/pb"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Validate(req *pb.CreateCharacterRequest) error {

	if req.TableUserId == 0 && req.SystemKey == 0 && req.PlayerName == "" {
		return fmt.Errorf("request body is empty")
	}
	if req.TableUserId == 0 {
		return ErrParamIsRequired("tableUserId", "string")
	}
	if req.SystemKey == 0 {
		return ErrParamIsRequired("systemKey", "const.SystemKey")
	}
	if req.CharacterName == "" {
		return ErrParamIsRequired("name", "string")
	}

	return nil
}
