package main

import (
	"flag"
	"fmt"
	"go-tools-api/controllers"
	"log"
	"os"

	"github.com/go-ini/ini"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	h := flag.Bool("h", false, "--help")
	conf := flag.String("c", "app.conf", "app.conf path")
	flag.Parse()
	if *h {
		flag.Usage()
	}
	_, err := os.Stat(*conf)
	if err != nil {
		fmt.Println("app.conf not found")
		os.Exit(-1)
	}

	cfg, err := ini.Load(*conf)
	if err != nil {
		log.Fatal("Fail to read file: ", err)
		os.Exit(-1)
	}

	host := cfg.Section("general").Key("host").String()
	port, err := cfg.Section("general").Key("port").Int()
	if err != nil {
		log.Fatal("Read conf error: ", err)
		os.Exit(-1)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	fileWriter, err := os.OpenFile("go-tools-api-running.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file go-tools-api-running.log failed:%v", err.Error())
	}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} - ${remote_ip} - ${method} - ${status} - ${uri}\n",
		Output: fileWriter,
	}))
	e.Use(middleware.Recover())

	controllers.ImagePath = cfg.Section("image").Key("path").String()

	// Routes
	e.GET("/rest/api/v1/ip", controllers.IP)

	e.POST("/rest/api/v1/telegram", controllers.Telegram)

	e.GET("/rest/api/v1/image", controllers.Image)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", host, port)))
}
