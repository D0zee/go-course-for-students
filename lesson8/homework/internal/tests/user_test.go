package tests

import (
	"github.com/stretchr/testify/assert"
	"homework8/internal/app"
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
