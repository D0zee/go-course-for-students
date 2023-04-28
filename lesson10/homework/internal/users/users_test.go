package users

import (
	"reflect"
	"testing"
)

type testCase struct {
	Id       int64
	Nickname string
	Email    string

	Expected *User
}

func TestCreatingUser(t *testing.T) {
	testCases := []testCase{{
		Id:       10,
		Nickname: "Nikolai",
		Email:    "work@yandex.ru",
		Expected: &User{
			Id:       10,
			Nickname: "Nikolai",
			Email:    "work@yandex.ru",
			Deleted:  false,
		},
	},
		{
			Id:       0,
			Nickname: "Иван",
			Email:    "русскаяраскладка@tinkoff.com",
			Expected: &User{
				Id:       0,
				Nickname: "Иван",
				Email:    "русскаяраскладка@tinkoff.com",
				Deleted:  false,
			},
		},

		{
			Id:       -1,
			Nickname: "",
			Email:    "",
			Expected: &User{
				Id:       -1,
				Nickname: "",
				Email:    "",
				Deleted:  false,
			},
		},
	}

	for _, test := range testCases {
		userFromFunc := New(test.Id, test.Nickname, test.Email)
		if !reflect.DeepEqual(userFromFunc, test.Expected) {
			t.Fatalf(`testCreatingUser: expect %v got %v`, test.Expected, userFromFunc)
		}
	}
}
