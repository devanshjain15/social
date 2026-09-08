package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "Ok",
		"env":     app.config.env,
		"version": VERSION,
	}

	if err := writeJson(w, http.StatusOK, data); err != nil {
		writeErrorJson(w, http.StatusInternalServerError, "An error occurred!")
	}
}
