package main

import (
	"flag"
	"fmt"
	"go-tools-api/m/controllers"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	host := flag.String("host", "localhost", "liscen host")
	port := flag.Int("port", 8085, "liscen port")
	imageDir := flag.String("image", "./images", "image folder(jpg type)")
	jsonValidityDays := flag.Int("jsonValidityDays", 31, "json validity days")
	jsonMaxMemorySize := flag.Int64("jsonMaxMemorySize", 32, "json max memory size")

	help := flag.Bool("h", false, "show help")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	fmt.Printf("Host: %s\n", *host)
	fmt.Printf("Port: %d\n", *port)
	fmt.Printf("Image Directory: %s\n", *imageDir)
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

	var files []string
	err := filepath.Walk(*imageDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil || len(files) < 2 {
		log.Fatal(err)
	}

	controllers.ImageList = files
	controllers.MaxMemory = *jsonMaxMemorySize << 20
	controllers.JsonValidityDays = *jsonValidityDays

	// Routes
	router.GET("/rest/api/v1/ip", controllers.IP)
	router.GET("/rest/api/v1/image", controllers.Image)
	router.GET("/rest/api/v1/anyjson", controllers.AnyJsonWithGet)
	router.POST("/rest/api/v1/anyjson", controllers.AnyJsonWithPost)

	// Start server
	router.Run(*host + ":" + strconv.Itoa(*port))
	// r.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", host, port)))
}
