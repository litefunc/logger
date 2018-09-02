package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type logs struct {
	Service string
	Ltime   string
	Lfile   string
	Lline   int
	Level   string
	Msg     string
}

const skip = 4

func genLtime() string {
	t := time.Now().UTC().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s", t)
}

func genLfile() string {
	_, file, _, _ := runtime.Caller(skip)
	pwd, _ := os.Getwd()
	file = strings.Replace(file, pwd+"/", "", -1)
	return file
}

func genLline() int {
	_, _, line, _ := runtime.Caller(skip)
	return line
}

func genMsg(msg []interface{}) string {
	var msgs []string
	for _, v := range msg {
		msgs = append(msgs, fmt.Sprintf("%+v", v))
	}
	return strings.Join(msgs, " ")
}

func genLog(level string, msg ...interface{}) logs {
	var log logs

	log.Service = defaultLogger.Service

	if defaultLogger.Flag&Ltime != 0 {
		log.Ltime = genLtime()
	}
	if defaultLogger.Flag&Lfile != 0 {
		log.Lfile = genLfile()
	}
	if defaultLogger.Flag&Lline != 0 {
		log.Lline = genLline()
	}
	log.Level = level
	log.Msg = genMsg(msg)

	return log
}

func (log logs) String() string {
	var msgs []string

	if log.Ltime != "" {
		msgs = append(msgs, log.Ltime)
	}
	if log.Service != "" {
		msgs = append(msgs, log.Service)
	}
	if log.Lfile != "" {
		msgs = append(msgs, fmt.Sprintf("%s:", log.Lfile))
	}
	if log.Lline != 0 {
		msgs = append(msgs, fmt.Sprintf("line:%d", log.Lline))
	}
	if log.Level != "" {
		msgs = append(msgs, fmt.Sprintf("| %s |", log.Level))
	}
	if log.Msg != "" {
		msgs = append(msgs, log.Msg)
	}

	return strings.Join(msgs, " ") + "\n"
}

func (log logs) FileString() string {
	var msgs []string

	if log.Ltime != "" {
		msgs = append(msgs, log.Ltime)
	}
	if log.Lfile != "" {
		msgs = append(msgs, fmt.Sprintf("%s:", log.Lfile))
	}
	if log.Lline != 0 {
		msgs = append(msgs, fmt.Sprintf("line:%d", log.Lline))
	}
	if log.Msg != "" {
		msgs = append(msgs, log.Msg)
	}

	return strings.Join(msgs, " ") + "\n"
}
