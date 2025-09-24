package token

import (
	"fmt"

	"github.com/GarotoCowboy/vttProject/api/grpc/proto/token/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Validate(req *pb.CreateTokenRequest) error {

	if req.TableId == 0 {
		return ErrParamIsRequired("tableId", "tableName")
	}

	return nil
}

func ValidadeAndBuildUpdateMap(req *pb.EditTokenRequest) (map[string]interface{}, error) {

	if req == nil || req.GetToken() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "the requisition and token cannot be null")
	}

	if req.GetToken().GetTokenId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "tokenID is necessary")
	}

	if req.GetToken().GetTableId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "tableID is necessary")
	}

	mask := req.GetMask()
	if mask == nil || len(mask.GetPaths()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "FieldMask is mandatory and must specify at least one field to update")
	}

	token := req.GetToken()
	updatesMap := make(map[string]interface{})

	for _, path := range mask.GetPaths() {
		switch path {
		case "name":
			updatesMap["name"] = token.GetName()
		case "image_url":
			updatesMap["imageUrl"] = token.GetImageUrl()
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unknown or not allowed field in mask: '%s'", path)
		}
	}

	if len(updatesMap) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no valid fields for update were provided in the mask")
	}
	return updatesMap, nil
}
