package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"homework9/internal/ports/grpc/proto"
	"testing"
	"time"
)

func TestGRRPCCreateUser(t *testing.T) {
	client, ctx := getGrpcClient(t)
	res, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	assert.Equal(t, "Oleg", res.Name)
}

func TestCreateAdGrpc(t *testing.T) {
	client, ctx := getGrpcClient(t)

	uResponse, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "ya@ya.ru"})
	assert.NoError(t, err)

	request := &proto.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: uResponse.Id,
	}
	response, err := client.CreateAd(ctx, request)
	assert.NoError(t, err)
	assert.Zero(t, response.Id)
	assert.Equal(t, response.Title, "hello")
	fmt.Println(response)
	assert.Equal(t, response.Text, "world")
	assert.Equal(t, response.AuthorId, int64(0))
	assert.False(t, response.Published)
	assert.True(t, isSameDate(response.CreationTime.AsTime(), time.Now()))
	assert.True(t, isSameDate(response.UpdateTime.AsTime(), time.Now()))
}
