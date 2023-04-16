package app

import (
	"context"
	"errors"
	"github.com/D0zee/advalidator"
	"homework8/internal/adapters/adrepo"
	"homework8/internal/ads"
)

var ErrInternal = errors.New("internal error")
var ErrAccess = errors.New("forbidden")
var ErrValidate = errors.New("not validated")

type App interface {
	CreateAd(ctx context.Context, title, text string, userId int64) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adId, userId int64, published bool) (*ads.Ad, error)
	UpdateAd(ctx context.Context, adId, userId int64, title, text string) (*ads.Ad, error)
}

type AdApp struct {
	Repo adrepo.Repository
}

func (a *AdApp) CreateAd(ctx context.Context, title, text string, userId int64) (*ads.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, ErrInternal
	default:
	}
	ad := ads.Ad{ID: a.Repo.GetNewId(), Title: title, Text: text, AuthorID: userId}
	if err := advalidator.Validate(ad); err != nil {
		return nil, ErrValidate
	}
	a.Repo.Insert(ad)
	return &ad, nil
}

func (a *AdApp) ChangeAdStatus(ctx context.Context, adId, userId int64, published bool) (*ads.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, ErrInternal
	default:
	}
	if !a.Repo.CheckAccess(adId, userId) {
		return nil, ErrAccess
	}
	ad := a.Repo.Get(adId, userId)
	ad.Published = published
	a.Repo.ReplaceById(ad, ad.ID, userId)
	return &ad, nil
}

func (a *AdApp) UpdateAd(ctx context.Context, adId, userId int64, title, text string) (*ads.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, ErrInternal
	default:
	}
	if !a.Repo.CheckAccess(adId, userId) {
		return nil, ErrAccess
	}
	ad := a.Repo.Get(adId, userId)
	ad.Title = title
	ad.Text = text
	if err := advalidator.Validate(ad); err != nil {
		return nil, ErrValidate
	}
	a.Repo.ReplaceById(ad, ad.ID, userId)
	return &ad, nil
}

func NewApp(repo adrepo.Repository) App {
	return &AdApp{Repo: repo}
}
