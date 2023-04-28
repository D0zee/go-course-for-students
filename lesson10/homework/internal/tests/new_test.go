package tests

import (
	"github.com/stretchr/testify/assert"
	"homework9/internal/app"
	"net/url"
	"strconv"
	"testing"
	"time"
)

// support of date are shown in basic tests

func TestGetAd(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("Oleg", "ya@ya.ru")
	assert.NoError(t, err)

	userId := uResponse.Data.ID

	ad, err := client.createAd(userId, "hello", "world")
	assert.NoError(t, err)

	response, err := client.getAd(0, userId)
	assert.NoError(t, err)
	assert.Equal(t, response, ad)

	// ad with this id isn't
	_, err = client.getAd(1, userId)
	assert.Error(t, err, app.ErrAvailabilityAd)

	// access permitted
	_, err = client.getAd(0, 122)
	assert.Error(t, err, app.ErrAccess)

}

func equalityOfAds(t *testing.T, ad1 adData, ad2 adResponse) {
	assert.Equal(t, ad1.ID, ad2.Data.ID)
	assert.Equal(t, ad1.Title, ad2.Data.Title)
	assert.Equal(t, ad1.Text, ad2.Data.Text)
	assert.Equal(t, ad1.AuthorID, ad2.Data.AuthorID)
	assert.Equal(t, ad1.Published, ad2.Data.Published)
}

func TestListAdsFilter(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("Oleg", "ya@ya.ru")
	assert.NoError(t, err)
	userId1 := uResponse.Data.ID

	response, err := client.createAd(userId1, "hello", "world")
	assert.NoError(t, err)

	ad1, err := client.changeAdStatus(userId1, response.Data.ID, true)
	assert.NoError(t, err)

	ad2, err := client.createAd(userId1, "best cat", "not for sale")
	assert.NoError(t, err)

	uResponse, err = client.createUser("Ivan", "you@ya.ru")
	assert.NoError(t, err)

	userId2 := uResponse.Data.ID

	ad3, err := client.createAd(userId2, "ad from IVAN", "it's my ad")
	assert.NoError(t, err)

	v := url.Values{}
	v.Add("author_id", strconv.Itoa(int(userId1)))
	queryString := v.Encode()

	ads, err := client.listAds(queryString)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
	equalityOfAds(t, ads.Data[0], ad1)
	equalityOfAds(t, ads.Data[1], ad2)

	v = url.Values{}
	currentTime := time.Now().UTC()
	dateStr := currentTime.Format(time.DateOnly)
	v.Add("time", dateStr)
	queryString = v.Encode()

	ads, err = client.listAds(queryString)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 3)
	equalityOfAds(t, ads.Data[0], ad1)
	equalityOfAds(t, ads.Data[1], ad2)
	equalityOfAds(t, ads.Data[2], ad3)

}

func TestAdsByTitle(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("Oleg", "ya@ya.ru")
	assert.NoError(t, err)
	userId1 := uResponse.Data.ID

	uResponse, err = client.createUser("Ivan", "aboba@ya.ru")
	assert.NoError(t, err)
	userId2 := uResponse.Data.ID

	ad1, err := client.createAd(userId1, "hello1", "world")
	assert.NoError(t, err)

	ad2, err := client.createAd(userId2, "hello2", "Hello spectators!")
	assert.NoError(t, err)

	ad3, err := client.createAd(userId2, "go is cool", "Hello spectators!")
	assert.NoError(t, err)

	helloAds, err := client.getAdsByTitle("hello")
	assert.NoError(t, err)
	assert.Len(t, helloAds.Data, 2)

	equalityOfAds(t, helloAds.Data[0], ad1)
	equalityOfAds(t, helloAds.Data[1], ad2)

	ad_, err := client.getAdsByTitle("go is cool")
	assert.NoError(t, err)
	assert.Len(t, ad_.Data, 1)
	equalityOfAds(t, ad_.Data[0], ad3)

}

func TestRemoveAd(t *testing.T) {
	client := getTestClient()

	uResponse, err := client.createUser("Oleg", "ya@ya.ru")
	assert.NoError(t, err)
	userId := uResponse.Data.ID

	adOriginal, err := client.createAd(userId, "hello1", "world")
	assert.NoError(t, err)

	// wrong user id
	_, err = client.removeAd(adOriginal.Data.ID, -1)
	assert.Error(t, err, app.ErrAccess)

	// wrong ad id
	_, err = client.removeAd(-1, userId)
	assert.Error(t, err, app.ErrAvailabilityAd)

	// ad is available
	ad, err := client.getAd(adOriginal.Data.ID, userId)
	assert.NoError(t, err)
	assert.Equal(t, ad, adOriginal)

	adRemoved, err := client.removeAd(adOriginal.Data.ID, userId)
	assert.NoError(t, err)
	assert.Equal(t, adRemoved, adOriginal)

	// get error if you want to get removed ad
	_, err = client.getAd(adOriginal.Data.ID, userId)
	assert.Error(t, err, app.ErrAvailabilityAd)

	// get error if we want to remove again
	_, err = client.removeAd(adOriginal.Data.ID, userId)
	assert.Error(t, err, app.ErrAvailabilityAd)
}
