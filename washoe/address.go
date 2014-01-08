package washoe

import ()

// TODO What identifies a unique address?

type Address struct {
	Id        int64   `json:"id"`
	Full      string  `json:"address"`
	Number    int     `json:"number"`
	Fraction  string  `json:"fraction"`
	Prefix    string  `json:"prefix"`
	Street    string  `json:"street"`
	Type      string  `json:"type"`
	Suffix    string  `json:"suffix"`
	Unit      string  `json:"unit"`
	City      string  `json:"city"` // aka Muni
	County    string  `json:"county"`
	State     string  `json:"state"` // Index 10
	Zip       int     `json:"zip"`
	Zip4      int     `json:"zip4"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TODO A create table statement for postgres

func (address *Address) String() string {
	return address.Full
}

func (address *Address) LatLong() (float64, float64) {
	return address.Latitude, address.Longitude
}

// Give each address in the array an id equal to its index + 1
func AutoAddressId(addresses []*Address) {
	for index, address := range addresses {
		address.Id = int64(index + 1)
	}
}
