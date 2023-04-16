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
	//gin.SetMode(gin.DebugMode)
	s := Server{port: port, app: gin.New()}
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
