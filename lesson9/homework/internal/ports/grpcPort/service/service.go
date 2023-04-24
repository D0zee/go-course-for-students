package service

import (
	"context"
	"errors"
	"homework9/internal/ads"
	"homework9/internal/app"
	"homework9/internal/ports/grpcPort/proto"
	"homework9/internal/users"
	"strings"
)

var ErrEmptyReq = errors.New("empty request")

type Service struct {
	App app.App
}

func NewService(app app.App) *Service {
	return &Service{App: app}
}

func (s *Service) CreateAd(ctx context.Context, request *proto.CreateAdRequest) (*proto.AdResponse, error) {
	ad, err := s.App.CreateAd(ctx, request.Title, request.Text, request.UserId)
	if err != nil {
		return nil, err
	}
	return AdSuccessResponse(ad), nil
}

func (s *Service) ChangeAdStatus(ctx context.Context, request *proto.ChangeAdStatusRequest) (*proto.AdResponse, error) {
	ad, err := s.App.ChangeAdStatus(ctx, request.AdId, request.UserId, request.Published)
	if err != nil {
		return nil, err
	}
	return AdSuccessResponse(ad), nil
}

func (s *Service) UpdateAd(ctx context.Context, request *proto.UpdateAdRequest) (*proto.AdResponse, error) {
	ad, err := s.App.UpdateAd(ctx, request.AdId, request.UserId, request.Title, request.Text)
	if err != nil {
		return nil, err
	}
	return AdSuccessResponse(ad), nil
}

func (s *Service) GetAd(ctx context.Context, request *proto.GetAdRequest) (*proto.AdResponse, error) {
	ad, err := s.App.GetAdById(ctx, request.AdId, request.UserId)
	if err != nil {
		return nil, err
	}
	return AdSuccessResponse(ad), nil
}

func (s *Service) DeleteAd(ctx context.Context, request *proto.DeleteAdRequest) (*proto.AdResponse, error) {
	ad, err := s.App.RemoveAd(ctx, request.AdId, request.UserId)
	if err != nil {
		return nil, err
	}
	return AdSuccessResponse(ad), nil
}

type adPredicate func(ad ads.Ad) bool

type adPredicates struct {
	predicates []adPredicate
}

func filter(Ads []ads.Ad, p adPredicates) []ads.Ad {
	var result []ads.Ad
	for _, ad := range Ads {
		passed := true
		for _, predicate := range p.predicates {
			if !predicate(ad) {
				passed = false
				break
			}
		}
		if passed {
			result = append(result, ad)
		}

	}
	return result
}

func (s *Service) ListAds(ctx context.Context, request *proto.ListAdRequest) (*proto.ListAdResponse, error) {
	var filterFunc adPredicates
	listAds := s.App.ListAds(ctx)
	if request.Time != nil {
		filterFunc.predicates = append(filterFunc.predicates, func(ad ads.Ad) bool {
			reqTime := request.Time.AsTime()
			adTime := ad.CreationTime
			return reqTime.Day() == adTime.Day() &&
				reqTime.Month() == adTime.Month() &&
				reqTime.Year() == adTime.Year()
		})
	}

	if request.UserId != nil {
		filterFunc.predicates = append(filterFunc.predicates, func(ad ads.Ad) bool {
			return ad.AuthorID == *request.UserId
		})
	}

	if request.Time == nil && request.UserId == nil {
		filterFunc.predicates = append(filterFunc.predicates, func(ad ads.Ad) bool {
			return ad.Published
		})
	}
	listAds = filter(listAds, filterFunc)
	return ListAdSuccessResponse(listAds), nil
}

func (s *Service) AdsByTitle(ctx context.Context, request *proto.AdsByTitleRequest) (*proto.ListAdResponse, error) {
	listAds := s.App.ListAds(ctx)
	title := request.Title

	var adsWithTitle []ads.Ad
	for _, ad := range listAds {
		if strings.HasPrefix(ad.Title, title) {
			adsWithTitle = append(adsWithTitle, ad)
		}
	}
	return ListAdSuccessResponse(adsWithTitle), nil
}

func (s *Service) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.UserResponse, error) {
	user, err := s.App.CreateUser(ctx, request.Name, request.Email)
	if err != nil {
		return nil, err
	}
	return UserSuccessResponse(user), nil
}

func (s *Service) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.UserResponse, error) {
	user, err := s.App.GetUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return UserSuccessResponse(user), nil
}

func (s *Service) ChangeUser(ctx context.Context, request *proto.ChangeUserRequest) (*proto.UserResponse, error) {
	if request.Email == nil && request.Nickname == nil {
		return nil, ErrEmptyReq
	}
	var user users.User
	var err error

	if request.Email != nil {
		if user, err = s.App.UpdateUser(ctx, request.Id, *request.Email, app.ChangeEmail); err != nil {
			return nil, err
		}
	}

	if request.Nickname != nil {
		if user, err = s.App.UpdateUser(ctx, request.Id, *request.Nickname, app.ChangeNickname); err != nil {
			return nil, err
		}
	}

	return UserSuccessResponse(user), nil
}

func (s *Service) DeleteUser(ctx context.Context, request *proto.DeleteUserRequest) (*proto.UserResponse, error) {
	user, err := s.App.RemoveUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return UserSuccessResponse(user), nil
}
