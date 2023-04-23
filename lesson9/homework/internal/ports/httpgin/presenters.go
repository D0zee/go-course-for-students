package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/ads"
	"homework9/internal/users"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	AuthorID     int64     `json:"author_id"`
	Published    bool      `json:"published"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type getAdRequest struct {
	UserID int64 `json:"user_id"`
}

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type changeUserRequest struct {
	Data string `json:"data"`
}

type filterQueryRequest struct {
	AuthorID int64     `form:"author_id,default=-1"`
	Time     time.Time `form:"time" time_format:"2006-01-02"`
}

type getAdsByTitleRequest struct {
	Title string `json:"title"`
}

func AdSuccessResponse(ad ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreationTime: ad.CreationTime,
			UpdateTime:   ad.UpdateTime,
		},
		"error": nil,
	}
}

func AdListSuccessResponse(ads []ads.Ad) *gin.H {
	var responses []adResponse
	for _, ad := range ads {
		responses = append(responses, adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreationTime: ad.CreationTime,
			UpdateTime:   ad.UpdateTime,
		})
	}
	return &gin.H{
		"data": responses,
	}
}

func ErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func UserSuccessResponse(u users.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			ID:       u.Id,
			Nickname: u.Nickname,
			Email:    u.Email,
		},
		"error": nil,
	}
}
