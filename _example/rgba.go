package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/octu0/wavelet"

	_ "embed"
)

//go:embed src.png
var testImgData []byte

var (
	testImg *image.RGBA
)

func init() {
	img, err := pngToRGBA(testImgData)
	if err != nil {
		panic(err.Error())
	}
	testImg = img
}

func main() {
	width, height := testImg.Rect.Dx(), testImg.Rect.Dy()

	highXR := make([][]float32, height)
	highXG := make([][]float32, height)
	highXB := make([][]float32, height)
	highXA := make([][]float32, height)
	lowXR := make([][]float32, height)
	lowXG := make([][]float32, height)
	lowXB := make([][]float32, height)
	lowXA := make([][]float32, height)

	for h := 0; h < height; h += 1 {
		r, g, b, a := make([]float32, width), make([]float32, width), make([]float32, width), make([]float32, width)
		for w := 0; w < width; w += 1 {
			c := testImg.RGBAAt(w, h)
			r[w] = float32(c.R)
			g[w] = float32(c.G)
			b[w] = float32(c.B)
			a[w] = float32(c.A)
		}
		highXR[h], lowXR[h] = wavelet.Haar(r)
		highXG[h], lowXG[h] = wavelet.Haar(g)
		highXB[h], lowXB[h] = wavelet.Haar(b)
		highXA[h], lowXA[h] = wavelet.Haar(a)
	}

	highYR := make([][]float32, width)
	highYG := make([][]float32, width)
	highYB := make([][]float32, width)
	highYA := make([][]float32, width)
	lowYR := make([][]float32, width)
	lowYG := make([][]float32, width)
	lowYB := make([][]float32, width)
	lowYA := make([][]float32, width)

	for w := 0; w < width; w += 1 {
		r, g, b, a := make([]float32, height), make([]float32, height), make([]float32, height), make([]float32, height)
		for h := 0; h < height; h += 1 {
			c := testImg.RGBAAt(w, h)
			r[h] = float32(c.R)
			g[h] = float32(c.G)
			b[h] = float32(c.B)
			a[h] = float32(c.A)
		}
		highYR[w], lowYR[w] = wavelet.Haar(r)
		highYG[w], lowYG[w] = wavelet.Haar(g)
		highYB[w], lowYB[w] = wavelet.Haar(b)
		highYA[w], lowYA[w] = wavelet.Haar(a)
	}

	intermidate := image.NewRGBA(image.Rect(0, 0, width/2, height/2))

	for h := 0; h < (height / 2); h += 1 {
		for w := 0; w < (width / 2); w += 1 {
			r := byte(wavelet.Clamp((highYR[w*2][h]+highXR[h*2][w])/2, 0, 255))
			g := byte(wavelet.Clamp((highYG[w*2][h]+highXG[h*2][w])/2, 0, 255))
			b := byte(wavelet.Clamp((highYB[w*2][h]+highXB[h*2][w])/2, 0, 255))
			a := byte(wavelet.Clamp((highYA[w*2][h]+highXA[h*2][w])/2, 0, 255))
			intermidate.SetRGBA(w, h, color.RGBA{
				R: r,
				G: g,
				B: b,
				A: a,
			})
		}
	}

	path1, err := saveImage(intermidate)
	if err != nil {
		panic(err)
	}
	println("intermidate", path1)

	inverse := image.NewRGBA(image.Rect(0, 0, width, height))
	for w := 0; w < width; w += 1 {
		yr := wavelet.InverseHaar(highYR[w], lowYR[w])
		yg := wavelet.InverseHaar(highYG[w], lowYG[w])
		yb := wavelet.InverseHaar(highYB[w], lowYB[w])
		ya := wavelet.InverseHaar(highYA[w], lowYA[w])
		for h := 0; h < height; h += 1 {
			xr := wavelet.InverseHaar(highXR[h], lowXR[h])
			xg := wavelet.InverseHaar(highXG[h], lowXG[h])
			xb := wavelet.InverseHaar(highXB[h], lowXB[h])
			xa := wavelet.InverseHaar(highXA[h], lowXA[h])

			inverse.SetRGBA(w, h, color.RGBA{
				R: byte((yr[h] + xr[w]) / 2),
				G: byte((yg[h] + xg[w]) / 2),
				B: byte((yb[h] + xb[w]) / 2),
				A: byte((ya[h] + xa[w]) / 2),
			})
		}
	}

	path2, err := saveImage(inverse)
	if err != nil {
		panic(err)
	}
	println("inverse", path2)

	ratio := float32(0.85)
	compress := image.NewRGBA(image.Rect(0, 0, width, height))

	for h := 0; h < height; h += 1 {
		wavelet.Threshold(highXR[h], ratio)
		wavelet.Threshold(highXG[h], ratio)
		wavelet.Threshold(highXB[h], ratio)
	}
	for w := 0; w < width; w += 1 {
		wavelet.Threshold(highYR[w], ratio)
		wavelet.Threshold(highYG[w], ratio)
		wavelet.Threshold(highYB[w], ratio)
	}
	for w := 0; w < width; w += 1 {
		yr := wavelet.InverseHaar(highYR[w], lowYR[w])
		yg := wavelet.InverseHaar(highYG[w], lowYG[w])
		yb := wavelet.InverseHaar(highYB[w], lowYB[w])
		ya := wavelet.InverseHaar(highYA[w], lowYA[w])
		for h := 0; h < height; h += 1 {
			xr := wavelet.InverseHaar(highXR[h], lowXR[h])
			xg := wavelet.InverseHaar(highXG[h], lowXG[h])
			xb := wavelet.InverseHaar(highXB[h], lowXB[h])
			xa := wavelet.InverseHaar(highXA[h], lowXA[h])

			compress.SetRGBA(w, h, color.RGBA{
				R: byte(wavelet.Clamp((yr[h]+xr[w])/2, 0, 255)),
				G: byte(wavelet.Clamp((yg[h]+xg[w])/2, 0, 255)),
				B: byte(wavelet.Clamp((yb[h]+xb[w])/2, 0, 255)),
				A: byte(wavelet.Clamp((ya[h]+xa[w])/2, 0, 255)),
			})
		}
	}

	path3, err := saveImage(compress)
	if err != nil {
		panic(err)
	}
	println("compress", path3)
}

func pngToRGBA(data []byte) (*image.RGBA, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if i, ok := img.(*image.RGBA); ok {
		return i, nil
	}

	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y += 1 {
		for x := b.Min.X; x < b.Max.X; x += 1 {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			rgba.Set(x, y, c)
		}
	}
	return rgba, nil
}

func saveImage(img *image.RGBA) (string, error) {
	out, err := os.CreateTemp("/tmp", "out*.png")
	if err != nil {
		return "", err
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		return "", err
	}
	return out.Name(), nil
}
