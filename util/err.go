package util

import "log"

func ErrHandler(message string, err error) {
	if err != nil {
		log.Panic(message, err)
	}
}