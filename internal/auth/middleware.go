package auth

import (
	"context"
	"github.com/wlcmtunknwndth/test_ozon/lib/slogAttr"
	"log/slog"
	"net/http"
)

const (
	usernameKey   = "user"
	adminKey      = "isadmin"
	registeredKey = "isregistered"
)

func (a *Auth) MiddlewareAuth() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "internal.auth.MiddlewareAuth"
			defer next.ServeHTTP(w, r)
			info, err := checkRequest(r)
			if err != nil {
				slog.Error("couldn't validate jwt token", slogAttr.SlogErr(op, err))
				//httpResponse.Write(w, http.StatusUnauthorized, unauthorized)
				r = r.WithContext(context.WithValue(r.Context(), registeredKey, false))
				return
			}

			ctx := context.WithValue(r.Context(), usernameKey, info.Username)
			ctx = context.WithValue(ctx, registeredKey, true)
			ctx = context.WithValue(ctx, adminKey, info.IsAdmin)
			r = r.WithContext(ctx)
		})
	}
}

func GetUsername(ctx context.Context) string {
	username, ok := ctx.Value(usernameKey).(string)
	if !ok {
		return ""
	}
	return username
}

func HasAdminRights(ctx context.Context) bool {
	ans, ok := ctx.Value(adminKey).(bool)
	if !ok {
		return false
	}
	return ans
}

func IsRegistered(ctx context.Context) bool {
	ans, ok := ctx.Value(registeredKey).(bool)
	if !ok {
		return false
	}
	return ans
}
