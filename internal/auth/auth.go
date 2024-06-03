package auth

import (
	"context"
	"encoding/json"
	"github.com/wlcmtunknwndth/test_ozon/lib/httpResponse"
	"github.com/wlcmtunknwndth/test_ozon/lib/slogAttr"

	"github.com/wlcmtunknwndth/test_ozon/graph/model"
	"log/slog"
	"net/http"
	"time"
)

const (
	badRequest          = "Bad request"
	internalServerError = "Internal server error"
	unauthorized        = "Unauthorized"
	authorized          = "Authorized"
	noEnoughPrivileges  = "Not enough privileges"
)

//go:generate mockery --name Storage

type Storage interface {
	GetPassword(context.Context, string) (string, error)
	CreateUser(context.Context, *model.NewUser) error
	IsAdmin(context.Context, string) (bool, error)
	DeleteUser(context.Context, string) error
}

type Auth struct {
	Db Storage
}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	const op = "auth.auth.Register"
	//corsSkip.EnableCors(w, r)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var usr model.NewUser

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		slog.Error("couldn't decode request: ", slogAttr.SlogErr(op, err))
		httpResponse.Write(w, http.StatusBadRequest, badRequest)
		return
	}

	if err = a.Db.CreateUser(ctx, &usr); err != nil {
		slog.Error("couldn't register user: ", slogAttr.SlogErr(op, err))
		httpResponse.Write(w, http.StatusInternalServerError, internalServerError)
		return
	}

	WriteNewToken(w, model.User{
		Username: usr.Username,
		Password: "",
		IsAdmin:  false,
	})

	httpResponse.Write(w, http.StatusOK, authorized)
	return
}

func (a *Auth) LogIn(w http.ResponseWriter, r *http.Request) {
	const op = "auth.auth.Login"
	//corsSkip.EnableCors(w, r)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	var usr model.User

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		httpResponse.Write(w, http.StatusBadRequest, badRequest)
		slog.Error("couldn't decode user: ", slogAttr.SlogErr(op, err))
		return
	}

	if len(usr.Password) < 4 || len(usr.Username) < 4 {
		httpResponse.Write(w, http.StatusUnauthorized, unauthorized)
		return
	}

	pass, err := a.Db.GetPassword(ctx, usr.Username)
	if err != nil || pass != usr.Password {
		slog.Error("couldn't get password from storage: ", slogAttr.SlogErr(op, err))
		httpResponse.Write(w, http.StatusUnauthorized, unauthorized)
		return
	}

	if usr.IsAdmin, err = a.Db.IsAdmin(ctx, usr.Username); err != nil {
		slog.Error("couldn't determine if user is admin: ", slogAttr.SlogErr(op, err))
		httpResponse.Write(w, http.StatusInternalServerError, internalServerError)
	}

	WriteNewToken(w, usr)
	httpResponse.Write(w, http.StatusOK, authorized)
	return
}

func (a *Auth) LogOut(w http.ResponseWriter, r *http.Request) {
	//corsSkip.EnableCors(w, r)
	http.SetCookie(w, &http.Cookie{
		Name:    access,
		Expires: time.Now(),
	})
	return
}

func (a *Auth) DeleteUser(w http.ResponseWriter, r *http.Request) {
	const op = "auth.auth.DeleteUser"
	//corsSkip.EnableCors(w, r)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if ok, err := IsAdmin(r); !ok {
		if err != nil {
			httpResponse.Write(w, http.StatusInternalServerError, internalServerError)
			slog.Error("couldn't determine if user is admin: ", slogAttr.SlogErr(op, err))
			return
		}
		httpResponse.Write(w, http.StatusForbidden, noEnoughPrivileges)
		return
	}

	type query struct {
		Username string `json:"username"`
	}

	var qry query
	err := json.NewDecoder(r.Body).Decode(&qry)
	if err != nil {
		httpResponse.Write(w, http.StatusBadRequest, badRequest)
		slog.Error("couldn't decode user: ", slogAttr.SlogErr(op, err))
		return
	}

	err = a.Db.DeleteUser(ctx, qry.Username)
	if err != nil {
		httpResponse.Write(w, http.StatusInternalServerError, internalServerError)
		slog.Error("couldn't delete user: ", slogAttr.SlogErr(op, err))
		return
	}
	return
}
