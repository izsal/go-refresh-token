package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// InitLogger initializes the log file with the current date
func InitLogger() *os.File {
	t := time.Now()
	logFileName := filepath.Join("logs", t.Format("2006-01-02")+".log")

	// Ensure the logs directory exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// Create or append to the log file
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	return file
}
