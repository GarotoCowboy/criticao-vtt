package bar

import (
	"fmt"
	"regexp"

	"github.com/GarotoCowboy/vttProject/api/grpc/pb/bar"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// regex to validate if the color is a hexadecimal
var hexColor = regexp.MustCompile("^#([a-fA-F0-9]){3}(([a-fA-F0-9]){3})?$")

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Validate(req *bar.CreateBarRequest) error {

	if req.TokenId == 0 {
		return ErrParamIsRequired("tokenId", "uint64")
	}
	if req.Name == "" {
		return ErrParamIsRequired("name", "string")
	}
	if !isHexColor(req.Color) {
		return ErrParamIsRequired("color", "hexadecimal")
	}

	return nil
}

func ValidadeAndBuildUpdateMap(req *bar.EditBarRequest) (map[string]interface{}, error) {

	if req == nil || req.GetBar() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "the requisition and bar id cannot be null")
	}

	if req.GetBar().GetBarId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "bar id is necessary")
	}

	if req.GetBar().GetTokenId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "token id is necessary")
	}

	if !isHexColor(req.GetBar().GetColor()) {
		return nil, status.Errorf(codes.InvalidArgument, "color is a hexadecimal value")
	}

	mask := req.GetMask()
	if mask == nil || len(mask.GetPaths()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "FieldMask is mandatory and must specify at least one field to update")
	}

	bar := req.GetBar()
	updatesMap := make(map[string]interface{})

	for _, path := range mask.GetPaths() {
		switch path {
		case "name":
			updatesMap["name"] = bar.GetName()
		case "value":
			updatesMap["value"] = bar.GetValue()
		case "max_value":
			updatesMap["maxValue"] = bar.GetMaxValue()
		case "color":
			updatesMap["color"] = bar.GetColor()
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unknown or not allowed field in mask: '%s'", path)
		}
	}

	if len(updatesMap) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no valid fields for update were provided in the mask")
	}
	return updatesMap, nil
}

func isHexColor(s string) bool {
	return hexColor.MatchString(s)
}
