package repo

import (
	"github.com/stretchr/testify/assert"
	"homework9/internal/ads"
	"testing"
)

type TestCase struct {
	repo Repository[ads.Ad]
}

func TestRepo(t *testing.T) {
	testCases := []TestCase{
		{NewMapAdRepo()},
		{NewSliceAdRepo()},
	}

	for _, test := range testCases {
		repo := test.repo
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

}

func FuzzRepo(f *testing.F) {

	// stress test for different degrees of load on service
	testcases := []int64{10, 100, 5000, 1000000}

	for _, tc := range testcases {
		f.Add(tc)
	}
	repo := NewMapAdRepo()

	go func() {
		for {
			obj := &ads.Ad{
				ID: repo.GetCurAvailableId(),
			}
			repo.Insert(obj)
		}
	}()

	f.Fuzz(func(t *testing.T, id int64) {
		_, contain := repo.Get(id)
		if contain && repo.GetCurAvailableId() <= id ||
			!contain && repo.GetCurAvailableId() > id {
			t.Fatalf("slice don't work right")
		}

	})
}
