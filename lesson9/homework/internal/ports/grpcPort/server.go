package grpcPort

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"homework9/internal/app"
	"homework9/internal/ports/grpcPort/proto"
	"homework9/internal/ports/grpcPort/service"
	"log"
	"net"
)

type GrpcServer struct {
	ctx    context.Context
	Addr   string
	Lis    net.Listener
	Server *grpc.Server
}

func NewGrpcServer(ctx context.Context, port string, app app.App) *GrpcServer {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serviceI := service.NewService(app)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(LoggerInterceptor),
		grpc.ChainUnaryInterceptor(PanicRecovery))
	proto.RegisterAdServiceServer(grpcServer, serviceI)

	return &GrpcServer{
		ctx:    ctx,
		Addr:   port,
		Lis:    lis,
		Server: grpcServer,
	}
}

func (s *GrpcServer) Run() error {
	log.Printf("starting grpc server, listening on %s\n", s.Addr)
	defer log.Printf("close grpc server listening on %s\n", s.Addr)

	errCh := make(chan error)

	defer func() {
		s.Server.GracefulStop()
		_ = s.Lis.Close()

		close(errCh)
	}()

	go func() {
		if err := s.Server.Serve(s.Lis); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	case err := <-errCh:
		return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
	}
}

func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("[GRPC] invoke method with name:", info.FullMethod)
	return handler(ctx, req)
}

func PanicRecovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("[GRPC] Was catch panic in method:" + info.FullMethod)
		}
	}()
	return handler(ctx, req)
}
