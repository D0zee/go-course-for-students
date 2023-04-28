package app

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"homework9/internal/adapters/repo"
	"homework9/internal/ads"
	"homework9/internal/users"
	"strings"
	"testing"
	"time"
)

type UserSuit struct {
	suite.Suite
	ctx      context.Context
	app      App
	userRepo repo.Repository[users.User]
	dftUser  users.User

	nickname string
	email    string
}

func (s *UserSuit) SetupTest() {
	s.ctx = context.Background()
	s.userRepo = repo.NewUserRepo()
	s.app = NewApp(repo.NewMapAdRepo(), s.userRepo)
	s.nickname = "aboba"
	s.nickname = "aboba"

	user, _ := s.app.CreateUser(s.ctx, s.nickname, s.email)
	s.dftUser = user
}

func (s *UserSuit) TestCreateUser() {
	userExpected := users.User{
		Id:       s.userRepo.GetCurAvailableId(),
		Nickname: s.nickname,
		Email:    s.email,
		Deleted:  false,
	}

	user, err := s.app.CreateUser(s.ctx, s.nickname, s.email)
	s.Nil(err)
	s.Equal(user, userExpected)
}

func (s *UserSuit) TestUpdateUser() {
	newNickname := "abobaNew"
	newEmail := "mailNew"

	updUser, err := s.app.UpdateUser(s.ctx, s.dftUser.Id, newNickname, ChangeNickname)
	s.Nil(err)
	s.dftUser.Nickname = newNickname
	s.Equal(updUser, s.dftUser)

	updUser, err = s.app.UpdateUser(s.ctx, s.dftUser.Id, newEmail, ChangeEmail)
	s.Nil(err)
	s.dftUser.Email = newEmail
	s.Equal(updUser, s.dftUser)

	emtUser, err := s.app.UpdateUser(s.ctx, 1000, newEmail, ChangeEmail)
	s.NotNil(err)
	s.Empty(emtUser)
}

func (s *UserSuit) TestGetUser() {
	user, err := s.app.GetUser(s.ctx, s.dftUser.Id)
	s.Nil(err)
	s.Equal(user, s.dftUser)

	emtUser, err := s.app.GetUser(s.ctx, 100)
	s.NotNil(err)
	s.Empty(emtUser)

	_, _ = s.app.RemoveUser(s.ctx, s.dftUser.Id)
	emtUser, err = s.app.GetUser(s.ctx, s.dftUser.Id)
	s.NotNil(err)
	s.Empty(emtUser)
}

func (s *UserSuit) TestWithCancelledCtx() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.app.CreateUser(ctx, s.nickname, s.email)
	s.NotNil(err)

	_, err = s.app.UpdateUser(ctx, s.dftUser.Id, s.nickname, ChangeNickname)
	s.NotNil(err)

	_, err = s.app.GetUser(ctx, s.dftUser.Id)
	s.NotNil(err)

	_, err = s.app.RemoveUser(ctx, s.dftUser.Id)
	s.NotNil(err)
}

func (s *UserSuit) TestRemoveUser() {
	user, err := s.app.RemoveUser(s.ctx, s.dftUser.Id)
	s.dftUser.Deleted = true
	s.Nil(err)
	s.Equal(s.dftUser, user)

	emtUser, err := s.app.RemoveUser(s.ctx, 1000)
	s.NotNil(err)
	s.Empty(emtUser)
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserSuit))
}

type adSuite struct {
	suite.Suite
	ctx     context.Context
	app     App
	adRepo  repo.Repository[ads.Ad]
	dftUser users.User
	dftAd   ads.Ad

	dftTitle string
	dftText  string
}

func (s *adSuite) SetupTest() {
	s.ctx = context.Background()
	s.adRepo = repo.NewMapAdRepo()
	s.app = NewApp(s.adRepo, repo.NewUserRepo())
	s.dftTitle = "hello"
	s.dftText = "world"

	nickname := "Иван"
	email := "hello@tinkoff"
	user, _ := s.app.CreateUser(s.ctx, nickname, email)
	s.dftUser = user

	ad, _ := s.app.CreateAd(s.ctx, s.dftTitle, s.dftText, s.dftUser.Id)
	s.dftAd = ad
}

func isSameDate(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func isSameAd(ad1, ad2 ads.Ad) {
	isSame := ad1.ID == ad2.ID && ad1.Title == ad2.Title &&
		ad1.Text == ad2.Text &&
		ad1.Published == ad2.Published &&
		ad1.Deleted == ad2.Deleted &&
		isSameDate(ad1.UpdateTime, ad2.UpdateTime) &&
		isSameDate(ad1.CreationTime, ad2.CreationTime)
	if !isSame {
		panic("not same ad")
	}
}

func (s *adSuite) TestAdCreate() {
	expAd := ads.Ad{
		ID:           s.adRepo.GetCurAvailableId(),
		Title:        s.dftTitle,
		Text:         s.dftText,
		AuthorID:     s.dftUser.Id,
		Published:    false,
		CreationTime: time.Now().UTC(),
		UpdateTime:   time.Now().UTC(),
	}

	ad, err := s.app.CreateAd(s.ctx, s.dftTitle, s.dftText, s.dftUser.Id)
	fmt.Println(ad)

	fmt.Println(expAd)
	s.Nil(err)
	isSameAd(ad, expAd)

	// creating ad from non-existent user
	ad, err = s.app.CreateAd(s.ctx, s.dftTitle, s.dftText, 100)
	s.NotNil(err)
	s.Empty(ad)

	// too long length of Title (0 < len <= 100)
	tooLongTitle := strings.Repeat("a", 101)
	emptyTitle := ""
	ad, err = s.app.CreateAd(s.ctx, emptyTitle, s.dftText, s.dftUser.Id)
	s.NotNil(err)
	s.Empty(ad)

	ad, err = s.app.CreateAd(s.ctx, tooLongTitle, s.dftText, s.dftUser.Id)
	s.NotNil(err)
	s.Empty(ad)

	// wrong length of Text
	tooLongText := strings.Repeat("a", 501)
	emptyText := ""
	ad, err = s.app.CreateAd(s.ctx, s.dftTitle, tooLongText, s.dftUser.Id)
	s.NotNil(err)
	s.Empty(ad)

	ad, err = s.app.CreateAd(s.ctx, s.dftTitle, emptyText, s.dftUser.Id)
	s.NotNil(err)
	s.Empty(ad)
}

func (s *adSuite) TestChangeAdStatus() {
	ad, err := s.app.ChangeAdStatus(s.ctx, s.dftAd.ID, s.dftUser.Id, true)
	s.Nil(err)
	s.dftAd.Published = true
	isSameAd(s.dftAd, ad)
}

func (s *adSuite) TestUpdateAd() {
	// change only in environment
	s.dftAd.Text = "new Text"
	s.dftAd.Title = "new Title"

	ad, err := s.app.UpdateAd(s.ctx, s.dftAd.ID, s.dftUser.Id, s.dftAd.Title, s.dftAd.Text)
	s.Nil(err)
	isSameAd(ad, s.dftAd)

	_, err = s.app.UpdateAd(s.ctx, s.dftAd.ID, 100, s.dftAd.Title, s.dftAd.Text)
	s.NotNil(err)

	// validation error
	emptyTitle := ""
	_, err = s.app.UpdateAd(s.ctx, s.dftAd.ID, s.dftUser.Id, emptyTitle, s.dftAd.Text)
	s.NotNil(err)
}

func (s *adSuite) TestGetAdById() {
	ad, err := s.app.GetAdById(s.ctx, s.dftAd.ID, s.dftUser.Id)
	s.Nil(err)
	isSameAd(ad, s.dftAd)
}

func (s *adSuite) TestRemoveAd() {
	ad, err := s.app.RemoveAd(s.ctx, s.dftAd.ID, s.dftUser.Id)
	s.Nil(err)
	s.dftAd.Deleted = true
	isSameAd(ad, s.dftAd)
}

func (s *adSuite) TestListAds() {
	anotherAd, err := s.app.CreateAd(s.ctx, s.dftTitle, s.dftText, s.dftUser.Id)
	s.Nil(err)

	ads := s.app.ListAds(s.ctx)
	isSameAd(ads[0], s.dftAd)
	isSameAd(ads[1], anotherAd)
}

func (s *adSuite) TestWithNonValidCtx() {
	nonValidCtx, cancel := context.WithCancel(s.ctx)
	cancel()

	_, err := s.app.CreateAd(nonValidCtx, s.dftTitle, s.dftText, s.dftUser.Id)
	s.NotNil(err)

	_, err = s.app.ChangeAdStatus(nonValidCtx, s.dftAd.ID, s.dftUser.Id, true)
	s.NotNil(err)

	_, err = s.app.UpdateAd(nonValidCtx, s.dftAd.ID, s.dftUser.Id, s.dftAd.Title, s.dftAd.Text)
	s.NotNil(err)

	_, err = s.app.GetAdById(nonValidCtx, s.dftAd.ID, s.dftUser.Id)
	s.NotNil(err)

	_, err = s.app.RemoveAd(nonValidCtx, s.dftAd.ID, s.dftUser.Id)
	s.NotNil(err)

	emptyAds := s.app.ListAds(nonValidCtx)
	s.Empty(emptyAds)
}

func (s *adSuite) TestAccessOnly() {
	// request from non-existence user
	empty, err := s.app.ChangeAdStatus(s.ctx, s.dftAd.ID, 100, true)
	s.NotNil(err)
	s.Empty(empty)

	// non-existence ad
	empty, err = s.app.ChangeAdStatus(s.ctx, 1000, s.dftUser.Id, true)
	s.NotNil(err)
	s.Empty(empty)

	// request on unavailable ad
	ad, err := s.app.RemoveAd(s.ctx, 100, s.dftUser.Id)
	s.NotNil(err)
	s.Empty(ad)

	// request from removed user
	_, err = s.app.RemoveUser(s.ctx, s.dftUser.Id)
	s.Nil(err)
	ad, err = s.app.ChangeAdStatus(s.ctx, s.dftAd.ID, s.dftUser.Id, true)
	s.NotNil(err)
	s.Empty(ad)

	// request from another user
	anotherUser, err := s.app.CreateUser(s.ctx, "Kolya", "email#tinkoff")
	s.Nil(err)
	ad, err = s.app.GetAdById(s.ctx, s.dftAd.ID, anotherUser.Id)
	s.Empty(ad)
	s.NotNil(err)
}

func TestAd(t *testing.T) {
	suite.Run(t, new(adSuite))
}

// two benchmarks which compare repo with map and repo with slice implementation

func BenchmarkCreateAdWithMapRepo(b *testing.B) {
	mr := repo.NewMapAdRepo()
	a := NewApp(mr, repo.NewUserRepo())
	ctx := context.Background()
	user, _ := a.CreateUser(ctx, "Иван", "tinkoff@com")

	for i := 0; i < b.N; i++ {
		_, _ = a.CreateAd(ctx, "hello", "world", user.Id)
	}
}

func BenchmarkCreateAdWithSliceRepo(b *testing.B) {
	mr := repo.NewSliceAdRepo()
	a := NewApp(mr, repo.NewUserRepo())
	ctx := context.Background()
	user, _ := a.CreateUser(ctx, "Иван", "tinkoff@com")

	for i := 0; i < b.N; i++ {
		_, _ = a.CreateAd(ctx, "hello", "world", user.Id)
	}
}
