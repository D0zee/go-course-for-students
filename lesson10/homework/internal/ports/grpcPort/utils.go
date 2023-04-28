package grpcPort

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework9/internal/adapters/repo"
	"homework9/internal/app"
	"homework9/internal/ports/grpcPort/proto"
	"homework9/internal/ports/grpcPort/service"
	"homework9/internal/users"
	"net"
	"testing"
	"time"
)

func getGrpcClient(t *testing.T, a app.App) (proto.AdServiceClient, context.Context) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(LoggerInterceptor),
		grpc.ChainUnaryInterceptor(PanicRecovery))
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := service.NewService(a)
	proto.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	return proto.NewAdServiceClient(conn), ctx
}

func GetGrpcClient(t *testing.T) (proto.AdServiceClient, context.Context) {
	a := app.NewApp(repo.NewMapAdRepo(), repo.NewUserRepo())
	return getGrpcClient(t, a)
}

func GetMockedGrpcClient(t *testing.T) (proto.AdServiceClient, context.Context) {
	a := &app.AppMock{}
	a.On("CreateUser", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(users.User{}, nil).Once()
	return getGrpcClient(t, a)
}
