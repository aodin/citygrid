package denver

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"github.com/aodin/citygrid"
)

var UnparsableAddress = errors.New("addresses: unparsable")
var UnexpectedLength = errors.New("addresses: unexpected length")

func ParseAddress(raw []string) (address *Address, err error) {
	if len(raw) != 21 {
		return nil, UnexpectedLength
	}

	// Strings
	address = &Address{
		Full:       raw[20],
		PreDir:     raw[9],
		Street:     raw[10],
		PostType:   raw[11],
		PostDir:    raw[12],
		UnitType:   raw[16],
		UnitNumber: raw[17],
	}
	// Convert the numbers
	address.Number, err = strconv.Atoi(raw[6])
	if err != nil {
		return nil, err
	}
	// TODO Apt numbers OR letters
	// if raw[17] != "" {
	// 	address.UnitNumber, err = strconv.Atoi(raw[17])
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// Convert the latitude and longitude
	address.Latitude, err = strconv.ParseFloat(raw[3], 64)
	if err != nil {
		return nil, err
	}
	address.Longitude, err = strconv.ParseFloat(raw[4], 64)
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
		return nil, err
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
		if address.Latitude != 0 && address.Longitude != 0 {
			addresses = append(addresses, address)
		}
	}
	return
}
