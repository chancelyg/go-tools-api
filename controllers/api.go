package controllers

import (
	"encoding/json"
	"fmt"
	"go-tools-api/m/models"
	"go-tools-api/m/utils"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var ImageList []string

// Handler
func IP(c *gin.Context) {
	c.JSON(200, successResponse("Query success", map[string]interface{}{"ip": c.ClientIP()}))
}

func Telegram(c *gin.Context) {
	var telegram models.TelegramData
	c.ShouldBindJSON(&telegram)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/SendMessage?text=%s&chat_id=%s&parse_mode=MarkdownV2", telegram.ApiKey, telegram.MsgText, telegram.ChatID)
	resp, gErr := http.Get(url)
	if gErr != nil {
		c.JSON(200, errorResponse(gErr.Error()))
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var telegramResponse models.TelegramResponseData
	json.Unmarshal(body, &telegramResponse)
	if !telegramResponse.Ok {
		c.JSON(200, errorResponse(telegramResponse.Description))
		return
	}
	c.JSON(200, successResponse("Send success", telegramResponse.Result))
}

var imageMap = map[string]int{}

func Image(c *gin.Context) {
	if len(ImageList) < 2 {
		c.JSON(200, errorResponse("The number for image collections is to small."))
		return
	}

	var imageData models.ImageData
	if err := c.BindQuery(&imageData); err != nil {
		c.JSON(200, errorResponse(err.Error()))
		return
	}

	maxSize := 4096

	if imageData.Width > maxSize || imageData.Height > maxSize {
		c.JSON(200, errorResponse("The resolution of the image to high."))
		return
	}

	imageIndex := rand.Intn(len(ImageList))

	if imageData.Id != "" {
		c.Writer.Header().Set("Etag", imageData.Id)
		if c.Writer.Header().Get("Etag") == c.GetHeader("If-None-Match") {
			c.Status(http.StatusNotModified)
			return
		}
		// Cache for 3 months.
		c.Writer.Header().Set("Cache-Control", "public, max-age=8035200")
		if value, exist := imageMap[imageData.Id]; exist {
			imageIndex = value
		} else {
			imageMap[imageData.Id] = imageIndex
		}
	}

	var filePath string
	var err error

	if imageData.Height != 0 && imageData.Width != 0 {
		filePath, err = utils.ClipImageWithTmpFile(ImageList[imageIndex], imageData.Width, imageData.Height)
		if err != nil {
			c.JSON(200, errorResponse(err.Error()))
			return
		}
		defer os.Remove(filePath)
	} else {
		filePath = ImageList[imageIndex]
	}

	c.File(filePath)
}

func errorResponse(msg string) *models.Response {
	return &models.Response{
		Status: 0,
		Msg:    msg,
		Data:   nil,
	}
}

func successResponse(msg string, data map[string]interface{}) *models.Response {
	return &models.Response{
		Status: 1,
		Msg:    msg,
		Data:   data,
	}

}
