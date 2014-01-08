package citygrid

import (
	"image/color"
	"math"
)

// HSVModel converts any color.Color to a HSV color.
var HSVModel = color.ModelFunc(hsvModel)

// HSV represents a cylindrical coordinate of points in an RGB color model.
//
// Values are in the range 0 to 1.
type HSV struct {
	H, S, V float64
}

// RGBA returns the alpha-premultiplied red, green, blue and alpha values
// for the HSV.
func (c HSV) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := HSVToRGB(c.H, c.S, c.V)
	return uint32(r) * 0x101, uint32(g) * 0x101, uint32(b) * 0x101, 0xffff
}

// hsvModel converts a color.Color to HSV.
func hsvModel(c color.Color) color.Color {
	if _, ok := c.(HSV); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	h, s, v := RGBToHSV(uint8(r>>8), uint8(g>>8), uint8(b>>8))
	return HSV{h, s, v}
}

// RGBToHSV converts an RGB triple to a HSV triple.
//
// Ported from http://goo.gl/Vg1h9
func RGBToHSV(r, g, b uint8) (h, s, v float64) {
	fR := float64(r) / 255
	fG := float64(g) / 255
	fB := float64(b) / 255
	max := math.Max(math.Max(fR, fG), fB)
	min := math.Min(math.Min(fR, fG), fB)
	d := max - min
	s, v = 0, max
	if max > 0 {
		s = d / max
	}
	if max == min {
		// Achromatic.
		h = 0
	} else {
		// Chromatic.
		switch max {
		case fR:
			h = (fG - fB) / d
			if fG < fB {
				h += 6
			}
		case fG:
			h = (fB-fR)/d + 2
		case fB:
			h = (fR-fG)/d + 4
		}
		h /= 6
	}
	return
}

// HSVToRGB converts an HSV triple to a RGB triple.
//
// Ported from http://goo.gl/Vg1h9
func HSVToRGB(h, s, v float64) (r, g, b uint8) {
	var fR, fG, fB float64
	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1.0 - s)
	q := v * (1.0 - f*s)
	t := v * (1.0 - (1.0-f)*s)
	switch int(i) % 6 {
	case 0:
		fR, fG, fB = v, t, p
	case 1:
		fR, fG, fB = q, v, p
	case 2:
		fR, fG, fB = p, v, t
	case 3:
		fR, fG, fB = p, q, v
	case 4:
		fR, fG, fB = t, p, v
	case 5:
		fR, fG, fB = v, p, q
	}
	r = uint8((fR * 255) + 0.5)
	g = uint8((fG * 255) + 0.5)
	b = uint8((fB * 255) + 0.5)
	return
}

// HSLModel converts any color.Color to a HSL color.
var HSLModel = color.ModelFunc(hslModel)

// HSL represents a cylindrical coordinate of points in an RGB color model.
//
// Values are in the range 0 to 1.
type HSL struct {
	H, S, L float64
}

// RGBA returns the alpha-premultiplied red, green, blue and alpha values
// for the HSL.
func (c HSL) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := HSLToRGB(c.H, c.S, c.L)
	return uint32(r) * 0x101, uint32(g) * 0x101, uint32(b) * 0x101, 0xffff
}

// hslModel converts a color.Color to HSL.
func hslModel(c color.Color) color.Color {
	if _, ok := c.(HSL); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	h, s, l := RGBToHSL(uint8(r>>8), uint8(g>>8), uint8(b>>8))
	return HSL{h, s, l}
}

// RGBToHSL converts an RGB triple to a HSL triple.
//
// Ported from http://goo.gl/Vg1h9
func RGBToHSL(r, g, b uint8) (h, s, l float64) {
	fR := float64(r) / 255
	fG := float64(g) / 255
	fB := float64(b) / 255
	max := math.Max(math.Max(fR, fG), fB)
	min := math.Min(math.Min(fR, fG), fB)
	l = (max + min) / 2
	if max == min {
		// Achromatic.
		h, s = 0, 0
	} else {
		// Chromatic.
		d := max - min
		if l > 0.5 {
			s = d / (2.0 - max - min)
		} else {
			s = d / (max + min)
		}
		switch max {
		case fR:
			h = (fG - fB) / d
			if fG < fB {
				h += 6
			}
		case fG:
			h = (fB-fR)/d + 2
		case fB:
			h = (fR-fG)/d + 4
		}
		h /= 6
	}
	return
}

// HSLToRGB converts an HSL triple to a RGB triple.
//
// Ported from http://goo.gl/Vg1h9
func HSLToRGB(h, s, l float64) (r, g, b uint8) {
	var fR, fG, fB float64
	if s == 0 {
		fR, fG, fB = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - s*l
		}
		p := 2*l - q
		fR = hueToRGB(p, q, h+1.0/3)
		fG = hueToRGB(p, q, h)
		fB = hueToRGB(p, q, h-1.0/3)
	}
	r = uint8((fR * 255) + 0.5)
	g = uint8((fG * 255) + 0.5)
	b = uint8((fB * 255) + 0.5)
	return
}

// hueToRGB is a helper function for HSLToRGB.
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6 {
		return p + (q-p)*6*t
	}
	if t < 0.5 {
		return q
	}
	if t < 2.0/3 {
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}