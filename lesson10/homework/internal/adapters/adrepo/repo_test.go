package adrepo

import (
	"github.com/stretchr/testify/assert"
	"homework9/internal/ads"
	"testing"
)

func TestRepo(t *testing.T) {
	repo := NewAdRepo()
	assert.NotNil(t, repo)

	// test only adRepo implementation because Repository is template struct

	userRepo := NewUserRepo()
	assert.NotNil(t, userRepo)

	id := repo.GetCurAvailableId()
	assert.Zero(t, id)

	ad := ads.Ad{
		ID:    id,
		Title: "Hello",
		Text:  "world",
	}

	repo.Insert(&ad)
	adFromRepo, contain := repo.Get(0)
	assert.True(t, contain)
	assert.Equal(t, *adFromRepo, ad)

	adFromRepo, contain = repo.Get(1)
	assert.False(t, contain)
	assert.Nil(t, adFromRepo)
}
