package main

import (
	"net/http"
	"social/internal/store"
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
