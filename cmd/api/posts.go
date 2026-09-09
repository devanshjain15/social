package main

import (
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type createPostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload createPostPayload
	if err := readJson(w, r, &payload); err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: validate

	newPost := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: change after auth
		UserId: 1,
	}

	// why this context?
	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, &newPost); err != nil {
		writeErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJson(w, http.StatusCreated, newPost); err != nil {
		writeErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "postId")
	postID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	ctx := r.Context()
	post, err := app.store.Posts.GetById(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeErrorJson(w, http.StatusBadRequest, err.Error())
			return
		default:
			writeErrorJson(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err := writeJson(w, http.StatusOK, post); err != nil {
		writeErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}
}
