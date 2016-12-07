package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
)

const (
	Accuracy = 200
)

type Hist []int

// GetHistogram computes a histogram from the image given as path.
// part into picture fragments, and decide color classification of the part.
func GetPartedHistogram(path string) (Hist, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// In current stage, use RGB classification.
	hist := make(Hist, 3, 3)

	rect := img.Bounds()
	// part into <partition> * <partition> fragments
	partition := Accuracy
	// omit fractions (max fraction < <partition>)
	pX := rect.Dx() / partition
	pY := rect.Dy() / partition

	for y := 0; y < partition; y++ {
		for x := 0; x < partition; x++ {
			RGBCheck(img, x*pX, y*pY, pX, pY, hist)
		}
	}

	return hist, nil
}

func GetHistogram(path string) (Hist, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// In current stage, use RGB classification.
	hist := make(Hist, 3, 3)

	rect := img.Bounds()
	for y := 0; y < rect.Dy(); y++ {
		for x := 0; x < rect.Dx(); x++ {
			pr, pg, pb, pa := img.At(x, y).RGBA()
			r, g, b := float2rgb(Unpremultiply(pr, pg, pb, pa))
			ClassifyRGB(r, g, b, hist)
		}
	}

	return hist, nil
}

// sx: start point x
// sy: start point y
// wx: x width
// wy: y width
// m:  array for histogram
func RGBCheck(img image.Image, sx, sy, wx, wy int, h Hist) {
	var r, g, b int
	for y := sy; y < sy+wy; y++ {
		for x := sx; x < sx+wx; x++ {
			pr, pg, pb, pa := img.At(x, y).RGBA()
			tr, tg, tb := float2rgb(Unpremultiply(pr, pg, pb, pa))
			r += tr
			g += tg
			b += tb
		}
	}

	ClassifyRGB(r, g, b, h)
}

// Premultiplied Alpha: Each RGB value is multiplied by alpha.
// Straight Alpha: Each RGB value is raw.
// Premultiplied Alpha -> Straight Alpha (Unpremultiplied)
func Unpremultiply(r, g, b, a uint32) (float64, float64, float64) {
	return float64(r) / float64(a), float64(g) / float64(a), float64(b) / float64(a)
}

// float2rgb converts 0.0 ~ 1.0 range color to 0 ~ 255 range color.
func float2rgb(r, g, b float64) (int, int, int) {
	return int(r * 255), int(g * 255), int(b * 255)
}

func ClassifyRGB(r, g, b int, h Hist) {
	if r > g {
		if r > b {
			h[0]++
		} else {
			h[2]++
		}
	} else {
		if g > b {
			h[1]++
		} else {
			h[2]++
		}
	}
}

func (h Hist) InfoShow() {
	fmt.Printf("R: %d\nG: %d\nB: %d\n", h[0], h[1], h[2])
	fmt.Println("Total Part:", h[0]+h[1]+h[2])
}
