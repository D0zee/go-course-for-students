package adrepo

import (
	"homework8/internal/ads"
)

type Repository interface {
	Insert(ad ads.Ad)
	Get(adId int64) ads.Ad
	GetCurAvailableId() int64
	ReplaceById(ad ads.Ad, adId, userId int64)
	CheckAccess(adId, userId int64) bool
}

type adRepo struct {
	CurAdId  int64
	AddById  map[int64]ads.Ad
	UserById map[int64]int64
}

func (m *adRepo) Insert(ad ads.Ad) {
	adId := ad.ID
	m.AddById[adId] = ad
	m.UserById[adId] = ad.AuthorID
	m.CurAdId++
}

func (m *adRepo) Get(adId int64) ads.Ad {
	return m.AddById[adId]
}

func (m *adRepo) ReplaceById(ad ads.Ad, adId, userId int64) {
	m.AddById[adId] = ad
}

func (m *adRepo) GetCurAvailableId() int64 {
	return m.CurAdId
}

func (m *adRepo) CheckAccess(adId, userId int64) bool {
	return m.UserById[adId] == userId
}

func New() Repository {
	return &adRepo{CurAdId: 0,
		AddById:  make(map[int64]ads.Ad),
		UserById: make(map[int64]int64)}
}
