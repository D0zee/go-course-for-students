package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"homework8/internal/app"
)

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, a app.App) Server {
	//gin.SetMode(gin.ReleaseMode)
	s := Server{port: port, app: gin.New()}
	s.app.POST("/api/v1/ads", createAd(a))
	//s.app.GET("/api/v1/ads", showListAds(a))
	s.app.PUT("/api/v1/ads/:id/status", changeAdStatus(a))
	s.app.PUT("/api/v1/ads/:id", updateAd(a))
	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
