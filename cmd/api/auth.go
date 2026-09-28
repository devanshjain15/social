package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"social/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=24"`
}

type UserWithToken struct {
	store.User
	token string
}

type CreateTokenPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var user store.User
	user.Username = payload.Username
	user.Email = payload.Email
	user.Password.Set(payload.Password)

	token := uuid.New().String()
	log.Print(token)
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	ctx := r.Context()

	if err := app.store.Users.CreateAndInvite(ctx, &user, hashedToken, app.config.mail.exp); err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.badRequestError(w, r, err)
		case errors.Is(err, store.ErrDuplicateUsername):
			app.badRequestError(w, r, err)
		default:
			app.badRequestError(w, r, err)
		}
		return
	}

	// userWithToken := UserWithToken{
	// 	User:  user,
	// 	token: token,
	// }

	// activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, token)
	// isProdEnv := app.config.env == "production"
	// vars := struct {
	// 	Username      string
	// 	ActivationURL string
	// }{
	// 	Username:      user.Username,
	// 	ActivationURL: activationURL,
	// }

	// err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	// if err != nil {
	// 	app.logger.Errorw("error sending welcome email", "error", err)
	// 	err = app.store.Users.Delete(ctx, userWithToken.Id)
	// 	if err != nil {
	// 		app.logger.Errorw("error deleting user", "error", err)
	// 	}

	// 	app.internalServerError(w, r, err)
	// 	return
	// }

	if err := jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	if err := app.store.Users.Activate(r.Context(), hashedToken); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// get the entered login and password
	var payload CreateTokenPayload
	if err := readJson(w, r, &payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// validate
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// get the user
	ctx := r.Context()
	user, err := app.store.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// generate token
	claims := jwt.RegisteredClaims{
		Issuer:    app.config.auth.tokenAuth.iss,
		Audience:  jwt.ClaimStrings{app.config.auth.tokenAuth.iss},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(app.config.auth.tokenAuth.exp)),
		Subject:   fmt.Sprint(user.Id),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	output := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	if err := jsonResponse(w, http.StatusCreated, output); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

