package auth

import (
	"context"
	"github.com/wlcmtunknwndth/test_ozon/lib/slogAttr"
	"log/slog"
	"net/http"
)

type key string

const (
	extUsername     key = "user"
	extIsAdmin      key = "isadmin"
	extIsRegistered key = "isregistered"
)

func (a *Auth) MiddlewareAuth() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "internal.auth.MiddlewareAuth"

			info, err := checkRequest(r)
			if err != nil {
				slog.Error("couldn't validate jwt token", slogAttr.SlogErr(op, err))
				//httpResponse.Write(w, http.StatusUnauthorized, unauthorized)
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), extIsRegistered, false)))
				return
			}

			ctx := context.WithValue(r.Context(), extIsRegistered, true)
			ctx = context.WithValue(ctx, extIsAdmin, info.IsAdmin)
			ctx = context.WithValue(ctx, extUsername, info.Username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUsername(ctx context.Context) string {
	data, ok := ctx.Value(extUsername).(string)
	if !ok {
		return ""
	}
	return data
}

func HasAdminRights(ctx context.Context) bool {
	data, ok := ctx.Value(extIsAdmin).(bool)
	if !ok {
		return false
	}
	return data
}

func IsRegistered(ctx context.Context) bool {
	data, ok := ctx.Value(extIsRegistered).(bool)
	if !ok {
		slog.Error("fucked")
		return false
	}
	return data
}
