package auth

import (
	"context"
	"github.com/wlcmtunknwndth/test_ozon/lib/slogAttr"
	"log/slog"
	"net/http"
)

type username string
type isAdmin string
type isRegistered string

const (
	user         = username("username")
	isadmin      = isAdmin("isadmin")
	isregistered = isRegistered("isregistered")
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
				r = r.WithContext(context.WithValue(r.Context(), isregistered, false))
				return
			}
			ctx := context.WithValue(r.Context(), isregistered, true)
			ctx = context.WithValue(r.Context(), isadmin, info.IsAdmin)
			ctx = context.WithValue(r.Context(), user, info.Username)

			r = r.WithContext(ctx)
		})
	}
}

func GetUsername(ctx context.Context) string {
	data, ok := ctx.Value(user).(string)
	if !ok {
		return ""
	}
	return data
}

func HasAdminRights(ctx context.Context) bool {
	data, ok := ctx.Value(isadmin).(bool)
	if !ok {
		return false
	}
	return data
}

func IsRegistered(ctx context.Context) bool {
	data, ok := ctx.Value(isregistered).(bool)
	if !ok {
		return false
	}
	return data
}
