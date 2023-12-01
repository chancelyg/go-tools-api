package dependencies

import (
	"bytes"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

func CompressImage(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	// Resize to width 1000 while preserving aspect ratio.
	m := resize.Resize(1000, 0, img, resize.Lanczos3)

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, m, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ClipImageWithTmpFile(sourcePath string, width, height int) (string, error) {
	// 读取源文件
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	srcImg, err := jpeg.Decode(srcFile)
	if err != nil {
		return "", err
	}

	// 裁剪图像
	dstImg := resize.Resize(uint(width), uint(height), srcImg, resize.Lanczos3)

	// 创建临时文件
	dstFile, err := os.CreateTemp("", "image-*.jpeg")
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	err = jpeg.Encode(dstFile, dstImg, nil)
	if err != nil {
		os.Remove(dstFile.Name())
		return "", err
	}

	// 返回临时文件的路径
	return dstFile.Name(), nil
}
