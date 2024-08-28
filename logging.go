package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func configureLogging() *os.File {
	logDirPath := "/tmp/xsdm"
	logFileName := "xsdm.log"
	logFilePath := filepath.Join(logDirPath, logFileName)

	err := os.MkdirAll(logDirPath, os.ModePerm)
	if err != nil {
		panicText := fmt.Sprintf("failed to create log directory: %s", err)
		panic(panicText)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panicText := fmt.Sprintf("failed to create log file: %s", err)
		panic(panicText)
	}

	log.SetFlags(log.LstdFlags | log.LUTC)
	log.SetOutput(logFile)
	return logFile
}
