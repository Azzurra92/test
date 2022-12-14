package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"sapienza/azzurra/wasaphoto/service/api/reqcontext"

	"sapienza/azzurra/wasaphoto/service/database"

	"github.com/julienschmidt/httprouter"
)

// create identifier for user login.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "Error decoding JSON",
		}
		ctx.Logger.WithError(err).Error("user: error decoding JSON")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	} else if !user.IsValid() {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "Error validating JSON",
		}
		ctx.Logger.Error("user: error validating JSON")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	// Create user in the DB
	dbuser, err := rt.db.CreateUser(user.ToDatabase())
	if errors.Is(err, database.ErrUserExists) {
		resp := ApiResponse{
			Code:    http.StatusConflict,
			Message: "user: The user alredy exists",
		}
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(resp)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("user: error creating user in DB")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.FromDatabase(dbuser)
	// Send the response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}
