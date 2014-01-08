package washoe

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"github.com/aodin/citygrid"
)

var UnparsableAddress = errors.New("Could not parse the given address")
var UnexpectedLength = errors.New("Unexpected Length")

func ParseAddress(raw []string) (address *Address, err error) {
	if len(raw) != 15 {
		return nil, UnexpectedLength
	}
	address = &Address{
		Full:     raw[0],
		Fraction: raw[2],
		Prefix:   raw[3],
		Street:   raw[4],
		Type:     raw[5],
		Suffix:   raw[6],
		Unit:     raw[7],
		City:     raw[8],
		County:   raw[9],
		State:    raw[10],
	}
	// Convert the number
	address.Number, err = strconv.Atoi(raw[1])
	if err != nil {
		return nil, err
	}
	// Convert the ZIP codes
	if len(raw[11]) > 0 {
		address.Zip, err = strconv.Atoi(raw[11])
		if err != nil {
			return nil, err
		}
	}
	if len(raw[12]) > 0 {
		address.Zip4, err = strconv.Atoi(raw[12])
		if err != nil {
			return nil, err
		}
	}
	// Convert the latitude and longitude
	address.Latitude, err = strconv.ParseFloat(raw[13], 64)
	if err != nil {
		return nil, err
	}
	address.Longitude, err = strconv.ParseFloat(raw[14], 64)
	if err != nil {
		return nil, err
	}
	// TODO Full return statement is not needed, but added for clarity
	return address, nil
}

func ParseAddresses(path string) (addresses []citygrid.Location, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)

	// Skip the header
	_, err = r.Read()
	if err != nil {
		return
	}
	// TODO Re-examine the csv.ReadAll() source
	for {
		line, err := r.Read()
		if err == io.EOF {
			return addresses, nil
		}
		if err != nil {
			return nil, err
		}
		address, err := ParseAddress(line)
		if err != nil {
			return nil, err
		}
		if address.Latitude >= 1 && address.Longitude <= 0 {
			addresses = append(addresses, address)
		}
	}
}
