package logger

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Logger represents an active logging object that generates lines ofoutput
type Logger struct {
	Flag      int
	Level     int
	SaveToDir string
	Service   string
	DB        *sql.DB
	LogFile   LogFile
}

const (
	Ltime = 1 << iota
	Lfile
	Lline
	Ltype
	Ldebug
)

const (
	LDebug = 1 << iota
	LInfo
	LWarn
	LError
	LPanic
	LFatal
	LHTTP
)

var defaultLogger = Logger{
	Flag:      Ltime | Lfile | Lline,
	Level:     LDebug | LInfo | LWarn | LError | LPanic | LFatal | LHTTP,
	SaveToDir: "",
	Service:   "",
	DB:        nil,
}

func Default() *Logger {
	return &defaultLogger
}

func SetLogger(flag, level int, db *sql.DB, save, ser string) {
	SetFlags(flag)
	SetLevel(level)
	SetDb(db)
	SetSaveToDir(save)
	SetService(ser)

	// use built-in log when logger not available
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// SetFlags represent developer can customize logger
func SetFlags(flag int) {
	defaultLogger.Flag = flag
}

func SetService(ser string) {
	defaultLogger.Service = ser
}

func SetLevel(level int) {
	defaultLogger.Level = level
}

func SetSaveToDir(dir string) {
	defaultLogger.SaveToDir = dir
	if dir != "" {
		defaultLogger.LogFile.listen()
		defaultLogger.LogFile.Files = make(map[string]*os.File)

		now := time.Now().UTC()
		y, m, d := now.Date()
		defaultLogger.LogFile.Today = fmt.Sprintf("%4d-%02d-%02d", y, int(m), d)

		// find tomorrow date
		y1, m1, d1 := now.Add(24 * time.Hour).UTC().Date()

		// find how much time it takes before tomorrow comes
		tomorrowStart := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
		interval := tomorrowStart.Sub(time.Now())

		time.AfterFunc(interval, func() {
			// after this interval file name changes
			defaultLogger.LogFile.Today = fmt.Sprintf("%4d-%02d-%02d", y1, int(m1), d1)

			// after first time file name changes, it changes every 24 hour
			ticker := time.Tick(24 * time.Hour)
			go func() {
				for t := range ticker {
					y, m, d := t.UTC().Date()
					defaultLogger.LogFile.Today = fmt.Sprintf("%4d-%02d-%02d", y, int(m), d)
				}
			}()
		})
	}
}

func printLog(c *color.Color, log string) {
	c.Println(fmt.Sprintf("%v", log))
}

func logging(level string, flag int, co color.Attribute, msg ...interface{}) {
	logs := genLog(level, msg...)
	if defaultLogger.DB != (nil) {
		saveToDb(defaultLogger.DB, logs)
	}

	if defaultLogger.SaveToDir != "" {
		defaultLogger.LogFile.writeToFile(strings.ToLower(level), logs.FileString())
	}

	if defaultLogger.Level&flag != 0 {
		newColor := color.New(co, color.Bold)
		printLog(newColor, logs.String())
	}
}

func loggingf(level string, flag int, co color.Attribute, format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog(level, s)
	if defaultLogger.DB != (nil) {
		saveToDb(defaultLogger.DB, logs)
	}

	if defaultLogger.SaveToDir != "" {
		defaultLogger.LogFile.writeToFile(strings.ToLower(level), logs.FileString())
	}

	if defaultLogger.Level&flag != 0 {
		newColor := color.New(co, color.Bold)
		printLog(newColor, logs.String())
	}
}

// don't write in this way because skip would be different between logging an logginf
// func loggingf(level string, flag int, co color.Attribute, format string, msg ...interface{}) {
// 	s := fmt.Sprintf(format, msg...)
// 	logging(level, flag, co, s)
// }

func Debug(msg ...interface{}) {
	logging("Debug", LDebug, color.FgHiBlue, msg...)
}

func Debugf(format string, msg ...interface{}) {
	loggingf("Debug", LDebug, color.FgHiBlue, format, msg...)
}

func Info(msg ...interface{}) {
	logging("Info", LInfo, color.FgHiGreen, msg...)
}

func Infof(format string, msg ...interface{}) {
	loggingf("Info", LInfo, color.FgHiGreen, format, msg...)
}

func Warn(msg ...interface{}) {
	logging("Warn", LWarn, color.FgHiMagenta, msg...)
}

func Warnf(format string, msg ...interface{}) {
	loggingf("Warn", LWarn, color.FgHiMagenta, format, msg...)
}

func Error(msg ...interface{}) {
	logging("Error", LError, color.FgHiRed, msg...)
}

func Errorf(format string, msg ...interface{}) {
	loggingf("Error", LError, color.FgHiRed, format, msg...)
}

func Panic(msg ...interface{}) {
	logging("Panic", LPanic, color.FgHiCyan, msg...)
	panic(msg)
}

func Panicf(format string, msg ...interface{}) {
	loggingf("Panic", LPanic, color.FgHiCyan, format, msg...)
	panic(msg)
}

func Fatal(msg ...interface{}) {
	msgs := append(msg, "os.Exit(1)")
	logging("Fatal", LFatal, color.FgHiYellow, msgs...)
	panic(msg)
}

func Fatalf(format string, msg ...interface{}) {
	msgs := append(msg, "os.Exit(1)")
	loggingf("Fatal", LFatal, color.FgHiYellow, format, msgs...)
	os.Exit(1)
}

func HTTP(msg ...interface{}) {
	logging("HTTP", LHTTP, color.FgHiCyan, msg...)
}

func HTTPf(format string, msg ...interface{}) {
	loggingf("HTTP", LHTTP, color.FgHiCyan, format, msg...)
}
