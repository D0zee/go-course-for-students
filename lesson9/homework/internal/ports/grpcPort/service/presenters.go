package service

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework9/internal/ads"
	"homework9/internal/ports/grpcPort/proto"
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
	var list []*proto.AdResponse
	for _, ad := range ads {
		list = append(list, AdSuccessResponse(ad))
	}
	return &proto.ListAdResponse{List: list}
}

func UserSuccessResponse(user users.User) *proto.UserResponse {
	return &proto.UserResponse{
		Id:    user.Id,
		Name:  user.Nickname,
		Email: user.Email,
	}
}
