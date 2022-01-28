package skin

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

var SkinWidths = map[int]int{
	64 * 32 * 4:   64,
	64 * 64 * 4:   64,
	128 * 128 * 4: 128,
}

var SkinHeights = map[int]int{
	64 * 32 * 4:   32,
	64 * 64 * 4:   64,
	128 * 128 * 4: 128,
}

type Skin struct {
	Username string
	Skin     string
}

func S(username string, skin string) *Skin {
	return &Skin{
		Skin:     skin,
		Username: username,
	}
}

func (s *Skin) ConvertToImage() (*image.RGBA, error) {
	data, err := base64.StdEncoding.DecodeString(s.Skin)
	dataString := string(data)

	skinSize := len(dataString)
	skinHeight := SkinHeights[skinSize]
	skinWidth := SkinWidths[skinSize]

	alphaImage := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: skinWidth, Y: skinHeight}})

	position := 0
	for y := 0; y < skinHeight; y++ {
		for x := 0; x < skinWidth; x++ {
			r := dataString[position]
			position++
			g := dataString[position]
			position++
			b := dataString[position]
			position++
			a := dataString[position]
			position++

			alphaImage.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	return alphaImage, err
}

func (s *Skin) SaveFullImage() error {
	workingDirectory, _ := os.Getwd()
	f, err := os.Create(workingDirectory + "\\images\\full\\" + s.Username + ".png")

	defer f.Close()

	w := bufio.NewWriter(f)
	skinImage, _ := s.ConvertToImage()

	err = png.Encode(w, skinImage)
	err = w.Flush()

	return err
}

func (s *Skin) SaveHeadImage() error {
	workingDirectory, _ := os.Getwd()
	f, err := os.Create(workingDirectory + "\\images\\head\\" + s.Username + ".png")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	skinImage, _ := s.ConvertToImage()
	headImage := skinImage.SubImage(image.Rectangle{
		// 7, 8 en 16, 15
		Min: image.Pt(8, 8),
		Max: image.Pt(16, 16),
	})

	err = png.Encode(w, headImage)

	err = w.Flush()

	return err
}

func (s *Skin) FullBytes() ([]byte, error) {

	imageStruct, err := s.ConvertToImage()
	buf := new(bytes.Buffer)
	err = png.Encode(buf, imageStruct)
	base64Image := buf.Bytes()
	return base64Image, err
}

func (s *Skin) HeadBytes() ([]byte, error) {

	imageStruct, err := s.ConvertToImage()

	/*
		data, err := base64.StdEncoding.DecodeString(s.Skin)
		dataString := string(data)
		skinSize := len(dataString)
		skinHeight := SkinHeights[skinSize]
		skinWidth := SkinWidths[skinSize]
	*/

	headImage := imageStruct.SubImage(image.Rectangle{
		// 7, 8 en 16, 15
		Min: image.Pt(8, 8),
		Max: image.Pt(16, 16),
	})

	buf := new(bytes.Buffer)
	err = png.Encode(buf, headImage)
	base64Image := buf.Bytes()

	return base64Image, err
}
