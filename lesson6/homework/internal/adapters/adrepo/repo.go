package adrepo

import (
	"errors"
	"homework6/internal/ads"
	"homework6/internal/app"
)

type MyRepo struct {
	curId    int64
	AddById  map[int64]*ads.Ad
	UserById map[int64]int64
}

func (m *MyRepo) Insert(ad *ads.Ad, userId int64) {
	ad.ID = m.curId
	m.AddById[m.curId] = ad
	m.UserById[m.curId] = userId
	m.curId++
}

func (m *MyRepo) Get(adId, userId int64) (*ads.Ad, error) {
	if m.UserById[adId] != userId {
		return nil, errors.New("forbidden")
	}
	return m.AddById[adId], nil
}

func (m *MyRepo) Remove(id int64) {
	delete(m.UserById, id)
}

func (m *MyRepo) GetNewId() int64 {
	return m.curId
}
func New() app.Repository {
	return &MyRepo{curId: 0, UserById: make(map[int64]int64), AddById: make(map[int64]*ads.Ad)}
}
