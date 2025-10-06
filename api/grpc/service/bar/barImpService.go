package bar

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/GarotoCowboy/vttProject/api/grpc/pb/bar"
	"github.com/GarotoCowboy/vttProject/api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *BarService) CreateBar(ctx context.Context, req *bar.CreateBarRequest) (*bar.CreateBarResponse, error) {
	//validate the codes
	if err := Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request body: %v", err.Error())
	}

	var token = models.Token{}

	//search the first token with inputted id
	if err := s.DB.Where("id = ?", req.GetTokenId()).First(&token).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid token id: %v", err.Error())
	}

	//create a barModel var
	var barModel = models.Bar{
		Name:     req.GetName(),
		MaxValue: req.GetMaxValue(),
		Value:    req.GetValue(),
		Color:    req.GetColor(),
		TokenID:  token.ID,
	}

	//create a barModel in DB
	if err := s.DB.Create(&barModel).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	s.Logger.InfoF("Token created: %v", barModel)

	//barModel response
	responseBar := &bar.Bar{
		Name:     barModel.Name,
		Color:    barModel.Color,
		Value:    barModel.Value,
		MaxValue: barModel.MaxValue,
		TokenId:  uint64(barModel.TokenID),
	}

	//if successful we return a response created before
	return &bar.CreateBarResponse{
		Bar: responseBar,
	}, nil

}

func (s *BarService) EditBar(ctx context.Context, req *bar.EditBarRequest) (*bar.EditBarResponse, error) {
	//Verify if the mask and values are valid
	updatesMap, err := ValidadeAndBuildUpdateMap(req)
	if err != nil {
		return nil, err
	}

	var barModel models.Bar
	var token models.Token

	if err := s.DB.Where("id = ?", req.GetBar().GetTokenId()).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "invalid token id: %v", "token id not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	//search bar by id
	if err := s.DB.Where("id = ?", req.GetBar().GetBarId()).First(&barModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "invalid bar id: %v", "bar id not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	if err := s.DB.Model(&barModel).Updates(updatesMap).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	responseBar := &bar.Bar{
		TokenId:  uint64(token.ID),
		BarId:    uint64(barModel.ID),
		Name:     barModel.Name,
		MaxValue: barModel.MaxValue,
		Value:    barModel.Value,
		Color:    barModel.Color,
	}

	return &bar.EditBarResponse{
		Bar: responseBar,
	}, nil
}
func (s *BarService) DeleteBar(ctx context.Context, req *bar.DeleteBarRequest) (*bar.DeleteBarResponse, error) {
	//Verify if the Ids are valid
	if req.BarId <= 0 && req.TokenId <= 0 {
		return &bar.DeleteBarResponse{
			Success:       false,
			Message:       "bar id or token id  must be greater than zero",
			MessageStatus: strconv.Itoa(http.StatusBadRequest),
		}, status.Errorf(codes.NotFound, "invalid request body: %v", req.TokenId)
	}

	var barModel models.Bar

	//Search first bar with inputted id
	if err := s.DB.Where("id = ? AND token_id = ?", req.BarId, req.TokenId).First(&barModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &bar.DeleteBarResponse{
				Success:       false,
				Message:       "token not found in the table",
				MessageStatus: strconv.Itoa(http.StatusNotFound),
			}, status.Errorf(codes.NotFound, "bar not found in the table: %v", req.TokenId)
		}
		return nil, err
	}
	//Delete bar with inputted id
	if err := s.DB.Delete(&barModel).Error; err != nil {
		return &bar.DeleteBarResponse{
			MessageStatus: "500",
			Message:       "cannot delete token",
			Success:       false,
		}, status.Errorf(codes.Internal, "cannot delete token: %v", err)
	}

	//sucess response(no content)
	return &bar.DeleteBarResponse{
		MessageStatus: "204",
		Message:       "no content",
		Success:       true,
	}, nil
}
func (s *BarService) GetBarsForToken(ctx context.Context, req *bar.GetBarForTokenRequest) (*bar.GetBarsForTokenResponse, error) {

	if req.TokenId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request body: %v", req.TokenId)
	}

	var bars []models.Bar

	if err := s.DB.Where("token_id= ?", req.GetTokenId()).Find(&bars).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid token id: %v", err.Error())
	}

	responseBar := make([]*bar.Bar, 0, len(bars))

	for _, barModelLop := range bars {
		responseBar = append(responseBar, &bar.Bar{
			BarId:    uint64(barModelLop.ID),
			TokenId:  uint64(barModelLop.TokenID),
			Name:     barModelLop.Name,
			Color:    barModelLop.Color,
			Value:    barModelLop.Value,
			MaxValue: barModelLop.MaxValue,
		})
	}

	return &bar.GetBarsForTokenResponse{
		Bars: responseBar,
	}, nil
}
