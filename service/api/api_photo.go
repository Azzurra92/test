package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"sapienza/azzurra/wasaphoto/service/api/reqcontext"
	"sapienza/azzurra/wasaphoto/service/database"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	photoId, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		resp := ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: "error parsing JSON",
		}
		ctx.Logger.WithError(err).Error("photo: Error parsing photoId")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	var req CommentRequest
	errR := json.NewDecoder(r.Body).Decode(&req)
	if errR != nil {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "Error decoding JSON",
		}
		ctx.Logger.WithError(errR).Error("comment: error decoding JSON")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	dbcomment, err := rt.db.CommentPhoto(userId, photoId, req.ToDatabase())
	if err != nil {
		resp := ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Constraints Failed",
		}
		ctx.Logger.WithError(errR).Error("comment: error retrieving data")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	u := User{
		Username: dbcomment.User.Username,
	}
	cr := CommentResponse{
		Id:       dbcomment.Id,
		From:     &u,
		Comment:  dbcomment.Comment,
		Datetime: dbcomment.Datetime,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(cr)
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId, err := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error parsing userId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	photoId, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error parsing photoId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbPhoto, err := rt.db.GetPhoto(userId, photoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: photo does not exist in db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := rt.db.DeletePhoto(userId, photoId); err != nil {
		ctx.Logger.WithError(err).Error("photo: Error deleting photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// deleting from filesystem
	photoPath := fmt.Sprintf("%s/%d/%s", rt.imagesFolder, userId, dbPhoto.PhotoUrl)
	if err := os.Remove(photoPath); err != nil {
		ctx.Logger.WithError(err).Warn("photo: could not remove photo from folder")
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId, err := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error parsing userId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	photoId, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error parsing photoId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	errD := rt.db.LikePhoto(userId, photoId)
	if errors.Is(errD, database.ErrLikesExists) {
		resp := ApiResponse{
			Code:    http.StatusConflict,
			Message: "The user has already liked",
		}
		ctx.Logger.WithError(errD).Error("like: Error The user has already liked")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(resp)
		return
	} else if errD != nil {
		resp := ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Constraint failed",
		}
		ctx.Logger.WithError(errD).Error("like: Error constraint failed")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId, err := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error parsing userId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	photoId, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error parsing photoId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commentId, err := strconv.ParseUint(ps.ByName("commentId"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error parsing photoId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := rt.db.DeleteComment(commentId, userId, photoId); err != nil {
		ctx.Logger.WithError(err).Error("photo: Error constraint failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId, err := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error parsing userId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	photoId, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error parsing photoId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := rt.db.DeleteLike(userId, photoId); err != nil {
		ctx.Logger.WithError(err).Error("photo: Error constraint failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		resp := ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "photo: Bad Request",
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error Retrieving the File")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	userid, err := strconv.ParseUint(ps.ByName("userId"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error parsing userId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		ctx.Logger.WithError(err).Error("photo: Error Creating the UUID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	imgid := uuid.String()

	// check user dir existence
	userDir := fmt.Sprintf("%s/%d", rt.imagesFolder, userid)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		if err := os.Mkdir(userDir, os.ModePerm); err != nil {
			ctx.Logger.WithError(err).Error("photo: .errors.upload_image.cannot_create_local_file")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
	}
	uniqueName := fmt.Sprintf("%s-%s", imgid, handler.Filename)
	tmpFileName := fmt.Sprintf("%s/%s", userDir, uniqueName)
	tmpFile, err := os.Create(tmpFileName)

	if err != nil {
		ctx.Logger.WithError(err).Error("photo: .errors.upload_image.cannot_create_local_file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()
	if _, err := io.Copy(tmpFile, file); err != nil {
		ctx.Logger.WithError(err).Error(".errors.upload_image.cannot_copy_to_file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p := Photo{
		Datetime: time.Now(),
		UUID:     imgid,
		UserId:   userid,
		Likes:    0,
		PhotoUrl: tmpFileName,
	}

	createdPhoto, err := rt.db.CreatePhoto(p.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error(".errors.upload_image.cannot_save_to_db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p.FromDatabase(createdPhoto)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(p)
}
