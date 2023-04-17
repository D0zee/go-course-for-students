package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/D0zee/advalidator"
	"homework8/internal/adapters/adrepo"
	"homework8/internal/ads"
	"homework8/internal/users"
	"log"
	"time"
)

var ErrInternal = errors.New("internal error")
var ErrAccess = errors.New("forbidden")
var ErrValidate = errors.New("not validated")
var ErrWrongUserId = errors.New("not contain user with this id")
var ErrAvailabilityAd = errors.New("ad with this id is not created")

type App interface {
	CreateAd(ctx context.Context, title, text string, userId int64) (ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adId, userId int64, published bool) (ads.Ad, error)
	UpdateAd(ctx context.Context, adId, userId int64, title, text string) (ads.Ad, error)
	GetAdById(ctx context.Context, adId, userId int64) (ads.Ad, error)
	Access(adId, userId int64) error

	CreateUser(ctx context.Context, nickname, email string) (users.User, error)
	UpdateUser(ctx context.Context, userId int64, nickname string, m Method) (users.User, error)
}

type AdApp struct {
	Repo     adrepo.Repository[ads.Ad]
	UserRepo adrepo.Repository[users.User]
}

func (a *AdApp) CreateAd(ctx context.Context, title, text string, userId int64) (ads.Ad, error) {
	select {
	case <-ctx.Done():
		log.Println("APP: troubles with context")
		return ads.Ad{}, ErrInternal
	default:
	}
	_, contain := a.UserRepo.Get(userId)
	if !contain {
		log.Println("APP: access error")
		return ads.Ad{}, ErrAccess
	}

	currentTime := time.Now()
	ad := ads.Ad{ID: a.Repo.GetCurAvailableId(), Title: title, Text: text,
		AuthorID: userId, CreationTime: currentTime, UpdateTime: currentTime}
	fmt.Println("AD:", ad)
	if err := advalidator.Validate(ad); err != nil {
		log.Println("APP: validation error")
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

	ad, contain := a.Repo.Get(adId)
	if !contain {
		return ErrWrongUserId
	}
	if ad.AuthorID != user.Id {
		return ErrAccess
	}
	return nil
}

func (a *AdApp) ChangeAdStatus(ctx context.Context, adId, userId int64, published bool) (ads.Ad, error) {
	select {
	case <-ctx.Done():
		return ads.Ad{}, ErrInternal
	default:
	}
	if err := a.Access(adId, userId); err != nil {
		return ads.Ad{}, err
	}
	ad, _ := a.Repo.Get(adId)
	ad.UpdateTime = time.Now()
	ad.Published = published
	return *ad, nil
}

func (a *AdApp) UpdateAd(ctx context.Context, adId, userId int64, title, text string) (ads.Ad, error) {
	select {
	case <-ctx.Done():
		return ads.Ad{}, ErrInternal
	default:
	}
	if err := a.Access(adId, userId); err != nil {
		return ads.Ad{}, err
	}
	ad, _ := a.Repo.Get(adId)

	newAd := *ad
	newAd.Title = title
	newAd.Text = text
	newAd.UpdateTime = time.Now()
	if err := advalidator.Validate(newAd); err != nil {
		return ads.Ad{}, ErrValidate
	}
	ad = &newAd
	return *ad, nil
}

func (a *AdApp) GetAdById(ctx context.Context, adId, userId int64) (ads.Ad, error) {
	select {
	case <-ctx.Done():
		return ads.Ad{}, ErrInternal
	default:
	}
	if err := a.Access(adId, userId); err != nil {
		return ads.Ad{}, err
	}
	ad, _ := a.Repo.Get(adId)
	return *ad, nil
}

func (a *AdApp) CreateUser(ctx context.Context, nickname, email string) (users.User, error) {
	select {
	case <-ctx.Done():
		return users.User{}, ErrInternal
	default:
	}
	userId := a.Repo.GetCurAvailableId()
	user := users.New(userId, nickname, email)
	// todo: validation of fields
	a.UserRepo.Insert(user)
	return *user, nil
}

type Method int64

const (
	ChangeEmail Method = iota
	ChangeNickname
)

func (a *AdApp) UpdateUser(ctx context.Context, userId int64, data string, m Method) (users.User, error) {
	select {
	case <-ctx.Done():
		return users.User{}, ErrInternal
	default:
	}
	user, contain := a.UserRepo.Get(userId)
	if !contain {
		return users.User{}, ErrWrongUserId
	}
	if m == ChangeEmail {
		user.Email = data
	} else if m == ChangeNickname {
		user.Nickname = data
	} else {
		panic(m)
	}
	return *user, nil
}

func NewApp(repo adrepo.Repository[ads.Ad], urepo adrepo.Repository[users.User]) App {
	return &AdApp{Repo: repo, UserRepo: urepo}
}
