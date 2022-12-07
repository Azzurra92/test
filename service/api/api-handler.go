package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.PUT("/user/:userId/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/user/:userId/following/:followingId", rt.wrap(rt.followUser))
	rt.router.DELETE("/user/:userId/following/:followingId", rt.wrap(rt.unfollowUser))
	rt.router.PUT("/user/:userId/ban/:userBanId", rt.wrap(rt.banUser))
	rt.router.DELETE("/user/:userId/ban/:userBanId", rt.wrap(rt.unbanUser))
	rt.router.GET("/user/:userId", rt.wrap(rt.getUserProfile))
	rt.router.GET("/streams", rt.wrap(rt.getMyStream))
	rt.router.POST("/photo", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/photo/:photoId", rt.wrap(rt.deletePhoto))
	rt.router.PUT("/photo/:photoId/like/:userId", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/photo/:photoId/like/:userId", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/photo/:photoId/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/photo/:photoId/comments/:commentId", rt.wrap(rt.uncommentPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
