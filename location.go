package citygrid

import (
	"log"
	"sort"
)

// All we need from addresses is their latitude and longitude
type Location interface {
	LatLong() (float64, float64)
}

// A simple implementation of the location interface
type Simple struct {
	Latitude  float64
	Longitude float64
}

func (s Simple) LatLong() (float64, float64) {
	return s.Latitude, s.Longitude
}

type Rect struct {
	minX float64
	minY float64
	maxX float64
	maxY float64
}

func Extrema(minX, minY, maxX, maxY float64) (rect Rect) {
	rect.minX = minX
	rect.minY = minY
	rect.maxX = maxX
	rect.maxY = maxY
	return
}

func LocationExtrema(locations []Location) (rect Rect) {
	var lat, long float64
	for i, location := range locations {
		// Always set the first location
		lat, long = location.LatLong()
		// Remember: lat is a y coordinate, long is an x!
		if i == 0 {
			rect.minX = long
			rect.maxX = long
			rect.minY = lat
			rect.maxY = lat
		} else {
			if long < rect.minX {
				rect.minX = long
			} else if long > rect.maxX {
				rect.maxX = long
			}
			if lat < rect.minY {
				rect.minY = lat
			} else if lat > rect.maxY {
				rect.maxY = lat
			}
		}
	}
	return
}

func LocationExtrema999(locations []Location) (rect Rect) {
	var lat, long float64
	var lats, longs []float64

	for _, location := range locations {
		// Always set the first location
		lat, long = location.LatLong()
		lats = append(lats, lat)
		longs = append(longs, long)
	}

	sort.Float64s(lats)
	sort.Float64s(longs)

	maxLatIndex := int(0.999 * float64(len(lats)))
	minLatIndex := len(lats) - maxLatIndex

	maxLongIndex := int(0.999 * float64(len(longs)))
	minLongIndex := len(longs) - maxLongIndex

	rect.minX = longs[minLongIndex]
	rect.maxX = longs[maxLongIndex]
	rect.minY = lats[minLatIndex]
	rect.maxY = lats[maxLatIndex]
	return
}

func LocationExtrema9999(locations []Location) (rect Rect) {
	var lat, long float64
	var lats, longs []float64

	for _, location := range locations {
		// Always set the first location
		lat, long = location.LatLong()
		lats = append(lats, lat)
		longs = append(longs, long)
	}

	sort.Float64s(lats)
	sort.Float64s(longs)

	maxLatIndex := int(0.9999 * float64(len(lats)))
	minLatIndex := len(lats) - maxLatIndex

	maxLongIndex := int(0.9999 * float64(len(longs)))
	minLongIndex := len(longs) - maxLongIndex

	rect.minX = longs[minLongIndex]
	rect.maxX = longs[maxLongIndex]
	rect.minY = lats[minLatIndex]
	rect.maxY = lats[maxLatIndex]
	return
}

type Histogram struct {
	Width  int
	Height int
	Blocks [][]int
	minX   float64
	minY   float64
	rangeX float64
	rangeY float64
	ratio  float64
}

// Uses the ratio of latitude and longitude at 45 deg N/S
func AutoHeightHistogram(width int, extrema Rect) *Histogram {
	// TODO These values must be non-zero

	histogram := &Histogram{}

	// Determine how many pixels should be used for the height
	histogram.rangeX = extrema.maxX - extrema.minX
	histogram.rangeY = extrema.maxY - extrema.minY

	histogram.minX = extrema.minX
	histogram.minY = extrema.minY

	ratio := histogram.rangeY / histogram.rangeX

	// All y coordinates will be multipled by this ratio
	// Since a single latitude span represents a span that is 41% larger than
	// a longitude span, we need 41% more pixels than the ratio would suggest
	histogram.ratio = 1.40944

	histogram.Width = width
	// TODO Rounding should be performed
	histogram.Height = int(float64(width) * ratio * histogram.ratio)

	histogram.Blocks = make([][]int, histogram.Height)
	for i, _ := range histogram.Blocks {
		histogram.Blocks[i] = make([]int, histogram.Width)
	}
	return histogram
}

func (h *Histogram) CountLocations(locations []Location) {
	var lat, long float64
	var x, y int

	width64 := float64(h.Width)
	height64 := float64(h.Height)

	for _, location := range locations {
		lat, long = location.LatLong()

		// Remember, lat is y, long is x
		x = int(((long - h.minX) / h.rangeX) * width64)

		// There will be some cropping of the y coordinate
		y = int(((lat - h.minY) / h.rangeY) * height64 * h.ratio)

		if x >= h.Width {
			// x = h.Width - 1
			continue
		} else if x < 0 {
			// x = 0
			continue
		}
		if y >= h.Height {
			// y = h.Height - 1
			continue
		} else if y < 0 {
			// y = 0
			continue
		}
		h.Blocks[y][x] += 1
	}
}

type Frequency []float64

func CreateMaxFrequency(h *Histogram) Frequency {
	return createFrequency(h, h.Width * h.Height - 1)
}

func CreateFrequency(h *Histogram, index int) Frequency {
	return createFrequency(h, h.Width * h.Height - 1 - index)
}

func createFrequency(h *Histogram, index int) (freq Frequency) {
	// Get the maximum, or optionally, use
	var max, v float64
	length := h.Height * h.Width
	values := make([]int, length)

	for r, row := range h.Blocks {
		for c, count := range row {
			values[(r * h.Width) + c] = count
		}
	}
	sort.Ints(values[0:])
	max = float64(values[index])
	log.Println("Mode:", max)

	freq = make([]float64, length)
	for r, row := range h.Blocks {
		for c, count := range row {
			v = float64(count) / max
			if v > 1.0 {
				v = 1.0
			}
			freq[(r * h.Width) + c] = v
		}
	}
	return
}
