package washoe

import (
	"testing"
)

var testFile = "./washoe_example.csv"

func TestParseAddresses(t *testing.T) {
	addresses, parseErr := ParseAddresses(testFile)
	if parseErr != nil {
		t.Fatal(parseErr)
	}
	if len(addresses) != 3 {
		t.Fatalf("Unexpected length of addresses: %d", len(addresses))
	}

	lat, long := addresses[0].LatLong()
	if lat != 39.632746 {
		t.Errorf("Unexpected latitude: %s", lat)
	}
	if long != -119.717986 {
		t.Errorf("Unexpected longitude: %s", long)
	}
}
