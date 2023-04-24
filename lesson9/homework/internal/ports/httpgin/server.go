package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"homework9/internal/app"
)

type Server struct {
	Addr    string
	Handler http.Handler
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	handler.Use(gin.Logger())
	handler.Use(Recovery)
	s := Server{Addr: port, Handler: handler}

	adsRoute := handler.Group("/api/v1/ads")
	adRouter(adsRoute, a)

	userRoute := handler.Group("api/v1/users")
	userRouter(userRoute, a)
	return s
}
