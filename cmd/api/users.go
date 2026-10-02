package main

import (
	"context"
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userCtx string

const userCtxKey userCtx = "user"
const targetUserCtxKey userCtx = "targetUser"

type FollowUserPayload struct {
	UserId int64 `json:"user_id"`
}

type UnfollowUserPayload struct {
	UserId int64 `json:"user_id"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getTargetUserFromCtx(r)
	if err := jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user to follow from the target user context (set by usersContextMiddleware)
	userToFollow := getTargetUserFromCtx(r)

	// Get the current authenticated user
	currentUser := getUserFromCtx(r)

	ctx := r.Context()
	if err := app.store.Followers.Follow(ctx, currentUser.Id, userToFollow.Id); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user to unfollow from the target user context (set by usersContextMiddleware)
	userToUnfollow := getTargetUserFromCtx(r)

	// Get the current authenticated user
	currentUser := getUserFromCtx(r)

	ctx := r.Context()
	if err := app.store.Followers.Unfollow(ctx, currentUser.Id, userToUnfollow.Id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "userId")
		userId, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetById(ctx, userId)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, targetUserCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	return r.Context().Value(userCtxKey).(*store.User)
}

func getTargetUserFromCtx(r *http.Request) *store.User {
	return r.Context().Value(targetUserCtxKey).(*store.User)
}
