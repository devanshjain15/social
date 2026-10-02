package main

import (
	"net/http"
	"social/internal/store"
)

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get current user from auth middleware
	currentUser := r.Context().Value(userCtxKey)
	if currentUser == nil {
		app.unauthroizedAccessError(w, r, nil)
		return
	}

	users, err := app.store.Users.GetAll(ctx)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Check follow status for each user
	for i := range users {
		if users[i].Id == currentUser.(*store.User).Id {
			continue // Skip current user
		}

		isFollowing, err := app.store.Followers.IsFollowing(ctx, currentUser.(*store.User).Id, users[i].Id)
		if err != nil {
			app.logger.Errorw("error checking follow status", "error", err)
			users[i].IsFollowing = false
		} else {
			users[i].IsFollowing = isFollowing
		}
	}

	if err := jsonResponse(w, http.StatusOK, users); err != nil {
		app.internalServerError(w, r, err)
	}
}
