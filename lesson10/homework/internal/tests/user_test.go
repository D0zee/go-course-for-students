package tests

import (
	"github.com/stretchr/testify/assert"
	"homework9/internal/app"
	"testing"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("aboba", "pushkin@ya.ru")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Nickname, "aboba")
	assert.Equal(t, response.Data.Email, "pushkin@ya.ru")
}

func TestChangeUser(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("aboba", "pushkin@ya.ru")
	assert.NoError(t, err)

	userId := uResponse.Data.ID

	response, err := client.updateUser(userId, "nickname", "Олег")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.ID, userId)
	assert.Equal(t, response.Data.Nickname, "Олег")
	assert.Equal(t, response.Data.Email, "pushkin@ya.ru")

	response, err = client.updateUser(userId, "email", "aoa@google.com")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.ID, userId)
	assert.Equal(t, response.Data.Nickname, "Олег")
	assert.Equal(t, response.Data.Email, "aoa@google.com")

	// wrong ID
	response, err = client.updateUser(1, "email", "aoa@google.com")
	assert.Error(t, err, app.ErrWrongUserId)
}

func TestGetUser(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("aboba", "pushkin@ya.ru")
	assert.NoError(t, err)

	userId := uResponse.Data.ID

	// get same user
	user, err := client.getUser(userId)
	assert.NoError(t, err)
	assert.Equal(t, user, uResponse)

	// get non existing user
	_, err = client.getUser(-1)
	assert.Error(t, err, app.ErrWrongUserId)
}

func TestRemoveUser(t *testing.T) {
	client := getTestClient()

	userOriginal, err := client.createUser("aboba", "pushkin@ya.ru")
	assert.NoError(t, err)

	// remove same user
	user, err := client.removeUser(userOriginal.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, userOriginal)

	// remove same user again is impossible
	_, err = client.removeUser(userOriginal.Data.ID)
	assert.Error(t, err, app.ErrWrongUserId)

	// get same user is impossible
	_, err = client.getUser(userOriginal.Data.ID)
	assert.Error(t, err, app.ErrWrongUserId)
}
