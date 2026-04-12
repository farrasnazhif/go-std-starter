package main

import (
	"net/http"
)

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Check if the API is running and healthy
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	error
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, "Fetched successfully", data); err != nil {
		app.internalServerError(w, r, err)
	}
}
