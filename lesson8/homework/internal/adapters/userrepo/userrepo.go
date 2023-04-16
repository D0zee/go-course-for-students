package userrepo

import (
	"homework8/internal/users"
)

type Repository interface {
	Insert(user *users.User)
	GetById(id int64) *users.User
	ContainsUserWithId(id int64) bool
	GetCurrentId() int64
}

type userRepo struct {
	currentId int64
	userById  map[int64]*users.User
}

func (s *userRepo) Insert(user *users.User) {
	s.userById[s.currentId] = user
	s.currentId++
}

func (s *userRepo) GetCurrentId() int64 {
	return s.currentId
}

func (s *userRepo) GetById(id int64) *users.User {
	return s.userById[id]
}

func (s *userRepo) ContainsUserWithId(id int64) bool {
	if _, contain := s.userById[id]; !contain {
		return false
	}
	return true
}

func New() Repository {
	return &userRepo{userById: make(map[int64]*users.User)}
}
