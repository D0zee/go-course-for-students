package tests

import (
	"github.com/stretchr/testify/assert"
	"homework9/internal/ports/grpc/proto"
	"testing"
)

func TestGRRPCCreateUser(t *testing.T) {
	client, ctx := getGrpcClient(t)
	res, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	assert.Equal(t, "Oleg", res.Name)
}
