package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sapienza/azzurra/wasaphoto/service/api/reqcontext"
	"sapienza/azzurra/wasaphoto/service/database"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userid, err := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	bannedid, err := strconv.ParseUint(ps.ByName("userBanId"), 10, 64)
	if err != nil {
		// The value was not uint64, reject the action indicating an error on the client side.
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = rt.db.BanUser(userid, bannedid)
	if errors.Is(err, database.ErrUserNotExists) {
		// The fountain (indicated by `id`) does not exist, reject the action indicating an error on the client side.
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id, err := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	if err != nil {
		jsonResp, _ := json.Marshal(ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "Error Parsing JSON",
		})
		ctx.Logger.WithError(err).Error("user: error parsing JSON")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
		return
	}
	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		jsonResp, _ := json.Marshal(ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "Error decoding JSON",
		})
		ctx.Logger.WithError(err).Error("user: error decoding JSON")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
		return
	} else if !updatedUser.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	updatedUser.ID = id
	dbuser, err := rt.db.UpdateUser(updatedUser.ToDatabase())
	if errors.Is(err, database.ErrUserNotExists) {
		jsonResp, _ := json.Marshal(ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "The user not exists",
		})
		ctx.Logger.WithError(err).Error("user: error user not exists")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonResp)
		return
	} else if err != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		// Note (2): we are adding the error and an additional field (`id`) to the log entry, so that we will receive
		// the identifier of the fountain that triggered the error.
		ctx.Logger.WithError(err).WithField("id", id).Error("can't update the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	updatedUser.FromDatabase(dbuser)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(updatedUser)
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
