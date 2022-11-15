package main

import (
	"log"
	"os"
)

var logFile, _ = os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
var logger = log.New(logFile, "log: ", log.Ldate|log.Ltime)

func writeLog(message string) {
	logger.Fatal(message)
}
