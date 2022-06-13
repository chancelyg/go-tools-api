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
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var ImagePath string

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

var imageIndex int

func Image(c *gin.Context) {
	var files []string

	err := filepath.Walk(ImagePath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil || len(files) < 2 {
		c.JSON(200, errorResponse("Image count < 2"))
		return
	}

	for {
		randomIndex := rand.Intn(len(files) - 1)
		if randomIndex != imageIndex {
			imageIndex = randomIndex
			break
		}
	}
	var imageData models.ImageData
	c.BindQuery(&imageData)
	if imageData.Height != 0 && imageData.Width != 0 {
		filePath, err := utils.ClipImageWithTmpFile(files[imageIndex], imageData.Width, imageData.Height)
		defer os.Remove(filePath)
		if err != nil {
			c.JSON(200, errorResponse(err.Error()))
			return
		}
		c.File(filePath)
		return
	}
	c.File(files[imageIndex])

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
