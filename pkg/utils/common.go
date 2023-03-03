package utils

import "log"

// FailOnError check error
// err: error
// msg: error message
// return: none
func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("Error happened, %s: %s", msg, err)
	}
}

// LogOnError check error
// err: error
// msg: error message
// return: none
func LogOnError(err error, msg string) {
	if err != nil {
		log.Printf("Error happened, %s: %s", msg, err)
	}
}
