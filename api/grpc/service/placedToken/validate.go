package placedToken

import (
	"fmt"

	"github.com/GarotoCowboy/vttProject/api/grpc/pb/placedToken"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

func Validate(req *placedToken.CreatePlacedTokenRequest) error {

	if req.SceneId == 0 && req.TokenId == 0 && req.PosX == 0 && req.PosY == 0 {
		return status.Errorf(codes.InvalidArgument, "mandatory fields not filled in")
	}

	if req.SceneId == 0 {
		return ErrParamIsRequired("scene_id", "uint64")
	}

	if req.TokenId == 0 {
		return ErrParamIsRequired("token_id", "uint64")
	}

	if req.PosY == 0 {
		return ErrParamIsRequired("pos_x", "int")
	}
	if req.PosY == 0 {
		return ErrParamIsRequired("pos_y", "int")
	}

	return nil
}

func ValidateAndBuildUpdateMap(req *placedToken.EditPlacedTokenRequest) (map[string]interface{}, error) {

	if req == nil || req.PlacedToken.SceneId == 0 || req.PlacedToken.TokenId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "the requisition and scene id and token id cannot be null")
	}

	if req.PlacedToken.SceneId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "scene id is necessary")
	}

	if req.PlacedToken.PlacedTokenId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "table id is necessary")
	}

	mask := req.GetMask()
	if mask == nil || len(mask.GetPaths()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "FieldMask is mandatory and must specify at least one field to update")
	}

	getPlacedToken := req.GetPlacedToken()
	updatesMap := make(map[string]interface{})

	for _, path := range mask.GetPaths() {
		switch path {
		case "layer":

			layerValue := getPlacedToken.GetLayer()

			if _, ok := placedToken.LayerType_name[int32(getPlacedToken.Layer)]; !ok {
				return nil, status.Errorf(codes.InvalidArgument, "layer type is invalid %d", layerValue)
			}
			updatesMap["layer"] = getPlacedToken.GetLayer()

		case "size":
			if getPlacedToken.GetSize() == 0 {
				return nil, status.Errorf(codes.InvalidArgument, "placedToken size cant't be empty")
			}
			updatesMap["size"] = getPlacedToken.GetSize()

		case "rotation":
			if getPlacedToken.GetRotation() == 0 {
				return nil, status.Errorf(codes.InvalidArgument, "placedToken rotation cant't be empty")
			}
			updatesMap["rotation"] = getPlacedToken.GetRotation()
		}

	}
	return updatesMap, nil
}

func MoveTokenValidate(req *placedToken.MoveTokenRequest) error {

	if req.SceneId == 0 && req.PlacedTokenId == 0 && req.PosX == 0 && req.PosY == 0 {
		return status.Errorf(codes.InvalidArgument, "mandatory fields not filled in")
	}

	if req.SceneId == 0 {
		return ErrParamIsRequired("scene_id", "uint64")
	}

	if req.PlacedTokenId == 0 {
		return ErrParamIsRequired("token_id", "uint64")
	}

	if req.PosY == 0 {
		return ErrParamIsRequired("pos_x", "int")
	}
	if req.PosY == 0 {
		return ErrParamIsRequired("pos_y", "int")
	}

	return nil
}
