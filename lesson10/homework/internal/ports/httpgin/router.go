package httpgin

import (
	"github.com/gin-gonic/gin"

	"homework9/internal/app"
)

func adRouter(r *gin.RouterGroup, a app.App) {
	r.POST("/", createAd(a))
	r.PUT("/:id/status", changeAdStatus(a))
	r.PUT("/:id", updateAd(a))

	r.GET("/:id", getAd(a))
	r.GET("/", getListAds(a))

	r.GET("/title", getAdsByTitle(a))
	r.DELETE("/:id", removeAd(a))
}

func userRouter(r *gin.RouterGroup, a app.App) {
	r.POST("/", createUser(a))
	r.PUT("/:id/nickname", changeUser(a, app.ChangeNickname))
	r.PUT("/:id/email", changeUser(a, app.ChangeEmail))

	r.GET("/:id", getUser(a))
	r.DELETE("/:id", removeUser(a))

}
