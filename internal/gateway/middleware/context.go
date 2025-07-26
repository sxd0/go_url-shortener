package middleware

import (
	"context"
	"errors"
)

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return 0, errors.New("user_id not found in context")
	}
	userID, ok := val.(uint)
	if !ok {
		return 0, errors.New("invalid user_id type in context")
	}
	return userID, nil
}
