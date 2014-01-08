package citygrid

import (
	"testing"
)

func TestLocation(t *testing.T) {
	wallyLand := Simple{33.809, -117.919}
	wallyWorld := Simple{28.418611, -81.581111}
	locations := []Location{wallyLand, wallyWorld}
	
	// Calculate the range of the lat and longs
	// TODO Is it faster to sort? 
	extrema := LocationExtrema(locations)
	if extrema.minX != wallyLand.Longitude {
		t.Errorf("Unexpected minimum longitude")
	}
	if extrema.minY != wallyWorld.Latitude {
		t.Errorf("Unexpected minimum latitude")
	}

	histogram := AutoHeightHistogram(512, extrema)
	histogram.CountLocations(locations)
	CreateMaxFrequency(histogram)
	CreateFrequency(histogram, 3)
}
