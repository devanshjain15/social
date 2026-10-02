package main

import (
	"net/http"
	"social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Get post from context (set by postsContextMiddleware)
	post := getPostFromCtx(r)

	// Get current user from auth middleware
	user := r.Context().Value(userCtxKey).(*store.User)

	var payload CreateCommentPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	comment := &store.Comment{
		Content: payload.Content,
		PostId:  post.Id,
		UserId:  user.Id,
		User: store.User{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	ctx := r.Context()
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
	}
}
