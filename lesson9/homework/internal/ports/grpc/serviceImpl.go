package grpc

import (
	"context"
	"errors"
	"homework9/internal/ads"
	"homework9/internal/app"
	"homework9/internal/ports/grpc/service"
	"homework9/internal/users"
	"strings"
)

var ErrCtx = errors.New("problems with ctx")
var ErrEmptyReq = errors.New("empty request")

type ServiceImpl struct {
	App app.App
}

func NewService(app app.App) *ServiceImpl {
	return &ServiceImpl{App: app}
}

func (s *ServiceImpl) CreateAd(ctx context.Context, request *service.CreateAdRequest) (*service.AdResponse, error) {
	ad, err := s.App.CreateAd(ctx, request.Title, request.Title, request.UserId)
	if err != nil {
		return nil, err
	}
	return AdResponse(ad), nil
}

func (s *ServiceImpl) ChangeAdStatus(ctx context.Context, request *service.ChangeAdStatusRequest) (*service.AdResponse, error) {
	ad, err := s.App.ChangeAdStatus(ctx, request.AdId, request.UserId, request.Published)
	if err != nil {
		return nil, err
	}
	return AdResponse(ad), nil
}

func (s *ServiceImpl) UpdateAd(ctx context.Context, request *service.UpdateAdRequest) (*service.AdResponse, error) {
	ad, err := s.App.UpdateAd(ctx, request.AdId, request.UserId, request.Title, request.Text)
	if err != nil {
		return nil, err
	}
	return AdResponse(ad), nil
}

func (s *ServiceImpl) GetAd(ctx context.Context, request *service.GetAdRequest) (*service.AdResponse, error) {
	ad, err := s.App.GetAdById(ctx, request.AdId, request.UserId)
	if err != nil {
		return nil, err
	}
	return AdResponse(ad), nil
}

func (s *ServiceImpl) DeleteAd(ctx context.Context, request *service.DeleteAdRequest) (*service.AdResponse, error) {
	ad, err := s.App.RemoveAd(ctx, request.AdId, request.UserId)
	if err != nil {
		return nil, err
	}
	return AdResponse(ad), nil
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

func (s *ServiceImpl) ListAds(ctx context.Context, request *service.ListAdRequest) (*service.ListAdResponse, error) {
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
	return ListAdResponse(listAds), nil
}

func (s *ServiceImpl) AdsByTitle(ctx context.Context, request *service.AdsByTitleRequest) (*service.ListAdResponse, error) {
	listAds := s.App.ListAds(ctx)
	title := request.Title

	var adsWithTitle []ads.Ad
	for _, ad := range listAds {
		if strings.HasPrefix(ad.Title, title) {
			adsWithTitle = append(adsWithTitle, ad)
		}
	}
	return ListAdResponse(listAds), nil
}

func (s *ServiceImpl) CreateUser(ctx context.Context, request *service.CreateUserRequest) (*service.UserResponse, error) {
	user, err := s.App.CreateUser(ctx, request.Name, request.Email)
	if err != nil {
		return nil, err
	}
	return UserResponse(user), nil
}

func (s *ServiceImpl) GetUser(ctx context.Context, request *service.GetUserRequest) (*service.UserResponse, error) {
	user, err := s.App.GetUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return UserResponse(user), nil
}

func (s *ServiceImpl) ChangeUser(ctx context.Context, request *service.ChangeUserRequest) (*service.UserResponse, error) {
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
		if user, err = s.App.UpdateUser(ctx, request.Id, *request.Email, app.ChangeEmail); err != nil {
			return nil, err
		}
	}

	return UserResponse(user), nil
}

func (s *ServiceImpl) DeleteUser(ctx context.Context, request *service.DeleteUserRequest) (*service.UserResponse, error) {
	user, err := s.App.RemoveUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return UserResponse(user), nil
}
