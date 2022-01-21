package util

import (
	"errors"
	"log"
)

var ErrBadType = errors.New("unable to find api request of that type")
var ErrNotFound = errors.New("no results found for your request")

func Check(err error) {
	if err != nil {
		log.Fatal("Unable to complete request due to error!", err)
	}
}
