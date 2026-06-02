package main

import (
	"net/http"

	"github.com/farrasnazhif/go-std-starter/internal/dto"
	"github.com/farrasnazhif/go-std-starter/internal/lib"
	"github.com/farrasnazhif/go-std-starter/internal/store/repositories"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterUserPayload
	if err := lib.ReadJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := lib.Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	result, err := app.authService.Register(r.Context(), payload.Username, payload.Email, payload.Password)
	if err != nil {
		switch err {
		case repositories.ErrDuplicateEmail, repositories.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := lib.JSONResponse(w, http.StatusCreated, "User registered. Please verify your email with the OTP sent.", result.User); err != nil {
		app.internalServerError(w, r, err)
	}
}
