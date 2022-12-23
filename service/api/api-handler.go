package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.PUT("/users/:userId/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/:userId/following/:followingId", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:userId/following/:followingId", rt.wrap(rt.unfollowUser))
	rt.router.PUT("/users/:userId/bans/:userBanId", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:userId/bans/:userBanId", rt.wrap(rt.unbanUser))
	rt.router.GET("/users/:userId", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users/:userId/streams", rt.wrap(rt.getMyStream))

	rt.router.POST("/users/:userId/photos", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/users/:userId/photos/:photoId", rt.wrap(rt.deletePhoto))
	rt.router.PUT("/users/:userId/photos/:photoId/likes", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:userId/photos/:photoId/likes", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/users/:userId/photos/:photoId/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/users/:userId/photos/:photoId/comments/:commentId", rt.wrap(rt.uncommentPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
