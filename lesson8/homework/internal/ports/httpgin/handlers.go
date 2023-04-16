package httpgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"

	"homework8/internal/app"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.ShouldBind(&reqBody)
		log.Println("BODY:", reqBody)

		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			log.Fatal(err)
		}

		ad, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)
		if errors.Is(err, app.ErrValidate) {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			log.Fatal(err)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			log.Fatal(err)
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
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			log.Fatal(err)
		}

		adID, err := strconv.Atoi(c.Param("id"))
		if c.Param("id") == "" || err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			log.Fatal(err)
		}

		ad, err := a.ChangeAdStatus(c, int64(adID), reqBody.UserID, reqBody.Published)
		log.Println("AD IN CHANGE STATUS:", ad)
		if errors.Is(err, app.ErrAccess) {
			c.JSON(http.StatusForbidden, AdErrorResponse(err))
			log.Fatal(err)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("id"))
		if c.Param("id") == "" || err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			log.Fatal(err)
		}

		ad, err := a.UpdateAd(c, int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if err != nil {
			if errors.Is(err, app.ErrValidate) {
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
				return
			}

			if errors.Is(err, app.ErrAccess) {
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
				return
			}

			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}
