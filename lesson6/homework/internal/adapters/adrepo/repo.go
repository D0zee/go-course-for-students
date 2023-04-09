package adrepo

import (
	"errors"
	"github.com/D0zee/advalidator"
	"homework6/internal/ads"
)

var ErrAccess = errors.New("forbidden")
var ErrValidate = errors.New("not validated")

type Repository interface {
	Insert(ad ads.Ad) error
	Get(adId, userId int64) (ads.Ad, error)
	GetNewId() int64
	ReplaceById(ad ads.Ad, adId, userId int64) error
}

type adRepo struct {
	curAdId  int64
	AddById  map[int64]ads.Ad
	UserById map[int64]int64
}

func (m *adRepo) Insert(ad ads.Ad) error {
	if err := advalidator.Validate(ad); err != nil {
		return ErrValidate
	}
	adId := ad.ID
	m.AddById[adId] = ad
	m.UserById[adId] = ad.AuthorID
	return nil
}

func (m *adRepo) Get(adId, userId int64) (ads.Ad, error) {
	if m.UserById[adId] != userId {
		return ads.Ad{}, ErrAccess
	}
	return m.AddById[adId], nil
}

func (m *adRepo) ReplaceById(ad ads.Ad, adId, userId int64) error {
	if m.UserById[adId] != userId {
		return ErrAccess
	}
	if err := advalidator.Validate(ad); err != nil {
		return ErrValidate
	}
	m.AddById[adId] = ad
	return nil
}

func (m *adRepo) GetNewId() int64 {
	oldId := m.curAdId
	m.curAdId++
	return oldId
}
func New() Repository {
	return &adRepo{curAdId: 0,
		AddById:  make(map[int64]ads.Ad),
		UserById: make(map[int64]int64)}
}
