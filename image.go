package citygrid

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

func Image(width, height, pixel int, freq Frequency) *image.RGBA {
	var maxL float64 = 0.9

	// Create a new blank image
	m := image.NewRGBA(image.Rect(0, 0, width * pixel, height * pixel))

	// Fill it with a single uniform color
	// black := color.RGBA{0, 0, 0, 255}
	// draw.Draw(m, m.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	var ix, iy int
	var v float64

	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 1 {

			v = freq[y * width + x]

			// Take the square root to blunt extremes
			if v != 0 {
				v = math.Sqrt(v)
			}

			// Square it to make the extremes more pronounced
			// v = v * v

			l := 1.0 - (v * maxL)
			// Base color HSL (0.5944, 0.92, l)
			// Convert the HSL value to rgb and switch the pixel
			// log.Println("Lightness:", freq[(y * 256) + x] / max, l)
			r, g, b := HSLToRGB(0.5944, 0.92, l)
			fg := color.RGBA{r, g, b, 255}

			ix = x * pixel
			// Invert the image vertically
			iy = (height - y - 1) * pixel
			block := image.Rect(ix, iy, ix + pixel, iy + pixel)
			draw.Draw(m, block, &image.Uniform{fg}, image.ZP, draw.Src)
		}
	}
	return m
}