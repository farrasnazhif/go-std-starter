package main

import (
	"net/http"

	"github.com/farrasnazhif/go-std-starter/internal/dto"
	"github.com/farrasnazhif/go-std-starter/internal/lib"
	"github.com/farrasnazhif/go-std-starter/internal/store/repositories"
)

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.RegisterUserRequest	true	"User credentials"
//	@Success		201		{object}	dto.UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterUserRequest
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

	userWithToken := dto.UserWithToken{
		User:  result.User,
		Token: result.Token,
	}

	if err := lib.JSONResponse(w, http.StatusCreated, "User registered successfully", userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}
