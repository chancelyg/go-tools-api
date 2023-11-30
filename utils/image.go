package utils

import (
	"bytes"
	"io"
	"os"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

func ClipImageWithTmpFile(sourcePath string, width, height int) (string, error) {
	// 读取源文件
	srcFile, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", err
	}

	// 解码 WebP 图像
	srcImg, err := webp.DecodeRGB(srcFile)
	if err != nil {
		return "", err
	}

	// 裁剪图像
	dstImg := resize.Resize(uint(width), uint(height), srcImg, resize.Lanczos3)

	// 创建临时文件
	dstFile, err := os.CreateTemp("", "image-*.webp")
	if err != nil {
		return "", err
	}

	// 将裁剪后的图像编码为 WebP 格式，然后写入临时文件
	dstData, err := webp.EncodeRGB(dstImg, 90)
	if err != nil {
		dstFile.Close()
		os.Remove(dstFile.Name())
		return "", err
	}
	_, err = io.Copy(dstFile, bytes.NewReader(dstData))
	if err != nil {
		dstFile.Close()
		os.Remove(dstFile.Name())
		return "", err
	}

	// 返回临时文件的路径
	return dstFile.Name(), nil
}
