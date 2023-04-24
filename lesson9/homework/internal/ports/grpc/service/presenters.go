package service

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework9/internal/ads"
	"homework9/internal/ports/grpc/proto"
	"homework9/internal/users"
)

func AdSuccessResponse(ad ads.Ad) *proto.AdResponse {
	return &proto.AdResponse{
		Id:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationTime: timestamppb.New(ad.CreationTime),
		UpdateTime:   timestamppb.New(ad.UpdateTime),
	}
}

func ListAdSuccessResponse(ads []ads.Ad) *proto.ListAdResponse {
	var response proto.ListAdResponse
	list := response.GetList()
	for _, ad := range ads {
		list = append(list, AdSuccessResponse(ad))
	}
	return &response
}

func UserSuccessResponse(user users.User) *proto.UserResponse {
	return &proto.UserResponse{
		Id:    user.Id,
		Name:  user.Nickname,
		Email: user.Email,
	}
}
