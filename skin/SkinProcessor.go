package skin

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
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

type Dimension struct {
	Width  int
	Height int
}

type Skin struct {
	Username   string
	Skin       string
	Dimensions Dimension
}

func S(username string, skin string) *Skin {
	data, _ := base64.StdEncoding.DecodeString(skin)
	skinSize := len(string(data))
	skinHeight := Heights[skinSize]
	skinWidth := Widths[skinSize]

	return &Skin{
		Skin:       skin,
		Username:   username,
		Dimensions: Dimension{skinWidth, skinHeight},
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

func (s *Skin) SaveFullImage(name string) (string, error) {
	workingDirectory, _ := os.Getwd()
	f, err := os.Create(workingDirectory + "/public/images/" + s.Username + "-" + name + ".png")

	defer f.Close()

	w := bufio.NewWriter(f)
	skinImage, _ := s.ConvertToImage()

	err = png.Encode(w, skinImage)
	err = w.Flush()

	return workingDirectory + "/public/images/" + s.Username + "-" + name + ".png", err
}

func (s *Skin) SaveHeadImage(name string) (string, error) {
	workingDirectory, _ := os.Getwd()
	f, err := os.Create(workingDirectory + "/public/images/" + s.Username + "-" + name + ".png")

	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(f)
	skinImage, _ := s.ConvertToImage()

	minP := s.Dimensions.Width / 8
	maxP := s.Dimensions.Width / 4

	headHeight := s.Dimensions.Width / 8

	overlayBeginPointX := int(float32(s.Dimensions.Width) * (5.0 / 8.0))
	overlayEndPointX := overlayBeginPointX + headHeight

	headImage := skinImage.SubImage(image.Rectangle{
		Min: image.Pt(minP, minP),
		Max: image.Pt(maxP, maxP),
	})

	overlayImage := skinImage.SubImage(image.Rectangle{
		Min: image.Pt(overlayBeginPointX, minP),
		Max: image.Pt(overlayEndPointX, minP+headHeight),
	})

	py := minP
	for y := minP; y < minP+headHeight; y++ {
		px := minP
		for x := overlayBeginPointX; x < overlayEndPointX; x++ {
			overlayColor := overlayImage.At(x, y)
			_, _, _, a := overlayColor.RGBA()
			if a == 0 {
				px++
				continue
			}
			headImage.(draw.Image).Set(px, py, overlayColor)
			px++
		}
		py++
	}

	err = png.Encode(w, headImage)

	err = w.Flush()

	f.Close()

	return workingDirectory + "/public/images/" + s.Username + "-" + name + ".png", err
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

	minP := s.Dimensions.Width / 8
	maxP := s.Dimensions.Width / 4

	headImage := imageStruct.SubImage(image.Rectangle{
		Min: image.Pt(minP, minP),
		Max: image.Pt(maxP, maxP),
	})

	buf := new(bytes.Buffer)
	err = png.Encode(buf, headImage)
	base64Image := buf.Bytes()

	return base64Image, err
}

func PseudoUuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
