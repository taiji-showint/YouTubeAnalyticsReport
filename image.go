package main

import (
	"image"
	"image/jpeg" // JPEGを読み書きする場合
	"os"
	"strings"
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
	// クロップのみ（アスペクト比を保持）
	// 左上(left, top)から右下(left+width, top+height)までの範囲を抽出
	type subImageIface interface {
		SubImage(r image.Rectangle) image.Image
	}

	if subImg, ok := img.(subImageIface); ok {
		return subImg.SubImage(image.Rect(left, top, left+width, top+height))
	}

	// SubImage に対応していない場合は RGBA にコピー
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height && top+y < img.Bounds().Max.Y; y++ {
		for x := 0; x < width && left+x < img.Bounds().Max.X; x++ {
			newImage.Set(x, y, img.At(left+x, top+y))
		}
	}
	return newImage
}

func trim_YT_Thumbnail(filename string) error {
	img, err := loadImage(filename)
	if err != nil {
		return err
	}

	// 画像サイズを取得
	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X

	// アスペクト比 16:9 を保持して trim
	// 高さから幅を計算（16:9）
	targetHeight := 315
	targetWidth := targetHeight * 16 / 9 // = 560

	// 左右中央からトリムするために、開始位置を計算
	startX := (imgWidth - targetWidth) / 2
	if startX < 0 {
		startX = 0
	}
	startY := 45

	img_trim := trimImage(img, startY, startX, targetWidth, targetHeight)
	filename_trim := strings.Replace(filename, ".jpg", "", -1) + "_trim.jpg"
	err = saveImage(filename_trim, img_trim)
	return err
}
