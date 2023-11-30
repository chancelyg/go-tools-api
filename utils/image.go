package utils

import (
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

func ClipImageWithTmpFile(sourcePath string, width, height int) (string, error) {
	// 读取源文件
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	// 解码 PNG 图像
	srcImg, err := png.Decode(srcFile)
	if err != nil {
		return "", err
	}

	// 裁剪图像
	dstImg := resize.Resize(uint(width), uint(height), srcImg, resize.Lanczos3)

	// 创建临时文件
	dstFile, err := os.CreateTemp("", "image-*.png")
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	// 将裁剪后的图像编码为 PNG 格式，然后写入临时文件
	err = png.Encode(dstFile, dstImg)
	if err != nil {
		os.Remove(dstFile.Name())
		return "", err
	}

	// 返回临时文件的路径
	return dstFile.Name(), nil
}
