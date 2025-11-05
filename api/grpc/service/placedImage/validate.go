package placedImage

import (
	"fmt"

	"github.com/GarotoCowboy/vttProject/api/grpc/pb/placedImage"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/placedToken"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

func Validate(req *placedImage.CreatePlacedImageRequest) error {

	if req.SceneId == 0 && req.ImageId == 0 && req.PosX == 0 && req.PosY == 0 {
		return status.Errorf(codes.InvalidArgument, "mandatory fields not filled in")
	}

	if req.SceneId == 0 {
		return ErrParamIsRequired("scene_id", "uint64")
	}

	if req.ImageId == 0 {
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

func ValidateAndBuildUpdateMap(req *placedImage.EditPlacedImageRequest) (map[string]interface{}, error) {

	if req == nil || req.PlacedImage.SceneId == 0 || req.PlacedImage.ImageId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "the requisition and scene id and token id cannot be null")
	}

	if req.PlacedImage.SceneId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "scene id is necessary")
	}

	if req.PlacedImage.PlacedImageId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "table id is necessary")
	}

	mask := req.GetMask()
	if mask == nil || len(mask.GetPaths()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "FieldMask is mandatory and must specify at least one field to update")
	}

	getPlacedImage := req.GetPlacedImage()
	updatesMap := make(map[string]interface{})

	for _, path := range mask.GetPaths() {
		switch path {
		case "layer":

			layerValue := getPlacedImage.GetLayer()

			if _, ok := placedToken.LayerType_name[int32(getPlacedImage.Layer)]; !ok {
				return nil, status.Errorf(codes.InvalidArgument, "layer type is invalid %d", layerValue)
			}
			updatesMap["layer"] = getPlacedImage.GetLayer()

		case "width":
			if getPlacedImage.GetWidth() == 0 {
				return nil, status.Errorf(codes.InvalidArgument, "placedToken width cant't be empty")
			}
			updatesMap["width"] = getPlacedImage.GetWidth()

		case "height":
			if getPlacedImage.GetHeight() == 0 {
				return nil, status.Errorf(codes.InvalidArgument, "placedToken height cant't be empty")
			}
			updatesMap["width"] = getPlacedImage.GetWidth()

		case "rotation":
			if getPlacedImage.GetRotation() == 0 {
				return nil, status.Errorf(codes.InvalidArgument, "placedToken rotation cant't be empty")
			}
			updatesMap["rotation"] = getPlacedImage.GetRotation()
		}

	}
	return updatesMap, nil
}

func MoveTokenValidate(req *placedImage.MoveImageRequest) error {

	if req.SceneId == 0 && req.PlacedImageId == 0 && req.PosX == 0 && req.PosY == 0 {
		return status.Errorf(codes.InvalidArgument, "mandatory fields not filled in")
	}

	if req.SceneId == 0 {
		return ErrParamIsRequired("scene_id", "uint64")
	}

	if req.PlacedImageId == 0 {
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
