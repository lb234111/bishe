package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

func InitLogger(modName string) *log.Logger {
	file := "./log_0413/" + time.Now().Format("2006010215") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Fatal(err)
	}
	prefix := fmt.Sprintf("[%s]", modName)
	loger := log.New(logFile, prefix, log.LstdFlags|log.Lshortfile|log.LUTC)
	return loger
}
