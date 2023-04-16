package tests

import (
	"github.com/stretchr/testify/assert"
	"homework8/internal/app"
	"testing"
)

func TestGetAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	response, err := client.getAd(0, 123)
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(123))
	assert.False(t, response.Data.Published)

	// ad with this id isn't
	_, err = client.getAd(1, 123)
	assert.Error(t, err, app.ErrAvailabilityAd)

	// access permitted
	_, err = client.getAd(0, 122)
	assert.Error(t, err, app.ErrAccess)

}
