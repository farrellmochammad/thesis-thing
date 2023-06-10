package logger

import (
	"log"
	"os"
	"time"
)

type MyLogger struct {
	file *os.File
}

func (l *MyLogger) Init(filename string) error {
	// Create or open the log file
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	l.file = file
	// Set the log output to the file
	log.SetOutput(l.file)
	return nil
}

func (l *MyLogger) Close() error {
	if l.file != nil {
		err := l.file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *MyLogger) Log(message string) {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00")
	log.Println(timestamp + message)
}
