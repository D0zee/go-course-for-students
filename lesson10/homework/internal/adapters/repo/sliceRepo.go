package repo

import (
	"homework9/internal/ads"
	"sync"
)

type sliceRepo[T any] struct {
	mx    *sync.RWMutex
	CurId int64
	store []*T
}

func (r *sliceRepo[T]) Insert(obj *T) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.store = append(r.store, obj)
	r.CurId++
}

func (r *sliceRepo[T]) Get(id int64) (*T, bool) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if id < 0 || id >= int64(len(r.store)) {
		return nil, false
	}

	return r.store[id], true
}

func (r *sliceRepo[T]) GetCurAvailableId() int64 {
	r.mx.RLock()
	defer r.mx.RUnlock()
	return r.CurId
}

func NewSliceAdRepo() Repository[ads.Ad] {
	return &sliceRepo[ads.Ad]{mx: &sync.RWMutex{}, CurId: 0}
}
