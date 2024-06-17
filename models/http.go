package models

import "time"

type Response struct {
	Status  int                    `json:"status"`
	Msg     string                 `json:"msg"`
	Data    map[string]interface{} `json:"data"`
	Version string                 `json:"version"`
}

type TelegramData struct {
	ApiKey  string `json:"apiKey"`
	ChatID  string `json:"chatId"`
	MsgText string `json:"msgText"`
}

type TelegramResponseData struct {
	Ok          bool                   `json:"ok"`
	Result      map[string]interface{} `json:"result"`
	ErrorCode   int                    `json:"error_code"`
	Description string                 `json:"description"`
}

type ImageData struct {
	Id     string `form:"id"`
	Width  int    `form:"width"`
	Height int    `form:"height"`
}

type AnyJsonData struct {
	Data   map[string]interface{}
	Expiry time.Time
}
