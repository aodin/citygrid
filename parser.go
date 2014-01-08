package citygrid

import (
	"os"
)

// Given a file of addresses, return a list of locations
type LocationParser interface {
	Parse(*os.File) ([]Location, error)
}