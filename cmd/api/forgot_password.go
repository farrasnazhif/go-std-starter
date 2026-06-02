package main

import (
	"net/http"

	"github.com/farrasnazhif/go-std-starter/internal/dto"
	"github.com/farrasnazhif/go-std-starter/internal/lib"
	"github.com/farrasnazhif/go-std-starter/internal/service"
	"github.com/farrasnazhif/go-std-starter/internal/store/repositories"
)

// Step 1: Send OTP for password reset
func (app *application) forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var payload dto.SendOTPPayload
	if err := lib.ReadJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := lib.Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if _, err := app.store.Users.GetByEmail(r.Context(), payload.Email); err != nil {
		switch err {
		case repositories.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.otpService.Send(r.Context(), payload.Email, service.PurposePasswordReset); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	lib.JSONResponse(w, http.StatusOK, "OTP sent to your email", nil)
}

// Step 2: Verify OTP, return reset token
func (app *application) verifyForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var payload dto.VerifyOTPPayload
	if err := lib.ReadJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := lib.Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	token, err := app.otpService.VerifyForReset(r.Context(), payload.Email, payload.Code)
	if err != nil {
		switch err {
		case service.ErrInvalidOTP:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	lib.JSONResponse(w, http.StatusOK, "OTP verified", map[string]string{"reset_token": token})
}

// Step 3: Reset password using reset token
func (app *application) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var payload dto.ResetPasswordPayload
	if err := lib.ReadJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := lib.Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.otpService.ResetPassword(r.Context(), payload.Token, payload.NewPassword); err != nil {
		switch err {
		case service.ErrInvalidResetToken:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	lib.JSONResponse(w, http.StatusOK, "Password reset successfully", nil)
}
