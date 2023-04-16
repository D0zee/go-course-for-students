package adrepo

import (
	"homework8/internal/ads"
	"homework8/internal/users"
)

type Repository interface {
	Insert(ad ads.Ad)
	Get(adId int64) ads.Ad
	GetCurAvailableId() int64
	ReplaceById(ad ads.Ad, adId, userId int64)
	CheckAccess(adId, userId int64) bool

	InsertUser(user *users.User)
	GetUserById(id int64) *users.User
	ContainsUserWithId(id int64) bool
	GetNextUserId() int64
}

type adRepo struct {
	CurAdId  int64
	AddById  map[int64]ads.Ad
	AdToUser map[int64]int64

	CurUserId int64
	IdToUser  map[int64]*users.User
}

func (m *adRepo) Insert(ad ads.Ad) {
	adId := ad.ID
	m.AddById[adId] = ad
	m.AdToUser[adId] = ad.AuthorID
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
	return m.AdToUser[adId] == userId
}

func (s *adRepo) InsertUser(user *users.User) {
	s.IdToUser[s.CurUserId] = user
	s.CurUserId++
}

func (s *adRepo) GetNextUserId() int64 {
	return s.CurUserId
}

func (s *adRepo) GetUserById(id int64) *users.User {
	return s.IdToUser[id]
}

func (s *adRepo) ContainsUserWithId(id int64) bool {
	if _, contain := s.IdToUser[id]; !contain {
		return false
	}
	return true
}

func New() Repository {
	return &adRepo{CurAdId: 0,
		AddById:  make(map[int64]ads.Ad),
		AdToUser: make(map[int64]int64)}
}
