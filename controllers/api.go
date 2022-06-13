package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"go-tools-api/models"

	"github.com/labstack/echo/v4"
)

var ImagePath string

// Handler
func IP(c echo.Context) error {
	r := &models.Response{
		Status: 1,
		Msg:    "查询成功",
		Data: map[string]interface{}{
			"ip": c.RealIP(),
		},
	}
	return c.JSON(http.StatusOK, r)
}

func Telegram(c echo.Context) (err error) {
	t := new(models.Telegram)
	if err = c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/SendMessage?text=%s&chat_id=%s&parse_mode=MarkdownV2", t.ApiKey, t.MsgText, t.ChatID)
	resp, gErr := http.Get(url)
	if gErr != nil {
		r := &models.Response{
			Status: 0,
			Msg:    gErr.Error(),
			Data:   nil,
		}
		return c.JSON(http.StatusOK, r)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var telegramResponse models.TelegramResponse
	json.Unmarshal(body, &telegramResponse)
	if !telegramResponse.Ok {
		r := &models.Response{
			Status: 0,
			Msg:    telegramResponse.Description,
			Data:   telegramResponse.Result,
		}
		return c.JSON(http.StatusOK, r)
	}
	r := &models.Response{
		Status: 1,
		Msg:    "消息发送成功",
		Data:   telegramResponse.Result,
	}
	return c.JSON(http.StatusOK, r)
}

var imageIndex int

func Image(c echo.Context) error {
	var files []string

	err := filepath.Walk(ImagePath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil || len(files) < 2 {
		r := &models.Response{
			Status: 0,
			Msg:    "Get image error",
			Data:   nil,
		}
		return c.JSON(http.StatusOK, r)
	}
	for {
		randomIndex := rand.Intn(len(files) - 1)
		if randomIndex != imageIndex {
			imageIndex = randomIndex
			break
		}
	}
	return c.File(files[imageIndex])

}
