package users

type User struct {
	Id       int64
	Nickname string
	Email    string
}

func New(id int64, nickname, email string) *User {
	return &User{id, nickname, email}
}
