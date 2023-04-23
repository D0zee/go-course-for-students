package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"homework9/internal/app"
)

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{port: port, app: gin.New()}
	s.app.Use(gin.Logger())
	s.app.Use(Recovery)
	adsRoute := s.app.Group("/api/v1/ads")
	adRouter(adsRoute, a)

	userRoute := s.app.Group("api/v1/users")
	userRouter(userRoute, a)

	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
