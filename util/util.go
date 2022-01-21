package util

import (
	"errors"
	"log"
)

var ErrNotFound = errors.New("cannot fetch fetch at this time. please try again")

func Check(err error) {
	if err != nil {
		log.Fatal("Unable to complete request due to error! ", err)
	}
}
