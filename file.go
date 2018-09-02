package logger

import (
	"fmt"
	"log"
	"os"
)

type LogFile struct {
	Today string
	Files map[string]*os.File
	write chan write
}

type write struct {
	file *os.File
	msg  string
}

func (logFile *LogFile) checkFile(level string) *os.File {
	var fileName string
	if defaultLogger.Service != "" {
		fileName = fmt.Sprintf(`%s/%s.%s.%s.log`, defaultLogger.SaveToDir, defaultLogger.Service, logFile.Today, level)
	} else {
		fileName = fmt.Sprintf(`%s/%s.%s.log`, defaultLogger.SaveToDir, logFile.Today, level)
	}

	// check if file exist. file may not exist due to 1: haven't created  2: been removed when service is running
	if _, err := os.Stat(fileName); err != nil {

		if os.IsNotExist(err) {
			f, err := os.Create(fileName)
			if err != nil {
				log.Println(err)
				return nil
			}
			f.Close()

			f1, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0777)

			if err != nil {
				log.Println(err)
				return nil
			}
			logFile.Files[level] = f1

		} else {
			log.Println(err)
			return nil
		}
	}

	// this value would be true when service restart regardless if file exists
	if logFile.Files[level] == nil {
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0777)

		if err != nil {
			log.Println(err)
			return nil
		}
		logFile.Files[level] = f
	}

	return logFile.Files[level]
}

// need pointer receiver
func (logFile *LogFile) listen() {
	logFile.write = make(chan write)
	go func() {
		for write := range logFile.write {
			if _, err := write.file.WriteString(write.msg); err != nil {
				log.Println(err)
			}
		}
	}()
}

func (logFile *LogFile) writeToFile(level string, str string) {

	if f := logFile.checkFile(level); f != nil {
		logFile.write <- write{file: f, msg: str}
	}
}
