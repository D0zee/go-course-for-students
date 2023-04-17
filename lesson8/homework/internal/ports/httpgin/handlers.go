package httpgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"homework8/internal/ads"
	//"homework8/internal/users"
	"log"
	"net/http"
	"strconv"

	"homework8/internal/app"
)

var ErrEmptyQueryParam = errors.New("error param is empty")

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.ShouldBindJSON(&reqBody)
		log.Println("BODY:", reqBody)

		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			log.Println("error with parse json:", err)
			return
		}

		ad, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)
		if errors.Is(err, app.ErrValidate) {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			log.Println("error validation:", err)
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			log.Println("error response IDK:", err)
			return
		}
		log.Println("AD:", ad)
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		if c.Param("id") == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse(ErrEmptyQueryParam))
			return
		}

		adID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, err := a.ChangeAdStatus(c, int64(adID), reqBody.UserID, reqBody.Published)
		log.Println("AD IN CHANGE STATUS:", ad)
		if errors.Is(err, app.ErrAccess) {
			c.JSON(http.StatusForbidden, ErrorResponse(err))
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		if c.Param("id") == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse(ErrEmptyQueryParam))
			return
		}

		adID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, err := a.UpdateAd(c, int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if err != nil {
			if errors.Is(err, app.ErrValidate) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
				return
			}

			if errors.Is(err, app.ErrAccess) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}

			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func getAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Param("id") == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse(ErrEmptyQueryParam))
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		var req getAdRequest
		if err = c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, err := a.GetAdById(c, int64(id), req.UserID)
		if err != nil {
			if errors.Is(err, app.ErrAccess) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
			if errors.Is(err, app.ErrAvailabilityAd) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// CreateUser - Method for creating user
func CreateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		user, err := a.CreateUser(c, req.Nickname, req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

// ChangeUser - Method for changing different fields of user structure
func ChangeUser(a app.App, m app.Method) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		if c.Param("id") == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse(ErrEmptyQueryParam))
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		user, err := a.UpdateUser(c, int64(id), req.Data, m)

		if err != nil {
			if errors.Is(err, app.ErrWrongUserId) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

type adPredicate func(ad ads.Ad) bool

func filter(Ads []ads.Ad, p adPredicate) []ads.Ad {
	var result []ads.Ad
	for _, ad := range Ads {
		if p(ad) {
			result = append(result, ad)
		}
	}
	return result
}

func getListAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		listAds := a.ListAds(c)
		var f adPredicate

		var queryParam filterQueryRequest

		if err := c.ShouldBindQuery(&queryParam); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		author := queryParam.AuthorID
		if author != -1 {
			f = func(ad ads.Ad) bool {
				return ad.AuthorID == author
			}
			listAds = filter(listAds, f)
		}

		date := queryParam.Time
		log.Println("AUTHOR:", author, " ", c.Query("author"))
		log.Println("TIME:", date, date.IsZero())

		if !date.IsZero() {
			f = func(ad ads.Ad) bool {
				adTime := ad.CreationTime
				return date.Day() == adTime.Day() &&
					date.Month() == adTime.Month() &&
					date.Year() == adTime.Year()
			}
			listAds = filter(listAds, f)
		}

		f = func(ad ads.Ad) bool {
			return ad.Published
		}
		listAds = filter(listAds, f)
		c.JSON(http.StatusOK, AdListSuccessResponse(listAds))
	}
}
