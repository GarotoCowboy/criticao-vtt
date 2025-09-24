package token

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/GarotoCowboy/vttProject/api/grpc/proto/token/pb"
	"github.com/GarotoCowboy/vttProject/api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *TokenService) CreateToken(ctx context.Context, req *pb.CreateTokenRequest) (*pb.CreateTokenResponse, error) {

	//validate the codes
	if err := Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request body: %v", err.Error())
	}

	var table = models.Table{}

	//search the first table with inputted id
	if err := s.DB.Where("id = ?", req.GetTableId()).First(&table).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid table id: %v", err.Error())
	}

	//create a token var
	var token = models.Token{
		Name:     req.GetName(),
		ImageURL: req.GetImageUrl(),
		Bars:     nil,
		TableID:  table.ID,
	}

	//create a token in DB
	if err := s.DB.Create(&token).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	s.Logger.InfoF("Token created: %v", token)

	//token response
	responseToken := &pb.Token{
		TokenId:  uint64(token.ID),
		TableId:  uint64(token.TableID),
		Name:     token.Name,
		ImageUrl: token.ImageURL,
	}

	//if successful we return a response created before
	return &pb.CreateTokenResponse{
		Token: responseToken,
	}, nil

}

func (s *TokenService) EditToken(ctx context.Context, req *pb.EditTokenRequest) (*pb.EditTokenResponse, error) {

	//Verify if the mask and values are valid
	updatesMap, err := ValidadeAndBuildUpdateMap(req)
	if err != nil {
		return nil, err
	}

	var token models.Token
	var table models.Table

	if err := s.DB.Where("id = ?", req.GetToken().GetTableId()).First(&table).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid table id: %v", "table id not found")
	}

	//search token by id
	if err := s.DB.Where("id = ?", req.GetToken().GetTokenId()).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "invalid token id: %v", "token id not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	if err := s.DB.Model(&token).Updates(updatesMap).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	responseToken := &pb.Token{
		TokenId:  uint64(token.ID),
		TableId:  uint64(token.TableID),
		Name:     token.Name,
		ImageUrl: token.ImageURL,
	}

	return &pb.EditTokenResponse{
		Token: responseToken,
	}, nil
}
func (s *TokenService) DeleteToken(ctx context.Context, req *pb.DeleteTokenRequest) (*pb.DeleteTokenResponse, error) {

	//Verify if the Ids are valid
	if req.TokenId <= 0 && req.TableId <= 0 {
		return &pb.DeleteTokenResponse{
			Success:       false,
			Message:       "TokenId or TableId must be greater than zero",
			MessageStatus: strconv.Itoa(http.StatusBadRequest),
		}, status.Errorf(codes.NotFound, "invalid request body: %v", req.TokenId)
	}

	var token = models.Token{}

	//Search first token with inputted id
	if err := s.DB.Where("id = ? AND table_id = ?", req.TokenId, req.TableId).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.DeleteTokenResponse{
				Success:       false,
				Message:       "token not found in the table",
				MessageStatus: strconv.Itoa(http.StatusNotFound),
			}, status.Errorf(codes.NotFound, "token not found in the table: %v", req.TokenId)
		}
		return nil, err
	}
	//Delete token with inputted id
	if err := s.DB.Delete(&token).Error; err != nil {
		return &pb.DeleteTokenResponse{
			MessageStatus: "500",
			Message:       "cannot delete token",
			Success:       false,
		}, status.Errorf(codes.Internal, "cannot delete token: %v", err)
	}

	//sucess response(no content)
	return &pb.DeleteTokenResponse{
		MessageStatus: "204",
		Message:       "no content",
	}, nil
}
func (s *TokenService) ListAllTokenInTable(ctx context.Context, req *pb.ListTokenRequest) (*pb.ListTokenResponse, error) {

	if req.TableId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request body: %v", req.TableId)
	}

	var tokens []models.Token

	if err := s.DB.Where("table_id = ?", req.GetTableId()).Find(&tokens).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid table id: %v", err.Error())
	}

	responseToken := make([]*pb.Token, 0, len(tokens))

	for _, token := range tokens {
		responseToken = append(responseToken, &pb.Token{
			TokenId:  uint64(token.ID),
			TableId:  uint64(token.TableID),
			Name:     token.Name,
			ImageUrl: token.ImageURL,
		})
	}

	return &pb.ListTokenResponse{
		Token: responseToken,
	}, nil
}
