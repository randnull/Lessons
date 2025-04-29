package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
)

func InitLogger() error {
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		err := os.Mkdir("./logs", 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile("./logs/users.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

func Info(msg string) {
	infoLogger.Println(msg)
}

func Error(msg string) {
	errorLogger.Println(msg)
}

func Debug(msg string) {
	debugLogger.Println(msg)
}
