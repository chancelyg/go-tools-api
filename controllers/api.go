package controllers

import (
	"fmt"
	"go-tools-api/m/dependencies"
	"go-tools-api/m/models"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var ImageList []string

var MaxMemory int64 = 8 << 20
var JsonValidityDays int = 31

var (
	dataStore   = make(map[string]models.AnyJsonData)
	mu          sync.Mutex
	currentSize int64
)

// Handler
func IP(c *gin.Context) {
	c.JSON(200, successResponse("Query success", map[string]interface{}{"ip": c.ClientIP()}))
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

	eTag := fmt.Sprintf("%s_%dx%d", imageData.Id, imageData.Width, imageData.Height)

	// Status not modified
	if imageData.Id != "" {
		if c.GetHeader("If-None-Match") == eTag {
			c.Status(http.StatusNotModified)
			return
		}
	}
	// Cache for 3 months.
	c.Writer.Header().Set("Etag", eTag)
	c.Writer.Header().Set("Cache-Control", "public, max-age=8035200")
	imageMap[eTag] = imageIndex

	var filePath string
	var err error

	if imageData.Height != 0 && imageData.Width != 0 {
		filePath, err = dependencies.ClipImageWithTmpFile(ImageList[imageIndex], imageData.Width, imageData.Height)
		if err != nil {
			c.JSON(200, errorResponse(err.Error()))
			return
		}
		defer os.Remove(filePath)
	} else {
		filePath = ImageList[imageIndex]
	}

	img, err := dependencies.CompressImage(filePath)
	if err != nil {
		c.JSON(200, errorResponse(err.Error()))
		return
	}

	c.Data(http.StatusOK, "image/jpeg", img)
}

func AnyJsonWithGet(c *gin.Context) {
	id := c.Query("id")

	mu.Lock()
	storedData, exists := dataStore[id]
	mu.Unlock()

	if !exists || time.Now().After(storedData.Expiry) {
		c.JSON(200, errorResponse("No data found or data expired"))
		return
	}

	c.JSON(http.StatusOK, storedData.Data)
}

func AnyJsonWithPost(c *gin.Context) {
	id := c.Query("id")

	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(200, errorResponse(err.Error()))
		return
	}

	newSize := dependencies.CalculateSize(jsonData)

	// 检查总数据大小
	mu.Lock()
	if currentSize+newSize > MaxMemory {
		mu.Unlock()
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Total stored data exceeds 128MB"})
		return
	}

	// 存储数据，设置过期时间为365天后
	if existingData, exists := dataStore[id]; exists {
		currentSize -= dependencies.CalculateSize(existingData.Data)
	}
	dataStore[id] = models.AnyJsonData{
		Data:   jsonData,
		Expiry: time.Now().Add(time.Duration(JsonValidityDays) * 24 * time.Hour),
	}
	currentSize += newSize
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": ""})
	c.JSON(200, successResponse("Data stored successfully", nil))
}

func errorResponse(msg string) *models.Response {
	return &models.Response{
		Status:  0,
		Msg:     msg,
		Data:    nil,
		Version: dependencies.Version,
	}
}

func successResponse(msg string, data map[string]interface{}) *models.Response {
	return &models.Response{
		Status:  1,
		Msg:     msg,
		Data:    data,
		Version: dependencies.Version,
	}

}
