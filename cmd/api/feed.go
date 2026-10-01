package main

import (
	"net/http"
	"social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	fq := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	userId := user.Id
	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(userId), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
