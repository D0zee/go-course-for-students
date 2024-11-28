package service_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework9/internal/app"
	"homework9/internal/ports/grpcPort"
	"homework9/internal/ports/grpcPort/proto"
	"homework9/internal/tests"
	"testing"
	"time"
)

func TestGRRPCCreateUser(t *testing.T) {
	client, ctx := grpcPort.GetGrpcClient(t)
	res, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	assert.Equal(t, "Oleg", res.Name)
}

func TestCreateAdGrpc(t *testing.T) {
	client, ctx := grpcPort.GetGrpcClient(t)

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
	assert.True(t, tests.IsSameTimes(response.CreationTime.AsTime(), time.Now().UTC()))
	assert.True(t, tests.IsSameTimes(response.UpdateTime.AsTime(), time.Now().UTC()))
}

func TestChangeAdStatusGrpc(t *testing.T) {
	client, ctx := grpcPort.GetGrpcClient(t)

	user, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "ya@ya.ru"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Ivan", Email: "ya@ya.ru"})
	assert.NoError(t, err)

	request := &proto.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	}
	ad, err := client.CreateAd(ctx, request)
	assert.NoError(t, err)

	response, err := client.ChangeAdStatus(ctx, &proto.ChangeAdStatusRequest{
		AdId:      ad.Id,
		UserId:    user.Id,
		Published: true,
	})
	assert.NoError(t, err)
	assert.True(t, response.Published)
	assert.True(t, tests.IsSameTimes(response.UpdateTime.AsTime(), time.Now().UTC()))

	// error when try to change ad status from another user
	_, err = client.ChangeAdStatus(ctx, &proto.ChangeAdStatusRequest{
		AdId:      ad.Id,
		UserId:    user2.Id,
		Published: true,
	})
	assert.Error(t, err)
}

func TestUpdateAdGrpc(t *testing.T) {
	client, ctx := grpcPort.GetGrpcClient(t)

	user, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "ya@ya.ru"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Ivan", Email: "ya@ya.ru"})
	assert.NoError(t, err)

	request := &proto.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	}
	ad, err := client.CreateAd(ctx, request)
	assert.NoError(t, err)

	response, err := client.UpdateAd(ctx, &proto.UpdateAdRequest{
		AdId:   ad.Id,
		UserId: user.Id,
		Title:  "not hello",
		Text:   "not world",
	})
	assert.NoError(t, err)
	assert.Equal(t, response.Title, "not hello")
	assert.Equal(t, response.Text, "not world")
	assert.True(t, tests.IsSameTimes(response.UpdateTime.AsTime(), time.Now().UTC()))

	// error when try to change ad status from another user
	_, err = client.UpdateAd(ctx, &proto.UpdateAdRequest{
		AdId:   ad.Id,
		UserId: user2.Id,
	})
	assert.Error(t, err)
}

func TestGetAdGrpc(t *testing.T) {
	client, ctx := grpcPort.GetGrpcClient(t)

	user, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "ya@ya.ru"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Ivan", Email: "ya@ya.ru"})
	assert.NoError(t, err)

	request := &proto.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	}
	adOriginal, err := client.CreateAd(ctx, request)
	assert.NoError(t, err)

	response, err := client.GetAd(ctx, &proto.GetAdRequest{
		AdId:   adOriginal.Id,
		UserId: user.Id,
	})
	assert.NoError(t, err)
	assert.Equal(t, response, adOriginal)

	// ad with this id isn't
	_, err = client.GetAd(ctx, &proto.GetAdRequest{
		AdId:   100000,
		UserId: user.Id,
	})
	assert.Error(t, err, app.ErrAvailabilityAd)

	// access permitted
	_, err = client.GetAd(ctx, &proto.GetAdRequest{
		AdId:   adOriginal.Id,
		UserId: user2.Id,
	})
	assert.Error(t, err, app.ErrAccess)
}

func TestDeleteAdGrpc(t *testing.T) {
	client, ctx := grpcPort.GetGrpcClient(t)

	user, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "ya@ya.ru"})
	assert.NoError(t, err)

	request := &proto.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	}
	adOriginal, err := client.CreateAd(ctx, request)
	assert.NoError(t, err)

	// wrong user id
	_, err = client.DeleteAd(ctx, &proto.DeleteAdRequest{
		AdId:   adOriginal.Id,
		UserId: -1,
	})
	assert.Error(t, err, app.ErrAccess)

	// wrong ad id
	_, err = client.DeleteAd(ctx, &proto.DeleteAdRequest{
		AdId:   -1,
		UserId: user.Id,
	})
	assert.Error(t, err, app.ErrAvailabilityAd)

	// ad is available
	ad, err := client.GetAd(ctx, &proto.GetAdRequest{
		AdId:   adOriginal.Id,
		UserId: user.Id,
	})
	assert.NoError(t, err)
	assert.Equal(t, ad, adOriginal)

	adRemoved, err := client.DeleteAd(ctx, &proto.DeleteAdRequest{
		AdId:   adOriginal.Id,
		UserId: user.Id,
	})
	assert.NoError(t, err)
	assert.Equal(t, adRemoved, adOriginal)

	// get error if you want to get removed ad
	_, err = client.GetAd(ctx, &proto.GetAdRequest{
		AdId:   adOriginal.Id,
		UserId: user.Id,
	})
	assert.Error(t, err, app.ErrAvailabilityAd)

	// get error if we want to remove again
	_, err = client.DeleteAd(ctx, &proto.DeleteAdRequest{
		AdId:   adOriginal.Id,
		UserId: user.Id,
	})
	assert.Error(t, err, app.ErrAvailabilityAd)
}

func isSameAd(t *testing.T, ad1, ad2 *proto.AdResponse) {
	assert.Equal(t, ad1.Id, ad2.Id)
	assert.Equal(t, ad1.Title, ad2.Title)
	assert.Equal(t, ad1.Text, ad2.Text)
	assert.Equal(t, ad1.AuthorId, ad2.AuthorId)
}

func TestListAdsFilterGrpc(t *testing.T) {
	client, ctx := grpcPort.GetGrpcClient(t)

	user, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "ya@ya.ru"})
	assert.NoError(t, err)

	ad1, err := client.CreateAd(ctx, &proto.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	})
	assert.NoError(t, err)

	ad1, err = client.ChangeAdStatus(ctx, &proto.ChangeAdStatusRequest{
		AdId:      ad1.Id,
		UserId:    user.Id,
		Published: true,
	})
	assert.NoError(t, err)

	ad2, err := client.CreateAd(ctx, &proto.CreateAdRequest{
		Title:  "best cat",
		Text:   "not for sale",
		UserId: user.Id,
	})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "IVAN", Email: "google@ya.ru"})
	assert.NoError(t, err)

	ad3, err := client.CreateAd(ctx, &proto.CreateAdRequest{
		Title:  "ad from Ivan",
		Text:   "empty",
		UserId: user2.Id,
	})
	assert.NoError(t, err)

	ads, err := client.ListAds(ctx, &proto.ListAdRequest{
		UserId: &user.Id,
	})
	assert.NoError(t, err)
	assert.Len(t, ads.List, 2)
	isSameAd(t, ads.List[0], ad1)
	isSameAd(t, ads.List[1], ad2)

	ads, err = client.ListAds(ctx, &proto.ListAdRequest{
		Time: timestamppb.Now(),
	})
	assert.NoError(t, err)
	assert.Len(t, ads.List, 3)
	isSameAd(t, ads.List[0], ad1)
	isSameAd(t, ads.List[1], ad2)
	isSameAd(t, ads.List[2], ad3)
}

func TestAdsByTitleGrpc(t *testing.T) {
	client, ctx := grpcPort.GetGrpcClient(t)

	user, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "ya@ya.ru"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Ivan", Email: "aboba@ya.ru"})
	assert.NoError(t, err)

	ad1, err := client.CreateAd(ctx, &proto.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: user.Id,
	})
	assert.NoError(t, err)

	ad2, err := client.CreateAd(ctx, &proto.CreateAdRequest{
		Title:  "hello2",
		Text:   "hello spectators!",
		UserId: user2.Id,
	})
	assert.NoError(t, err)

	ad3, err := client.CreateAd(ctx, &proto.CreateAdRequest{
		Title:  "go is cool",
		Text:   "hello spectators!",
		UserId: user2.Id,
	})
	assert.NoError(t, err)

	helloAds, err := client.AdsByTitle(ctx, &proto.AdsByTitleRequest{Title: "hello"})
	assert.NoError(t, err)
	assert.Len(t, helloAds.List, 2)
	isSameAd(t, helloAds.List[0], ad1)
	isSameAd(t, helloAds.List[1], ad2)

	ad, err := client.AdsByTitle(ctx, &proto.AdsByTitleRequest{Title: "go"})
	assert.NoError(t, err)
	assert.Len(t, ad.List, 1)
	isSameAd(t, ad.List[0], ad3)
}

func TestChangeUserGrpc(t *testing.T) {
	client, ctx := grpcPort.GetGrpcClient(t)

	user, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "pushkin@ya.ru"})
	assert.NoError(t, err)

	newName := "Олег"
	updatedUser, err := client.ChangeUser(ctx, &proto.ChangeUserRequest{
		Id:       user.Id,
		Nickname: &newName,
		Email:    nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Id, user.Id)
	assert.Equal(t, updatedUser.Name, newName)
	assert.Equal(t, updatedUser.Email, user.Email)

	newEmail := "oao@gmail.com"
	updatedUser, err = client.ChangeUser(ctx, &proto.ChangeUserRequest{
		Id:       user.Id,
		Nickname: nil,
		Email:    &newEmail,
	})
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Id, user.Id)
	assert.Equal(t, updatedUser.Name, newName)
	assert.Equal(t, updatedUser.Email, newEmail)

	//wrong ID
	_, err = client.ChangeUser(ctx, &proto.ChangeUserRequest{
		Id:       1000,
		Nickname: nil,
		Email:    &newEmail,
	})
	assert.Error(t, err, app.ErrWrongUserId)
}

func TestGetUserGrpc(t *testing.T) {
	client, ctx := grpcPort.GetGrpcClient(t)

	user, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "pushkin@ya.ru"})
	assert.NoError(t, err)

	sameUser, err := client.GetUser(ctx, &proto.GetUserRequest{Id: user.Id})
	assert.NoError(t, err)
	assert.Equal(t, user, sameUser)

	// get non existing user
	_, err = client.GetUser(ctx, &proto.GetUserRequest{Id: 1000})
	assert.Error(t, err, app.ErrWrongUserId)
}

func TestRemoveUserGrpc(t *testing.T) {
	client, ctx := grpcPort.GetGrpcClient(t)

	user, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "pushkin@ya.ru"})
	assert.NoError(t, err)

	// remove same user
	removedUser, err := client.DeleteUser(ctx, &proto.DeleteUserRequest{Id: user.Id})
	assert.NoError(t, err)
	assert.Equal(t, user, removedUser)

	// remove same user again is impossible
	_, err = client.DeleteUser(ctx, &proto.DeleteUserRequest{Id: user.Id})
	assert.Error(t, err, app.ErrWrongUserId)

	// get same user is impossible
	_, err = client.GetUser(ctx, &proto.GetUserRequest{Id: user.Id})
	assert.Error(t, err, app.ErrWrongUserId)
}
