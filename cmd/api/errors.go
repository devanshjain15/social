package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal Server Error Occurred", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeErrorJson(w, http.StatusInternalServerError, "Server encountered a problem, We're fixing it!")
}

func (app *application) forbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Forbidden", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeErrorJson(w, http.StatusForbidden, "Forbidden")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Bad Request Error Occurred", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeErrorJson(w, http.StatusBadRequest, "Bad request by user")
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Resource not found Error Occurred", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeErrorJson(w, http.StatusNotFound, "Resource that user looking for not found!")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Conflict Error Occurred", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeErrorJson(w, http.StatusConflict, "Not Found!")
}

func (app *application) unauthroizedAccessError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Unauthorized Error Occurred", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeErrorJson(w, http.StatusUnauthorized, "Unauthorized")
}

func (app *application) unauthorizedBasicError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Unauthorized Error Occurred", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="Protected Area"`)
	writeErrorJson(w, http.StatusUnauthorized, "Unauthorized")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryAfter)
	writeErrorJson(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
