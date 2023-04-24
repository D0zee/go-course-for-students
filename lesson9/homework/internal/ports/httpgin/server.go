package httpgin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"homework9/internal/app"
)

type ServerWithCtx struct {
	ctx    context.Context
	Server *http.Server
}

func (s *ServerWithCtx) Run() error {
	log.Printf("starting http server, listening on %s\n", s.Server.Addr)
	defer log.Printf("close http server listening on %s\n", s.Server.Addr)

	errCh := make(chan error)

	defer func() {
		shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.Server.Shutdown(shCtx); err != nil {
			log.Printf("can't close http server listening on %s: %s", s.Server.Addr, err.Error())
		}

		close(errCh)
	}()

	go func() {
		if err := s.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	case err := <-errCh:
		return fmt.Errorf("http server can't listen and serve requests: %w", err)
	}
}

func NewHTTPServer(ctx context.Context, port string, a app.App) *ServerWithCtx {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	handler.Use(gin.Logger())
	handler.Use(Recovery)

	adsRoute := handler.Group("/api/v1/ads")
	adRouter(adsRoute, a)

	userRoute := handler.Group("api/v1/users")
	userRouter(userRoute, a)

	return &ServerWithCtx{ctx, &http.Server{Addr: port, Handler: handler}}
}
