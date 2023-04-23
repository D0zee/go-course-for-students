package grpc

import (
	"homework9/internal/ads"
	"homework9/internal/ports/grpc/service"
	"homework9/internal/users"
)

func AdResponse(ad ads.Ad) *service.AdResponse {
	return &service.AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorId:  ad.AuthorID,
		Published: ad.Published}
}

func ListAdResponse(ads []ads.Ad) *service.ListAdResponse {
	var response service.ListAdResponse
	list := response.GetList()
	for _, ad := range ads {
		list = append(list, AdResponse(ad))
	}
	return &response
}

func UserResponse(user users.User) *service.UserResponse {
	return &service.UserResponse{
		Id:    user.Id,
		Name:  user.Nickname,
		Email: user.Email,
	}
}
