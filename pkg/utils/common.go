package utils

import (
	"fmt"
	"log"
	"os"
)

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

// ShowTips show tips
func ShowTips() {
	log.Printf("[error] Parse parmas error, Please check your input params!")
	fmt.Printf("Welcome to use, you can type ./sugar-agent -h to show help message." +
		"\nUsages: ./sugar-agent -user guest -password guest -host localhost -port 5672 -exchange-name device_exchange " +
		"-device-id 6\n")
	os.Exit(4)
}
