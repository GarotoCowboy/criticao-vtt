package utils

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

//func NewContextWithUserID(ctx context.Context, userID uint) context.Context {
//	return context.WithValue(ctx, userIDKey, userID)
//}

func PickUserIdJWT(ctx context.Context) (uint, error) {

	userIDFromCtx := ctx.Value("user_id")
	userID, ok := userIDFromCtx.(uint)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, "tableUser not authenticated")
	}
	return userID, nil
}
