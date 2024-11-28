package service_test

import (
	"context"
	"github.com/stretchr/testify/suite"
	"homework9/internal/ports/grpcPort"
	"homework9/internal/ports/grpcPort/proto"
	"testing"
)

type UserSuit struct {
	suite.Suite
	t               *testing.T
	client          proto.AdServiceClient
	clientWithError proto.AdServiceClient
	ctx             context.Context
}

func (s *UserSuit) SetupTest() {
	s.ctx = context.Background()
	s.t = &testing.T{}
	s.client, _ = grpcPort.GetMockedGrpcClient(s.t, grpcPort.HappyCase)
	s.clientWithError, _ = grpcPort.GetMockedGrpcClient(s.t, grpcPort.ErrorCase)
}

func (s *UserSuit) TestUserCreate() {
	req := &proto.CreateUserRequest{
		Name:  "aboba",
		Email: "bebrovna",
	}

	_, err := s.client.CreateUser(s.ctx, req)
	s.Nil(err)

	_, err = s.clientWithError.CreateUser(s.ctx, req)
	s.NotNil(err)
}

func (s *UserSuit) TestUserUpdate() {
	data := "Ivan"
	req := &proto.ChangeUserRequest{
		Id:       0,
		Nickname: &data,
		Email:    &data,
	}

	_, err := s.client.ChangeUser(s.ctx, req)
	s.Nil(err)

	emptyReq := &proto.ChangeUserRequest{}
	_, err = s.client.ChangeUser(s.ctx, emptyReq)
	s.NotNil(err)

	_, err = s.clientWithError.ChangeUser(s.ctx, req)
	s.NotNil(err)
}

func (s *UserSuit) TestGetUser() {
	req := &proto.GetUserRequest{Id: 0}

	_, err := s.client.GetUser(s.ctx, req)
	s.Nil(err)

	_, err = s.clientWithError.GetUser(s.ctx, req)
	s.NotNil(err)
}

func (s *UserSuit) TestRemoveUser() {
	req := &proto.DeleteUserRequest{Id: 0}

	_, err := s.client.DeleteUser(s.ctx, req)
	s.Nil(err)

	_, err = s.clientWithError.DeleteUser(s.ctx, req)
	s.NotNil(err)
}

func (s *UserSuit) TestCreateAd() {
	req := &proto.CreateAdRequest{}

	_, err := s.client.CreateAd(s.ctx, req)
	s.Nil(err)

	_, err = s.clientWithError.CreateAd(s.ctx, req)
	s.NotNil(err)
}

func (s *UserSuit) TestUpdateAd() {
	req := &proto.UpdateAdRequest{}

	_, err := s.client.UpdateAd(s.ctx, req)
	s.Nil(err)

	_, err = s.clientWithError.UpdateAd(s.ctx, req)
	s.NotNil(err)
}

func (s *UserSuit) TestChangeAdStatue() {
	req := &proto.ChangeAdStatusRequest{}

	_, err := s.client.ChangeAdStatus(s.ctx, req)
	s.Nil(err)

	_, err = s.clientWithError.ChangeAdStatus(s.ctx, req)
	s.NotNil(err)
}

func (s *UserSuit) TestGetAd() {
	req := &proto.GetAdRequest{}

	_, err := s.client.GetAd(s.ctx, req)
	s.Nil(err)

	_, err = s.clientWithError.GetAd(s.ctx, req)
	s.NotNil(err)
}

func (s *UserSuit) TestRemoveAd() {
	req := &proto.DeleteAdRequest{}

	_, err := s.client.DeleteAd(s.ctx, req)
	s.Nil(err)

	_, err = s.clientWithError.DeleteAd(s.ctx, req)
	s.NotNil(err)
}

func (s *UserSuit) TestListAds() {
	req := &proto.ListAdRequest{}

	// error will be never
	_, err := s.client.ListAds(s.ctx, req)
	s.Nil(err)

	_, err = s.clientWithError.ListAds(s.ctx, req)
	s.Nil(err)
}

func (s *UserSuit) TestAdsByTitle() {
	req := &proto.AdsByTitleRequest{}

	_, err := s.client.AdsByTitle(s.ctx, req)
	s.Nil(err)

	_, err = s.clientWithError.AdsByTitle(s.ctx, req)
	s.Nil(err)
}

func TestGrpcWithMockedApp(t *testing.T) {
	suite.Run(t, new(UserSuit))
}
