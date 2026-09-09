package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal Server Error Occurred => Method: %s Path: %s Error: %s", r.Method, r.URL.Path, err)
	writeErrorJson(w, http.StatusInternalServerError, "Server encountered a problem, We're fixing it!")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad Request Error Occurred => Method: %s Path: %s Error: %s", r.Method, r.URL.Path, err)
	writeErrorJson(w, http.StatusBadRequest, "Bad request by user")
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Resource Not Found Error Occurred => Method: %s Path: %s Error: %s", r.Method, r.URL.Path, err)
	writeErrorJson(w, http.StatusNotFound, "Resource that user looking for not found!")
}
