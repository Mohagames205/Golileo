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

func (s *Skin) SaveFullImage() (string, error) {
	workingDirectory, _ := os.Getwd()
	f, err := os.Create(workingDirectory + "/public/images/" + s.Username + "-full.png")

	defer f.Close()

	w := bufio.NewWriter(f)
	skinImage, _ := s.ConvertToImage()

	err = png.Encode(w, skinImage)
	err = w.Flush()

	return workingDirectory + "/public/images/" + s.Username + "-full.png", err
}

func (s *Skin) SaveFullFrontImage() (string, error) {
	workingDirectory, _ := os.Getwd()
	f, err := os.Create(workingDirectory + "/public/images/" + s.Username + "-full_front.png")

	defer f.Close()

	w := bufio.NewWriter(f)
	skinImage, _ := s.ConvertToImage()

	fmt.Println(skinImage.Bounds())

	legLeftFront := image.Rectangle{
		Min: image.Pt(20, 52),
		Max: image.Pt(24, 64),
	}

	legRightFront := image.Rectangle{
		Min: image.Pt(4, 20),
		Max: image.Pt(8, 32),
	}

	armLeftFront := image.Rectangle{
		Min: image.Pt(36, 52),
		Max: image.Pt(40, 64),
	}

	armRightFront := image.Rectangle{
		Min: image.Pt(44, 20),
		Max: image.Pt(48, 32),
	}

	bodyFront := image.Rectangle{
		Min: image.Pt(20, 20),
		Max: image.Pt(28, 32),
	}

	minP := s.Dimensions.Width / 8
	maxP := s.Dimensions.Width / 4

	headFront := image.Rectangle{
		Min: image.Pt(minP, minP),
		Max: image.Pt(maxP, maxP),
	}

	headFront = image.Rectangle{
		Min: image.Pt(8, 8),
		Max: image.Pt(16, 16),
	}

	legLeftFrontImage := skinImage.SubImage(legLeftFront)
	legRightFrontImage := skinImage.SubImage(legRightFront)

	armRightFrontImage := skinImage.SubImage(armRightFront)
	armLeftFrontImage := skinImage.SubImage(armLeftFront)

	bodyFrontImage := skinImage.SubImage(bodyFront)
	headFrontImage := skinImage.SubImage(headFront)

	alphaImage := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: armLeftFront.Dx()*2 + bodyFront.Dx(), Y: legLeftFront.Dy() + bodyFront.Dy() + headFront.Dy()}})

	// linker arm
	leftArmY := armLeftFront.Min.Y
	for y := headFront.Dy(); y < headFront.Dy()+armLeftFront.Dy(); y++ {
		leftArmX := armLeftFront.Min.X
		for x := 0; x < armLeftFront.Dx(); x++ {
			alphaImage.Set(x, y, armLeftFrontImage.At(leftArmX, leftArmY))
			leftArmX++
		}
		leftArmY++
	}

	// rechterarm
	rightArmY := armRightFront.Min.Y
	for y := headFront.Dy(); y < headFront.Dy()+armLeftFront.Dy(); y++ {
		rightArmX := armRightFront.Min.X
		for x := armLeftFront.Dx() + bodyFront.Dx(); x < armLeftFront.Dx()*2+bodyFront.Dx(); x++ {
			alphaImage.Set(x, y, armRightFrontImage.At(rightArmX, rightArmY))
			rightArmX++
		}
		rightArmY++
	}

	// linker been
	leftLegY := legLeftFront.Min.Y
	for y := headFront.Dy() + bodyFront.Dy(); y < headFront.Dy()+bodyFront.Dy()+legLeftFront.Dy(); y++ {
		leftLegX := legLeftFront.Min.X
		for x := armLeftFront.Dx(); x < armLeftFront.Dx()+legLeftFront.Dx(); x++ {
			alphaImage.Set(x, y, legLeftFrontImage.At(leftLegX, leftLegY))
			leftLegX++
		}
		leftLegY++
	}

	// rechter been
	rightLegY := legRightFront.Min.Y
	for y := headFront.Dy() + bodyFront.Dy(); y < headFront.Dy()+bodyFront.Dy()+legLeftFront.Dy(); y++ {
		rightLegX := legRightFront.Min.X
		for x := armLeftFront.Dx() + legLeftFront.Dx(); x < armLeftFront.Dx()+legLeftFront.Dx()*2; x++ {
			alphaImage.Set(x, y, legRightFrontImage.At(rightLegX, rightLegY))
			rightLegX++
		}
		rightLegY++
	}

	// head
	headY := headFront.Min.Y
	for y := 0; y < headFront.Dy(); y++ {
		headX := headFront.Min.X
		for x := armLeftFront.Dx(); x < armLeftFront.Dx()+bodyFront.Dx(); x++ {
			alphaImage.Set(x, y, headFrontImage.At(headX, headY))
			headX++
		}
		headY++
	}

	// body
	bodyY := bodyFront.Min.Y
	for y := headFront.Dy(); y < headFront.Dy()+bodyFront.Dy(); y++ {
		bodyX := bodyFront.Min.X
		for x := armLeftFront.Dx(); x < armLeftFront.Dx()+bodyFront.Dx(); x++ {
			alphaImage.Set(x, y, bodyFrontImage.At(bodyX, bodyY))
			bodyX++
		}
		bodyY++
	}

	err = png.Encode(w, alphaImage)
	err = w.Flush()

	return workingDirectory + "/public/images/" + s.Username + "-full_front.png", err
}

func (s *Skin) SaveHeadImage() (string, error) {
	workingDirectory, _ := os.Getwd()
	f, err := os.Create(workingDirectory + "/public/images/" + s.Username + "-head.png")

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

	return workingDirectory + "/public/images/" + s.Username + "-head.png", err
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
