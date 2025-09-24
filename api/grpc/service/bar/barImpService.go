package bar

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/GarotoCowboy/vttProject/api/grpc/proto/bar/pb"
	"github.com/GarotoCowboy/vttProject/api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *BarService) CreateBar(ctx context.Context, req *pb.CreateBarRequest) (*pb.CreateBarResponse, error) {
	//validate the codes
	if err := Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request body: %v", err.Error())
	}

	var token = models.Token{}

	//search the first token with inputted id
	if err := s.DB.Where("id = ?", req.GetTokenId()).First(&token).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid token id: %v", err.Error())
	}

	//create a bar var
	var bar = models.Bar{
		Name:     req.GetName(),
		MaxValue: req.GetMaxValue(),
		Value:    req.GetValue(),
		Color:    req.GetColor(),
		TokenID:  token.ID,
	}

	//create a bar in DB
	if err := s.DB.Create(&bar).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	s.Logger.InfoF("Token created: %v", bar)

	//bar response
	responseBar := &pb.Bar{
		Name:     bar.Name,
		Color:    bar.Color,
		Value:    bar.Value,
		MaxValue: bar.MaxValue,
		TokenId:  uint64(bar.TokenID),
	}

	//if successful we return a response created before
	return &pb.CreateBarResponse{
		Bar: responseBar,
	}, nil

}

func (s *BarService) EditBar(ctx context.Context, req *pb.EditBarRequest) (*pb.EditBarResponse, error) {
	//Verify if the mask and values are valid
	updatesMap, err := ValidadeAndBuildUpdateMap(req)
	if err != nil {
		return nil, err
	}

	var bar models.Bar
	var token models.Token

	if err := s.DB.Where("id = ?", req.GetBar().GetTokenId()).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "invalid token id: %v", "token id not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	//search bar by id
	if err := s.DB.Where("id = ?", req.GetBar().GetBarId()).First(&bar).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "invalid bar id: %v", "bar id not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	if err := s.DB.Model(&bar).Updates(updatesMap).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	responseBar := &pb.Bar{
		TokenId:  uint64(token.ID),
		BarId:    uint64(bar.ID),
		Name:     bar.Name,
		MaxValue: bar.MaxValue,
		Value:    bar.Value,
		Color:    bar.Color,
	}

	return &pb.EditBarResponse{
		Bar: responseBar,
	}, nil
}
func (s *BarService) DeleteBar(ctx context.Context, req *pb.DeleteBarRequest) (*pb.DeleteBarResponse, error) {
	//Verify if the Ids are valid
	if req.BarId <= 0 && req.TokenId <= 0 {
		return &pb.DeleteBarResponse{
			Success:       false,
			Message:       "bar id or token id  must be greater than zero",
			MessageStatus: strconv.Itoa(http.StatusBadRequest),
		}, status.Errorf(codes.NotFound, "invalid request body: %v", req.TokenId)
	}

	var bar models.Bar

	//Search first bar with inputted id
	if err := s.DB.Where("id = ? AND token_id = ?", req.BarId, req.TokenId).First(&bar).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.DeleteBarResponse{
				Success:       false,
				Message:       "token not found in the table",
				MessageStatus: strconv.Itoa(http.StatusNotFound),
			}, status.Errorf(codes.NotFound, "bar not found in the table: %v", req.TokenId)
		}
		return nil, err
	}
	//Delete bar with inputted id
	if err := s.DB.Delete(&bar).Error; err != nil {
		return &pb.DeleteBarResponse{
			MessageStatus: "500",
			Message:       "cannot delete token",
			Success:       false,
		}, status.Errorf(codes.Internal, "cannot delete token: %v", err)
	}

	//sucess response(no content)
	return &pb.DeleteBarResponse{
		MessageStatus: "204",
		Message:       "no content",
		Success:       true,
	}, nil
}
func (s *BarService) GetBarsForToken(ctx context.Context, req *pb.GetBarForTokenRequest) (*pb.GetBarsForTokenResponse, error) {

	if req.TokenId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request body: %v", req.TokenId)
	}

	var bars []models.Bar

	if err := s.DB.Where("token_id= ?", req.GetTokenId()).Find(&bars).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid token id: %v", err.Error())
	}

	responseBar := make([]*pb.Bar, 0, len(bars))

	for _, bar := range bars {
		responseBar = append(responseBar, &pb.Bar{
			BarId:    uint64(bar.ID),
			TokenId:  uint64(bar.TokenID),
			Name:     bar.Name,
			Color:    bar.Color,
			Value:    bar.Value,
			MaxValue: bar.MaxValue,
		})
	}

	return &pb.GetBarsForTokenResponse{
		Bars: responseBar,
	}, nil
}
