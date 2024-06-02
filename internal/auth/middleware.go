package auth

import (
	"context"
	"github.com/wlcmtunknwndth/test_ozon/lib/httpResponse"
	"github.com/wlcmtunknwndth/test_ozon/lib/slogAttr"
	"log/slog"
	"net/http"
)

const (
	usernameKey = "user"
	adminKey    = "isadmin"
)

func (a *Auth) MiddlewareAuth() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "internal.auth.MiddlewareAuth"
			c, err := r.Cookie(access)
			if err != nil || c == nil {
				httpResponse.Write(w, http.StatusUnauthorized, statusUnauthorized)
				a.LogIn(w, r)
				return
			}
			info, err := checkRequest(r)
			if err != nil {
				slog.Error("couldn't validate jwt token", slogAttr.SlogErr(op, err))
				httpResponse.Write(w, http.StatusUnauthorized, unauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), usernameKey, info.Username)
			ctx = context.WithValue(ctx, adminKey, info.IsAdmin)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
