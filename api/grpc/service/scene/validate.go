package scene

import (
	"fmt"

	"github.com/GarotoCowboy/vttProject/api/grpc/pb/scene"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

func Validate(req *scene.CreateSceneRequest) error {

	if req.TableId == 0 && req.Name == "" && req.Width == 0 && req.Height == 0 && req.GridType == 0 {
		return status.Errorf(codes.InvalidArgument, "mandatory fields not filled in")
	}

	if req.TableId == 0 {
		return ErrParamIsRequired("tableId", "uint64")
	}

	if req.Name == "" {
		return ErrParamIsRequired("name", "string")
	}

	if req.Width == 0 {
		return ErrParamIsRequired("width", "uint64")
	}
	if req.Height == 0 {
		return ErrParamIsRequired("height", "uint64")
	}
	//if req.GridType < 0 && req.GridType > 1 {
	//	return ErrParamIsRequired("gridType", "consts.GridType")
	//}

	return nil
}

func ValidateAndBuildUpdateMap(req *scene.EditSceneRequest) (map[string]interface{}, error) {

	if req == nil || req.Scene.SceneId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "the requisition and scene id cannot be null")
	}

	if req.Scene.SceneId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "scene id is necessary")
	}

	if req.Scene.TableId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "table id is necessary")
	}

	mask := req.GetMask()
	if mask == nil || len(mask.GetPaths()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "FieldMask is mandatory and must specify at least one field to update")
	}

	getScene := req.GetScene()
	updatesMap := make(map[string]interface{})

	for _, path := range mask.GetPaths() {
		switch path {
		case "name":
			if getScene.GetName() == "" {
				return nil, status.Errorf(codes.InvalidArgument, "scene name cant't be empty")
			}
			updatesMap["name"] = getScene.GetName()

		case "height":
			if getScene.GetHeight() == 0 {
				return nil, status.Errorf(codes.InvalidArgument, "scene height cant't be empty")
			}
			updatesMap["height"] = getScene.GetHeight()

		case "width":
			if getScene.GetWidth() == 0 {
				return nil, status.Errorf(codes.InvalidArgument, "scene width cant't be empty")
			}
			updatesMap["width"] = getScene.GetWidth()

		case "background_image_path":
			if getScene.GetBackgroundImagePath() == "" {
				return nil, status.Errorf(codes.InvalidArgument, "scene background_image_path cant't be empty")
			}
			updatesMap["background_image_path"] = getScene.GetBackgroundImagePath()
		case "is_visible":
			if getScene.IsVisible != nil {
				updatesMap["is_visible"] = getScene.IsVisible.GetValue()
			} else {
				return nil, status.Errorf(codes.InvalidArgument, "field 'is_visible' is in the mask but value is null'")
			}
			//case "grid_type":
			//	{
			//		if getScene.GridType < 0 && getScene.GridType > 1 {
			//			return nil, status.Errorf(codes.InvalidArgument, "gridType  be equal to 0 or GridTyper equal to 1")
			//		}
			//		updatesMap["grid_type"] = getScene.GridType
			//}

		}

	}
	return updatesMap, nil
}
