package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"social/internal/store"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) tokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Try to get token from Authorization header first
		header := r.Header.Get("Authorization")
		if header != "" {
			parts := strings.Split(header, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// If no header token, try to get from cookie
		if token == "" {
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token = cookie.Value
			}
		}

		// If still no token, return unauthorized
		if token == "" {
			app.unauthroizedAccessError(w, r, fmt.Errorf("No Authorization header or cookie!"))
			return
		}

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthroizedAccessError(w, r, err)
			return
		}
		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userIdString, _ := claims.GetSubject()
		userId, err := strconv.ParseInt(userIdString, 10, 64)
		if err != nil {
			app.unauthroizedAccessError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.store.Users.GetById(ctx, userId)
		if err != nil {
			app.unauthroizedAccessError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) basicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			// check the credentials
			username := app.config.auth.basicAuth.username
			password := app.config.auth.basicAuth.password

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		if post.UserId == user.Id {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		allowed, err := app.checkRolePrecedence(ctx, user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		app.logger.Info(user.Id)

		if !allowed {
			app.forbiddenError(w, r, fmt.Errorf("User cannot perform this action lower role level!"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
