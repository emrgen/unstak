package x

import "context"

func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "userID", userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}

func ContextWithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, "token", token)
}

func TokenFromContext(ctx context.Context) (string, bool) {
	return ctx.Value("token").(string), true
}
