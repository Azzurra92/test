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
	userid, erru := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	bannedid, errb := strconv.ParseUint(ps.ByName("userBanId"), 10, 64)
	if erru != nil || errb != nil {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "error parsing JSON",
		}
		ctx.Logger.WithError(erru).Error("user: error parsing JSON")
		ctx.Logger.WithError(errb).Error("user: error parsing JSON")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	errB := rt.db.BanUser(userid, bannedid)
	if errB != nil {
		resp := ApiResponse{
			Code:    http.StatusConflict,
			Message: "Constraint failed",
		}
		ctx.Logger.WithError(errB).Error("user: error constraint failed")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	errF := rt.db.DeleteFollowerUser(userid, bannedid)
	if errF != nil {
		resp := ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: "unfollowed user",
		}
		ctx.Logger.WithError(errF).Error("user: the banned user does not follow the user who banned him")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userid, erru := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	followingId, errb := strconv.ParseUint(ps.ByName("followingId"), 10, 64)
	if erru != nil || errb != nil {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "error parsing JSON",
		}
		ctx.Logger.WithError(erru).Error("user: error parsing JSON")
		ctx.Logger.WithError(errb).Error("user: error parsing JSON")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	err := rt.db.FollowerUser(userid, followingId)
	if err != nil {
		resp := ApiResponse{
			Code:    http.StatusConflict,
			Message: "Constraint failed",
		}
		ctx.Logger.WithError(err).Error("user: error constraint failed")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId, err := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	if err != nil {
		resp := ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: "error parsing JSON",
		}
		ctx.Logger.WithError(err).Error("stream: Error parsing userId")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	photos, err := rt.db.GetStream(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("stream: Error getting photos")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stream := Stream{
		Photos: make([]Photo, 0),
	}

	for _, p := range photos {
		apiPhoto := Photo{}
		apiPhoto.FromDatabase(p)
		stream.Photos = append(stream.Photos, apiPhoto)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(stream)
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id, err := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	if err != nil {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "error parsing JSON",
		}
		ctx.Logger.WithError(err).Error("user: error parsing JSON")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	profiledb, err := rt.db.GetUserProfile(id)
	if errors.Is(err, database.ErrUserNotExists) {
		resp := ApiResponse{
			Code:    http.StatusConflict,
			Message: "user: The user not exists",
		}
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(resp)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("Profile: Error getting profile")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u := User{
		Username: profiledb.User.Username,
	}
	profile := Profile{
		User:      &u,
		Post:      profiledb.Post,
		Follower:  profiledb.Follower,
		Following: profiledb.Following,
		Photos:    make([]Photo, 0),
	}

	for _, p := range profiledb.Photos {
		apiPhoto := Photo{}
		apiPhoto.FromDatabase(p)
		profile.Photos = append(profile.Photos, apiPhoto)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(profile)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id, err := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	if err != nil {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "Error Parsing JSON",
		}
		ctx.Logger.WithError(err).Error("user: error parsing JSON")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "Error decoding JSON",
		}
		ctx.Logger.WithError(err).Error("user: error decoding JSON")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	} else if !updatedUser.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	updatedUser.ID = id
	dbuser, err := rt.db.UpdateUser(updatedUser.ToDatabase())
	if errors.Is(err, database.ErrUserNotExists) {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "The user not exists",
		}
		ctx.Logger.WithError(err).Error("user: error user not exists")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(resp)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).WithField("id", id).Error("can't update the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	updatedUser.FromDatabase(dbuser)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(updatedUser)
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userid, erru := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	bannedid, errb := strconv.ParseUint(ps.ByName("userBanId"), 10, 64)
	if erru != nil || errb != nil {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "error parsing JSON",
		}
		ctx.Logger.WithError(erru).Error("user: error parsing JSON")
		ctx.Logger.WithError(errb).Error("user: error parsing JSON")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	if userid == bannedid {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "error parsing JSON",
		}
		ctx.Logger.WithError(erru).Error("user: bannerId cannot be equal to userid")
		ctx.Logger.WithError(errb).Error("user: bannerId cannot be equal to userid")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	err := rt.db.DeleteBan(userid, bannedid)
	if err != nil {
		resp := ApiResponse{
			Code:    http.StatusConflict,
			Message: "Constraint failed",
		}
		ctx.Logger.WithError(err).Error("user: error constraint failed")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userid, erru := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	followingId, errb := strconv.ParseUint(ps.ByName("followingId"), 10, 64)
	if erru != nil || errb != nil {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "error parsing JSON",
		}
		ctx.Logger.WithError(erru).Error("user: error parsing JSON")
		ctx.Logger.WithError(errb).Error("user: error parsing JSON")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	err := rt.db.FollowerUser(userid, followingId)
	if err != nil {
		resp := ApiResponse{
			Code:    http.StatusConflict,
			Message: "Constraint failed",
		}
		ctx.Logger.WithError(err).Error("user: error constraint failed")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
