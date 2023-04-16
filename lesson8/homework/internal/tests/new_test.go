package tests

import (
	"github.com/stretchr/testify/assert"
	"homework8/internal/app"
	"testing"
)

func TestGetAd(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("Oleg", "ya@ya.ru")
	assert.NoError(t, err)

	userId := uResponse.Data.ID

	_, err = client.createAd(userId, "hello", "world")
	assert.NoError(t, err)

	response, err := client.getAd(0, userId)
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, userId)
	assert.False(t, response.Data.Published)

	// ad with this id isn't
	_, err = client.getAd(1, userId)
	assert.Error(t, err, app.ErrAvailabilityAd)

	// access permitted
	_, err = client.getAd(0, 122)
	assert.Error(t, err, app.ErrAccess)

}
