package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("Oleg", "ya@ya.ru")
	assert.NoError(t, err)

	response, err := client.createAd(uResponse.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(0))
	assert.False(t, response.Data.Published)
	assert.True(t, IsSameTimes(response.Data.CreationTime, time.Now().UTC()))
	assert.True(t, IsSameTimes(response.Data.UpdateTime, time.Now().UTC()))
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("Oleg", "ya@ya.ru")
	assert.NoError(t, err)
	userId := uResponse.Data.ID

	response, err := client.createAd(userId, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(userId, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)
	assert.True(t, IsSameTimes(response.Data.UpdateTime, time.Now().UTC()))

	response, err = client.changeAdStatus(userId, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
	assert.True(t, IsSameTimes(response.Data.UpdateTime, time.Now().UTC()))

	response, err = client.changeAdStatus(userId, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("Oleg", "ya@ya.ru")
	assert.NoError(t, err)

	userId := uResponse.Data.ID

	response, err := client.createAd(userId, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(userId, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
	assert.True(t, IsSameTimes(response.Data.UpdateTime, time.Now().UTC()))
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("Oleg", "ya@ya.ru")
	assert.NoError(t, err)

	userId := uResponse.Data.ID

	response, err := client.createAd(userId, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(userId, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(userId, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds("")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}
