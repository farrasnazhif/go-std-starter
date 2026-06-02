package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/farrasnazhif/go-std-starter/internal/lib"
	"github.com/farrasnazhif/go-std-starter/internal/store/models"
	"github.com/farrasnazhif/go-std-starter/internal/store/repositories"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := lib.JSONResponse(w, http.StatusOK, "User retrieved successfully", user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		user, err := app.store.Users.GetByID(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, repositories.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *models.User {
	user, _ := r.Context().Value(userCtx).(*models.User)
	return user
}
