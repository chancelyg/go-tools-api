package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Response struct {
	Status int                    `json:"status"`
	Msg    string                 `json:"msg"`
	Data   map[string]interface{} `json:"data"`
}

type Telegram struct {
	ApiKey  string `json:"apiKey"`
	ChatID  string `json:"chatId"`
	MsgText string `json:"msgText"`
}

type TelegramResponse struct {
	Ok          bool                   `json:"ok"`
	Result      map[string]interface{} `json:"result"`
	ErrorCode   int                    `json:"error_code"`
	Description string                 `json:"description"`
}

var (
	h    bool
	host string
	port int
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&host, "host", "127.0.0.1", "listen host")
	flag.IntVar(&port, "port", 1323, "listen port")
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
	}
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	fileWriter, err := os.OpenFile("ta.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file ta.log failed:%v", err.Error())
	}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} - ${remote_ip} - ${method} - ${status} - ${uri}\n",
		Output: fileWriter,
	}))
	e.Use(middleware.Recover())

	// Routes
	e.GET("/rest/api/v1/ip", ip)

	e.POST("/rest/api/v1/telegram", telegram)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", host, port)))
}

// Handler
func ip(c echo.Context) error {
	r := &Response{
		Status: 1,
		Msg:    "查询成功",
		Data: map[string]interface{}{
			"ip": c.RealIP(),
		},
	}
	return c.JSON(http.StatusOK, r)
}

func telegram(c echo.Context) (err error) {
	t := new(Telegram)
	if err = c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/SendMessage?text=%s&chat_id=%s&parse_mode=MarkdownV2", t.ApiKey, t.MsgText, t.ChatID)
	resp, gErr := http.Get(url)
	if gErr != nil {
		r := &Response{
			Status: 0,
			Msg:    gErr.Error(),
			Data:   nil,
		}
		return c.JSON(http.StatusOK, r)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var telegramResponse TelegramResponse
	json.Unmarshal(body, &telegramResponse)
	if !telegramResponse.Ok {
		r := &Response{
			Status: 0,
			Msg:    telegramResponse.Description,
			Data:   telegramResponse.Result,
		}
		return c.JSON(http.StatusOK, r)
	}
	r := &Response{
		Status: 1,
		Msg:    "消息发送成功",
		Data:   telegramResponse.Result,
	}
	return c.JSON(http.StatusOK, r)

}
