package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"

	"github.com/octu0/wavelet"
	"github.com/pkg/errors"

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

	lowXR := make([][]float32, height)
	lowXG := make([][]float32, height)
	lowXB := make([][]float32, height)
	lowXA := make([][]float32, height)
	highXR := make([][]float32, height)
	highXG := make([][]float32, height)
	highXB := make([][]float32, height)
	highXA := make([][]float32, height)

	for h := 0; h < height; h += 1 {
		r, g, b, a := make([]float32, width), make([]float32, width), make([]float32, width), make([]float32, width)
		for w := 0; w < width; w += 1 {
			c := testImg.RGBAAt(w, h)
			r[w] = float32(c.R)
			g[w] = float32(c.G)
			b[w] = float32(c.B)
			a[w] = float32(c.A)
		}
		lowXR[h], highXR[h] = wavelet.Haar(r)
		lowXG[h], highXG[h] = wavelet.Haar(g)
		lowXB[h], highXB[h] = wavelet.Haar(b)
		lowXA[h], highXA[h] = wavelet.Haar(a)
	}

	lowYR := make([][]float32, width)
	lowYG := make([][]float32, width)
	lowYB := make([][]float32, width)
	lowYA := make([][]float32, width)
	highYR := make([][]float32, width)
	highYG := make([][]float32, width)
	highYB := make([][]float32, width)
	highYA := make([][]float32, width)

	for w := 0; w < width; w += 1 {
		r, g, b, a := make([]float32, height), make([]float32, height), make([]float32, height), make([]float32, height)
		for h := 0; h < height; h += 1 {
			c := testImg.RGBAAt(w, h)
			r[h] = float32(c.R)
			g[h] = float32(c.G)
			b[h] = float32(c.B)
			a[h] = float32(c.A)
		}
		lowYR[w], highYR[w] = wavelet.Haar(r)
		lowYG[w], highYG[w] = wavelet.Haar(g)
		lowYB[w], highYB[w] = wavelet.Haar(b)
		lowYA[w], highYA[w] = wavelet.Haar(a)
	}

	intermidate := image.NewRGBA(image.Rect(0, 0, width/2, height/2))

	for h := 0; h < (height / 2); h += 1 {
		for w := 0; w < (width / 2); w += 1 {
			r := byte(wavelet.Clamp((lowYR[w*2][h]+lowXR[h*2][w])/2, 0, 255))
			g := byte(wavelet.Clamp((lowYG[w*2][h]+lowXG[h*2][w])/2, 0, 255))
			b := byte(wavelet.Clamp((lowYB[w*2][h]+lowXB[h*2][w])/2, 0, 255))
			a := byte(wavelet.Clamp((lowYA[w*2][h]+lowXA[h*2][w])/2, 0, 255))
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
		yr := wavelet.InverseHaar(lowYR[w], highYR[w])
		yg := wavelet.InverseHaar(lowYG[w], highYG[w])
		yb := wavelet.InverseHaar(lowYB[w], highYB[w])
		ya := wavelet.InverseHaar(lowYA[w], highYA[w])
		for h := 0; h < height; h += 1 {
			xr := wavelet.InverseHaar(lowXR[h], highXR[h])
			xg := wavelet.InverseHaar(lowXG[h], highXG[h])
			xb := wavelet.InverseHaar(lowXB[h], highXB[h])
			xa := wavelet.InverseHaar(lowXA[h], highXA[h])

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
	substract := image.NewRGBA(image.Rect(0, 0, width, height))

	for h := 0; h < height; h += 1 {
		wavelet.Threshold(lowXR[h], ratio)
		wavelet.Threshold(lowXG[h], ratio)
		wavelet.Threshold(lowXB[h], ratio)
	}
	for w := 0; w < width; w += 1 {
		wavelet.Threshold(lowYR[w], ratio)
		wavelet.Threshold(lowYG[w], ratio)
		wavelet.Threshold(lowYB[w], ratio)
	}
	for w := 0; w < width; w += 1 {
		yr := wavelet.InverseHaar(lowYR[w], highYR[w])
		yg := wavelet.InverseHaar(lowYG[w], highYG[w])
		yb := wavelet.InverseHaar(lowYB[w], highYB[w])
		ya := wavelet.InverseHaar(lowYA[w], highYA[w])
		for h := 0; h < height; h += 1 {
			xr := wavelet.InverseHaar(lowXR[h], highXR[h])
			xg := wavelet.InverseHaar(lowXG[h], highXG[h])
			xb := wavelet.InverseHaar(lowXB[h], highXB[h])
			xa := wavelet.InverseHaar(lowXA[h], highXA[h])

			substract.SetRGBA(w, h, color.RGBA{
				R: byte(wavelet.Clamp((yr[h]+xr[w])/2, 0, 255)),
				G: byte(wavelet.Clamp((yg[h]+xg[w])/2, 0, 255)),
				B: byte(wavelet.Clamp((yb[h]+xb[w])/2, 0, 255)),
				A: byte(wavelet.Clamp((ya[h]+xa[w])/2, 0, 255)),
			})
		}
	}

	path3, err := saveImage(substract)
	if err != nil {
		panic(err)
	}
	println("substract", path3)

	dataR := make([]byte, height*width)
	dataG := make([]byte, height*width)
	dataB := make([]byte, height*width)
	dataA := make([]byte, height*width)
	i := 0
	for h := 0; h < (height / 2); h += 1 {
		for w := 0; w < (width / 2); w += 1 {
			r := byte(wavelet.Clamp((lowYR[w*2][h]+lowXR[h*2][w])/2, 0, 255))
			g := byte(wavelet.Clamp((lowYG[w*2][h]+lowXG[h*2][w])/2, 0, 255))
			b := byte(wavelet.Clamp((lowYB[w*2][h]+lowXB[h*2][w])/2, 0, 255))
			a := byte(wavelet.Clamp((lowYA[w*2][h]+lowXA[h*2][w])/2, 0, 255))
			dataR[i] = r
			dataG[i] = g
			dataB[i] = b
			dataA[i] = a
			i += 1
		}
	}

	bufR := bytes.NewBuffer(nil)
	bufG := bytes.NewBuffer(nil)
	bufB := bytes.NewBuffer(nil)
	bufA := bytes.NewBuffer(nil)

	r := &doubleRunlength{}
	sizeR, err := r.encode(bufR, dataR)
	if err != nil {
		panic(err)
	}
	sizeG, err := r.encode(bufG, dataG)
	if err != nil {
		panic(err)
	}
	sizeB, err := r.encode(bufB, dataB)
	if err != nil {
		panic(err)
	}
	sizeA, err := r.encode(bufA, dataA)
	if err != nil {
		panic(err)
	}

	out := bytes.NewBuffer(nil)
	if err := binary.Write(out, binary.BigEndian, uint64(width)); err != nil {
		panic(err)
	}
	if err := binary.Write(out, binary.BigEndian, uint64(height)); err != nil {
		panic(err)
	}
	if err := binary.Write(out, binary.BigEndian, uint64(sizeR)); err != nil {
		panic(err)
	}
	out.Write(bufR.Bytes())
	if err := binary.Write(out, binary.BigEndian, uint64(sizeG)); err != nil {
		panic(err)
	}
	out.Write(bufG.Bytes())
	if err := binary.Write(out, binary.BigEndian, uint64(sizeB)); err != nil {
		panic(err)
	}
	out.Write(bufB.Bytes())
	if err := binary.Write(out, binary.BigEndian, uint64(sizeA)); err != nil {
		panic(err)
	}
	out.Write(bufA.Bytes())

	fmt.Printf(
		"compressed %dKB\nby png %3.4f%% \nby rgba %3.4f%% \n",
		out.Len()/1024,
		(1-(float64(out.Len())/float64(len(testImgData))))*100,
		(1-(float64(out.Len())/float64(width*height*4)))*100,
	)

	encoded := bytes.NewReader(out.Bytes())
	eWidth := uint64(0)
	if err := binary.Read(encoded, binary.BigEndian, &eWidth); err != nil {
		panic(err)
	}
	eHeight := uint64(0)
	if err := binary.Read(encoded, binary.BigEndian, &eHeight); err != nil {
		panic(err)
	}
	eSizeR := uint64(0)
	if err := binary.Read(encoded, binary.BigEndian, &eSizeR); err != nil {
		panic(err)
	}
	encodedR := make([]byte, eSizeR)
	if _, err := encoded.Read(encodedR); err != nil {
		panic(err)
	}
	eSizeG := uint64(0)
	if err := binary.Read(encoded, binary.BigEndian, &eSizeG); err != nil {
		panic(err)
	}
	encodedG := make([]byte, eSizeG)
	if _, err := encoded.Read(encodedG); err != nil {
		panic(err)
	}
	eSizeB := uint64(0)
	if err := binary.Read(encoded, binary.BigEndian, &eSizeB); err != nil {
		panic(err)
	}
	encodedB := make([]byte, eSizeB)
	if _, err := encoded.Read(encodedB); err != nil {
		panic(err)
	}
	eSizeA := uint64(0)
	if err := binary.Read(encoded, binary.BigEndian, &eSizeA); err != nil {
		panic(err)
	}
	encodedA := make([]byte, eSizeA)
	if _, err := encoded.Read(encodedA); err != nil {
		panic(err)
	}

	decodedBufR := bytes.NewBuffer(nil)
	decodedBufG := bytes.NewBuffer(nil)
	decodedBufB := bytes.NewBuffer(nil)
	decodedBufA := bytes.NewBuffer(nil)
	if err := r.decode(decodedBufR, bytes.NewReader(encodedR)); err != nil {
		panic(err)
	}
	if err := r.decode(decodedBufG, bytes.NewReader(encodedG)); err != nil {
		panic(err)
	}
	if err := r.decode(decodedBufB, bytes.NewReader(encodedB)); err != nil {
		panic(err)
	}
	if err := r.decode(decodedBufA, bytes.NewReader(encodedA)); err != nil {
		panic(err)
	}

	dR := decodedBufR.Bytes()
	dG := decodedBufG.Bytes()
	dB := decodedBufB.Bytes()
	dA := decodedBufA.Bytes()

	decoded := image.NewRGBA(image.Rect(0, 0, int(eWidth), int(eHeight)))
	pos := 0
	for h := 0; h < int(eHeight); h += 2 {
		for w := 0; w < int(eWidth); w += 1 {
			decoded.SetRGBA(w, h, color.RGBA{
				R: dR[pos],
				G: dG[pos],
				B: dB[pos],
				A: dA[pos],
			})
			decoded.SetRGBA(w, h+1, color.RGBA{
				R: dR[pos],
				G: dG[pos],
				B: dB[pos],
				A: dA[pos],
			})
			pos += 1
		}
	}

	path4, err := saveImage(decoded)
	if err != nil {
		panic(err)
	}
	println("decoded", path4)
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

type doubleRunlength struct {
}

func (r *doubleRunlength) encode(out io.Writer, data []byte) (int64, error) {
	currentValue := data[0]
	currentLength := byte(1)
	written := int64(0)

	for i := 1; i < len(data); i += 1 {
		flush := true
		if currentValue == data[i] {
			flush = false
			currentLength += 1
			if 127 <= currentLength {
				flush = true
			}
		}

		if flush {
			size, err := out.Write([]byte{currentLength * 2, currentValue})
			if err != nil {
				return 0, errors.Wrapf(err, "failed to write data:%d %v", currentLength, currentValue)
			}
			currentLength = 1
			currentValue = data[i]
			written += int64(size)
		}
	}
	size, err := out.Write([]byte{currentLength * 2, currentValue})
	if err != nil {
		return 0, errors.Wrapf(err, "failed to write data:%d %v", currentLength, currentValue)
	}
	written += int64(size)
	return written, nil
}

func (r *doubleRunlength) decode(out io.Writer, in io.Reader) error {
	buf := make([]byte, 2)
	for {
		_, err := in.Read(buf[0:2])
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return errors.WithStack(err)
		}

		length := buf[0]
		values := make([]byte, length)
		for i := byte(0); i < length; i += 1 {
			values[i] = buf[1]
		}
		if _, err := out.Write(values); err != nil {
			return errors.Wrapf(err, "failed to decoded value")
		}
	}
	return nil
}
