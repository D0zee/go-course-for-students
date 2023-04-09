package adrepo

import (
	"errors"
	"homework6/internal/ads"
)

var ErrAccess = errors.New("forbidden")

type Repository interface {
	Insert(ad *ads.Ad, userId int64)
	Get(adId, userId int64) (*ads.Ad, error)
	GetNewId() int64
}

type adRepo struct {
	curAdId  int64
	AddById  map[int64]*ads.Ad
	UserById map[int64]int64
}

func (m *adRepo) Insert(ad *ads.Ad, userId int64) {
	ad.ID = m.curAdId
	m.curAdId++
	m.AddById[ad.ID] = ad
	m.UserById[ad.ID] = userId
}

func (m *adRepo) Get(adId, userId int64) (*ads.Ad, error) {
	if m.UserById[adId] != userId {
		return nil, ErrAccess
	}
	return m.AddById[adId], nil
}

func (m *adRepo) GetNewId() int64 {
	return m.curAdId
}
func New() Repository {
	return &adRepo{curAdId: 0,
		AddById:  make(map[int64]*ads.Ad),
		UserById: make(map[int64]int64)}
}
