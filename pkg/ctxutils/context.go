package ctxutils

import "context"

type userDataKey struct{}
type UserData struct {
	UserId int32
	Email  string
	Role   string
}

func CtxWithUserData(ctx context.Context, data UserData) context.Context {
	return context.WithValue(ctx, userDataKey{}, data)
}

func GetUserDataFromCtx(ctx context.Context) UserData {
	return ctx.Value(userDataKey{}).(UserData)
}
