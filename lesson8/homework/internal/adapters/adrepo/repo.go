package adrepo

import (
	"homework8/internal/ads"
	"homework8/internal/users"
)

type Repository[T any] interface {
	Insert(object *T)
	Get(id int64) (*T, bool)
	GetCurAvailableId() int64
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
