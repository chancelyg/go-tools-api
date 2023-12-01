package controllers

import (
	"fmt"
	"go-tools-api/m/dependencies"
	"go-tools-api/m/models"
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
