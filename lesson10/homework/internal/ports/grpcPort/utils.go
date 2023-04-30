package grpcPort

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework9/internal/adapters/repo"
	"homework9/internal/ads"
	"homework9/internal/app"
	"homework9/internal/mocks"
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

type Mode int64

const (
	HappyCase = iota
	ErrorCase
)

func GetMockedGrpcClient(t *testing.T, mode Mode) (proto.AdServiceClient, context.Context) {
	a := &mocks.App{}

	err := errors.New("smth")
	if mode == HappyCase {
		err = nil
	}

	a.On("CreateUser", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(users.User{}, err).Once()
	a.On("UpdateUser", mock.Anything, mock.AnythingOfType("int64"),
		mock.AnythingOfType("string"), mock.Anything).
		Return(users.User{}, err).Once()
	a.On("GetUser", mock.Anything, mock.AnythingOfType("int64")).
		Return(users.User{}, err).Once()
	a.On("RemoveUser", mock.Anything, mock.AnythingOfType("int64")).
		Return(users.User{}, err).Once()

	a.On("CreateAd", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(ads.Ad{}, err).Once()
	a.On("UpdateAd", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(ads.Ad{}, err).Once()
	a.On("ChangeAdStatus", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("bool")).
		Return(ads.Ad{}, err).Once()
	a.On("RemoveAd", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
		Return(ads.Ad{}, err).Once()
	a.On("GetAdById", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
		Return(ads.Ad{}, err).Once()
	a.On("ListAds", mock.Anything).
		Return([]ads.Ad{}, err).Once()
	return getGrpcClient(t, a)
}
