package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func isSameDate(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

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
	assert.True(t, isSameDate(response.Data.CreationTime, time.Now()))
	assert.True(t, isSameDate(response.Data.UpdateTime, time.Now()))
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
	assert.True(t, isSameDate(response.Data.UpdateTime, time.Now()))

	response, err = client.changeAdStatus(userId, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
	assert.True(t, isSameDate(response.Data.UpdateTime, time.Now()))

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
	assert.True(t, isSameDate(response.Data.UpdateTime, time.Now()))
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
