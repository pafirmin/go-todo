package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/pafirmin/go-todo/internal/data"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	creds := &data.Credentials{}

	err := app.readJSON(w, r, creds)
	if err != nil {
		app.badRequest(w, err.Error())
		return
	}

	u, err := app.models.Users.Authenticate(creds)
	if err != nil {
		app.unauthorized(w)
		return
	}

	exp := time.Now().Add(5 * time.Minute)
	accessToken, err := app.jwtService.Sign(u.ID, exp)
	if err != nil {
		app.serverError(w, err)
		return
	}

	exp = time.Now().Add(7 * 24 * time.Hour)
	refreshToken, err := app.models.Tokens.New(u.ID, exp, data.ScopeRefresh)
	if err != nil {
		app.serverError(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Plaintext,
		Expires:  exp,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	app.writeJSON(w, http.StatusOK, responsePayload{"access_token": accessToken, "user": u})
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if nil == err {
		err := app.models.Tokens.Delete(cookie.Value)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	cookie = &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	app.writeJSON(w, http.StatusOK, responsePayload{"message": "successfully logged out"})
}

func (app *application) logoutEverywhere(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	err := app.models.Tokens.DeleteForUser(data.ScopeRefresh, claims.UserID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	app.writeJSON(w, http.StatusOK, responsePayload{"message": "successfully logged out"})
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			app.unauthorized(w)
			return
		default:
			app.serverError(w, err)
			return
		}
	}

	u, err := app.models.Users.GetByToken(data.ScopeRefresh, cookie.Value)
	if err != nil {
		app.infoLog.Print(err)
		app.unauthorized(w)

		return
	}

	exp := time.Now().Add(5 * time.Minute)
	token, err := app.jwtService.Sign(u.ID, exp)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, responsePayload{"access_token": token, "user": u})
}

func (app *application) guestLogin(w http.ResponseWriter, r *http.Request) {
	u, err := app.models.Users.GetByEmail("guest@example.com")
	if err != nil {
		app.unauthorized(w)
		return
	}

	exp := time.Now().Add(24 * time.Hour)
	token, err := app.jwtService.Sign(u.ID, exp)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, responsePayload{"access_token": token, "user": u})
}
