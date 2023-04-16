package adrepo

import (
	"homework8/internal/ads"
	"homework8/internal/users"
)

type Repository[T any] interface {
	Insert(object *T)
	Get(id int64) (*T, bool)
	GetCurAvailableId() int64
	//ReplaceById(object T, id int64)
}

type myRepo[T any] struct {
	CurId   int64
	IdToObj map[int64]*T
}

func (r *myRepo[T]) Insert(obj *T) {
	r.IdToObj[r.CurId] = obj
	r.CurId++
}

func (r *myRepo[T]) Get(id int64) (*T, bool) {
	value, contain := r.IdToObj[id]
	if !contain {
		return nil, false
	}
	return value, true
}

func (r *myRepo[T]) GetCurAvailableId() int64 {
	return r.CurId
}

func NewAdRepo() Repository[ads.Ad] {
	return &myRepo[ads.Ad]{CurId: 0,
		IdToObj: make(map[int64]*ads.Ad)}
}

func NewUserRepo() Repository[users.User] {
	return &myRepo[users.User]{CurId: 0,
		IdToObj: make(map[int64]*users.User)}
}

//type userRepo struct {
//	CurUserId int64
//	IdToUser  map[int64]users.User
//}
//
//func (r *userRepo) Insert(user users.User) {
//	r.IdToUser[r.CurUserId] = user
//	r.CurUserId++
//}
//
//func (r *userRepo) Get(id int64) (users.User, bool) {
//	return r.IdToUser[id]
//}
//
//func (r *userRepo) GetCurAvailableId() int64 {
//	return r.CurUserId
//}
//
//func (r *userRepo) ReplaceById(user users.User, id int64) {
//	r.IdToUser[id] = user
//}

//func (r *adRepo) ReplaceById(ad ads.Ad, adId int64) {
//	r.AddById[adId] = ad
//}
//
//func (r *adRepo) GetCurAvailableId() int64 {
//	return r.CurAdId
//}
//
//type adRepo struct {
//	CurAdId int64
//	AddById map[int64]ads.Ad
//}
//
//func (r *adRepo) Insert(ad ads.Ad) {
//	adId := ad.ID
//	r.AddById[adId] = ad
//	r.CurAdId++
//}
//
//func (r *adRepo) Get(adId int64) (ads.Ad, bool) {
//	value, contain := r.AddById[adId]
//	if !contain {
//		return ads.Ad{}, false
//	}
//	return value, true
//}
//
//func (r *adRepo) ReplaceById(ad ads.Ad, adId int64) {
//	r.AddById[adId] = ad
//}
//
//func (r *adRepo) GetCurAvailableId() int64 {
//	return r.CurAdId
//}
