package app

import (
	"context"
	"homework6/internal/ads"
)

type App interface {
	CreateAd(ctx *context.Context, title, text string, userId int64) *ads.Ad
	ChangeAdStatus(ctx *context.Context, adId, userId int64, published bool) (*ads.Ad, error)
	UpdateAd(ctx *context.Context, adId, userId int64, title, text string) (*ads.Ad, error)
}

type MyApp struct {
	repo Repository
}

func (app *MyApp) CreateAd(ctx *context.Context, title, text string, userId int64) *ads.Ad {
	ad := ads.Ad{Title: title, Text: text, AuthorID: userId, Published: false}
	app.repo.Insert(&ad, userId)
	return &ad
}

func (app *MyApp) ChangeAdStatus(ctx *context.Context, adId, userId int64, published bool) (*ads.Ad, error) {
	ad, err := app.repo.Get(adId, userId)
	if err != nil {
		return nil, err
	}
	ad.Published = published
	return ad, nil
}

func (app *MyApp) UpdateAd(ctx *context.Context, adId, userId int64, title, text string) (*ads.Ad, error) {
	ad, err := app.repo.Get(adId, userId)
	if err != nil {
		return nil, err
	}
	ad.Title = title
	ad.Text = text
	return ad, nil
}

type Repository interface {
	Insert(ad *ads.Ad, userId int64)
	Get(adId, userId int64) (*ads.Ad, error)
	GetNewId() int64
	Remove(id int64)
}

func NewApp(repo Repository) App {
	return &MyApp{repo}
}
