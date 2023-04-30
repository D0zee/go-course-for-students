package repo

import (
	"homework9/internal/ads"
	"homework9/internal/users"
	"sync"
)

type Repository[T any] interface {
	Insert(object *T)
	Get(id int64) (*T, bool)
	GetCurAvailableId() int64
}

type mapRepo[T any] struct {
	mx      *sync.RWMutex
	CurId   int64
	IdToObj map[int64]*T
}

func (r *mapRepo[T]) Insert(obj *T) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.IdToObj[r.CurId] = obj
	r.CurId++
}

func (r *mapRepo[T]) Get(id int64) (*T, bool) {
	r.mx.Lock()
	defer r.mx.Unlock()

	value, contain := r.IdToObj[id]
	if !contain {
		return nil, false
	}

	return value, true
}

func (r *mapRepo[T]) GetCurAvailableId() int64 {
	r.mx.RLock()
	defer r.mx.RUnlock()
	return r.CurId
}

func NewMapAdRepo() Repository[ads.Ad] {
	return &mapRepo[ads.Ad]{mx: &sync.RWMutex{}, CurId: 0,
		IdToObj: make(map[int64]*ads.Ad)}
}

func NewUserRepo() Repository[users.User] {
	return &mapRepo[users.User]{mx: &sync.RWMutex{}, CurId: 0,
		IdToObj: make(map[int64]*users.User)}
}
