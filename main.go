package main

import (
	"flag"
	"fmt"
	"go-tools-api/m/controllers"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
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

	router := gin.Default()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	controllers.ImagePath = cfg.Section("image").Key("path").String()

	// Routes
	router.GET("/rest/api/v1/ip", controllers.IP)

	router.POST("/rest/api/v1/telegram", controllers.Telegram)

	router.GET("/rest/api/v1/image", controllers.Image)

	// Start server
	router.Run(host + ":" + strconv.Itoa(port))
	// r.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", host, port)))
}
