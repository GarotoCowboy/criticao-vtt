package libraryImg

import (
	"bytes"
	"fmt"
	"image"
	"net/http"

	"github.com/GarotoCowboy/vttProject/api/grpc/proto/imageLibrary/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

func ValidateUploadInit(initMsg *pb.UploadInit) error {

	if initMsg == nil {
		return fmt.Errorf("The first message needs be initialization (init)")
	}

	if initMsg.TableId == 0 {
		return ErrParamIsRequired("tableId", "uint64")
	}

	if initMsg.Name == "" {
		return ErrParamIsRequired("name", "string")
	}

	return nil
}

func validateImageContent(imageBytes []byte) (*image.Config, string, error) {

	if len(imageBytes) == 0 {
		return nil, "", ErrParamIsRequired("imageBytes", "[]bytes")
	}

	imgConfig, _, err := image.DecodeConfig(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, "", ErrParamIsRequired("Image", "PNG or JPEG or JPG")
	}

	contentType := http.DetectContentType(imageBytes)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" {
		return nil, "", ErrParamIsRequired("Image", "PNG or JPEG or JPG")
	}

	return &imgConfig, contentType, nil

}

func ValidadeAndBuildUpdateMap(req *pb.EditImageRequest) (map[string]interface{}, error) {

	if req == nil || req.Image.ImageId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "the requisition and image id cannot be null")
	}

	if req.Image.ImageId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "image id is necessary")
	}

	if req.Image.TableId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "table id is necessary")
	}

	mask := req.GetMask()
	if mask == nil || len(mask.GetPaths()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "FieldMask is mandatory and must specify at least one field to update")
	}

	getImage := req.GetImage()
	updatesMap := make(map[string]interface{})

	for _, path := range mask.GetPaths() {
		switch path {
		case "name":
			if getImage.GetName() == "" {
				return nil, status.Errorf(codes.InvalidArgument, "image name cant't be empty")
			}
			updatesMap["name"] = getImage.GetName()
		}
	}
	return updatesMap, nil
}
