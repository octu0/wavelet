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

	lowR := make([][]float32, height)
	lowG := make([][]float32, height)
	lowB := make([][]float32, height)
	lowA := make([][]float32, height)

	intermediateR := make([][]float32, height)
	intermediateG := make([][]float32, height)
	intermediateB := make([][]float32, height)
	intermediateA := make([][]float32, height)
	for h := 0; h < height; h += 1 {
		r, g, b, a := make([]float32, width), make([]float32, width), make([]float32, width), make([]float32, width)
		for w := 0; w < width; w += 1 {
			c := testImg.RGBAAt(w, h)
			r[w] = float32(c.R)
			g[w] = float32(c.G)
			b[w] = float32(c.B)
			a[w] = float32(c.A)
		}
		intermediateR[h], lowR[h] = wavelet.Wavelet(r)
		intermediateG[h], lowG[h] = wavelet.Wavelet(g)
		intermediateB[h], lowB[h] = wavelet.Wavelet(b)
		intermediateA[h], lowA[h] = wavelet.Wavelet(a)
	}

	intermidate := image.NewRGBA(image.Rect(0, 0, width/2, height/2))

	for h := 0; h < (height / 2); h += 1 {
		for w := 0; w < (width / 2); w += 1 {
			intermidate.SetRGBA(w, h, color.RGBA{
				R: byte(wavelet.Clamp(intermediateR[h*2][w], 0, 255)),
				G: byte(wavelet.Clamp(intermediateG[h*2][w], 0, 255)),
				B: byte(wavelet.Clamp(intermediateB[h*2][w], 0, 255)),
				A: byte(wavelet.Clamp(intermediateA[h*2][w], 0, 255)),
			})
		}
	}

	path1, err := saveImage(intermidate)
	if err != nil {
		panic(err)
	}
	println("intermidate", path1)

	inverse := image.NewRGBA(image.Rect(0, 0, width, height))
	for h := 0; h < height; h += 1 {
		r := wavelet.Inverse(intermediateR[h], lowR[h])
		g := wavelet.Inverse(intermediateG[h], lowG[h])
		b := wavelet.Inverse(intermediateB[h], lowB[h])
		a := wavelet.Inverse(intermediateA[h], lowA[h])
		for w := 0; w < width; w += 1 {
			inverse.SetRGBA(w, h, color.RGBA{
				R: byte(r[w]),
				G: byte(g[w]),
				B: byte(b[w]),
				A: byte(a[w]),
			})
		}
	}

	path2, err := saveImage(inverse)
	if err != nil {
		panic(err)
	}
	println("inverse", path2)
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
