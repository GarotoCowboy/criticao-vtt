package token

import (
	"context"
	"errors"

	"github.com/GarotoCowboy/vttProject/api/grpc/events"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/token"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (s *TokenService) CreateToken(ctx context.Context, req *token.CreateTokenRequest) (*token.CreateTokenResponse, error) {

	//validate the codes
	if err := Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request body: %v", err.Error())
	}

	var tableModel = models.Table{}

	//create a token var

	var tokenModel = models.Token{}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		//search the first tableModel with inputted id
		if err := tx.WithContext(ctx).Where("id = ?", req.GetTableId()).First(&tableModel).Error; err != nil {
			return status.Errorf(codes.NotFound, "invalid tableModel id: %v", err.Error())
		}

		tokenModel.Name = req.GetName()
		tokenModel.ImageURL = req.GetImageUrl()
		tokenModel.TableID = uint(req.TableId)

		//create a token in DB
		if err := tx.WithContext(ctx).Create(&tokenModel).Error; err != nil {
			return status.Errorf(codes.Internal, "internal error: %v", err.Error())
		}

		s.Logger.InfoF("Token created: %v", tokenModel)

		return nil
	})
	if err != nil {
		s.Logger.WarningF("Error creating token: %v", err.Error())
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	//token response
	responseToken := &token.Token{
		TokenId:  uint64(tokenModel.ID),
		TableId:  uint64(tokenModel.TableID),
		Name:     tokenModel.Name,
		ImageUrl: tokenModel.ImageURL,
	}

	s.Logger.InfoF("synchronizing this new event")
	event := events.NewCreateTokenEvent(responseToken)

	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(tokenModel.TableID), event)

	s.Logger.InfoF("token created with sucess")
	//if successful we return a response created before
	return &token.CreateTokenResponse{
		Token: responseToken,
	}, nil

}

func (s *TokenService) EditToken(ctx context.Context, req *token.EditTokenRequest) (*token.EditTokenResponse, error) {
	s.Logger.InfoF("GRPC requisition to EditToken")
	//Verify if the mask and values are valid
	updatesMap, err := ValidadeAndBuildUpdateMap(req)
	if err != nil {
		return nil, err
	}

	var tokenModel models.Token
	var tableModel models.Table

	s.Logger.InfoF("Checking if token and table exists")
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.DB.WithContext(ctx).Where("id = ?", req.GetToken().GetTableId()).First(&tableModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.WarningF("Table %s not found", req.GetToken().GetTableId())
				return status.Errorf(codes.NotFound, "invalid tableModel id: %v", "tableModel id not found")
			}
			s.Logger.InfoF("Internal error")
			return status.Errorf(codes.Internal, "internal error: %v", err.Error())
		}

		//search tokenModel by id
		if err := s.DB.WithContext(ctx).Where("id = ?", req.GetToken().GetTokenId()).First(&tokenModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.WarningF("Token %s not found", req.GetToken().GetTokenId())
				return status.Errorf(codes.NotFound, "invalid tokenModel id: %v", "tokenModel id not found")
			}
			return status.Errorf(codes.Internal, "internal error: %v", err.Error())
		}

		s.Logger.InfoF("editing token")
		if err := s.DB.WithContext(ctx).Model(&tokenModel).Updates(updatesMap).Error; err != nil {
			return status.Errorf(codes.Internal, "internal error: %v", err.Error())
		}

		return nil
	})
	if err != nil {
		s.Logger.WarningF("Error editing token: %v", err.Error())
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	responseToken := &token.Token{
		TokenId:  uint64(tokenModel.ID),
		TableId:  uint64(tokenModel.TableID),
		Name:     tokenModel.Name,
		ImageUrl: tokenModel.ImageURL,
	}

	s.Logger.InfoF("synchronizing this new event")
	event := events.NewUpdatedTokenEvent(responseToken)
	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(tokenModel.TableID), event)

	s.Logger.InfoF("token updated with sucess")

	return &token.EditTokenResponse{
		Token: responseToken,
	}, nil
}
func (s *TokenService) DeleteToken(ctx context.Context, req *token.DeleteTokenRequest) (*token.DeleteTokenResponse, error) {

	//Verify if the Ids are valid
	if req.TokenId <= 0 && req.TableId <= 0 {
		return nil, status.Errorf(codes.NotFound, "invalid request body: %v", req.TokenId)
	}

	var tokenModel = models.Token{}

	result := s.DB.WithContext(ctx).Where("id = ? AND table_id = ?", req.TokenId, req.TableId).Delete(&tokenModel)

	if result.Error != nil {
		s.Logger.WarningF("Token %s or table_id %s not found", req.TokenId, req.TableId)
		return nil, status.Errorf(codes.NotFound, "invalid request body: %v", req.TokenId)
	}

	if result.RowsAffected == 0 {
		s.Logger.WarningF("No token found to delete with id %d and table_id %d", req.TokenId, req.TableId)
		return nil, status.Errorf(codes.NotFound, "token not found or does not belong to the specified table")
	}

	s.Logger.InfoF("Sycronizing this new Event")
	event := events.NewDeleteTokenEvent(req.TableId, req.TokenId)
	s.Broker.Publish(pubSubSyncConst.TableSync, req.TableId, event)

	s.Logger.InfoF("GRPC Requisition to Delete Scene finished...")

	//sucess response(no content)
	return &token.DeleteTokenResponse{
		Empty: &emptypb.Empty{},
	}, nil
}
func (s *TokenService) ListAllTokenInTable(ctx context.Context, req *token.ListTokenRequest) (*token.ListTokenResponse, error) {

	if req.TableId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request body: %v", req.TableId)
	}

	var tokens []models.Token

	if err := s.DB.WithContext(ctx).Where("table_id = ?", req.GetTableId()).Find(&tokens).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid table id: %v", err.Error())
	}

	responseToken := make([]*token.Token, 0, len(tokens))

	for _, tokenLop := range tokens {
		responseToken = append(responseToken, &token.Token{
			TokenId:  uint64(tokenLop.ID),
			TableId:  uint64(tokenLop.TableID),
			Name:     tokenLop.Name,
			ImageUrl: tokenLop.ImageURL,
		})
	}

	return &token.ListTokenResponse{
		Tokens: responseToken,
	}, nil
}
