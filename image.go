package main

import (
	"image"
	"image/jpeg" // JPEGを読み書きする場合
	"os"
	"strings"

	"golang.org/x/image/draw"
)

// 画像を読み取るための関数。
// ファイルパスを指定すると、画像データを返してくれる。
func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// 画像を保存する関数。
// 保存先のパスと画像データを渡すと保存してくれる。
func saveImage(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = jpeg.Encode(f, img, &jpeg.Options{
		Quality: 80, // JPEGのクオリティ設定。省略するとjpeg.DefaultQualityの値（75）が使われる。
	})
	return err
}

func trimImage(img image.Image, top, left, width, height int) image.Image {
	// 新しい画像を用意する
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))

	// 左上(left, top)から右下(left+width, top+height)までの範囲を、新しい画像にコピーする
	draw.BiLinear.Scale(newImage, newImage.Bounds(), img, image.Rect(left, top, left+width, top+height), draw.Over, nil)

	return newImage
}

func trim_YT_Thumbnail(filename string) error {
	img, err := loadImage(filename)
	if err != nil {
		return err
	}
	img_trim := trimImage(img, 45, 0, 480, 315)
	filename_trim := strings.Replace(filename, ".jpg", "", -1) + "_trim.jpg"
	err = saveImage(filename_trim, img_trim)
	return err
}
