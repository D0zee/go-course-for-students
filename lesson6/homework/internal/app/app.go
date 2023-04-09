package app

import (
	"context"
	"errors"
	"fmt"
	"homework6/internal/adapters/adrepo"
	"homework6/internal/ads"
)

var ErrInternal = errors.New("internal error")

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
	ad := ads.Ad{Title: title, Text: text, AuthorID: userId}
	fmt.Println(title)
	fmt.Println(text)
	err := a.Repo.Insert(&ad, userId)
	if err != nil {
		return nil, err
	}
	return &ad, nil
}

func (a *AdApp) ChangeAdStatus(ctx context.Context, adId, userId int64, published bool) (*ads.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, ErrInternal
	default:
	}
	ad, err := a.Repo.Get(adId, userId)
	if err != nil {
		return nil, err
	}
	ad.Published = published
	return ad, nil
}

func (a *AdApp) UpdateAd(ctx context.Context, adId, userId int64, title, text string) (*ads.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, ErrInternal
	default:
	}
	ad, err := a.Repo.Get(adId, userId)
	if err != nil {
		return nil, err
	}
	ad.Title = title
	ad.Text = text
	return ad, nil
}

func NewApp(repo adrepo.Repository) App {
	return &AdApp{Repo: repo}
}
