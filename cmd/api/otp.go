package main

import (
	"net/http"

	"github.com/farrasnazhif/go-std-starter/internal/dto"
	"github.com/farrasnazhif/go-std-starter/internal/lib"
	"github.com/farrasnazhif/go-std-starter/internal/service"
	"github.com/farrasnazhif/go-std-starter/internal/store/repositories"
)

func (app *application) sendOTPHandler(w http.ResponseWriter, r *http.Request) {
	var payload dto.SendOTPPayload
	if err := lib.ReadJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := lib.Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Ensure user exists and is not already active
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case repositories.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if user.IsActive {
		app.badRequestResponse(w, r, service.ErrUserActive)
		return
	}

	if err := app.otpService.Send(r.Context(), payload.Email); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	lib.JSONResponse(w, http.StatusOK, "OTP sent successfully", nil)
}

func (app *application) verifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var payload dto.VerifyOTPPayload
	if err := lib.ReadJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := lib.Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.otpService.Verify(r.Context(), payload.Email, payload.Code); err != nil {
		switch err {
		case service.ErrInvalidOTP:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	lib.JSONResponse(w, http.StatusOK, "Account activated successfully", nil)
}
