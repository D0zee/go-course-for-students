package httpgin

import (
	"github.com/gin-gonic/gin"

	"homework8/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.POST("/ads", createAd(a))
	//s.app.GET("/api/v1/ads", showListAds(a))
	r.PUT("/ads/:id/status", changeAdStatus(a))
	r.PUT("/ads/:id", updateAd(a))
	r.POST("/users", CreateUser(a))
	r.PUT("/users/:id/nickname", ChangeUser(a, method(changeNickname)))
	r.PUT("/users/:id/email", ChangeUser(a, method(changeEmail)))
}
