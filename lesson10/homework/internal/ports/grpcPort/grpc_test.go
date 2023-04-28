package grpcPort

import (
	"github.com/stretchr/testify/assert"
	"homework9/internal/ports/grpcPort/proto"
	"testing"
)

func Test(t *testing.T) {
	client, ctx := GetMockedGrpcClient(t)
	_, err := client.CreateUser(ctx, &proto.CreateUserRequest{
		Name:  "aboba",
		Email: "bebrovna",
	})
	assert.Nil(t, err)
	//fmt.Println(user)

}
