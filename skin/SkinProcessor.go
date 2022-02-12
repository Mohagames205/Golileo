package skin

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

var Widths = map[int]int{
	64 * 32 * 4:   64,
	64 * 64 * 4:   64,
	128 * 128 * 4: 128,
}

var Heights = map[int]int{
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
	skinHeight := Heights[skinSize]
	skinWidth := Widths[skinSize]

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

func (s *Skin) SaveFullImage() (string, error) {
	workingDirectory, _ := os.Getwd()
	uuid := pseudo_uuid()
	f, err := os.Create(workingDirectory + "/images/" + uuid + ".png")

	defer f.Close()

	w := bufio.NewWriter(f)
	skinImage, _ := s.ConvertToImage()

	err = png.Encode(w, skinImage)
	err = w.Flush()

	return uuid, err
}

func (s *Skin) SaveHeadImage() (string, error) {
	workingDirectory, _ := os.Getwd()
	uuid := pseudo_uuid()
	f, err := os.Create(workingDirectory + "/images/" + uuid + ".png")

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

	return uuid, err
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

func pseudo_uuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
