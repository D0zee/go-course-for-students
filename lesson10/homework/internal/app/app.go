package app

import (
	"context"
	"errors"
	"github.com/D0zee/advalidator"
	"homework9/internal/adapters/repo"
	"homework9/internal/ads"
	"homework9/internal/users"
	"time"
)

var ErrInternal = errors.New("internal error")
var ErrAccess = errors.New("forbidden")
var ErrValidate = errors.New("not validated")
var ErrWrongUserId = errors.New("not contain user with this id")
var ErrAvailabilityAd = errors.New("not contain ad with this id")

type App interface {
	CreateAd(ctx context.Context, title, text string, userId int64) (ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adId, userId int64, published bool) (ads.Ad, error)
	UpdateAd(ctx context.Context, adId, userId int64, title, text string) (ads.Ad, error)
	GetAdById(ctx context.Context, adId, userId int64) (ads.Ad, error)
	RemoveAd(ctx context.Context, adId, userId int64) (ads.Ad, error)
	Access(adId, userId int64) error

	ListAds(ctx context.Context) []ads.Ad

	CreateUser(ctx context.Context, nickname, email string) (users.User, error)
	UpdateUser(ctx context.Context, userId int64, nickname string, m Method) (users.User, error)
	RemoveUser(ctx context.Context, userId int64) (users.User, error)
	GetUser(ctx context.Context, userId int64) (users.User, error)
}

func contextEnd(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

type AdApp struct {
	Repo     repo.Repository[ads.Ad]
	UserRepo repo.Repository[users.User]
}

func (a *AdApp) CreateAd(ctx context.Context, title, text string, userId int64) (ads.Ad, error) {
	if contextEnd(ctx) {
		return ads.Ad{}, ErrInternal
	}
	_, contain := a.UserRepo.Get(userId)
	if !contain {
		return ads.Ad{}, ErrAccess
	}

	currentTime := time.Now().UTC()
	ad := ads.Ad{ID: a.Repo.GetCurAvailableId(), Title: title, Text: text,
		AuthorID: userId, CreationTime: currentTime, UpdateTime: currentTime}
	if err := advalidator.Validate(ad); err != nil {
		return ads.Ad{}, ErrValidate
	}
	a.Repo.Insert(&ad)
	return ad, nil
}

func (a *AdApp) Access(adId, userId int64) error {
	user, contain := a.UserRepo.Get(userId)
	if !contain {
		return ErrAccess
	}

	if user.Deleted {
		return ErrWrongUserId
	}

	ad, contain := a.Repo.Get(adId)
	if !contain || ad.Deleted {
		return ErrAvailabilityAd
	}
	if ad.AuthorID != user.Id {
		return ErrAccess
	}
	return nil
}

func (a *AdApp) ChangeAdStatus(ctx context.Context, adId, userId int64, published bool) (ads.Ad, error) {
	if contextEnd(ctx) {
		return ads.Ad{}, ErrInternal
	}
	if err := a.Access(adId, userId); err != nil {
		return ads.Ad{}, err
	}
	ad, _ := a.Repo.Get(adId)
	ad.UpdateTime = time.Now().UTC()
	ad.Published = published
	return *ad, nil
}

func (a *AdApp) UpdateAd(ctx context.Context, adId, userId int64, title, text string) (ads.Ad, error) {
	if contextEnd(ctx) {
		return ads.Ad{}, ErrInternal
	}
	if err := a.Access(adId, userId); err != nil {
		return ads.Ad{}, err
	}
	ad, _ := a.Repo.Get(adId)

	newAd := *ad
	newAd.Title = title
	newAd.Text = text
	newAd.UpdateTime = time.Now().UTC()
	if err := advalidator.Validate(newAd); err != nil {
		return ads.Ad{}, ErrValidate
	}
	ad = &newAd
	return *ad, nil
}

func (a *AdApp) GetAdById(ctx context.Context, adId, userId int64) (ads.Ad, error) {
	if contextEnd(ctx) {
		return ads.Ad{}, ErrInternal
	}
	if err := a.Access(adId, userId); err != nil {
		return ads.Ad{}, err
	}
	ad, _ := a.Repo.Get(adId)
	return *ad, nil
}

func (a *AdApp) RemoveAd(ctx context.Context, adId, userId int64) (ads.Ad, error) {
	if contextEnd(ctx) {
		return ads.Ad{}, ErrInternal
	}
	if err := a.Access(adId, userId); err != nil {
		return ads.Ad{}, err
	}
	ad, _ := a.Repo.Get(adId)
	ad.Deleted = true
	return *ad, nil
}

func (a *AdApp) ListAds(ctx context.Context) []ads.Ad {
	if contextEnd(ctx) {
		return []ads.Ad{}
	}
	var result []ads.Ad
	for i := int64(0); i < a.Repo.GetCurAvailableId(); i++ {
		ad, _ := a.Repo.Get(i)
		if !ad.Deleted {
			result = append(result, *ad)
		}
	}
	return result
}

func (a *AdApp) CreateUser(ctx context.Context, nickname, email string) (users.User, error) {
	if contextEnd(ctx) {
		return users.User{}, ErrInternal
	}
	userId := a.UserRepo.GetCurAvailableId()
	user := users.New(userId, nickname, email)
	a.UserRepo.Insert(user)
	return *user, nil
}

type Method int64

const (
	ChangeEmail Method = iota
	ChangeNickname
)

func (a *AdApp) UpdateUser(ctx context.Context, userId int64, data string, m Method) (users.User, error) {
	if contextEnd(ctx) {
		return users.User{}, ErrInternal
	}
	user, contain := a.UserRepo.Get(userId)
	if !contain || user.Deleted {
		return users.User{}, ErrWrongUserId
	}
	if m == ChangeEmail {
		user.Email = data
	} else if m == ChangeNickname {
		user.Nickname = data
	}
	return *user, nil
}

func (a *AdApp) RemoveUser(ctx context.Context, userId int64) (users.User, error) {
	if contextEnd(ctx) {
		return users.User{}, ErrInternal
	}
	user, contain := a.UserRepo.Get(userId)
	if !contain || user.Deleted {
		return users.User{}, ErrWrongUserId
	}
	user.Deleted = true
	return *user, nil
}

func (a *AdApp) GetUser(ctx context.Context, userId int64) (users.User, error) {
	if contextEnd(ctx) {
		return users.User{}, ErrInternal
	}
	user, contain := a.UserRepo.Get(userId)
	if !contain || user.Deleted {
		return users.User{}, ErrWrongUserId
	}
	return *user, nil
}

func NewApp(repo repo.Repository[ads.Ad], urepo repo.Repository[users.User]) App {
	return &AdApp{Repo: repo, UserRepo: urepo}
}
