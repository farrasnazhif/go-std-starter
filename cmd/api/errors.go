package main

import (
	"net/http"

	"github.com/farrasnazhif/go-std-starter/internal/lib"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err)
	lib.WriteJSONError(w, http.StatusInternalServerError, "Internal server error", err.Error())
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err)
	lib.WriteJSONError(w, http.StatusBadRequest, "Validation failed", err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found", "method", r.Method, "path", r.URL.Path, "error", err)
	lib.WriteJSONError(w, http.StatusNotFound, "Resource not found", err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("conflict response", "method", r.Method, "path", r.URL.Path, "error", err)
	lib.WriteJSONError(w, http.StatusConflict, "Resource conflicted", err.Error())
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnf("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	lib.WriteJSONError(w, http.StatusTooManyRequests, "Rate limit exceeded", "you have made too many requests, please try again later")
}
